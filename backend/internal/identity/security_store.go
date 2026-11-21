package identity

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (store *Store) CreateAuthChallenge(ctx context.Context, userID, purpose string, tokenHash []byte, expiresAt time.Time) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin auth challenge: %w", err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `
		UPDATE auth_challenges SET consumed_at = now()
		WHERE user_id = $1 AND purpose = $2 AND consumed_at IS NULL
	`, userID, purpose); err != nil {
		return fmt.Errorf("expire previous auth challenge: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO auth_challenges (user_id, purpose, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`, userID, purpose, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("insert auth challenge: %w", err)
	}
	return tx.Commit(ctx)
}

func (store *Store) AuthChallengeByToken(ctx context.Context, token, purpose string) (AuthChallenge, error) {
	digest := sha256.Sum256([]byte(token))
	var challenge AuthChallenge
	err := store.pool.QueryRow(ctx, `
		SELECT challenge.id::text, challenge.purpose, challenge.pending_secret_encrypted,
		       challenge.attempt_count, challenge.expires_at,
		       users.id::text, users.email::text, users.display_name,
		       COALESCE(users.password_hash, ''), users.status, users.platform_role,
		       (users.mfa_enabled_at IS NOT NULL), users.mfa_secret_encrypted
		FROM auth_challenges AS challenge
		JOIN users ON users.id = challenge.user_id
		WHERE challenge.token_hash = $1
		  AND challenge.purpose = $2
		  AND challenge.consumed_at IS NULL
		  AND challenge.expires_at > now()
		  AND challenge.attempt_count < 8
		  AND users.status = 'active'
	`, digest[:], purpose).Scan(
		&challenge.ID,
		&challenge.Purpose,
		&challenge.PendingSecretEncrypted,
		&challenge.AttemptCount,
		&challenge.ExpiresAt,
		&challenge.User.ID,
		&challenge.User.Email,
		&challenge.User.DisplayName,
		&challenge.User.PasswordHash,
		&challenge.User.Status,
		&challenge.User.PlatformRole,
		&challenge.User.MFAEnabled,
		&challenge.MFASecretEncrypted,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return AuthChallenge{}, ErrNotFound
	}
	if err != nil {
		return AuthChallenge{}, fmt.Errorf("resolve auth challenge: %w", err)
	}
	return challenge, nil
}

func (store *Store) SetAuthChallengeSecret(ctx context.Context, challengeID string, encrypted []byte) error {
	command, err := store.pool.Exec(ctx, `
		UPDATE auth_challenges SET pending_secret_encrypted = $2
		WHERE id = $1 AND consumed_at IS NULL AND expires_at > now() AND purpose = 'mfa_enroll'
	`, challengeID, encrypted)
	if err != nil {
		return fmt.Errorf("set challenge secret: %w", err)
	}
	if command.RowsAffected() != 1 {
		return ErrNotFound
	}
	return nil
}

func (store *Store) RecordAuthChallengeFailure(ctx context.Context, challengeID string) error {
	_, err := store.pool.Exec(ctx, `
		UPDATE auth_challenges
		SET attempt_count = LEAST(attempt_count + 1, 8)
		WHERE id = $1 AND consumed_at IS NULL
	`, challengeID)
	if err != nil {
		return fmt.Errorf("record challenge failure: %w", err)
	}
	return nil
}

func (store *Store) ConsumeAuthChallenge(ctx context.Context, challengeID string) error {
	command, err := store.pool.Exec(ctx, `
		UPDATE auth_challenges SET consumed_at = now()
		WHERE id = $1 AND consumed_at IS NULL AND expires_at > now() AND attempt_count < 8
	`, challengeID)
	if err != nil {
		return fmt.Errorf("consume auth challenge: %w", err)
	}
	if command.RowsAffected() != 1 {
		return ErrNotFound
	}
	return nil
}

func (store *Store) EnableMFA(ctx context.Context, userID string, encryptedSecret []byte, lastUsedStep int64, recoveryHashes [][]byte) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin MFA enrollment: %w", err)
	}
	defer tx.Rollback(ctx)
	command, err := tx.Exec(ctx, `
		UPDATE users SET mfa_secret_encrypted = $2, mfa_enabled_at = now(), mfa_last_used_step = $3, updated_at = now()
		WHERE id = $1 AND mfa_enabled_at IS NULL
	`, userID, encryptedSecret, lastUsedStep)
	if err != nil {
		return fmt.Errorf("enable MFA: %w", err)
	}
	if command.RowsAffected() != 1 {
		return ErrConflict
	}
	if _, err := tx.Exec(ctx, `DELETE FROM mfa_recovery_codes WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("clear MFA recovery codes: %w", err)
	}
	for _, hash := range recoveryHashes {
		if _, err := tx.Exec(ctx, `INSERT INTO mfa_recovery_codes (user_id, code_hash) VALUES ($1, $2)`, userID, hash); err != nil {
			return fmt.Errorf("insert MFA recovery code: %w", err)
		}
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id)
		VALUES ($1, 'auth.mfa.enabled', 'user', $1)
	`, userID); err != nil {
		return fmt.Errorf("audit MFA enrollment: %w", err)
	}
	return tx.Commit(ctx)
}

func (store *Store) ConsumeRecoveryCode(ctx context.Context, userID string, codeHash []byte) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin recovery code consumption: %w", err)
	}
	defer tx.Rollback(ctx)
	command, err := tx.Exec(ctx, `
		UPDATE mfa_recovery_codes SET used_at = now()
		WHERE user_id = $1 AND code_hash = $2 AND used_at IS NULL
	`, userID, codeHash)
	if err != nil {
		return fmt.Errorf("consume recovery code: %w", err)
	}
	if command.RowsAffected() != 1 {
		return ErrNotFound
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id)
		VALUES ($1, 'auth.mfa.recovery_code.used', 'user', $1)
	`, userID); err != nil {
		return fmt.Errorf("audit recovery code: %w", err)
	}
	return tx.Commit(ctx)
}

func (store *Store) ConsumeTOTP(ctx context.Context, userID string, step int64) error {
	command, err := store.pool.Exec(ctx, `
		UPDATE users SET mfa_last_used_step = $2, updated_at = now()
		WHERE id = $1 AND mfa_enabled_at IS NOT NULL
		  AND (mfa_last_used_step IS NULL OR mfa_last_used_step < $2)
	`, userID, step)
	if err != nil {
		return fmt.Errorf("consume TOTP step: %w", err)
	}
	if command.RowsAffected() != 1 {
		return ErrConflict
	}
	return nil
}

func (store *Store) CreatePasswordReset(ctx context.Context, userID string, tokenHash []byte, expiresAt time.Time) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin password reset: %w", err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `
		UPDATE password_reset_tokens SET used_at = now()
		WHERE user_id = $1 AND used_at IS NULL
	`, userID); err != nil {
		return fmt.Errorf("invalidate password resets: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("insert password reset: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id)
		VALUES ($1, 'auth.password_reset.requested', 'user', $1)
	`, userID); err != nil {
		return fmt.Errorf("audit password reset request: %w", err)
	}
	return tx.Commit(ctx)
}

func (store *Store) ResetPassword(ctx context.Context, token, passwordHash string) error {
	digest := sha256.Sum256([]byte(token))
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin password reset completion: %w", err)
	}
	defer tx.Rollback(ctx)
	var resetID, userID string
	err = tx.QueryRow(ctx, `
		SELECT id::text, user_id::text
		FROM password_reset_tokens
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > now()
		FOR UPDATE
	`, digest[:]).Scan(&resetID, &userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("resolve password reset: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE password_reset_tokens SET used_at = now() WHERE id = $1`, resetID); err != nil {
		return fmt.Errorf("consume password reset: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		UPDATE users SET password_hash = $2, password_changed_at = now(), updated_at = now()
		WHERE id = $1 AND status = 'active'
	`, userID, passwordHash); err != nil {
		return fmt.Errorf("update reset password: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE sessions SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`, userID); err != nil {
		return fmt.Errorf("revoke reset sessions: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE auth_challenges SET consumed_at = now() WHERE user_id = $1 AND consumed_at IS NULL`, userID); err != nil {
		return fmt.Errorf("revoke reset challenges: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id)
		VALUES ($1, 'auth.password_reset.completed', 'user', $1)
	`, userID); err != nil {
		return fmt.Errorf("audit password reset: %w", err)
	}
	return tx.Commit(ctx)
}

func (store *Store) ChangePassword(ctx context.Context, userID, passwordHash string) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin password change: %w", err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `
		UPDATE users SET password_hash = $2, password_changed_at = now(), updated_at = now()
		WHERE id = $1 AND status = 'active'
	`, userID, passwordHash); err != nil {
		return fmt.Errorf("change password: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE sessions SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`, userID); err != nil {
		return fmt.Errorf("revoke changed password sessions: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO audit_log (actor_user_id, action, resource_type, resource_id)
		VALUES ($1, 'auth.password.changed', 'user', $1)
	`, userID); err != nil {
		return fmt.Errorf("audit password change: %w", err)
	}
	return tx.Commit(ctx)
}
