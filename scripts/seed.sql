-- Seed data for local development.
-- Use the API for realistic user creation and password hashing in most cases.

INSERT INTO users (id, email, phone, password_hash, role, first_name, last_name)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'parent@example.com', '+263700000001', 'seeded-password-hash', 'parent', 'Phase', 'Parent')
ON CONFLICT (email) DO NOTHING;
