package handlers

import (
	"construction-backend/internal/middleware"
	"construction-backend/internal/models"
	"construction-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	senderID, _ := middleware.GetUserIDFromContext(c)

	var req struct {
		ReceiverID string  `json:"receiverId" binding:"required"`
		Content    string  `json:"content" binding:"required"`
		BookingID  *string `json:"bookingId,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
		return
	}

	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    req.Content,
		IsRead:     false,
	}

	if req.BookingID != nil && *req.BookingID != "" {
		bookingID, err := uuid.Parse(*req.BookingID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
			return
		}
		message.BookingID = &bookingID
	}

	if err := h.messageService.SendMessage(message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": message, "status": "Message sent successfully"})
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID1, _ := middleware.GetUserIDFromContext(c)

	userID2, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	messages, err := h.messageService.GetMessagesBetweenUsers(userID1, userID2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages, "count": len(messages)})
}

func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID, _ := middleware.GetUserIDFromContext(c)

	conversations, err := h.messageService.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations, "count": len(conversations)})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	if err := h.messageService.MarkAsRead(messageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
}

func (h *MessageHandler) GetBookingMessages(c *gin.Context) {
	bookingID, err := uuid.Parse(c.Param("bookingId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	messages, err := h.messageService.GetMessagesByBookingID(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages, "count": len(messages)})
}

func (h *MessageHandler) MarkBookingMessagesAsRead(c *gin.Context) {
	userID, _ := middleware.GetUserIDFromContext(c)
	bookingID, err := uuid.Parse(c.Param("bookingId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	if err := h.messageService.MarkBookingMessagesAsRead(bookingID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Messages marked as read"})
}