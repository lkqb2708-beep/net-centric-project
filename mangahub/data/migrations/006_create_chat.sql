-- 006_create_chat.sql
CREATE TABLE IF NOT EXISTS chat_rooms (
    id          UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(20)  NOT NULL DEFAULT 'general'
                    CHECK (type IN ('general','manga')),
    manga_id    UUID         REFERENCES manga(id) ON DELETE SET NULL,
    description TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_chat_rooms_type     ON chat_rooms(type);
CREATE INDEX IF NOT EXISTS idx_chat_rooms_manga_id ON chat_rooms(manga_id);

CREATE TABLE IF NOT EXISTS chat_messages (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id    UUID        NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_chat_msgs_room_id    ON chat_messages(room_id);
CREATE INDEX IF NOT EXISTS idx_chat_msgs_user_id    ON chat_messages(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_msgs_created_at ON chat_messages(created_at DESC);

-- Insert default chat rooms
INSERT INTO chat_rooms (name, type, description) VALUES
    ('General', 'general', 'General manga discussion'),
    ('Announcements', 'general', 'Platform announcements and updates'),
    ('Recommendations', 'general', 'Recommend manga to fellow readers')
ON CONFLICT DO NOTHING;
