-- Phase One: Create teacher_profiles table
CREATE TABLE teacher_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    school_name VARCHAR(255),
    subjects_taught TEXT,
    bio TEXT,
    profile_complete BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_teacher_profiles_user_id ON teacher_profiles(user_id);
