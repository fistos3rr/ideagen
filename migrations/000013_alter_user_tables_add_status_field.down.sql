ALTER TABLE user_ideas DROP COLUMN IF EXISTS status;
DROP INDEX IF EXISTS idx_status_user_ideas;
