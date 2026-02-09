package handlers

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"construction-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PromotionHandler struct {
	promotionService service.PromotionService
	workerRepo       *repository.WorkerRepository
}

func NewPromotionHandler(promotionService service.PromotionService, workerRepo *repository.WorkerRepository) *PromotionHandler {
	return &PromotionHandler{
		promotionService: promotionService,
		workerRepo:       workerRepo,
	}
}

// GetPromotionPricing returns all available promotion options
// @Summary Get promotion pricing
// @Description Get all available promotion pricing options
// @Tags promotions
// @Produce json
// @Success 200 {array} models.PromotionPricing
// @Router /promotions/pricing [get]
func (h *PromotionHandler) GetPromotionPricing(c *gin.Context) {
	pricing, err := h.promotionService.GetPromotionPricing()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch promotion pricing"})
		return
	}
	c.JSON(http.StatusOK, pricing)
}

// GetTopWorkers returns promoted workers
// @Summary Get top promoted workers
// @Description Get workers with active promotions sorted by promotion level
// @Tags promotions
// @Produce json
// @Param limit query int false "Number of workers to return" default(10)
// @Success 200 {array} models.Worker
// @Router /promotions/top-workers [get]
func (h *PromotionHandler) GetTopWorkers(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		_, _ = c.GetQuery("limit")
	}

	workers, err := h.promotionService.GetTopWorkers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top workers"})
		return
	}
	c.JSON(http.StatusOK, workers)
}

// PromoteWorker promotes a worker (admin only)
// @Summary Promote a worker
// @Description Promote a worker with a specific promotion type (admin only)
// @Tags promotions
// @Accept json
// @Produce json
// @Param workerId path string true "Worker ID"
// @Param request body models.PromoteWorkerRequest true "Promotion details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /admin/workers/{workerId}/promote [post]
func (h *PromotionHandler) PromoteWorker(c *gin.Context) {
	workerID := c.Param("workerId")
	parsedID, err := uuid.Parse(workerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	var req models.PromoteWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.promotionService.PromoteWorker(parsedID, req.PromotionType, req.DurationDays); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Worker promoted successfully"})
}

// CancelPromotion cancels a worker's promotion (admin only)
// @Summary Cancel worker promotion
// @Description Cancel an active promotion for a worker (admin only)
// @Tags promotions
// @Produce json
// @Param workerId path string true "Worker ID"
// @Success 200 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /admin/workers/{workerId}/cancel-promotion [post]
func (h *PromotionHandler) CancelPromotion(c *gin.Context) {
	workerID := c.Param("workerId")
	parsedID, err := uuid.Parse(workerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	if err := h.promotionService.CancelPromotion(parsedID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel promotion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion cancelled successfully"})
}

// GetPromotionHistory returns promotion history for a worker
// @Summary Get promotion history
// @Description Get all promotion history for a specific worker
// @Tags promotions
// @Produce json
// @Param workerId path string true "Worker ID"
// @Success 200 {array} models.PromotionHistory
// @Router /workers/{workerId}/promotion-history [get]
func (h *PromotionHandler) GetPromotionHistory(c *gin.Context) {
	workerID := c.Param("workerId")
	parsedID, err := uuid.Parse(workerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	history, err := h.promotionService.GetPromotionHistory(parsedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch promotion history"})
		return
	}
	c.JSON(http.StatusOK, history)
}

// CreatePromotionRequest creates a new promotion request from a worker
// @Summary Request promotion
// @Description Worker requests a promotion (worker only)
// @Tags promotions
// @Accept json
// @Produce json
// @Param request body models.PromoteWorkerRequest true "Promotion request details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /workers/request-promotion [post]
func (h *PromotionHandler) CreatePromotionRequest(c *gin.Context) {
	userIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is a worker
	userType, _ := c.Get("userType")
	if userType != "worker" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can request promotions"})
		return
	}

	// Get worker profile
	worker, err := h.workerRepo.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Worker profile not found"})
		return
	}

	var req struct {
		PromotionType string `json:"promotionType" binding:"required"`
		DurationDays  int    `json:"durationDays" binding:"required,min=7,max=365"`
		Message       string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.promotionService.CreatePromotionRequest(worker.ID, req.PromotionType, req.DurationDays, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create promotion request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion request submitted successfully"})
}

// GetPromotionRequests returns all promotion requests (admin only)
// @Summary Get promotion requests
// @Description Get all promotion requests (admin only)
// @Tags promotions
// @Produce json
// @Param status query string false "Filter by status (pending/approved/rejected)"
// @Success 200 {array} models.PromotionRequest
// @Router /admin/promotion-requests [get]
func (h *PromotionHandler) GetPromotionRequests(c *gin.Context) {
	status := c.Query("status")
	requests, err := h.promotionService.GetPromotionRequests(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch promotion requests"})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// GetWorkerPromotionRequests returns promotion requests for current worker
// @Summary Get my promotion requests
// @Description Get promotion requests for the authenticated worker
// @Tags promotions
// @Produce json
// @Success 200 {array} models.PromotionRequest
// @Router /workers/my-promotion-requests [get]
func (h *PromotionHandler) GetWorkerPromotionRequests(c *gin.Context) {
	userIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is a worker
	userType, _ := c.Get("userType")
	if userType != "worker" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can access this"})
		return
	}

	// Get worker profile
	worker, err := h.workerRepo.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Worker profile not found"})
		return
	}

	requests, err := h.promotionService.GetWorkerPromotionRequests(worker.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch promotion requests"})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// ApprovePromotionRequest approves a promotion request (admin only)
// @Summary Approve promotion request
// @Description Approve a worker's promotion request (admin only)
// @Tags promotions
// @Accept json
// @Produce json
// @Param requestId path string true "Request ID"
// @Param notes body object false "Admin notes"
// @Success 200 {object} gin.H
// @Router /admin/promotion-requests/{requestId}/approve [post]
func (h *PromotionHandler) ApprovePromotionRequest(c *gin.Context) {
	requestID := c.Param("requestId")
	parsedRequestID, err := uuid.Parse(requestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	adminIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminID, ok := adminIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin ID"})
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	c.ShouldBindJSON(&req)

	if err := h.promotionService.ApprovePromotionRequest(parsedRequestID, adminID, req.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion request approved and worker promoted"})
}

// RejectPromotionRequest rejects a promotion request (admin only)
// @Summary Reject promotion request
// @Description Reject a worker's promotion request (admin only)
// @Tags promotions
// @Accept json
// @Produce json
// @Param requestId path string true "Request ID"
// @Param notes body object true "Rejection reason"
// @Success 200 {object} gin.H
// @Router /admin/promotion-requests/{requestId}/reject [post]
func (h *PromotionHandler) RejectPromotionRequest(c *gin.Context) {
	requestID := c.Param("requestId")
	parsedRequestID, err := uuid.Parse(requestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	adminIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminID, ok := adminIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin ID"})
		return
	}

	var req struct {
		Notes string `json:"notes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rejection reason is required"})
		return
	}

	if err := h.promotionService.RejectPromotionRequest(parsedRequestID, adminID, req.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion request rejected"})
}
