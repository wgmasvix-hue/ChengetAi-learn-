CREATE TABLE learner_profiles (
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id            UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    education_level    VARCHAR(100),
    grade              VARCHAR(50),
    preferred_language VARCHAR(50) DEFAULT 'en',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_learner_profiles_user_id ON learner_profiles(user_id);
