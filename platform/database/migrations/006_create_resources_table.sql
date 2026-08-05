-- Create resources table for knowledge items
CREATE TABLE IF NOT EXISTS resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dspace_id VARCHAR(255) UNIQUE NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    creator_id UUID NOT NULL REFERENCES users(id),
    resource_type VARCHAR(50), -- e.g., 'textbook', 'paper', 'video', 'quiz'
    curriculum VARCHAR(100),
    subject VARCHAR(100),
    country VARCHAR(2),
    language VARCHAR(5) DEFAULT 'en',
    grade_level VARCHAR(50),
    publication_date DATE,
    file_url TEXT,
    file_hash VARCHAR(255),
    file_size BIGINT,
    mime_type VARCHAR(100),
    views_count INTEGER DEFAULT 0,
    downloads_count INTEGER DEFAULT 0,
    bookmarks_count INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'published',
    verified_at TIMESTAMP,
    indexed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_synced_at TIMESTAMP,

    CONSTRAINT status_check CHECK (status IN ('draft', 'published', 'archived', 'flagged'))
);

-- Create resource content table (extracted text, OCR results)
CREATE TABLE IF NOT EXISTS resource_content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    content_type VARCHAR(50), -- 'raw_text', 'ocr_text', 'extracted_metadata'
    content TEXT NOT NULL,
    extracted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    confidence_score DECIMAL(3, 2), -- For OCR quality
    processing_time_ms INTEGER
);

-- Create resource embeddings table
CREATE TABLE IF NOT EXISTS resource_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    embedding_model VARCHAR(100), -- e.g., 'sentence-transformers/all-MiniLM-L6-v2'
    vector VECTOR(384), -- Dimension depends on model used
    chunk_index INTEGER DEFAULT 0,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (resource_id, chunk_index)
);

-- Create resource tags table
CREATE TABLE IF NOT EXISTS resource_tags (
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,
    tag_type VARCHAR(50), -- 'subject', 'topic', 'skill'
    PRIMARY KEY (resource_id, tag)
);

-- Create resource ratings table
CREATE TABLE IF NOT EXISTS resource_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER,
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (resource_id, user_id),
    CONSTRAINT rating_check CHECK (rating >= 1 AND rating <= 5)
);

-- Create indexes
CREATE INDEX idx_resources_dspace_id ON resources(dspace_id);
CREATE INDEX idx_resources_creator_id ON resources(creator_id);
CREATE INDEX idx_resources_status ON resources(status);
CREATE INDEX idx_resources_curriculum ON resources(curriculum);
CREATE INDEX idx_resources_subject ON resources(subject);
CREATE INDEX idx_resources_country ON resources(country);
CREATE INDEX idx_resources_created_at ON resources(created_at);
CREATE INDEX idx_resource_content_resource_id ON resource_content(resource_id);
CREATE INDEX idx_resource_embeddings_resource_id ON resource_embeddings(resource_id);
CREATE INDEX idx_resource_tags_tag ON resource_tags(tag);
CREATE INDEX idx_resource_ratings_resource_id ON resource_ratings(resource_id);
CREATE INDEX idx_resource_ratings_user_id ON resource_ratings(user_id);

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_resources_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER resources_updated_at_trigger
BEFORE UPDATE ON resources
FOR EACH ROW
EXECUTE FUNCTION update_resources_updated_at();
