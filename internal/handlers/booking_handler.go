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

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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