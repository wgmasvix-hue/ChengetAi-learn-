-- Tutor Sessions and AI Learning Interactions

-- Create tutor sessions table
CREATE TABLE IF NOT EXISTS tutor_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learner_profile_id UUID NOT NULL REFERENCES learner_profiles(id) ON DELETE CASCADE,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    learning_objective_id UUID REFERENCES learning_objectives(id) ON DELETE SET NULL,

    -- Session tracking
    session_type VARCHAR(50),
    session_title VARCHAR(255),
    description TEXT,

    -- Mastery tracking
    mastery_before DECIMAL(5, 2),
    mastery_after DECIMAL(5, 2),

    -- Timing
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    duration_minutes INT,

    -- Session status
    status VARCHAR(50) DEFAULT 'active',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT status CHECK (status IN ('active', 'completed', 'abandoned'))
);

CREATE INDEX idx_tutor_sessions_learner_profile_id ON tutor_sessions(learner_profile_id);
CREATE INDEX idx_tutor_sessions_subject_id ON tutor_sessions(subject_id);
CREATE INDEX idx_tutor_sessions_topic_id ON tutor_sessions(topic_id);
CREATE INDEX idx_tutor_sessions_status ON tutor_sessions(status);
CREATE INDEX idx_tutor_sessions_started_at ON tutor_sessions(started_at);

-- Create tutor messages table
CREATE TABLE IF NOT EXISTS tutor_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tutor_session_id UUID NOT NULL REFERENCES tutor_sessions(id) ON DELETE CASCADE,

    -- Message metadata
    role VARCHAR(50) NOT NULL,
    message_type VARCHAR(50),
    content TEXT NOT NULL,

    -- Message classification
    intent VARCHAR(100),
    confidence DECIMAL(3, 2),

    -- Message tagging
    uses_hint BOOLEAN DEFAULT FALSE,
    includes_explanation BOOLEAN DEFAULT FALSE,
    includes_example BOOLEAN DEFAULT FALSE,
    includes_question BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT role CHECK (role IN ('learner', 'tutor', 'system')),
    CONSTRAINT message_type CHECK (message_type IN ('question', 'answer', 'explanation', 'hint', 'feedback', 'assessment', 'system', 'practice'))
);

CREATE INDEX idx_tutor_messages_tutor_session_id ON tutor_messages(tutor_session_id);
CREATE INDEX idx_tutor_messages_role ON tutor_messages(role);
CREATE INDEX idx_tutor_messages_message_type ON tutor_messages(message_type);
CREATE INDEX idx_tutor_messages_created_at ON tutor_messages(created_at);

-- Create tutor session feedback table
CREATE TABLE IF NOT EXISTS session_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tutor_session_id UUID NOT NULL UNIQUE REFERENCES tutor_sessions(id) ON DELETE CASCADE,

    -- Learner feedback
    clarity_rating INT,
    helpfulness_rating INT,
    pacing_rating INT,
    engagement_rating INT,
    comments TEXT,

    -- Tutor performance metrics
    misconceptions_addressed INT,
    learning_objectives_met INT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT clarity_rating CHECK (clarity_rating >= 1 AND clarity_rating <= 5),
    CONSTRAINT helpfulness_rating CHECK (helpfulness_rating >= 1 AND helpfulness_rating <= 5),
    CONSTRAINT pacing_rating CHECK (pacing_rating >= 1 AND pacing_rating <= 5),
    CONSTRAINT engagement_rating CHECK (engagement_rating >= 1 AND engagement_rating <= 5)
);

CREATE INDEX idx_session_feedback_tutor_session_id ON session_feedback(tutor_session_id);

-- Create conversation insights table (for analytics)
CREATE TABLE IF NOT EXISTS session_insights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tutor_session_id UUID NOT NULL REFERENCES tutor_sessions(id) ON DELETE CASCADE,

    -- Detected issues
    misconceptions TEXT,
    knowledge_gaps TEXT,
    strengths TEXT,
    recommended_next_steps TEXT,

    -- Session quality metrics
    question_count INT,
    explanation_count INT,
    hint_count INT,
    feedback_count INT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_session_insights_tutor_session_id ON session_insights(tutor_session_id);

-- Create trigger for tutor_sessions updated_at
CREATE OR REPLACE FUNCTION update_tutor_sessions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tutor_sessions_updated_at_trigger
BEFORE UPDATE ON tutor_sessions
FOR EACH ROW
EXECUTE FUNCTION update_tutor_sessions_updated_at();
