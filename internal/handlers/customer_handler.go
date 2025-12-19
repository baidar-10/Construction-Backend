package handlers

import (
	"construction-backend/internal/middleware"
	"construction-backend/internal/service"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	Service *service.CustomerService
}

func NewCustomerHandler(s *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{Service: s}
}

func (h *CustomerHandler) ToggleFavorite(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var input struct {
		WorkerID string `json:"worker_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workerUUID, err := uuid.Parse(input.WorkerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	if err := h.Service.ToggleFavorite(userID, workerUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update favorites"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Favorites updated"})
}

func (h *CustomerHandler) GetFavorites(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	favorites, err := h.Service.GetFavorites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
		return
	}
	c.JSON(http.StatusOK, favorites)
}

// GetProfile (Specific customer view)
func (h *CustomerHandler) GetProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	customer, err := h.Service.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// === Router-compatible wrapper methods ===

func (h *CustomerHandler) GetCustomerProfile(c *gin.Context) {
	// Just reuse GetProfile
	h.GetProfile(c)
}

func (h *CustomerHandler) GetFavoriteWorkers(c *gin.Context) {
	// Return favorites for the customer ID in path
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	favorites, err := h.Service.GetFavorites(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
		return
	}
	c.JSON(http.StatusOK, favorites)
}

func (h *CustomerHandler) AddFavoriteWorker(c *gin.Context) {
	idStr := c.Param("id")
	custID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}
	var input struct {
		WorkerID string `json:"worker_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workerUUID, err := uuid.Parse(input.WorkerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}
	if err := h.Service.ToggleFavorite(custID, workerUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Favorite added"})
}

func (h *CustomerHandler) RemoveFavoriteWorker(c *gin.Context) {
	idStr := c.Param("id")
	custID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}
	workerIDStr := c.Param("workerId")
	workerID, err := uuid.Parse(workerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}
	if err := h.Service.RemoveFavorite(custID, workerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed"})
}

func (h *CustomerHandler) UpdateCustomerProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *CustomerHandler) GetBookingHistory(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}