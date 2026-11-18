package identity

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("identity not found")
	ErrEmailConflict = errors.New("email is already registered")
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (store *Store) FindUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := store.pool.QueryRow(ctx, `
		SELECT id::text, email::text, display_name, password_hash, status, platform_role
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.Status,
		&user.PlatformRole,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("find user by email: %w", err)
	}
	return user, nil
}

func (store *Store) Memberships(ctx context.Context, userID string) ([]Membership, error) {
	rows, err := store.pool.Query(ctx, `
		SELECT tenant_id::text, tenants.name, tenant_users.role
		FROM tenant_users
		JOIN tenants ON tenants.id = tenant_users.tenant_id
		WHERE tenant_users.user_id = $1
		  AND tenant_users.status = 'active'
		  AND tenants.status = 'active'
		ORDER BY tenants.name
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("query memberships: %w", err)
	}
	defer rows.Close()

	memberships, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Membership])
	if err != nil {
		return nil, fmt.Errorf("collect memberships: %w", err)
	}
	return memberships, nil
}

func (store *Store) CreateSession(
	ctx context.Context,
	user User,
	tokenHash []byte,
	csrfHash []byte,
	expiresAt time.Time,
) (Session, error) {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return Session{}, fmt.Errorf("begin session transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var session Session
	session.User = user
	session.CSRFHash = csrfHash
	session.ExpiresAt = expiresAt
	err = tx.QueryRow(ctx, `
		INSERT INTO sessions (user_id, token_hash, csrf_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text
	`, user.ID, tokenHash, csrfHash, expiresAt).Scan(&session.ID)
	if err != nil {
		return Session{}, fmt.Errorf("insert session: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id)
		VALUES ($1, 'auth.login.succeeded', 'session', $2)
	`, user.ID, session.ID); err != nil {
		return Session{}, fmt.Errorf("audit successful login: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Session{}, fmt.Errorf("commit session transaction: %w", err)
	}
	return session, nil
}

func (store *Store) SessionByToken(ctx context.Context, token string) (Session, error) {
	tokenHash := sha256.Sum256([]byte(token))
	var session Session
	err := store.pool.QueryRow(ctx, `
		UPDATE sessions
		SET last_seen_at = now()
		FROM users
		WHERE sessions.token_hash = $1
		  AND sessions.user_id = users.id
		  AND sessions.revoked_at IS NULL
		  AND sessions.expires_at > now()
		  AND users.status = 'active'
		RETURNING sessions.id::text,
		          sessions.csrf_hash,
		          sessions.expires_at,
		          users.id::text,
		          users.email::text,
		          users.display_name,
		          users.password_hash,
		          users.status,
		          users.platform_role
	`, tokenHash[:]).Scan(
		&session.ID,
		&session.CSRFHash,
		&session.ExpiresAt,
		&session.User.ID,
		&session.User.Email,
		&session.User.DisplayName,
		&session.User.PasswordHash,
		&session.User.Status,
		&session.User.PlatformRole,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Session{}, ErrNotFound
	}
	if err != nil {
		return Session{}, fmt.Errorf("resolve session: %w", err)
	}

	session.Memberships, err = store.Memberships(ctx, session.User.ID)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

func (store *Store) RevokeSession(ctx context.Context, session Session) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin revoke transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE sessions
		SET revoked_at = now()
		WHERE id = $1 AND revoked_at IS NULL
	`, session.ID); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id)
		VALUES ($1, 'auth.logout', 'session', $2)
	`, session.User.ID, session.ID); err != nil {
		return fmt.Errorf("audit logout: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit revoke transaction: %w", err)
	}
	return nil
}

func (store *Store) CreatePlatformAdmin(ctx context.Context, email, displayName, passwordHash string) (User, error) {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return User{}, fmt.Errorf("begin admin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var user User
	err = tx.QueryRow(ctx, `
		INSERT INTO users (email, display_name, password_hash, platform_role)
		VALUES ($1, $2, $3, 'platform_admin')
		RETURNING id::text, email::text, display_name, password_hash, status, platform_role
	`, email, displayName, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.Status,
		&user.PlatformRole,
	)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			return User{}, ErrEmailConflict
		}
		return User{}, fmt.Errorf("insert platform admin: %w", err)
	}

	metadata, err := json.Marshal(map[string]string{"source": "apu-admin"})
	if err != nil {
		return User{}, fmt.Errorf("encode bootstrap audit metadata: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id, metadata)
		VALUES ($1, 'identity.platform_admin.created', 'user', $1, $2)
	`, user.ID, string(metadata)); err != nil {
		return User{}, fmt.Errorf("audit platform admin creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, fmt.Errorf("commit admin transaction: %w", err)
	}
	return user, nil
}
