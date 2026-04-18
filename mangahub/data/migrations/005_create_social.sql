-- 005_create_social.sql
CREATE TABLE IF NOT EXISTS reviews (
    id         UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    manga_id   UUID         NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    rating     NUMERIC(4,2) NOT NULL CHECK (rating >= 0 AND rating <= 10),
    content    TEXT         NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, manga_id)
);

CREATE INDEX IF NOT EXISTS idx_reviews_manga_id ON reviews(manga_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user_id  ON reviews(user_id);

CREATE TABLE IF NOT EXISTS friends (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id  UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status     VARCHAR(20) NOT NULL DEFAULT 'pending'
                   CHECK (status IN ('pending','accepted','blocked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, friend_id)
);

CREATE INDEX IF NOT EXISTS idx_friends_user_id   ON friends(user_id);
CREATE INDEX IF NOT EXISTS idx_friends_friend_id ON friends(friend_id);
CREATE INDEX IF NOT EXISTS idx_friends_status    ON friends(status);

CREATE TABLE IF NOT EXISTS activity_feed (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action_type VARCHAR(50) NOT NULL,
    manga_id    UUID        REFERENCES manga(id) ON DELETE SET NULL,
    metadata    JSONB       NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_activity_user_id    ON activity_feed(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_created_at ON activity_feed(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_manga_id   ON activity_feed(manga_id);
