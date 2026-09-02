ALTER TABLE user_ideas ADD COLUMN status smallint NOT NULL DEFAULT 1;
CREATE INDEX idx_status_user_ideas ON user_ideas(status);
