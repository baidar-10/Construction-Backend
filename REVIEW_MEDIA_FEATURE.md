# Review Photo/Video Upload Feature

## Summary

Successfully added the ability for customers to attach photos and videos when leaving reviews for workers.

## Changes Made

### Backend Changes

1. **Database Migration** (`migrations/011_add_review_media.sql`)
   - Added `media_urls` column to `reviews` table
   - Type: TEXT[] (array of URLs)
   - Stores paths to uploaded images and videos

2. **Model Update** (`internal/models/models.go`)
   - Added `MediaURLs []string` field to `Review` struct
   - Maps to `media_urls` database column

3. **Handler Update** (`internal/handlers/review_handler.go`)
   - Enhanced `CreateReview` to handle both JSON and multipart form data
   - Supports file uploads with validation:
     - Allowed formats: jpg, jpeg, png, gif, webp, mp4, mov, avi, webm
     - Maximum file size: 10MB per file
     - Maximum files: 5 per review
   - Files are saved to `/uploads/reviews/` directory
   - Backward compatible with existing JSON-only requests

### Frontend Changes

1. **Review Modal** (`components/booking/ReviewModal.jsx`)
   - Added file upload UI with drag-and-drop support
   - Live preview of selected images/videos
   - Shows file count (X/5)
   - Individual file removal buttons
   - File type icons (image/video indicators)
   - Sends files via FormData with multipart/form-data

2. **Review Service** (`api/reviewService.js`)
   - Added `createReviewWithMedia()` method
   - Handles multipart form data uploads
   - Sets proper Content-Type header

3. **Translations** (en.json, ru.json, kk.json)
   - `review.photos`: "Photos/Videos"
   - `review.uploadMedia`: "Click to upload photos or videos"
   - `review.maxSize`: "Max 10MB per file, up to 5 files"
   - `review.invalidFiles`: Error message for rejected files

## How It Works

### For Customers:

1. Complete a job and click "Leave a Review"
2. Rate the worker (1-5 stars)
3. Write optional comment
4. **NEW:** Click upload area to add photos/videos
5. Preview shows thumbnails with remove buttons
6. Submit review with media attachments

### File Validation:

- ✅ Images: .jpg, .jpeg, .png, .gif, .webp
- ✅ Videos: .mp4, .mov, .avi, .webm
- ✅ Max size: 10MB per file
- ✅ Max count: 5 files total
- ❌ Other formats rejected with warning

### Storage:

- Files saved in: `Construction-Backend/uploads/reviews/`
- Filename format: `{timestamp}_{uuid}.{ext}`
- URLs stored as array in database
- Example: `["/uploads/reviews/1707753600_a1b2c3d4.jpg"]`

## Migration Applied

```sql
ALTER TABLE reviews ADD COLUMN media_urls TEXT[] DEFAULT '{}';
```

Migration successfully applied to database.

## Testing

To test the feature:

1. Log in as a customer
2. Complete a job
3. Leave a review
4. Upload 1-5 photos or videos
5. Submit and verify files are saved

## Next Steps (Optional Enhancements)

- Display review media on worker profile pages
- Add image/video viewer/gallery modal
- Compress images before upload
- Add client-side image editing (crop, rotate)
- Support for more video formats
- Progress indicator for large file uploads
