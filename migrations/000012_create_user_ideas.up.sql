CREATE TABLE IF NOT EXISTS user_ideas (
    user_id BIGINT NOT NULL,
    idea_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, idea_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (idea_id) REFERENCES ideas(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_ideas_user_id ON user_ideas(user_id);
CREATE INDEX idx_user_ideas_idea_id ON user_ideas(idea_id);
