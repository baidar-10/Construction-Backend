package handlers

import (
	"net/http"
	"strconv"

	"construction-backend/internal/models"
	"construction-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VerificationHandler struct {
	service *service.VerificationService
}

func NewVerificationHandler(service *service.VerificationService) *VerificationHandler {
	return &VerificationHandler{service: service}
}

// UploadDocument handles document upload for verification
// @Summary Upload verification document
// @Description Upload an identity document for verification
// @Tags verification
// @Accept multipart/form-data
// @Produce json
// @Param document formance file true "Document file (JPEG, PNG, or PDF)"
// @Param documentType formData string true "Document type" Enums(passport, id_card, driver_license)
// @Success 201 {object} models.VerificationDocument
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/verification/upload [post]
func (h *VerificationHandler) UploadDocument(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse multipart form
	file, header, err := c.Request.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	documentType := c.PostForm("documentType")
	if documentType == "" {
		documentType = "id_card"
	}

	// Validate document type
	validTypes := map[string]bool{"passport": true, "id_card": true, "driver_license": true}
	if !validTypes[documentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document type"})
		return
	}

	doc, err := h.service.UploadDocument(c.Request.Context(), userID.(uuid.UUID), file, header, documentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// GetMyDocuments retrieves the current user's verification documents
// @Summary Get my verification documents
// @Description Get all verification documents for the authenticated user
// @Tags verification
// @Produce json
// @Success 200 {array} models.VerificationDocument
// @Failure 401 {object} map[string]string
// @Router /api/verification/my-documents [get]
func (h *VerificationHandler) GetMyDocuments(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	docs, err := h.service.GetUserDocuments(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, docs)
}

// GetVerificationStatus retrieves the current user's verification status
// @Summary Get verification status
// @Description Get the verification status for the authenticated user
// @Tags verification
// @Produce json
// @Success 200 {object} models.VerificationStatusResponse
// @Failure 401 {object} map[string]string
// @Router /api/verification/status [get]
func (h *VerificationHandler) GetVerificationStatus(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	status, err := h.service.GetVerificationStatus(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch verification status"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// DeleteDocument deletes a verification document
// @Summary Delete verification document
// @Description Delete a pending or rejected verification document
// @Tags verification
// @Param id path string true "Document ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/verification/{id} [delete]
func (h *VerificationHandler) DeleteDocument(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	err = h.service.DeleteDocument(c.Request.Context(), docID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// ============ ADMIN ENDPOINTS ============

// GetAllVerifications retrieves all verification documents (admin only)
// @Summary Get all verifications (Admin)
// @Description Get all verification documents with optional status filter
// @Tags admin
// @Produce json
// @Param status query string false "Filter by status (pending, approved, rejected, all)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/admin/verifications [get]
func (h *VerificationHandler) GetAllVerifications(c *gin.Context) {
	status := c.DefaultQuery("status", "pending")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	docs, total, err := h.service.GetAllDocuments(status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch verifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents":  docs,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

// GetDocumentURL generates a presigned URL for viewing a document (admin only)
// @Summary Get document URL (Admin)
// @Description Generate a presigned URL to view a verification document
// @Tags admin
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/admin/verifications/{id}/url [get]
func (h *VerificationHandler) GetDocumentURL(c *gin.Context) {
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	url, err := h.service.GetDocumentURL(c.Request.Context(), docID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// ApproveVerification approves a verification document (admin only)
// @Summary Approve verification (Admin)
// @Description Approve a pending verification document
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Param body body models.ReviewVerificationRequest false "Review details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/admin/verifications/{id}/approve [post]
func (h *VerificationHandler) ApproveVerification(c *gin.Context) {
	adminID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req models.ReviewVerificationRequest
	_ = c.ShouldBindJSON(&req) // Comment is optional

	err = h.service.ApproveDocument(c.Request.Context(), docID, adminID.(uuid.UUID), req.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification approved successfully"})
}

// RejectVerification rejects a verification document (admin only)
// @Summary Reject verification (Admin)
// @Description Reject a pending verification document
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Param body body models.ReviewVerificationRequest true "Review details with rejection reason"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/admin/verifications/{id}/reject [post]
func (h *VerificationHandler) RejectVerification(c *gin.Context) {
	adminID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req models.ReviewVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment is required when rejecting"})
		return
	}

	if req.Comment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a reason for rejection"})
		return
	}

	err = h.service.RejectDocument(c.Request.Context(), docID, adminID.(uuid.UUID), req.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification rejected"})
}
