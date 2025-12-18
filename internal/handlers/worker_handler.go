package handlers

import (
	"construction-backend/internal/models"
	"construction-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkerHandler struct {
	Service *service.WorkerService
}

func NewWorkerHandler(s *service.WorkerService) *WorkerHandler {
	return &WorkerHandler{Service: s}
}

func (h *WorkerHandler) GetWorker(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	worker, err := h.Service.GetWorkerByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Worker not found"})
		return
	}
	c.JSON(http.StatusOK, worker)
}

func (h *WorkerHandler) SearchWorkers(c *gin.Context) {
	query := c.Query("q")
	skill := c.Query("skill")
	workers, err := h.Service.SearchWorkers(query, skill)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}
	c.JSON(http.StatusOK, workers)
}

func (h *WorkerHandler) UpdateProfile(c *gin.Context) {
	var profile models.WorkerProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Note: In a real app, you should verify the user ID from the token matches the profile
	if err := h.Service.UpdateProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated", "data": profile})
}