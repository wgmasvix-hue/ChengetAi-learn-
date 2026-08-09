-- Mastery Tracking and Assessment Tables

-- Create questions table
CREATE TABLE IF NOT EXISTS questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    skill_id UUID REFERENCES skills(id) ON DELETE SET NULL,
    learning_objective_id UUID REFERENCES learning_objectives(id) ON DELETE SET NULL,

    -- Question content
    prompt TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL,
    difficulty_level VARCHAR(50),
    estimated_time_seconds INT,

    -- Source tracking
    source_id UUID,
    source_type VARCHAR(100),

    -- Quality tracking
    review_status VARCHAR(50) DEFAULT 'draft',
    ai_generated BOOLEAN DEFAULT FALSE,
    explanation TEXT,
    marking_scheme TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT question_type CHECK (question_type IN ('multiple_choice', 'short_answer', 'numeric', 'true_false', 'matching', 'structured', 'essay')),
    CONSTRAINT difficulty_level CHECK (difficulty_level IN ('easy', 'medium', 'hard')),
    CONSTRAINT review_status CHECK (review_status IN ('draft', 'under_review', 'approved', 'rejected'))
);

CREATE INDEX idx_questions_subject_id ON questions(subject_id);
CREATE INDEX idx_questions_topic_id ON questions(topic_id);
CREATE INDEX idx_questions_skill_id ON questions(skill_id);
CREATE INDEX idx_questions_difficulty_level ON questions(difficulty_level);
CREATE INDEX idx_questions_question_type ON questions(question_type);
CREATE INDEX idx_questions_review_status ON questions(review_status);

-- Create question answers table (for multiple choice, matching, etc.)
CREATE TABLE IF NOT EXISTS question_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    answer_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    explanation TEXT,
    sequence_order INT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_question_answers_question_id ON question_answers(question_id);
CREATE INDEX idx_question_answers_is_correct ON question_answers(is_correct);

-- Create mastery tracking table
CREATE TABLE IF NOT EXISTS learner_mastery (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learner_profile_id UUID NOT NULL REFERENCES learner_profiles(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    skill_id UUID REFERENCES skills(id) ON DELETE CASCADE,

    -- Mastery metrics
    mastery_score DECIMAL(5, 2) DEFAULT 0.00,
    confidence_score DECIMAL(5, 2) DEFAULT 0.00,

    -- Attempt tracking
    total_attempts INT DEFAULT 0,
    correct_attempts INT DEFAULT 0,
    average_time_seconds INT,

    -- Progress tracking
    last_attempted_at TIMESTAMP,
    estimated_completion_date DATE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (learner_profile_id, topic_id, COALESCE(skill_id, '00000000-0000-0000-0000-000000000000')),
    CONSTRAINT mastery_score CHECK (mastery_score >= 0 AND mastery_score <= 100),
    CONSTRAINT confidence_score CHECK (confidence_score >= 0 AND confidence_score <= 100)
);

CREATE INDEX idx_learner_mastery_learner_profile_id ON learner_mastery(learner_profile_id);
CREATE INDEX idx_learner_mastery_topic_id ON learner_mastery(topic_id);
CREATE INDEX idx_learner_mastery_skill_id ON learner_mastery(skill_id);
CREATE INDEX idx_learner_mastery_mastery_score ON learner_mastery(mastery_score);

-- Create attempt history table
CREATE TABLE IF NOT EXISTS attempt_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learner_profile_id UUID NOT NULL REFERENCES learner_profiles(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,

    -- Response tracking
    learner_answer TEXT,
    is_correct BOOLEAN,
    score INT,
    max_score INT,
    time_taken_seconds INT,
    hints_used INT DEFAULT 0,

    -- Attempt metadata
    attempt_number INT,
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT is_correct CHECK (is_correct IN (true, false))
);

CREATE INDEX idx_attempt_history_learner_profile_id ON attempt_history(learner_profile_id);
CREATE INDEX idx_attempt_history_question_id ON attempt_history(question_id);
CREATE INDEX idx_attempt_history_topic_id ON attempt_history(topic_id);
CREATE INDEX idx_attempt_history_attempted_at ON attempt_history(attempted_at);
CREATE INDEX idx_attempt_history_is_correct ON attempt_history(is_correct);

-- Create assessments table
CREATE TABLE IF NOT EXISTS assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    assessment_type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INT,
    total_marks INT,
    pass_marks INT,

    -- Composition
    number_of_questions INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT assessment_type CHECK (assessment_type IN ('diagnostic', 'formative', 'summative', 'practice', 'mock_exam'))
);

CREATE INDEX idx_assessments_topic_id ON assessments(topic_id);
CREATE INDEX idx_assessments_subject_id ON assessments(subject_id);
CREATE INDEX idx_assessments_assessment_type ON assessments(assessment_type);

-- Create assessment questions junction
CREATE TABLE IF NOT EXISTS assessment_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    marks INT,
    sequence_order INT,

    PRIMARY KEY (assessment_id, question_id)
);

CREATE INDEX idx_assessment_questions_assessment_id ON assessment_questions(assessment_id);
CREATE INDEX idx_assessment_questions_question_id ON assessment_questions(question_id);

-- Create learner assessment attempts table
CREATE TABLE IF NOT EXISTS learner_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learner_profile_id UUID NOT NULL REFERENCES learner_profiles(id) ON DELETE CASCADE,
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,

    -- Performance
    score INT,
    max_score INT,
    percentage_score DECIMAL(5, 2),
    passed BOOLEAN,

    -- Timing
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    time_taken_seconds INT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT passed CHECK (passed IN (true, false))
);

CREATE INDEX idx_learner_assessments_learner_profile_id ON learner_assessments(learner_profile_id);
CREATE INDEX idx_learner_assessments_assessment_id ON learner_assessments(assessment_id);
CREATE INDEX idx_learner_assessments_completed_at ON learner_assessments(completed_at);
CREATE INDEX idx_learner_assessments_passed ON learner_assessments(passed);

-- Create trigger for questions and assessments updated_at
CREATE OR REPLACE FUNCTION update_question_tables_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER questions_updated_at_trigger
BEFORE UPDATE ON questions
FOR EACH ROW
EXECUTE FUNCTION update_question_tables_updated_at();

CREATE TRIGGER assessments_updated_at_trigger
BEFORE UPDATE ON assessments
FOR EACH ROW
EXECUTE FUNCTION update_question_tables_updated_at();

CREATE TRIGGER learner_mastery_updated_at_trigger
BEFORE UPDATE ON learner_mastery
FOR EACH ROW
EXECUTE FUNCTION update_question_tables_updated_at();
