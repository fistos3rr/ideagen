CREATE TABLE IF NOT EXISTS ideas (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    types text[] NOT NULL,
    message text NOT NULL
);
