ALTER TABLE workers
ADD COLUMN IF NOT EXISTS payment_type VARCHAR(20),
ADD COLUMN IF NOT EXISTS currency VARCHAR(10) DEFAULT 'KZT';

CREATE INDEX IF NOT EXISTS idx_workers_payment_type ON workers(payment_type);
CREATE INDEX IF NOT EXISTS idx_workers_currency ON workers(currency);
