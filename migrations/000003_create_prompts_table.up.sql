CREATE TABLE IF NOT EXISTS prompts (
    id bigserial PRIMARY KEY,
    type_id bigint NOT NULL REFERENCES types(id) ON DELETE RESTRICT,
    text text NOT NULL
);

CREATE INDEX idx_prompts_type_id ON prompts(type_id);