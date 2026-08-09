-- Phase One: Create learner_profiles table
CREATE TABLE learner_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    education_level VARCHAR(100),
    grade VARCHAR(50),
    preferred_language VARCHAR(5) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    profile_complete BOOLEAN DEFAULT FALSE,
    onboarding_complete BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_learner_profiles_user_id ON learner_profiles(user_id);
CREATE INDEX idx_learner_profiles_education_level ON learner_profiles(education_level);
