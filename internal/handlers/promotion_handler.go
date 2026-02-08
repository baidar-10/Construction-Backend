package handlers

import (
	"construction-backend/internal/models"
	"construction-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PromotionHandler struct {
	promotionService service.PromotionService
}

func NewPromotionHandler(promotionService service.PromotionService) *PromotionHandler {
	return &PromotionHandler{
		promotionService: promotionService,
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
