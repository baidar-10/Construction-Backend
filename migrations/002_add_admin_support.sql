-- Add admin user type support
-- This migration adds the ability for users to be admins

-- First, check if we need to update the user_type check constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_user_type_check;
ALTER TABLE users ADD CONSTRAINT users_user_type_check CHECK (user_type IN ('customer', 'worker', 'admin'));

-- Create a default admin user (password: admin123)
-- Password hash for 'admin123' using bcrypt
INSERT INTO users (id, email, password_hash, full_name, phone_number, user_type, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'admin@stroyhub.com',
    '$2a$10$YourBcryptHashHere', -- You'll need to generate this properly
    'System Administrator',
    '+77001234567',
    'admin',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;
