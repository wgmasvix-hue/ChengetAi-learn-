-- Recommendation history table: Tracks generated recommendations
CREATE TABLE recommendation_history (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id),
    resource_id VARCHAR(255) NOT NULL REFERENCES resources(id),
    score DECIMAL(3, 2) NOT NULL, -- 0.0 to 1.0
    rank INTEGER NOT NULL, -- Position in recommendation list
    strategy VARCHAR(50) NOT NULL CHECK (strategy IN ('collaborative', 'content_based', 'hybrid')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Recommendation feedback table: Captures user feedback on recommendations
CREATE TABLE recommendation_feedback (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id),
    resource_id VARCHAR(255) NOT NULL REFERENCES resources(id),
    helpful BOOLEAN NOT NULL,
    reason TEXT, -- Why they found it (un)helpful
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- User engagement metrics table: Aggregates recommendation effectiveness
CREATE TABLE user_engagement_metrics (
    user_id VARCHAR(255) PRIMARY KEY REFERENCES users(id),
    total_recommendations_received INTEGER DEFAULT 0,
    recommendations_clicked INTEGER DEFAULT 0,
    recommendations_completed INTEGER DEFAULT 0,
    average_satisfaction DECIMAL(3, 2),
    preferred_strategy VARCHAR(50),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indices for efficient queries
CREATE INDEX idx_recommendation_history_user ON recommendation_history(user_id, created_at DESC);
CREATE INDEX idx_recommendation_history_resource ON recommendation_history(resource_id);
CREATE INDEX idx_recommendation_history_strategy ON recommendation_history(strategy);
CREATE INDEX idx_recommendation_history_rank ON recommendation_history(rank);

CREATE INDEX idx_recommendation_feedback_user ON recommendation_feedback(user_id, created_at DESC);
CREATE INDEX idx_recommendation_feedback_resource ON recommendation_feedback(resource_id);
CREATE INDEX idx_recommendation_feedback_helpful ON recommendation_feedback(helpful);
CREATE INDEX idx_recommendation_feedback_created ON recommendation_feedback(created_at DESC);

-- Composite index for finding recommendations for a user within a time period
CREATE INDEX idx_recommendation_history_user_time ON recommendation_history(user_id, created_at DESC, strategy);

-- Trigger for updated_at on engagement metrics
CREATE TRIGGER update_user_engagement_metrics_updated_at
BEFORE UPDATE ON user_engagement_metrics
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
