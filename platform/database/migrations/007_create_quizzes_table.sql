-- Quizzes table: Stores quiz definitions
CREATE TABLE quizzes (
    id VARCHAR(255) PRIMARY KEY,
    resource_id VARCHAR(255) NOT NULL REFERENCES resources(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    difficulty_level VARCHAR(50) NOT NULL CHECK (difficulty_level IN ('easy', 'medium', 'hard')),
    estimated_time INTEGER NOT NULL, -- minutes
    passing_score INTEGER NOT NULL DEFAULT 70, -- 0-100
    questions JSONB NOT NULL, -- Array of Question objects
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Quiz responses table: Stores user quiz submissions and scores
CREATE TABLE quiz_responses (
    id VARCHAR(255) PRIMARY KEY,
    quiz_id VARCHAR(255) NOT NULL REFERENCES quizzes(id),
    user_id VARCHAR(255) NOT NULL REFERENCES users(id),
    responses JSONB NOT NULL, -- Array of Response objects with answers
    score INTEGER NOT NULL, -- Points earned
    percentage DECIMAL(5, 2) NOT NULL, -- Percentage score
    time_spent INTEGER NOT NULL, -- Seconds
    submitted_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Quiz attempts index: Track when users attempt quizzes
CREATE INDEX idx_quiz_responses_user_quiz ON quiz_responses(user_id, quiz_id, submitted_at DESC);
CREATE INDEX idx_quiz_responses_quiz ON quiz_responses(quiz_id);
CREATE INDEX idx_quiz_responses_user ON quiz_responses(user_id);
CREATE INDEX idx_quiz_responses_submitted ON quiz_responses(submitted_at DESC);

-- Quizzes index
CREATE INDEX idx_quizzes_resource ON quizzes(resource_id);
CREATE INDEX idx_quizzes_difficulty ON quizzes(difficulty_level);

-- Trigger for updated_at
CREATE TRIGGER update_quizzes_updated_at
BEFORE UPDATE ON quizzes
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
