package handlers

import (
	"construction-backend/internal/models"
	"construction-backend/internal/service"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	workerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	// Check if it's a multipart form (with files) or JSON
	contentType := c.GetHeader("Content-Type")
	var review *models.Review
	
	if strings.Contains(contentType, "multipart/form-data") {
		// Handle form data with files
		bookingID, err := uuid.Parse(c.PostForm("bookingId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
			return
		}

		customerID, err := uuid.Parse(c.PostForm("customerId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}

		rating := 0
		fmt.Sscanf(c.PostForm("rating"), "%d", &rating)
		if rating < 1 || rating > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
			return
		}

		review = &models.Review{
			BookingID:  bookingID,
			WorkerID:   workerID,
			CustomerID: customerID,
			Rating:     rating,
			Comment:    c.PostForm("comment"),
			MediaURLs:  models.StringArray{},
		}

		// Handle file uploads
		form, err := c.MultipartForm()
		if err == nil && form.File["media"] != nil {
			files := form.File["media"]
			
			// Create uploads directory if it doesn't exist
			uploadDir := "./uploads/reviews"
			os.MkdirAll(uploadDir, 0755)

			var mediaURLs []string
			for _, file := range files {
				// Validate file type
				ext := strings.ToLower(filepath.Ext(file.Filename))
				validExts := map[string]bool{
					".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
					".mp4": true, ".mov": true, ".avi": true, ".webm": true,
				}
				
				if !validExts[ext] {
					continue
				}

				// Generate unique filename
				filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], ext)
				filepath := filepath.Join(uploadDir, filename)

				// Save file
				if err := c.SaveUploadedFile(file, filepath); err != nil {
					continue
				}

				// Store relative path
				mediaURLs = append(mediaURLs, "/uploads/reviews/"+filename)
			}
			review.MediaURLs = models.StringArray(mediaURLs)
		}
	} else {
		// Handle JSON request (backward compatibility)
		var req struct {
			BookingID  string `json:"bookingId" binding:"required"`
			CustomerID string `json:"customerId" binding:"required"`
			Rating     int    `json:"rating" binding:"required,min=1,max=5"`
			Comment    string `json:"comment"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bookingID, err := uuid.Parse(req.BookingID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
			return
		}

		customerID, err := uuid.Parse(req.CustomerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}

		review = &models.Review{
			BookingID:  bookingID,
			WorkerID:   workerID,
			CustomerID: customerID,
			Rating:     req.Rating,
			Comment:    req.Comment,
			MediaURLs:  models.StringArray{},
		}
	}

	if err := h.reviewService.CreateReview(review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"review": review, "message": "Review created successfully"})
}

func (h *ReviewHandler) GetWorkerReviews(c *gin.Context) {
	workerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	reviews, err := h.reviewService.GetWorkerReviews(workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews, "count": len(reviews)})
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.reviewService.UpdateReview(id, req.Rating, req.Comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.reviewService.DeleteReview(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}