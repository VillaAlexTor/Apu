-- +goose Up
ALTER TABLE users
    ADD COLUMN mfa_secret_encrypted BYTEA,
    ADD COLUMN mfa_enabled_at TIMESTAMPTZ,
    ADD COLUMN mfa_last_used_step BIGINT;

ALTER TABLE sessions
    ADD COLUMN mfa_verified_at TIMESTAMPTZ;

CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    token_hash BYTEA NOT NULL UNIQUE CHECK (octet_length(token_hash) = 32),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (expires_at > created_at)
);

CREATE INDEX password_reset_tokens_active_idx
    ON password_reset_tokens (user_id, expires_at)
    WHERE used_at IS NULL;

CREATE TABLE auth_challenges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    purpose TEXT NOT NULL CHECK (purpose IN ('mfa_enroll', 'mfa_verify')),
    token_hash BYTEA NOT NULL UNIQUE CHECK (octet_length(token_hash) = 32),
    pending_secret_encrypted BYTEA,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    attempt_count SMALLINT NOT NULL DEFAULT 0 CHECK (attempt_count BETWEEN 0 AND 8),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (expires_at > created_at)
);

CREATE INDEX auth_challenges_active_idx
    ON auth_challenges (user_id, purpose, expires_at)
    WHERE consumed_at IS NULL;

CREATE TABLE mfa_recovery_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    code_hash BYTEA NOT NULL UNIQUE CHECK (octet_length(code_hash) = 32),
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX mfa_recovery_codes_active_idx
    ON mfa_recovery_codes (user_id)
    WHERE used_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS mfa_recovery_codes;
DROP TABLE IF EXISTS auth_challenges;
DROP TABLE IF EXISTS password_reset_tokens;
ALTER TABLE sessions DROP COLUMN IF EXISTS mfa_verified_at;
ALTER TABLE users
    DROP COLUMN IF EXISTS mfa_last_used_step,
    DROP COLUMN IF EXISTS mfa_enabled_at,
    DROP COLUMN IF EXISTS mfa_secret_encrypted;
