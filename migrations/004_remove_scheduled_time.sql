-- Migration: Remove scheduled_time column from bookings
-- Description: Remove the scheduled_time column as it's no longer needed

ALTER TABLE bookings DROP COLUMN IF EXISTS scheduled_time;
