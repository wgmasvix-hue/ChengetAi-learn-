-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create role_permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Create user_roles junction table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);

-- Create indexes
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- Insert default roles
INSERT INTO roles (id, name, description) VALUES
    (gen_random_uuid(), 'admin', 'Administrator with full access'),
    (gen_random_uuid(), 'teacher', 'Teacher who creates and shares educational content'),
    (gen_random_uuid(), 'student', 'Student who accesses educational resources'),
    (gen_random_uuid(), 'institution', 'Institution administrator')
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (id, name, description, resource, action) VALUES
    (gen_random_uuid(), 'users.create', 'Create new users', 'users', 'create'),
    (gen_random_uuid(), 'users.read', 'Read user profiles', 'users', 'read'),
    (gen_random_uuid(), 'users.update', 'Update user profiles', 'users', 'update'),
    (gen_random_uuid(), 'users.delete', 'Delete user accounts', 'users', 'delete'),
    (gen_random_uuid(), 'resources.create', 'Create educational resources', 'resources', 'create'),
    (gen_random_uuid(), 'resources.read', 'Read resources', 'resources', 'read'),
    (gen_random_uuid(), 'resources.update', 'Update resources', 'resources', 'update'),
    (gen_random_uuid(), 'resources.delete', 'Delete resources', 'resources', 'delete'),
    (gen_random_uuid(), 'wallet.read', 'View wallet balance', 'wallet', 'read'),
    (gen_random_uuid(), 'wallet.withdraw', 'Request withdrawals', 'wallet', 'withdraw'),
    (gen_random_uuid(), 'analytics.read', 'View analytics', 'analytics', 'read'),
    (gen_random_uuid(), 'admin.access', 'Access admin panel', 'admin', 'access')
ON CONFLICT (name) DO NOTHING;
