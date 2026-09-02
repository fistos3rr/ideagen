ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'User';
CREATE INDEX idx_role_users ON users(role);
