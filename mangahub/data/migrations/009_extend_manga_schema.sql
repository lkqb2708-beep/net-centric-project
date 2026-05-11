-- 009_extend_manga_schema.sql
-- Extends manga table with format, franchise, reading source, and provenance fields.

ALTER TABLE manga
  ADD COLUMN IF NOT EXISTS format            VARCHAR(20)  NOT NULL DEFAULT 'manga'
    CHECK (format IN ('manga','manhwa','manhua','light_novel','one_shot')),
  ADD COLUMN IF NOT EXISTS franchise         VARCHAR(255) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS franchise_part    INT          NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS reading_url       TEXT         NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS reading_source    VARCHAR(100) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS reading_region    VARCHAR(50)  NOT NULL DEFAULT 'worldwide',
  ADD COLUMN IF NOT EXISTS meta_source       VARCHAR(100) NOT NULL DEFAULT 'manual',
  ADD COLUMN IF NOT EXISTS cover_source      VARCHAR(100) NOT NULL DEFAULT 'manual',
  ADD COLUMN IF NOT EXISTS source_confidence VARCHAR(10)  NOT NULL DEFAULT 'high'
    CHECK (source_confidence IN ('high','medium','low','unverified'));

CREATE INDEX IF NOT EXISTS idx_manga_format    ON manga(format);
CREATE INDEX IF NOT EXISTS idx_manga_franchise ON manga(franchise);
