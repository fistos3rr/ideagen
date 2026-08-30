CREATE TABLE IF NOT EXISTS ideas (
    id bigserial PRIMARY KEY,
    type_id bigint NOT NULL REFERENCES types(id) ON DELETE RESTRICT,
    text text NOT NULL
);

CREATE INDEX idx_ideas_type_id ON ideas(type_id);
CREATE INDEX IF NOT EXISTS idx_text_ideas ON ideas USING GIN (to_tsvector('simple', text));
