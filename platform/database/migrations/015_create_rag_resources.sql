-- RAG (Retrieval-Augmented Generation) Resources and Knowledge Base

-- Create resources table
CREATE TABLE IF NOT EXISTS rag_resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content TEXT,

    -- Resource metadata
    author VARCHAR(255),
    publisher VARCHAR(255),
    source_type VARCHAR(100),
    source_id VARCHAR(255),
    source_url TEXT,

    -- Curriculum mapping
    curriculum_id UUID REFERENCES curricula(id) ON DELETE SET NULL,
    education_level_id UUID REFERENCES education_levels(id) ON DELETE SET NULL,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,

    -- Content properties
    language VARCHAR(5) DEFAULT 'en',
    content_type VARCHAR(100),
    publication_date DATE,
    license VARCHAR(255),
    copyright_notice TEXT,

    -- Quality tracking
    review_status VARCHAR(50) DEFAULT 'draft',
    quality_score DECIMAL(3, 2),
    citation_count INT DEFAULT 0,
    is_approved BOOLEAN DEFAULT FALSE,

    -- Technical metadata
    word_count INT,
    estimated_read_time_minutes INT,
    reading_level VARCHAR(50),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT review_status CHECK (review_status IN ('draft', 'under_review', 'approved', 'rejected')),
    CONSTRAINT source_type CHECK (source_type IN ('dspace', 'dare', 'openstax', 'gutenberg', 'institutional', 'teacher_created', 'research', 'textbook', 'other')),
    CONSTRAINT content_type CHECK (content_type IN ('text', 'pdf', 'image', 'video', 'audio', 'interactive'))
);

CREATE INDEX idx_rag_resources_subject_id ON rag_resources(subject_id);
CREATE INDEX idx_rag_resources_topic_id ON rag_resources(topic_id);
CREATE INDEX idx_rag_resources_review_status ON rag_resources(review_status);
CREATE INDEX idx_rag_resources_is_approved ON rag_resources(is_approved);
CREATE INDEX idx_rag_resources_source_type ON rag_resources(source_type);

-- Create resource chunks (for RAG - breaking large documents into searchable chunks)
CREATE TABLE IF NOT EXISTS resource_chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL REFERENCES rag_resources(id) ON DELETE CASCADE,

    -- Chunk content
    chunk_text TEXT NOT NULL,
    chunk_index INT,

    -- Chunk position in document
    start_position INT,
    end_position INT,

    -- Vector embedding (for semantic search)
    embedding vector(1536),

    -- Metadata
    chunk_type VARCHAR(50),
    keywords TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chunk_type CHECK (chunk_type IN ('paragraph', 'section', 'example', 'exercise', 'definition'))
);

CREATE INDEX idx_resource_chunks_resource_id ON resource_chunks(resource_id);
CREATE INDEX idx_resource_chunks_embedding ON resource_chunks USING ivfflat (embedding vector_cosine_ops);

-- Create resource citations/sources table (for tracking where information comes from)
CREATE TABLE IF NOT EXISTS resource_citations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_chunk_id UUID NOT NULL REFERENCES resource_chunks(id) ON DELETE CASCADE,

    citation_text VARCHAR(255),
    citation_source VARCHAR(255),
    page_number INT,
    url TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_resource_citations_resource_chunk_id ON resource_citations(resource_chunk_id);

-- Create resource ratings table (for quality assessment)
CREATE TABLE IF NOT EXISTS resource_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL REFERENCES rag_resources(id) ON DELETE CASCADE,
    rater_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Rating dimensions
    accuracy_rating INT,
    relevance_rating INT,
    clarity_rating INT,
    completeness_rating INT,
    overall_rating INT,

    comments TEXT,
    rated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (resource_id, rater_id),
    CONSTRAINT accuracy_rating CHECK (accuracy_rating >= 1 AND accuracy_rating <= 5),
    CONSTRAINT relevance_rating CHECK (relevance_rating >= 1 AND relevance_rating <= 5),
    CONSTRAINT clarity_rating CHECK (clarity_rating >= 1 AND clarity_rating <= 5),
    CONSTRAINT completeness_rating CHECK (completeness_rating >= 1 AND completeness_rating <= 5),
    CONSTRAINT overall_rating CHECK (overall_rating >= 1 AND overall_rating <= 5)
);

CREATE INDEX idx_resource_ratings_resource_id ON resource_ratings(resource_id);
CREATE INDEX idx_resource_ratings_rater_id ON resource_ratings(rater_id);

-- Create resource usage statistics table
CREATE TABLE IF NOT EXISTS resource_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL REFERENCES rag_resources(id) ON DELETE CASCADE,

    view_count INT DEFAULT 0,
    citation_count INT DEFAULT 0,
    assignment_count INT DEFAULT 0,

    last_viewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (resource_id)
);

CREATE INDEX idx_resource_usage_resource_id ON resource_usage(resource_id);

-- Create RAG queries log (for monitoring and improving search)
CREATE TABLE IF NOT EXISTS rag_queries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learner_profile_id UUID REFERENCES learner_profiles(id) ON DELETE SET NULL,
    teacher_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Query information
    query_text TEXT NOT NULL,
    query_type VARCHAR(50),
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,

    -- Results
    results_returned INT DEFAULT 0,
    top_result_id UUID REFERENCES rag_resources(id) ON DELETE SET NULL,
    query_quality_feedback VARCHAR(50),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT query_type CHECK (query_type IN ('search', 'tutor_lookup', 'teacher_lookup', 'auto_suggest')),
    CONSTRAINT query_quality_feedback CHECK (query_quality_feedback IN ('helpful', 'unhelpful', 'partially_helpful'))
);

CREATE INDEX idx_rag_queries_learner_profile_id ON rag_queries(learner_profile_id);
CREATE INDEX idx_rag_queries_teacher_id ON rag_queries(teacher_id);
CREATE INDEX idx_rag_queries_created_at ON rag_queries(created_at);

-- Create ai_generated_answers table (for logging AI responses)
CREATE TABLE IF NOT EXISTS ai_generated_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tutor_session_id UUID REFERENCES tutor_sessions(id) ON DELETE SET NULL,
    query_id UUID REFERENCES rag_queries(id) ON DELETE SET NULL,

    -- Query and answer
    query_text TEXT NOT NULL,
    answer_text TEXT NOT NULL,

    -- Source tracking
    resources_used JSONB,
    resource_ids UUID[],

    -- Quality metrics
    relevance_score DECIMAL(3, 2),
    confidence_score DECIMAL(3, 2),
    contains_citation BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ai_generated_answers_tutor_session_id ON ai_generated_answers(tutor_session_id);
CREATE INDEX idx_ai_generated_answers_created_at ON ai_generated_answers(created_at);

-- Create trigger for rag_resources updated_at
CREATE OR REPLACE FUNCTION update_rag_resources_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER rag_resources_updated_at_trigger
BEFORE UPDATE ON rag_resources
FOR EACH ROW
EXECUTE FUNCTION update_rag_resources_updated_at();

CREATE TRIGGER resource_usage_updated_at_trigger
BEFORE UPDATE ON resource_usage
FOR EACH ROW
EXECUTE FUNCTION update_rag_resources_updated_at();

-- Create pgvector extension if not exists
CREATE EXTENSION IF NOT EXISTS vector;
