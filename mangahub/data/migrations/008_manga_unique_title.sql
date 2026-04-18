-- 008_manga_unique_title.sql
-- Add a unique constraint to title to allow for upserting seed data.
ALTER TABLE manga ADD CONSTRAINT manga_title_unique UNIQUE (title);
