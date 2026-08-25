CREATE TABLE IF NOT EXISTS types (
  id bigserial PRIMARY KEY,
  name text NOT NULL UNIQUE
);
