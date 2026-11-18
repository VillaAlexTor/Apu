package identity

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (store *Store) PlatformSummary(ctx context.Context) (PlatformSummary, error) {
	var summary PlatformSummary
	err := store.pool.QueryRow(ctx, `
		SELECT
			(SELECT count(*) FROM tenants),
			(SELECT count(*) FROM tenants WHERE status = 'active'),
			(SELECT count(*) FROM users WHERE status = 'active'),
			(SELECT count(*) FROM users WHERE status = 'invited'),
			(SELECT count(*) FROM users WHERE status = 'suspended'),
			(SELECT count(*) FROM user_invitations WHERE accepted_at IS NULL AND revoked_at IS NULL AND expires_at > now())
	`).Scan(
		&summary.TenantsTotal,
		&summary.TenantsActive,
		&summary.UsersActive,
		&summary.UsersInvited,
		&summary.UsersSuspended,
		&summary.PendingInvitations,
	)
	if err != nil {
		return PlatformSummary{}, fmt.Errorf("query platform summary: %w", err)
	}
	return summary, nil
}

func (store *Store) PlatformTenants(ctx context.Context) ([]TenantRecord, error) {
	rows, err := store.pool.Query(ctx, `
		SELECT tenants.id::text,
		       tenants.slug::text,
		       tenants.name,
		       tenants.status,
		       owner_user.display_name,
		       owner_user.email::text,
		       count(DISTINCT tenant_users.user_id)::int,
		       count(DISTINCT sessions.id)::int,
		       tenants.created_at
		FROM tenants
		LEFT JOIN tenant_users ON tenant_users.tenant_id = tenants.id
		LEFT JOIN users ON users.id = tenant_users.user_id
		LEFT JOIN sessions ON sessions.user_id = users.id
		  AND sessions.revoked_at IS NULL
		  AND sessions.expires_at > now()
		LEFT JOIN LATERAL (
			SELECT owner.display_name, owner.email
			FROM tenant_users owner_membership
			JOIN users owner ON owner.id = owner_membership.user_id
			WHERE owner_membership.tenant_id = tenants.id
			  AND owner_membership.role = 'tenant_owner'
			ORDER BY owner_membership.created_at
			LIMIT 1
		) owner_user ON true
		GROUP BY tenants.id, owner_user.display_name, owner_user.email
		ORDER BY tenants.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query platform tenants: %w", err)
	}
	defer rows.Close()

	records := make([]TenantRecord, 0)
	for rows.Next() {
		var record TenantRecord
		if err := rows.Scan(
			&record.ID,
			&record.Slug,
			&record.Name,
			&record.Status,
			&record.OwnerName,
			&record.OwnerEmail,
			&record.UserCount,
			&record.ActiveSessions,
			&record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan platform tenant: %w", err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate platform tenants: %w", err)
	}
	return records, nil
}

func (store *Store) CreateTenant(ctx context.Context, actor User, name, slug string) (TenantRecord, error) {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return TenantRecord{}, fmt.Errorf("begin create tenant: %w", err)
	}
	defer tx.Rollback(ctx)

	var tenant TenantRecord
	err = tx.QueryRow(ctx, `
		INSERT INTO tenants (name, slug)
		VALUES ($1, $2)
		RETURNING id::text, slug::text, name, status, created_at
	`, name, slug).Scan(&tenant.ID, &tenant.Slug, &tenant.Name, &tenant.Status, &tenant.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return TenantRecord{}, ErrConflict
		}
		return TenantRecord{}, fmt.Errorf("insert tenant: %w", err)
	}

	metadata, _ := json.Marshal(map[string]string{"name": name, "slug": slug})
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, tenant_id, action, resource_type, resource_id, metadata)
		VALUES ($1, $2, 'platform.tenant.created', 'tenant', $2, $3)
	`, actor.ID, tenant.ID, string(metadata)); err != nil {
		return TenantRecord{}, fmt.Errorf("audit tenant creation: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return TenantRecord{}, fmt.Errorf("commit tenant creation: %w", err)
	}
	return tenant, nil
}

func (store *Store) SetTenantStatus(ctx context.Context, actor User, tenantID, status string) error {
	return store.updateStatusWithAudit(ctx, actor, "tenants", tenantID, status, "platform.tenant.status_changed", tenantID)
}

func (store *Store) PlatformUsers(ctx context.Context) ([]PlatformUser, error) {
	rows, err := store.pool.Query(ctx, `
		SELECT users.id::text,
		       users.email::text,
		       users.display_name,
		       users.status,
		       users.platform_role,
		       users.created_at,
		       count(DISTINCT sessions.id)::int
		FROM users
		LEFT JOIN sessions ON sessions.user_id = users.id
		  AND sessions.revoked_at IS NULL
		  AND sessions.expires_at > now()
		GROUP BY users.id
		ORDER BY users.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query platform users: %w", err)
	}
	defer rows.Close()

	users := make([]PlatformUser, 0)
	for rows.Next() {
		var user PlatformUser
		if err := rows.Scan(&user.ID, &user.Email, &user.DisplayName, &user.Status, &user.PlatformRole, &user.CreatedAt, &user.ActiveSessions); err != nil {
			return nil, fmt.Errorf("scan platform user: %w", err)
		}
		user.Memberships, err = store.Memberships(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate platform users: %w", err)
	}
	return users, nil
}

func (store *Store) CreateInvitation(
	ctx context.Context,
	actor User,
	tenantID, email, displayName, role string,
	tokenHash []byte,
	expiresAt time.Time,
) (InvitationResult, error) {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return InvitationResult{}, fmt.Errorf("begin invitation: %w", err)
	}
	defer tx.Rollback(ctx)

	var tenantStatus string
	if err := tx.QueryRow(ctx, `SELECT status FROM tenants WHERE id = $1 FOR UPDATE`, tenantID).Scan(&tenantStatus); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InvitationResult{}, ErrNotFound
		}
		return InvitationResult{}, fmt.Errorf("load invitation tenant: %w", err)
	}
	if tenantStatus == "deactivated" {
		return InvitationResult{}, ErrConflict
	}

	var user PlatformUser
	var passwordHash *string
	err = tx.QueryRow(ctx, `
		SELECT id::text, email::text, display_name, password_hash, status, platform_role, created_at
		FROM users WHERE email = $1 FOR UPDATE
	`, email).Scan(&user.ID, &user.Email, &user.DisplayName, &passwordHash, &user.Status, &user.PlatformRole, &user.CreatedAt)
	requiresActivation := false
	if errors.Is(err, pgx.ErrNoRows) {
		err = tx.QueryRow(ctx, `
			INSERT INTO users (email, display_name, password_hash, status)
			VALUES ($1, $2, NULL, 'invited')
			RETURNING id::text, email::text, display_name, status, platform_role, created_at
		`, email, displayName).Scan(&user.ID, &user.Email, &user.DisplayName, &user.Status, &user.PlatformRole, &user.CreatedAt)
		requiresActivation = true
	} else if err == nil {
		requiresActivation = passwordHash == nil
		if user.Status == "deactivated" {
			return InvitationResult{}, ErrConflict
		}
	} else {
		return InvitationResult{}, fmt.Errorf("load invitation user: %w", err)
	}
	if err != nil {
		return InvitationResult{}, fmt.Errorf("insert invitation user: %w", err)
	}

	membershipStatus := "active"
	if requiresActivation {
		membershipStatus = "invited"
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO tenant_users (tenant_id, user_id, role, status)
		VALUES ($1, $2, $3, $4)
	`, tenantID, user.ID, role, membershipStatus); err != nil {
		if isUniqueViolation(err) {
			return InvitationResult{}, ErrConflict
		}
		return InvitationResult{}, fmt.Errorf("insert invitation membership: %w", err)
	}

	if requiresActivation {
		if _, err := tx.Exec(ctx, `
			UPDATE user_invitations
			SET revoked_at = now()
			WHERE user_id = $1 AND accepted_at IS NULL AND revoked_at IS NULL
		`, user.ID); err != nil {
			return InvitationResult{}, fmt.Errorf("revoke prior invitation: %w", err)
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO user_invitations (user_id, tenant_id, token_hash, created_by, expires_at)
			VALUES ($1, $2, $3, $4, $5)
		`, user.ID, tenantID, tokenHash, actor.ID, expiresAt); err != nil {
			return InvitationResult{}, fmt.Errorf("insert invitation: %w", err)
		}
	}

	metadata, _ := json.Marshal(map[string]string{"email": email, "role": role})
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, tenant_id, action, resource_type, resource_id, metadata)
		VALUES ($1, $2, 'platform.user.invited', 'user', $3, $4)
	`, actor.ID, tenantID, user.ID, string(metadata)); err != nil {
		return InvitationResult{}, fmt.Errorf("audit invitation: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return InvitationResult{}, fmt.Errorf("commit invitation: %w", err)
	}
	return InvitationResult{User: user, RequiresActivation: requiresActivation}, nil
}

func (store *Store) ActivateInvitation(ctx context.Context, token, passwordHash string) error {
	tokenDigest := sha256.Sum256([]byte(token))
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin invitation activation: %w", err)
	}
	defer tx.Rollback(ctx)

	var invitationID, userID, tenantID string
	err = tx.QueryRow(ctx, `
		SELECT id::text, user_id::text, tenant_id::text
		FROM user_invitations
		WHERE token_hash = $1
		  AND accepted_at IS NULL
		  AND revoked_at IS NULL
		  AND expires_at > now()
		FOR UPDATE
	`, tokenDigest[:]).Scan(&invitationID, &userID, &tenantID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("resolve invitation: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE users
		SET password_hash = $2, status = 'active', password_changed_at = now(), updated_at = now()
		WHERE id = $1
	`, userID, passwordHash); err != nil {
		return fmt.Errorf("activate invited user: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		UPDATE tenant_users SET status = 'active', updated_at = now()
		WHERE user_id = $1 AND status = 'invited'
	`, userID); err != nil {
		return fmt.Errorf("activate invited memberships: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE user_invitations SET accepted_at = now() WHERE id = $1`, invitationID); err != nil {
		return fmt.Errorf("complete invitation: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, tenant_id, action, resource_type, resource_id)
		VALUES ($1, $2, 'identity.invitation.accepted', 'user', $1)
	`, userID, tenantID); err != nil {
		return fmt.Errorf("audit invitation activation: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit invitation activation: %w", err)
	}
	return nil
}

func (store *Store) SetUserStatus(ctx context.Context, actor User, userID, status string) error {
	if actor.ID == userID && status != "active" {
		return ErrConflict
	}
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin user status change: %w", err)
	}
	defer tx.Rollback(ctx)

	command, err := tx.Exec(ctx, `
		UPDATE users
		SET status = $2, updated_at = now()
		WHERE id = $1
		  AND ($2 <> 'active' OR password_hash IS NOT NULL)
	`, userID, status)
	if err != nil {
		return fmt.Errorf("update user status: %w", err)
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	if status != "active" {
		if _, err := tx.Exec(ctx, `UPDATE sessions SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`, userID); err != nil {
			return fmt.Errorf("revoke suspended user sessions: %w", err)
		}
	}
	metadata, _ := json.Marshal(map[string]string{"status": status})
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id, metadata)
		VALUES ($1, 'platform.user.status_changed', 'user', $2, $3)
	`, actor.ID, userID, string(metadata)); err != nil {
		return fmt.Errorf("audit user status: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit user status: %w", err)
	}
	return nil
}

func (store *Store) RevokeUserSessions(ctx context.Context, actor User, userID string) (int64, error) {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin session revocation: %w", err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, userID).Scan(&exists); err != nil {
		return 0, fmt.Errorf("check session user: %w", err)
	}
	if !exists {
		return 0, ErrNotFound
	}
	command, err := tx.Exec(ctx, `UPDATE sessions SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`, userID)
	if err != nil {
		return 0, fmt.Errorf("revoke user sessions: %w", err)
	}
	count := command.RowsAffected()
	metadata, _ := json.Marshal(map[string]int64{"sessions": count})
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id, metadata)
		VALUES ($1, 'platform.user.sessions_revoked', 'user', $2, $3)
	`, actor.ID, userID, string(metadata)); err != nil {
		return 0, fmt.Errorf("audit session revocation: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit session revocation: %w", err)
	}
	return count, nil
}

func (store *Store) SetMembershipRole(ctx context.Context, actor User, tenantID, userID, role string) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin membership role change: %w", err)
	}
	defer tx.Rollback(ctx)
	command, err := tx.Exec(ctx, `
		UPDATE tenant_users SET role = $3, updated_at = now()
		WHERE tenant_id = $1 AND user_id = $2
	`, tenantID, userID, role)
	if err != nil {
		return fmt.Errorf("update membership role: %w", err)
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	metadata, _ := json.Marshal(map[string]string{"role": role})
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, tenant_id, action, resource_type, resource_id, metadata)
		VALUES ($1, $2, 'platform.membership.role_changed', 'user', $3, $4)
	`, actor.ID, tenantID, userID, string(metadata)); err != nil {
		return fmt.Errorf("audit membership role: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit membership role: %w", err)
	}
	return nil
}

func (store *Store) PlatformAudit(ctx context.Context, limit int) ([]AuditEvent, error) {
	rows, err := store.pool.Query(ctx, `
		SELECT audit_log.id,
		       users.email::text,
		       tenants.name,
		       audit_log.action,
		       audit_log.resource_type,
		       audit_log.resource_id::text,
		       audit_log.metadata,
		       audit_log.occurred_at
		FROM audit_log
		LEFT JOIN users ON users.id = audit_log.actor_user_id
		LEFT JOIN tenants ON tenants.id = audit_log.tenant_id
		ORDER BY audit_log.occurred_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("query platform audit: %w", err)
	}
	defer rows.Close()
	events := make([]AuditEvent, 0)
	for rows.Next() {
		var event AuditEvent
		if err := rows.Scan(
			&event.ID,
			&event.ActorEmail,
			&event.TenantName,
			&event.Action,
			&event.ResourceType,
			&event.ResourceID,
			&event.Metadata,
			&event.OccurredAt,
		); err != nil {
			return nil, fmt.Errorf("scan audit event: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit events: %w", err)
	}
	return events, nil
}

func (store *Store) updateStatusWithAudit(
	ctx context.Context,
	actor User,
	table, resourceID, status, action, tenantID string,
) error {
	if table != "tenants" {
		return errors.New("unsupported status table")
	}
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin status update: %w", err)
	}
	defer tx.Rollback(ctx)
	command, err := tx.Exec(ctx, `UPDATE tenants SET status = $2, updated_at = now() WHERE id = $1`, resourceID, status)
	if err != nil {
		return fmt.Errorf("update tenant status: %w", err)
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	metadata, _ := json.Marshal(map[string]string{"status": status})
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, tenant_id, action, resource_type, resource_id, metadata)
		VALUES ($1, $2, $3, 'tenant', $2, $4)
	`, actor.ID, tenantID, action, string(metadata)); err != nil {
		return fmt.Errorf("audit tenant status: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit status update: %w", err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var postgresError interface{ SQLState() string }
	return errors.As(err, &postgresError) && postgresError.SQLState() == "23505"
}
