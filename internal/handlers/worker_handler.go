package handlers

import (
	"construction-backend/internal/models"
	"construction-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorkerHandler struct {
	Service *service.WorkerService
}

func NewWorkerHandler(s *service.WorkerService) *WorkerHandler {
	return &WorkerHandler{Service: s}
}

func (h *WorkerHandler) GetWorker(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	worker, err := h.Service.GetWorkerByID(id)
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

func (h *WorkerHandler) GetAllWorkers(c *gin.Context) {
	filters := make(map[string]interface{})
	if specialty := c.Query("specialty"); specialty != "" {
		filters["specialty"] = specialty
	}
	if location := c.Query("location"); location != "" {
		filters["location"] = location
	}
	// parse more filters as needed
	workers, err := h.Service.ListWorkers(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list workers"})
		return
	}
	c.JSON(http.StatusOK, workers)
}

func (h *WorkerHandler) GetWorkerByID(c *gin.Context) {
	// reuse GetWorker behavior
	h.GetWorker(c)
}

func (h *WorkerHandler) FilterWorkers(c *gin.Context) {
	// alias for search with query & skill
	h.SearchWorkers(c)
}

func (h *WorkerHandler) UpdateWorker(c *gin.Context) {
	h.UpdateProfile(c)
}

func (h *WorkerHandler) UpdateProfile(c *gin.Context) {
	var profile models.Worker
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Note: In a real app, verify the user ID from the token matches the profile
	if err := h.Service.UpdateProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated", "data": profile})
}

func (h *WorkerHandler) AddPortfolio(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}