-- 004_create_library_and_history.sql
CREATE TABLE IF NOT EXISTS library_entries (
    id              UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    manga_id        UUID        NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    status          VARCHAR(20) NOT NULL DEFAULT 'plan_to_read'
                        CHECK (status IN ('reading','completed','plan_to_read','on_hold','dropped')),
    current_chapter INT         NOT NULL DEFAULT 0,
    current_volume  INT         NOT NULL DEFAULT 0,
    rating          NUMERIC(4,2),
    notes           TEXT        NOT NULL DEFAULT '',
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, manga_id)
);

CREATE INDEX IF NOT EXISTS idx_library_user_id  ON library_entries(user_id);
CREATE INDEX IF NOT EXISTS idx_library_manga_id ON library_entries(manga_id);
CREATE INDEX IF NOT EXISTS idx_library_status   ON library_entries(status);
CREATE INDEX IF NOT EXISTS idx_library_updated  ON library_entries(updated_at DESC);

CREATE TABLE IF NOT EXISTS reading_history (
    id             UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id        UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    manga_id       UUID        NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    chapter_number INT         NOT NULL DEFAULT 0,
    volume_number  INT         NOT NULL DEFAULT 0,
    read_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_history_user_id ON reading_history(user_id);
CREATE INDEX IF NOT EXISTS idx_history_manga_id ON reading_history(manga_id);
CREATE INDEX IF NOT EXISTS idx_history_read_at  ON reading_history(read_at DESC);
