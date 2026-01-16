-- Migration: Add support for open bookings
-- Description: Adds is_open field to bookings table and makes worker_id nullable

-- Make worker_id nullable for open bookings
ALTER TABLE bookings ALTER COLUMN worker_id DROP NOT NULL;

-- Add is_open field
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS is_open BOOLEAN DEFAULT FALSE;

-- Create index for faster querying of open bookings
CREATE INDEX IF NOT EXISTS idx_bookings_open ON bookings(is_open, status, worker_id) 
WHERE is_open = true AND status = 'pending' AND worker_id IS NULL;
