package handlers

import (
	"construction-backend/internal/middleware"
	"construction-backend/internal/models"
	"construction-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a booking request for a worker (customers only)
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Booking true "Booking details"
// @Success 201 {object} map[string]interface{} "Booking created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - only customers can create bookings"
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "message": err.Error()})
		return
	}

	// Ensure authenticated user is a customer and set as the booking's customer
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Authentication required"})
		return
	}
	userType, _ := middleware.GetUserTypeFromContext(c)
	if userType != "customer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "Only customers can create bookings"})
		return
	}

	if err := h.bookingService.CreateBookingForUser(userID, &booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create booking", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"booking": booking, "message": "Booking created successfully"})
}

func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user type from context to determine which bookings to fetch
	userType, _ := c.Get("userType")
	bookings, err := h.bookingService.GetUserBookings(userID, userType.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings, "count": len(bookings)})
}

func (h *BookingHandler) GetWorkerBookings(c *gin.Context) {
	workerID, err := uuid.Parse(c.Param("workerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	bookings, err := h.bookingService.GetWorkerBookings(workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings, "count": len(bookings)})
}

func (h *BookingHandler) UpdateBookingStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.bookingService.UpdateBookingStatus(id, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated successfully"})
}
func (h *BookingHandler) AcceptBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	// Ensure authenticated user is a worker
	userType, _ := middleware.GetUserTypeFromContext(c)
	if userType != "worker" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can accept bookings"})
		return
	}

	if err := h.bookingService.AcceptBooking(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking accepted successfully"})
}

func (h *BookingHandler) DeclineBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	// Ensure authenticated user is a worker
	userType, _ := middleware.GetUserTypeFromContext(c)
	if userType != "worker" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can decline bookings"})
		return
	}

	if err := h.bookingService.DeclineBooking(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking declined successfully"})
}
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	if err := h.bookingService.CancelBooking(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}

// CompleteBooking godoc
// @Summary Complete a booking
// @Description Mark a booking as completed by the worker
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID"
// @Success 200 {object} map[string]interface{} "Booking completed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Forbidden - only workers can complete bookings"
// @Router /bookings/{id}/complete [put]
func (h *BookingHandler) CompleteBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	// Ensure authenticated user is a worker
	userType, _ := middleware.GetUserTypeFromContext(c)
	if userType != "worker" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can complete bookings"})
		return
	}

	if err := h.bookingService.CompleteBooking(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking completed successfully"})
}

// GetOpenBookings godoc
// @Summary Get all open bookings
// @Description Get all bookings that are open for any worker to accept
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of open bookings"
// @Failure 403 {object} map[string]interface{} "Forbidden - only workers can view open bookings"
// @Router /bookings/open [get]
func (h *BookingHandler) GetOpenBookings(c *gin.Context) {
// Ensure authenticated user is a worker
userType, _ := middleware.GetUserTypeFromContext(c)
if userType != "worker" {
c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can view open bookings"})
return
}

bookings, err := h.bookingService.GetOpenBookings()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{"bookings": bookings, "count": len(bookings)})
}

// ClaimOpenBooking godoc
// @Summary Claim an open booking
// @Description Worker claims an open booking and becomes assigned to it
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID"
// @Param request body map[string]string true "Worker ID"
// @Success 200 {object} map[string]interface{} "Booking claimed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Forbidden - only workers can claim bookings"
// @Router /bookings/{id}/claim [put]
func (h *BookingHandler) ClaimOpenBooking(c *gin.Context) {
id, err := uuid.Parse(c.Param("id"))
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
return
}

// Ensure authenticated user is a worker
userType, _ := middleware.GetUserTypeFromContext(c)
if userType != "worker" {
c.JSON(http.StatusForbidden, gin.H{"error": "Only workers can claim bookings"})
return
}

var req struct {
WorkerID string `json:"workerId" binding:"required"`
}

if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
return
}

workerID, err := uuid.Parse(req.WorkerID)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
return
}

if err := h.bookingService.ClaimOpenBooking(id, workerID); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{"message": "Booking claimed successfully"})
}
