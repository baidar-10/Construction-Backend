-- Migration: Add verification fields to users

ALTER TABLE users ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS verification_code VARCHAR(10);
ALTER TABLE users ADD COLUMN IF NOT EXISTS code_expires_at TIMESTAMP;
