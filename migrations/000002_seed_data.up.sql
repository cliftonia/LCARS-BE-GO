-- Seed sample data for development
-- Password hash is for 'admin123' using bcrypt cost 12
INSERT INTO users (id, name, email, password_hash) VALUES
    ('00000000-0000-0000-0000-000000000001', 'Admin User', 'admin@subspace.dev', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYbYvJ6Y7bm'),
    ('00000000-0000-0000-0000-000000000002', 'Test User', 'test@subspace.dev', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYbYvJ6Y7bm')
ON CONFLICT (email) DO NOTHING;

INSERT INTO messages (user_id, content, kind) VALUES
    ('00000000-0000-0000-0000-000000000001', 'Welcome to Subspace!', 'info'),
    ('00000000-0000-0000-0000-000000000001', 'This is your first message in the system.', 'info'),
    ('00000000-0000-0000-0000-000000000002', 'Hello from Test User!', 'info')
ON CONFLICT DO NOTHING;
