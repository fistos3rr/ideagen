CREATE TABLE IF NOT EXISTS prompts (
    id bigserial PRIMARY KEY,
    type_id bigint NOT NULL REFERENCES types(id) ON DELETE RESTRICT,
    text text NOT NULL
);

CREATE INDEX idx_prompts_type_id ON prompts(type_id);
CREATE INDEX IF NOT EXISTS idx_text_prompts ON prompts USING GIN (to_tsvector('simple', text));
ALTER TABLE prompts ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT true;
