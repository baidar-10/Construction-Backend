-- Migration: Add booking applications table
-- This allows workers to apply for open bookings

CREATE TABLE IF NOT EXISTS booking_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    worker_id UUID NOT NULL REFERENCES workers(id) ON DELETE CASCADE,
    message TEXT,
    proposed_price DECIMAL(10, 2),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(booking_id, worker_id)
);

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_booking_applications_booking ON booking_applications(booking_id);
CREATE INDEX IF NOT EXISTS idx_booking_applications_worker ON booking_applications(worker_id);
CREATE INDEX IF NOT EXISTS idx_booking_applications_status ON booking_applications(status);

-- Modify bookings table to support open bookings (worker_id can be null)
ALTER TABLE bookings ALTER COLUMN worker_id DROP NOT NULL;

-- Add status 'open' for bookings without assigned worker
-- Update check constraint if exists
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_status_check;
ALTER TABLE bookings ADD CONSTRAINT bookings_status_check 
    CHECK (status IN ('pending', 'open', 'accepted', 'in_progress', 'completed', 'cancelled'));
