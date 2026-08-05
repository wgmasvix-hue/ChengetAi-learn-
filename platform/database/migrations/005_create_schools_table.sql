-- Create schools (institutions) table
CREATE TABLE IF NOT EXISTS schools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    country VARCHAR(2),
    region VARCHAR(100),
    city VARCHAR(100),
    website VARCHAR(255),
    email VARCHAR(255),
    phone_number VARCHAR(20),
    logo_url TEXT,
    status VARCHAR(50) DEFAULT 'active',
    subscription_tier VARCHAR(50) DEFAULT 'free',
    subscription_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT status_check CHECK (status IN ('active', 'suspended', 'inactive')),
    CONSTRAINT tier_check CHECK (subscription_tier IN ('free', 'basic', 'professional', 'enterprise'))
);

-- Create school admins junction table
CREATE TABLE IF NOT EXISTS school_admins (
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'admin',
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (school_id, user_id)
);

-- Create school students junction table
CREATE TABLE IF NOT EXISTS school_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    grade_level VARCHAR(50),
    enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active',

    UNIQUE (school_id, user_id),
    CONSTRAINT status_check CHECK (status IN ('active', 'inactive', 'suspended'))
);

-- Create school curricula table
CREATE TABLE IF NOT EXISTS school_curricula (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    curriculum_type VARCHAR(100), -- e.g., 'zimsec', 'cambridge', 'ib'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (school_id, name)
);

-- Create indexes
CREATE INDEX idx_schools_country ON schools(country);
CREATE INDEX idx_schools_status ON schools(status);
CREATE INDEX idx_schools_subscription_tier ON schools(subscription_tier);
CREATE INDEX idx_school_admins_school_id ON school_admins(school_id);
CREATE INDEX idx_school_admins_user_id ON school_admins(user_id);
CREATE INDEX idx_school_enrollments_school_id ON school_enrollments(school_id);
CREATE INDEX idx_school_enrollments_user_id ON school_enrollments(user_id);
CREATE INDEX idx_school_enrollments_status ON school_enrollments(status);
CREATE INDEX idx_school_curricula_school_id ON school_curricula(school_id);

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_schools_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER schools_updated_at_trigger
BEFORE UPDATE ON schools
FOR EACH ROW
EXECUTE FUNCTION update_schools_updated_at();
