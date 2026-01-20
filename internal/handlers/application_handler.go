package handlers

import (
	"construction-backend/internal/middleware"
	"construction-backend/internal/models"
	"construction-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApplicationHandler struct {
	applicationService *service.ApplicationService
	workerService      *service.WorkerService
}

func NewApplicationHandler(applicationService *service.ApplicationService, workerService *service.WorkerService) *ApplicationHandler {
	return &ApplicationHandler{
		applicationService: applicationService,
		workerService:      workerService,
	}
}

// CreateApplication godoc
// @Summary Create a job application
// @Description Worker applies to an open booking
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BookingApplication true "Application details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /applications [post]
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userType, _ := middleware.GetUserTypeFromContext(c)
	if userType != "worker" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can apply to bookings"})
		return
	}

	// Get worker profile
	worker, err := h.workerService.GetOrCreateWorkerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get worker profile"})
		return
	}

	var application models.BookingApplication
	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	application.WorkerID = worker.ID

	if err := h.applicationService.CreateApplication(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"application": application, "message": "Application submitted successfully"})
}

// GetBookingApplications godoc
// @Summary Get applications for a booking
// @Description Get all applications for a specific booking (customer only)
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param bookingId path string true "Booking ID"
// @Success 200 {object} map[string]interface{}
// @Router /applications/booking/{bookingId} [get]
func (h *ApplicationHandler) GetBookingApplications(c *gin.Context) {
	bookingID, err := uuid.Parse(c.Param("bookingId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	applications, err := h.applicationService.GetBookingApplications(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

// GetWorkerApplications godoc
// @Summary Get worker's applications
// @Description Get all applications submitted by the authenticated worker
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /applications/my [get]
func (h *ApplicationHandler) GetWorkerApplications(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	worker, err := h.workerService.GetOrCreateWorkerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get worker profile"})
		return
	}

	applications, err := h.applicationService.GetWorkerApplications(worker.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

// AcceptApplication godoc
// @Summary Accept an application
// @Description Customer accepts a worker's application
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param applicationId path string true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Router /applications/{applicationId}/accept [put]
func (h *ApplicationHandler) AcceptApplication(c *gin.Context) {
	applicationID, err := uuid.Parse(c.Param("applicationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	userType, _ := middleware.GetUserTypeFromContext(c)
	if userType != "customer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only customers can accept applications"})
		return
	}

	if err := h.applicationService.AcceptApplication(applicationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application accepted successfully"})
}

// RejectApplication godoc
// @Summary Reject an application
// @Description Customer rejects a worker's application
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param applicationId path string true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Router /applications/{applicationId}/reject [put]
func (h *ApplicationHandler) RejectApplication(c *gin.Context) {
	applicationID, err := uuid.Parse(c.Param("applicationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	userType, _ := middleware.GetUserTypeFromContext(c)
	if userType != "customer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only customers can reject applications"})
		return
	}

	if err := h.applicationService.RejectApplication(applicationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application rejected successfully"})
}
