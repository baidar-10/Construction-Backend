-- Migration: Add support for multiple images in portfolio items

-- Add image_urls array column to portfolio_items table
ALTER TABLE portfolio_items ADD COLUMN IF NOT EXISTS image_urls TEXT[];

-- Create index for image_urls
CREATE INDEX IF NOT EXISTS idx_portfolio_items_image_urls ON portfolio_items USING GIN (image_urls);

-- Comment
COMMENT ON COLUMN portfolio_items.image_urls IS 'Array of image URLs for portfolio item (supports multiple images)';
