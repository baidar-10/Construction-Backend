package handlers

import (
	"construction-backend/internal/models"
	"construction-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	Service *service.MessageService
}

func NewMessageHandler(s *service.MessageService) *MessageHandler {
	return &MessageHandler{Service: s}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	senderID, _ := c.Get("userID")
	var msg models.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg.SenderID = senderID.(uint)
	if err := h.Service.SendMessage(&msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

func (h *MessageHandler) GetConversation(c *gin.Context) {
	userID, _ := c.Get("userID")
	otherUserID, _ := strconv.Atoi(c.Param("userId"))

	messages, err := h.Service.GetConversation(userID.(uint), uint(otherUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}