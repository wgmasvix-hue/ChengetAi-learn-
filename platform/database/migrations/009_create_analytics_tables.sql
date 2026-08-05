-- Analytics events table: Stores all platform events
CREATE TABLE analytics_events (
    id BIGSERIAL PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL, -- view, download, bookmark, quiz_submit, etc
    user_id VARCHAR(255) REFERENCES users(id),
    resource_id VARCHAR(255) REFERENCES resources(id),
    school_id VARCHAR(255) REFERENCES schools(id),
    action VARCHAR(255),
    metadata JSONB, -- Additional event context
    session_id VARCHAR(255), -- Session identifier
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Analytics reports table: Stores generated reports
CREATE TABLE analytics_reports (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- platform, user, resource, school
    period VARCHAR(50) NOT NULL, -- daily, weekly, monthly
    metrics JSONB NOT NULL, -- Report metrics and KPIs
    generated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Event aggregations table: Pre-computed hourly metrics
CREATE TABLE event_aggregations (
    id BIGSERIAL PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,
    hour_bucket TIMESTAMP NOT NULL, -- Hourly bucket for aggregation
    count INTEGER DEFAULT 0,
    unique_users INTEGER DEFAULT 0,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_type, hour_bucket)
);

-- User session tracking table: Tracks user sessions
CREATE TABLE user_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id),
    session_id VARCHAR(255) NOT NULL UNIQUE,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    session_length INTEGER, -- Duration in seconds
    events_count INTEGER DEFAULT 0,
    last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indices for efficient analytics queries
CREATE INDEX idx_analytics_events_event_type ON analytics_events(event_type);
CREATE INDEX idx_analytics_events_user_id ON analytics_events(user_id, timestamp DESC);
CREATE INDEX idx_analytics_events_resource_id ON analytics_events(resource_id, timestamp DESC);
CREATE INDEX idx_analytics_events_timestamp ON analytics_events(timestamp DESC);
CREATE INDEX idx_analytics_events_session_id ON analytics_events(session_id);
CREATE INDEX idx_analytics_events_user_resource ON analytics_events(user_id, resource_id, event_type);
CREATE INDEX idx_analytics_events_school_id ON analytics_events(school_id);

-- Composite index for common queries
CREATE INDEX idx_analytics_events_composite ON analytics_events(event_type, timestamp DESC, user_id);

CREATE INDEX idx_analytics_reports_type ON analytics_reports(type);
CREATE INDEX idx_analytics_reports_created ON analytics_reports(created_at DESC);

CREATE INDEX idx_event_aggregations_hour ON event_aggregations(hour_bucket DESC);
CREATE INDEX idx_event_aggregations_event_type ON event_aggregations(event_type);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id, started_at DESC);
CREATE INDEX idx_user_sessions_session_id ON user_sessions(session_id);
CREATE INDEX idx_user_sessions_started_at ON user_sessions(started_at DESC);

-- Partition analytics events by month for large datasets
-- This improves query performance on time-range queries
-- SELECT create_hypertable('analytics_events', 'timestamp', if_not_exists => TRUE);
