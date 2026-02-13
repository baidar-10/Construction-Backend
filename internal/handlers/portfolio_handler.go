package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"construction-backend/internal/middleware"
	"construction-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB sets the database instance for the portfolio handler
func SetDB(database *gorm.DB) {
	db = database
}

// UploadPortfolioItem handles worker portfolio upload
// @Summary Upload portfolio item
// @Description Worker uploads a portfolio item with image
// @Tags Portfolio
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param title formData string true "Portfolio item title"
// @Param description formData string false "Portfolio item description"
// @Param image formData file true "Portfolio image"
// @Success 201 {object} models.Portfolio
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/workers/portfolio [post]
func UploadPortfolioItem(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user is a worker
	var worker models.Worker
	if err := db.Where("user_id = ?", userID).First(&worker).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can upload portfolio items"})
		return
	}

	// Get form data
	title := c.PostForm("title")
	description := c.PostForm("description")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	// Handle multiple file uploads
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one image is required"})
		return
	}

	// Validate max 5 images
	if len(files) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 5 images allowed"})
		return
	}

	var imageURLs []string
	var primaryImageURL string

	for i, file := range files {
		// Validate file type
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG, PNG, and WebP images are allowed"})
			return
		}

		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Each image must be less than 5MB"})
			return
		}

		// Generate unique filename
		filename := fmt.Sprintf("portfolio_%s_%d_%d%s", worker.ID, time.Now().Unix(), i, ext)
		filePath := fmt.Sprintf("./uploads/portfolio/%s", filename)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		imageURL := fmt.Sprintf("/uploads/portfolio/%s", filename)
		imageURLs = append(imageURLs, imageURL)

		// First image is the primary image
		if i == 0 {
			primaryImageURL = imageURL
		}
	}

	// Create portfolio item
	portfolio := models.Portfolio{
		WorkerID:    worker.ID,
		Title:       title,
		Description: description,
		ImageURL:    primaryImageURL,
		ImageURLs:   imageURLs,
		Status:      "pending",
	}

	if err := db.Create(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create portfolio item"})
		return
	}

	c.JSON(http.StatusCreated, portfolio)
}

// GetWorkerPortfolio gets portfolio items for a specific worker
// @Summary Get worker portfolio
// @Description Get all approved portfolio items for a worker (or all for own portfolio)
// @Tags Portfolio
// @Produce json
// @Param id path string true "Worker ID"
// @Param Authorization header string false "Bearer token (optional, required to see own pending items)"
// @Success 200 {array} models.Portfolio
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/workers/{id}/portfolio [get]
func GetWorkerPortfolio(c *gin.Context) {
	workerID := c.Param("id")
	workerUUID, err := uuid.Parse(workerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	log.Printf("[DEBUG] GetWorkerPortfolio called for worker: %s", workerID)

	var portfolio []models.Portfolio
	query := db.Where("worker_id = ?", workerUUID)

	// If user is viewing their own portfolio, show all items
	// Otherwise, only show approved items
	userID, exists := middleware.GetUserIDFromContext(c)
	isOwner := false
	if exists {
		log.Printf("[DEBUG] User ID from context: %s", userID)
		var worker models.Worker
		if err := db.Where("user_id = ? AND id = ?", userID, workerUUID).First(&worker).Error; err == nil {
			isOwner = true
			log.Printf("[DEBUG] User is owner of this worker profile")
		} else {
			log.Printf("[DEBUG] User is NOT owner. Error: %v", err)
		}
	} else {
		log.Printf("[DEBUG] No user ID in context (not authenticated)")
	}

	if !isOwner {
		log.Printf("[DEBUG] Filtering to only show approved items")
		query = query.Where("status = ?", "approved")
	} else {
		log.Printf("[DEBUG] Showing ALL portfolio items (owner view)")
	}

	if err := query.Order("created_at DESC").Find(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portfolio"})
		return
	}

	log.Printf("[DEBUG] Returning %d portfolio items", len(portfolio))
	c.JSON(http.StatusOK, portfolio)
}

// GetPendingPortfolio gets all pending portfolio items for admin review
// @Summary Get pending portfolio items
// @Description Admin endpoint to get all pending portfolio items
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Portfolio
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/portfolio/pending [get]
func GetPendingPortfolio(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user is admin
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil || user.UserType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var portfolio []models.Portfolio
	if err := db.Where("status = ?", "pending").
		Preload("Worker").
		Order("created_at ASC").
		Find(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending portfolio items"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// ApprovePortfolioItem approves a portfolio item
// @Summary Approve portfolio item
// @Description Admin endpoint to approve a portfolio item
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Portfolio item ID"
// @Success 200 {object} models.Portfolio
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/portfolio/{id}/approve [put]
func ApprovePortfolioItem(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user is admin
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil || user.UserType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	portfolioID := c.Param("id")
	portfolioUUID, err := uuid.Parse(portfolioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid portfolio ID"})
		return
	}

	var portfolio models.Portfolio
	if err := db.Where("id = ?", portfolioUUID).First(&portfolio).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio item not found"})
		return
	}

	now := time.Now()
	portfolio.Status = "approved"
	portfolio.ApprovedAt = &now
	portfolio.ApprovedBy = &userID
	portfolio.RejectionReason = ""

	if err := db.Save(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve portfolio item"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// RejectPortfolioItem rejects a portfolio item
// @Summary Reject portfolio item
// @Description Admin endpoint to reject a portfolio item
// @Tags Admin
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Portfolio item ID"
// @Param reason body object true "Rejection reason"
// @Success 200 {object} models.Portfolio
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/portfolio/{id}/reject [put]
func RejectPortfolioItem(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user is admin
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil || user.UserType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	portfolioID := c.Param("id")
	portfolioUUID, err := uuid.Parse(portfolioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid portfolio ID"})
		return
	}

	var portfolio models.Portfolio
	if err := db.Where("id = ?", portfolioUUID).First(&portfolio).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio item not found"})
		return
	}

	var input struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	portfolio.Status = "rejected"
	portfolio.RejectionReason = input.Reason
	portfolio.ApprovedAt = nil
	portfolio.ApprovedBy = nil

	if err := db.Save(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject portfolio item"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// DeletePortfolioItem deletes a portfolio item
// @Summary Delete portfolio item
// @Description Worker can delete their own portfolio item
// @Tags Portfolio
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Portfolio item ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/workers/portfolio/{id} [delete]
func DeletePortfolioItem(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user is a worker
	var worker models.Worker
	if err := db.Where("user_id = ?", userID).First(&worker).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can delete portfolio items"})
		return
	}

	portfolioID := c.Param("id")
	portfolioUUID, err := uuid.Parse(portfolioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid portfolio ID"})
		return
	}

	var portfolio models.Portfolio
	if err := db.Where("id = ? AND worker_id = ?", portfolioUUID, worker.ID).First(&portfolio).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio item not found"})
		return
	}

	if err := db.Delete(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete portfolio item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Portfolio item deleted successfully"})
}
