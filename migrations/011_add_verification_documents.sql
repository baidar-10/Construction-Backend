-- Migration: Add verification documents table for user identity verification

-- Create enum for document types
DO $$ BEGIN
    CREATE TYPE document_type AS ENUM ('passport', 'id_card', 'driver_license');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create enum for verification status
DO $$ BEGIN
    CREATE TYPE verification_status AS ENUM ('pending', 'approved', 'rejected');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create verification_documents table
CREATE TABLE IF NOT EXISTS verification_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type document_type NOT NULL DEFAULT 'id_card',
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    mime_type VARCHAR(100) NOT NULL DEFAULT 'image/jpeg',
    status verification_status NOT NULL DEFAULT 'pending',
    admin_id UUID REFERENCES users(id),
    admin_comment TEXT,
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_verification_documents_user_id ON verification_documents(user_id);
CREATE INDEX IF NOT EXISTS idx_verification_documents_status ON verification_documents(status);
CREATE INDEX IF NOT EXISTS idx_verification_documents_created_at ON verification_documents(created_at);

-- Add is_identity_verified column to users table if not exists
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_identity_verified BOOLEAN DEFAULT FALSE;

-- Create index for identity verification
CREATE INDEX IF NOT EXISTS idx_users_identity_verified ON users(is_identity_verified);

-- Comment on table
COMMENT ON TABLE verification_documents IS 'Stores user identity verification documents for admin review';
