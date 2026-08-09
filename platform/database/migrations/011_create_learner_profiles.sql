-- Learner Profiles and Learning Preferences

-- Create learner profiles table
CREATE TABLE IF NOT EXISTS learner_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    education_level_id UUID NOT NULL REFERENCES education_levels(id),
    curriculum_id UUID NOT NULL REFERENCES curricula(id),
    school_id UUID REFERENCES schools(id) ON DELETE SET NULL,
    preferred_language VARCHAR(5) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    date_of_birth DATE,
    grade VARCHAR(50),

    -- Learning preferences
    learning_pace VARCHAR(50),
    preferred_content_types TEXT,
    notification_frequency VARCHAR(50),

    -- Profile status
    profile_complete BOOLEAN DEFAULT FALSE,
    onboarding_complete BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT learning_pace CHECK (learning_pace IN ('slow', 'normal', 'fast', 'self-paced'))
);

CREATE INDEX idx_learner_profiles_user_id ON learner_profiles(user_id);
CREATE INDEX idx_learner_profiles_education_level_id ON learner_profiles(education_level_id);
CREATE INDEX idx_learner_profiles_curriculum_id ON learner_profiles(curriculum_id);
CREATE INDEX idx_learner_profiles_school_id ON learner_profiles(school_id);

-- Create learner subject enrollment
CREATE TABLE IF NOT EXISTS learner_subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learner_profile_id UUID NOT NULL REFERENCES learner_profiles(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active',

    UNIQUE (learner_profile_id, subject_id),
    CONSTRAINT status CHECK (status IN ('active', 'completed', 'dropped', 'paused'))
);

CREATE INDEX idx_learner_subjects_learner_profile_id ON learner_subjects(learner_profile_id);
CREATE INDEX idx_learner_subjects_subject_id ON learner_subjects(subject_id);

-- Create learner topic enrollment (tracks which topics are started/completed)
CREATE TABLE IF NOT EXISTS learner_topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learner_profile_id UUID NOT NULL REFERENCES learner_profiles(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'not_started',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (learner_profile_id, topic_id),
    CONSTRAINT status CHECK (status IN ('not_started', 'in_progress', 'completed', 'abandoned'))
);

CREATE INDEX idx_learner_topics_learner_profile_id ON learner_topics(learner_profile_id);
CREATE INDEX idx_learner_topics_topic_id ON learner_topics(topic_id);
CREATE INDEX idx_learner_topics_status ON learner_topics(status);

-- Create learning history/activity log
CREATE TABLE IF NOT EXISTS learning_activity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learner_profile_id UUID NOT NULL REFERENCES learner_profiles(id) ON DELETE CASCADE,
    activity_type VARCHAR(100) NOT NULL,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    duration_minutes INT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT activity_type CHECK (activity_type IN ('lesson', 'practice', 'quiz', 'assessment', 'tutor_session', 'project_work', 'reading'))
);

CREATE INDEX idx_learning_activity_learner_profile_id ON learning_activity(learner_profile_id);
CREATE INDEX idx_learning_activity_activity_type ON learning_activity(activity_type);
CREATE INDEX idx_learning_activity_created_at ON learning_activity(created_at);
CREATE INDEX idx_learning_activity_subject_id ON learning_activity(subject_id);

-- Create trigger for learner_profiles updated_at
CREATE OR REPLACE FUNCTION update_learner_profiles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER learner_profiles_updated_at_trigger
BEFORE UPDATE ON learner_profiles
FOR EACH ROW
EXECUTE FUNCTION update_learner_profiles_updated_at();

CREATE TRIGGER learner_topics_updated_at_trigger
BEFORE UPDATE ON learner_topics
FOR EACH ROW
EXECUTE FUNCTION update_learner_profiles_updated_at();
