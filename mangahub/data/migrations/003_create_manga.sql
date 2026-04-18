-- 003_create_manga.sql
CREATE TABLE IF NOT EXISTS manga (
    id              UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    title           VARCHAR(500) NOT NULL,
    author          VARCHAR(255) NOT NULL DEFAULT '',
    artist          VARCHAR(255) NOT NULL DEFAULT '',
    genres          TEXT[]       NOT NULL DEFAULT '{}',
    status          VARCHAR(20)  NOT NULL DEFAULT 'ongoing' CHECK (status IN ('ongoing','completed','hiatus','cancelled')),
    chapter_count   INT          NOT NULL DEFAULT 0,
    volume_count    INT          NOT NULL DEFAULT 0,
    description     TEXT         NOT NULL DEFAULT '',
    cover_url       TEXT         NOT NULL DEFAULT '',
    year            INT,
    rating          NUMERIC(4,2) NOT NULL DEFAULT 0,
    popularity_rank INT          NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_manga_title      ON manga(title);
CREATE INDEX IF NOT EXISTS idx_manga_author     ON manga(author);
CREATE INDEX IF NOT EXISTS idx_manga_status     ON manga(status);
CREATE INDEX IF NOT EXISTS idx_manga_year       ON manga(year);
CREATE INDEX IF NOT EXISTS idx_manga_rating     ON manga(rating DESC);
CREATE INDEX IF NOT EXISTS idx_manga_popularity ON manga(popularity_rank);
CREATE INDEX IF NOT EXISTS idx_manga_title_trgm ON manga USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_manga_genres     ON manga USING GIN (genres);
