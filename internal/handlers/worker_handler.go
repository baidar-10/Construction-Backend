package handlers

import (
	"construction-backend/internal/middleware"
	"construction-backend/internal/service"
	"log"
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
		log.Printf("GetAllWorkers error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to fetch workers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workers": workers, "count": len(workers)})
}

func (h *WorkerHandler) GetWorkerByID(c *gin.Context) {
	param := c.Param("id")
	log.Printf("GetWorkerByID called with param=%s\n", param)
	id, err := uuid.Parse(param)
	if err != nil {
		log.Printf("GetWorkerByID: invalid uuid param=%s err=%v\n", param, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID", "message": "Invalid worker ID"})
		return
	}

	worker, err := h.workerService.GetWorkerByID(id)
	if err != nil {
		log.Printf("GetWorkerByID: worker not found id=%s err=%v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Worker not found", "message": "Worker not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"worker": worker})
}

// GetWorkerByUserID finds a worker profile by the associated user ID
// This endpoint requires authentication and only allows the owner (same user) to access their worker profile.
func (h *WorkerHandler) GetWorkerByUserID(c *gin.Context) {
	param := c.Param("userId")
	log.Printf("GetWorkerByUserID called with param=%s headers.Authorization=%s\n", param, c.GetHeader("Authorization"))

	userID, err := uuid.Parse(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "message": "Invalid user ID"})
		return
	}

	// Ensure requester is authenticated
	authUserID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		log.Printf("GetWorkerByUserID: unauthorized request for userId=%s\n", userID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Authentication required"})
		return
	}

	// Only allow users to access their own worker profile
	if authUserID != userID {
		log.Printf("GetWorkerByUserID: forbidden access - authUser=%s requested=%s\n", authUserID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "You are not authorized to access this resource"})
		return
	}

	worker, err := h.workerService.GetOrCreateWorkerByUserID(userID)
	if err != nil {
		log.Printf("GetOrCreateWorkerByUserID error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to get or create worker profile"})
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
		log.Printf("SearchWorkers error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to search workers"})
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
		log.Printf("FilterBySkill error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to filter workers by skill"})
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