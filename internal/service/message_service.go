package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"

	"github.com/google/uuid"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
}

func NewMessageService(messageRepo *repository.MessageRepository) *MessageService {
	return &MessageService{messageRepo: messageRepo}
}

func (s *MessageService) SendMessage(message *models.Message) error {
	return s.messageRepo.Create(message)
}

func (s *MessageService) GetMessagesBetweenUsers(userID1, userID2 uuid.UUID) ([]models.Message, error) {
	return s.messageRepo.FindBetweenUsers(userID1, userID2)
}

func (s *MessageService) GetConversations(userID uuid.UUID) ([]map[string]interface{}, error) {
	return s.messageRepo.GetConversations(userID)
}

func (s *MessageService) MarkAsRead(messageID uuid.UUID) error {
	return s.messageRepo.MarkAsRead(messageID)
}

func (s *MessageService) MarkAllAsRead(senderID, receiverID uuid.UUID) error {
	return s.messageRepo.MarkAllAsRead(senderID, receiverID)
}