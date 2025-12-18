package handlers

import (
	"construction-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	Service *service.CustomerService
}

func NewCustomerHandler(s *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{Service: s}
}

func (h *CustomerHandler) ToggleFavorite(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input struct {
		WorkerID uint `json:"worker_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.ToggleFavorite(userID.(uint), input.WorkerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update favorites"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Favorites updated"})
}

func (h *CustomerHandler) GetFavorites(c *gin.Context) {
	userID, _ := c.Get("userID")
	favorites, err := h.Service.GetFavorites(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
		return
	}
	c.JSON(http.StatusOK, favorites)
}

// GetProfile (Specific customer view)
func (h *CustomerHandler) GetProfile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	customer, err := h.Service.GetProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}