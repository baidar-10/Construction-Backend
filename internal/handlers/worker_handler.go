package handlers

import (
	"construction-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorkerHandler struct {
	workerService *service.WorkerService
}

func NewWorkerHandler(workerService *service.WorkerService) *WorkerHandler {
	return &WorkerHandler{workerService: workerService}
}

func (h *WorkerHandler) GetAllWorkers(c *gin.Context) {
	filters := make(map[string]interface{})

	if specialty := c.Query("specialty"); specialty != "" {
		filters["specialty"] = specialty
	}
	if location := c.Query("location"); location != "" {
		filters["location"] = location
	}
	if minRate := c.Query("minRate"); minRate != "" {
		if rate, err := strconv.ParseFloat(minRate, 64); err == nil {
			filters["minRate"] = rate
		}
	}
	if maxRate := c.Query("maxRate"); maxRate != "" {
		if rate, err := strconv.ParseFloat(maxRate, 64); err == nil {
			filters["maxRate"] = rate
		}
	}
	if availability := c.Query("availability"); availability != "" {
		filters["availability"] = availability
	}

	workers, err := h.workerService.GetAllWorkers(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workers": workers, "count": len(workers)})
}

func (h *WorkerHandler) GetWorkerByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	worker, err := h.workerService.GetWorkerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Worker not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"worker": worker})
}

func (h *WorkerHandler) SearchWorkers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	workers, err := h.workerService.SearchWorkers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workers": workers, "count": len(workers)})
}

func (h *WorkerHandler) FilterWorkers(c *gin.Context) {
	skill := c.Query("skill")
	if skill == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Skill parameter is required"})
		return
	}

	workers, err := h.workerService.FilterBySkill(skill)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workers": workers, "count": len(workers)})
}

func (h *WorkerHandler) UpdateWorker(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	worker, err := h.workerService.UpdateWorker(id, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"worker": worker, "message": "Worker updated successfully"})
}

func (h *WorkerHandler) AddPortfolio(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageURL    string `json:"imageUrl" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.workerService.AddPortfolio(id, req.Title, req.Description, req.ImageURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Portfolio item added successfully"})
}