CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          TEXT PRIMARY KEY,                 -- UUID (jti) – уникальный идентификатор токена
    user_id     INTEGER NOT NULL,                 -- ссылка на пользователя
    token_hash  TEXT NOT NULL UNIQUE,             -- SHA-256 хеш refresh-токена (защита от утечки БД)
    expires_at  TIMESTAMP WITH TIME ZONE NOT NULL, -- время истечения
    revoked     BOOLEAN NOT NULL DEFAULT FALSE,   -- принудительный отзыв (выход, смена пароля)
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


ALTER TABLE refresh_tokens
    ADD CONSTRAINT fk_refresh_tokens_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
