-- Add media fields to reviews table
ALTER TABLE reviews 
ADD COLUMN media_urls TEXT[] DEFAULT '{}';

-- Add comment for documentation
COMMENT ON COLUMN reviews.media_urls IS 'Array of URLs for photos/videos attached to the review';
