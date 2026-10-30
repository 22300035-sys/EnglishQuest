-- EnglishQuest initial schema: users and refresh tokens.
-- Forward-only initial migration. Run explicitly after reviewing the target.

BEGIN;

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(254) NOT NULL UNIQUE CHECK (email = lower(email)),
    password_hash  VARCHAR(255) NOT NULL,
    first_name      VARCHAR(100) NOT NULL DEFAULT '',
    last_name       VARCHAR(100) NOT NULL DEFAULT '',
    cefr_level      VARCHAR(4) NOT NULL DEFAULT 'A1' CHECK (cefr_level IN ('A1','A2','B1','B2','C1','C2')),
    role            VARCHAR(20) NOT NULL DEFAULT 'STUDENT' CHECK (role IN ('STUDENT','ADMIN')),
    timezone        VARCHAR(100) NOT NULL DEFAULT 'UTC',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_expiry ON refresh_tokens (expires_at);
COMMIT;
