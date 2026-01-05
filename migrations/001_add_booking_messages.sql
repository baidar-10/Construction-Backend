-- Migration to add booking_id to messages table and update booking status
-- Run this migration if you have an existing database

-- Add booking_id column to messages table
ALTER TABLE messages ADD COLUMN IF NOT EXISTS booking_id UUID REFERENCES bookings(id) ON DELETE CASCADE;

-- Update booking status constraint to include 'accepted' and 'declined'
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_status_check;
ALTER TABLE bookings ADD CONSTRAINT bookings_status_check 
    CHECK (status IN ('pending', 'accepted', 'declined', 'confirmed', 'in_progress', 'completed', 'cancelled'));

-- Create index on booking_id for better query performance
CREATE INDEX IF NOT EXISTS idx_messages_booking_id ON messages(booking_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender_receiver ON messages(sender_id, receiver_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
