-- Migration: Add worker promotion system
-- Allows admins to promote workers as TOP/Featured

ALTER TABLE workers ADD COLUMN IF NOT EXISTS is_promoted BOOLEAN DEFAULT FALSE;
ALTER TABLE workers ADD COLUMN IF NOT EXISTS promotion_type VARCHAR(50) DEFAULT 'none' CHECK (promotion_type IN ('none', 'featured', 'top', 'premium'));
ALTER TABLE workers ADD COLUMN IF NOT EXISTS promotion_expires_at TIMESTAMP;
ALTER TABLE workers ADD COLUMN IF NOT EXISTS promotion_payment_date TIMESTAMP;
ALTER TABLE workers ADD COLUMN IF NOT EXISTS promotion_price DECIMAL(10, 2);

-- Create promotion history table
CREATE TABLE IF NOT EXISTS promotion_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    worker_id UUID REFERENCES workers(id) ON DELETE CASCADE,
    promotion_type VARCHAR(50) NOT NULL CHECK (promotion_type IN ('featured', 'top', 'premium')),
    payment_amount DECIMAL(10, 2),
    duration_days INTEGER,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'expired', 'cancelled')),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for promotions
CREATE INDEX IF NOT EXISTS idx_workers_is_promoted ON workers(is_promoted);
CREATE INDEX IF NOT EXISTS idx_workers_promotion_type ON workers(promotion_type);
CREATE INDEX IF NOT EXISTS idx_workers_promotion_expires ON workers(promotion_expires_at);
CREATE INDEX IF NOT EXISTS idx_promotion_history_worker_id ON promotion_history(worker_id);
CREATE INDEX IF NOT EXISTS idx_promotion_history_status ON promotion_history(status);

-- Create trigger to update promotion_history updated_at
CREATE TRIGGER update_promotion_history_updated_at BEFORE UPDATE ON promotion_history FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Promotion pricing table
CREATE TABLE IF NOT EXISTS promotion_pricing (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    promotion_type VARCHAR(50) UNIQUE NOT NULL CHECK (promotion_type IN ('featured', 'top', 'premium')),
    price_per_day DECIMAL(10, 2) NOT NULL,
    min_duration_days INTEGER DEFAULT 7,
    max_duration_days INTEGER DEFAULT 365,
    description TEXT,
    display_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default promotion pricing
INSERT INTO promotion_pricing (promotion_type, price_per_day, min_duration_days, max_duration_days, description, display_order)
VALUES 
    ('featured', 100.00, 7, 30, 'Featured in search results', 1),
    ('top', 250.00, 7, 30, 'Top position in search results', 2),
    ('premium', 500.00, 7, 30, 'Premium position with badge', 3)
ON CONFLICT DO NOTHING;

CREATE TRIGGER update_promotion_pricing_updated_at BEFORE UPDATE ON promotion_pricing FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
