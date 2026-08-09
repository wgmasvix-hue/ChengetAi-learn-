-- Class Management Tables

-- Create classes table (for teacher-led groups)
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    education_level_id UUID NOT NULL REFERENCES education_levels(id),
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,

    -- Class info
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    description TEXT,
    max_learners INT,

    -- Class status
    status VARCHAR(50) DEFAULT 'active',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT status CHECK (status IN ('active', 'completed', 'archived'))
);

CREATE INDEX idx_classes_school_id ON classes(school_id);
CREATE INDEX idx_classes_teacher_id ON classes(teacher_id);
CREATE INDEX idx_classes_subject_id ON classes(subject_id);
CREATE INDEX idx_classes_status ON classes(status);

-- Create class members junction table
CREATE TABLE IF NOT EXISTS class_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    learner_profile_id UUID NOT NULL REFERENCES learner_profiles(id) ON DELETE CASCADE,

    -- Member role
    role VARCHAR(50) DEFAULT 'member',
    enrollment_status VARCHAR(50) DEFAULT 'active',

    enrolled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    enrollment_date TIMESTAMP,

    UNIQUE (class_id, learner_profile_id),
    CONSTRAINT role CHECK (role IN ('member', 'group_leader')),
    CONSTRAINT enrollment_status CHECK (enrollment_status IN ('active', 'inactive', 'dropped'))
);

CREATE INDEX idx_class_members_class_id ON class_members(class_id);
CREATE INDEX idx_class_members_learner_profile_id ON class_members(learner_profile_id);
CREATE INDEX idx_class_members_enrollment_status ON class_members(enrollment_status);

-- Create class lessons table
CREATE TABLE IF NOT EXISTS class_lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,

    -- Lesson info
    title VARCHAR(255) NOT NULL,
    description TEXT,
    lesson_type VARCHAR(50),
    duration_minutes INT,

    -- Content
    content TEXT,
    resources_json JSONB,

    -- Scheduling
    scheduled_date DATE,
    scheduled_time TIME,
    status VARCHAR(50) DEFAULT 'draft',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT lesson_type CHECK (lesson_type IN ('lecture', 'practice', 'assessment', 'discussion', 'project')),
    CONSTRAINT status CHECK (status IN ('draft', 'scheduled', 'completed', 'cancelled'))
);

CREATE INDEX idx_class_lessons_class_id ON class_lessons(class_id);
CREATE INDEX idx_class_lessons_topic_id ON class_lessons(topic_id);
CREATE INDEX idx_class_lessons_status ON class_lessons(status);
CREATE INDEX idx_class_lessons_scheduled_date ON class_lessons(scheduled_date);

-- Create class assignments table
CREATE TABLE IF NOT EXISTS class_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    assessment_id UUID REFERENCES assessments(id) ON DELETE SET NULL,

    -- Assignment info
    title VARCHAR(255) NOT NULL,
    description TEXT,
    assignment_type VARCHAR(50),

    -- Deadlines
    due_date DATE,
    due_time TIME,
    late_submission_allowed BOOLEAN DEFAULT FALSE,
    late_submission_days INT,

    -- Status
    status VARCHAR(50) DEFAULT 'active',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT assignment_type CHECK (assignment_type IN ('homework', 'classwork', 'project', 'quiz', 'test')),
    CONSTRAINT status CHECK (status IN ('draft', 'active', 'closed', 'cancelled'))
);

CREATE INDEX idx_class_assignments_class_id ON class_assignments(class_id);
CREATE INDEX idx_class_assignments_topic_id ON class_assignments(topic_id);
CREATE INDEX idx_class_assignments_status ON class_assignments(status);
CREATE INDEX idx_class_assignments_due_date ON class_assignments(due_date);

-- Create assignment submissions table
CREATE TABLE IF NOT EXISTS assignment_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_assignment_id UUID NOT NULL REFERENCES class_assignments(id) ON DELETE CASCADE,
    learner_profile_id UUID NOT NULL REFERENCES learner_profiles(id) ON DELETE CASCADE,

    -- Submission content
    submission_text TEXT,
    submission_files JSONB,
    notes TEXT,

    -- Grading
    score INT,
    max_score INT,
    feedback TEXT,
    graded_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Timing
    submitted_at TIMESTAMP,
    submitted_late BOOLEAN DEFAULT FALSE,
    graded_at TIMESTAMP,

    submission_status VARCHAR(50) DEFAULT 'draft',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (class_assignment_id, learner_profile_id),
    CONSTRAINT submission_status CHECK (submission_status IN ('draft', 'submitted', 'graded', 'returned'))
);

CREATE INDEX idx_assignment_submissions_class_assignment_id ON assignment_submissions(class_assignment_id);
CREATE INDEX idx_assignment_submissions_learner_profile_id ON assignment_submissions(learner_profile_id);
CREATE INDEX idx_assignment_submissions_submitted_at ON assignment_submissions(submitted_at);
CREATE INDEX idx_assignment_submissions_graded_at ON assignment_submissions(graded_at);
CREATE INDEX idx_assignment_submissions_submission_status ON assignment_submissions(submission_status);

-- Create class announcements table
CREATE TABLE IF NOT EXISTS class_announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    priority VARCHAR(50) DEFAULT 'normal',

    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT priority CHECK (priority IN ('low', 'normal', 'high', 'urgent'))
);

CREATE INDEX idx_class_announcements_class_id ON class_announcements(class_id);
CREATE INDEX idx_class_announcements_teacher_id ON class_announcements(teacher_id);
CREATE INDEX idx_class_announcements_published_at ON class_announcements(published_at);

-- Create class analytics view helper table
CREATE TABLE IF NOT EXISTS class_analytics_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL UNIQUE REFERENCES classes(id) ON DELETE CASCADE,

    -- Performance metrics
    average_mastery DECIMAL(5, 2),
    learners_count INT,
    assignments_completed INT,
    assessment_average DECIMAL(5, 2),

    -- Common issues
    common_misconceptions TEXT,
    weak_topics TEXT,

    -- Last updated
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_class_analytics_cache_class_id ON class_analytics_cache(class_id);

-- Create trigger for class tables updated_at
CREATE OR REPLACE FUNCTION update_class_tables_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER classes_updated_at_trigger
BEFORE UPDATE ON classes
FOR EACH ROW
EXECUTE FUNCTION update_class_tables_updated_at();

CREATE TRIGGER class_lessons_updated_at_trigger
BEFORE UPDATE ON class_lessons
FOR EACH ROW
EXECUTE FUNCTION update_class_tables_updated_at();

CREATE TRIGGER class_assignments_updated_at_trigger
BEFORE UPDATE ON class_assignments
FOR EACH ROW
EXECUTE FUNCTION update_class_tables_updated_at();

CREATE TRIGGER assignment_submissions_updated_at_trigger
BEFORE UPDATE ON assignment_submissions
FOR EACH ROW
EXECUTE FUNCTION update_class_tables_updated_at();

CREATE TRIGGER class_announcements_updated_at_trigger
BEFORE UPDATE ON class_announcements
FOR EACH ROW
EXECUTE FUNCTION update_class_tables_updated_at();
