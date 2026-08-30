CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    email text UNIQUE NOT NULL,
    password_hash bytea NOT NULL
);

CREATE UNIQUE INDEX users_email_lower_idx ON users (LOWER(email));
