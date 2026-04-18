-- 007_create_notifications_settings_logs.sql
CREATE TABLE IF NOT EXISTS notifications (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       VARCHAR(50) NOT NULL,
    title      TEXT        NOT NULL,
    body       TEXT        NOT NULL DEFAULT '',
    is_read    BOOLEAN     NOT NULL DEFAULT FALSE,
    metadata   JSONB       NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id    ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read    ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);

CREATE TABLE IF NOT EXISTS user_settings (
    id                 UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id            UUID        NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    notification_prefs JSONB       NOT NULL DEFAULT '{"email":true,"push":true,"chapter_release":true,"friend_activity":true}',
    theme              VARCHAR(20) NOT NULL DEFAULT 'dark',
    language           VARCHAR(10) NOT NULL DEFAULT 'en',
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS server_logs (
    id         BIGSERIAL   PRIMARY KEY,
    level      VARCHAR(10) NOT NULL,
    message    TEXT        NOT NULL,
    context    JSONB       NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_server_logs_level      ON server_logs(level);
CREATE INDEX IF NOT EXISTS idx_server_logs_created_at ON server_logs(created_at DESC);

CREATE TABLE IF NOT EXISTS backups (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename   TEXT        NOT NULL,
    size_bytes BIGINT      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_backups_user_id ON backups(user_id);
