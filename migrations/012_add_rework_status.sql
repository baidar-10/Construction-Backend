-- Migration: Add rework_required status to verification_status enum

-- Add new enum value for rework_required
ALTER TYPE verification_status ADD VALUE IF NOT EXISTS 'rework_required' BEFORE 'rejected';
