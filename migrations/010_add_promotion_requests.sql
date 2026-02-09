-- Create promotion requests table
CREATE TABLE IF NOT EXISTS promotion_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    worker_id UUID NOT NULL REFERENCES workers(id) ON DELETE CASCADE,
    promotion_type VARCHAR(50) NOT NULL,
    duration_days INTEGER NOT NULL CHECK (duration_days >= 7 AND duration_days <= 365),
    message TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    admin_notes TEXT,
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_promotion_requests_worker_id ON promotion_requests(worker_id);
CREATE INDEX IF NOT EXISTS idx_promotion_requests_status ON promotion_requests(status);
CREATE INDEX IF NOT EXISTS idx_promotion_requests_created_at ON promotion_requests(created_at DESC);

-- Update timestamp trigger
CREATE OR REPLACE FUNCTION update_promotion_requests_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER promotion_requests_updated_at
    BEFORE UPDATE ON promotion_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_promotion_requests_updated_at();
