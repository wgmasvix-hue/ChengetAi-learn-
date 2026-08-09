-- Curriculum Structure Tables
-- Hierarchical curriculum model: Curriculum -> Education Level -> Subject -> Syllabus -> Unit -> Topic -> Learning Objective -> Skill

-- Create curricula table
CREATE TABLE IF NOT EXISTS curricula (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    country VARCHAR(2) NOT NULL,
    description TEXT,
    version VARCHAR(50),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT curriculum_status CHECK (status IN ('active', 'archived', 'draft'))
);

CREATE INDEX idx_curricula_country ON curricula(country);
CREATE INDEX idx_curricula_status ON curricula(status);

-- Create education levels table
CREATE TABLE IF NOT EXISTS education_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    curriculum_id UUID NOT NULL REFERENCES curricula(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50),
    description TEXT,
    sequence_order INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (curriculum_id, code)
);

CREATE INDEX idx_education_levels_curriculum_id ON education_levels(curriculum_id);

-- Create subjects table
CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    education_level_id UUID NOT NULL REFERENCES education_levels(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50),
    description TEXT,
    sequence_order INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (education_level_id, code)
);

CREATE INDEX idx_subjects_education_level_id ON subjects(education_level_id);

-- Create syllabuses table
CREATE TABLE IF NOT EXISTS syllabuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration_weeks INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_syllabuses_subject_id ON syllabuses(subject_id);

-- Create units table (within syllabuses)
CREATE TABLE IF NOT EXISTS units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    syllabus_id UUID NOT NULL REFERENCES syllabuses(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sequence_order INT,
    duration_weeks INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_units_syllabus_id ON units(syllabus_id);

-- Create topics table
CREATE TABLE IF NOT EXISTS topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sequence_order INT,
    estimated_hours INT,
    difficulty_level VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT difficulty_level CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced'))
);

CREATE INDEX idx_topics_unit_id ON topics(unit_id);
CREATE INDEX idx_topics_subject_id ON topics(subject_id);

-- Create learning objectives table
CREATE TABLE IF NOT EXISTS learning_objectives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    objective_type VARCHAR(50),
    bloom_level VARCHAR(50),
    sequence_order INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT bloom_level CHECK (bloom_level IN ('remember', 'understand', 'apply', 'analyze', 'evaluate', 'create'))
);

CREATE INDEX idx_learning_objectives_topic_id ON learning_objectives(topic_id);

-- Create skills table
CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learning_objective_id UUID NOT NULL REFERENCES learning_objectives(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    skill_category VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_skills_learning_objective_id ON skills(learning_objective_id);

-- Create trigger for updated_at on curriculum tables
CREATE OR REPLACE FUNCTION update_curriculum_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER curricula_updated_at_trigger
BEFORE UPDATE ON curricula
FOR EACH ROW
EXECUTE FUNCTION update_curriculum_updated_at();

CREATE TRIGGER subjects_updated_at_trigger
BEFORE UPDATE ON subjects
FOR EACH ROW
EXECUTE FUNCTION update_curriculum_updated_at();

CREATE TRIGGER syllabuses_updated_at_trigger
BEFORE UPDATE ON syllabuses
FOR EACH ROW
EXECUTE FUNCTION update_curriculum_updated_at();

CREATE TRIGGER units_updated_at_trigger
BEFORE UPDATE ON units
FOR EACH ROW
EXECUTE FUNCTION update_curriculum_updated_at();

CREATE TRIGGER topics_updated_at_trigger
BEFORE UPDATE ON topics
FOR EACH ROW
EXECUTE FUNCTION update_curriculum_updated_at();

CREATE TRIGGER learning_objectives_updated_at_trigger
BEFORE UPDATE ON learning_objectives
FOR EACH ROW
EXECUTE FUNCTION update_curriculum_updated_at();

CREATE TRIGGER skills_updated_at_trigger
BEFORE UPDATE ON skills
FOR EACH ROW
EXECUTE FUNCTION update_curriculum_updated_at();
