package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"time"
	"github.com/google/uuid"
)

type MessageService struct {
	Repo *repository.MessageRepository
}

// Simple Mock if repository isn't fully defined yet
func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{Repo: repo}
}

func (s *MessageService) SendMessage(msg *models.Message) error {
	msg.CreatedAt = time.Now()
	msg.IsRead = false
	return s.Repo.Create(msg)
}

func (s *MessageService) GetConversation(user1, user2 uuid.UUID) ([]models.Message, error) {
	return s.Repo.FindBetweenUsers(user1, user2)
}

func (s *MessageService) MarkAsRead(msgID uuid.UUID) error {
	return s.Repo.MarkAsRead(msgID)
}

func (s *MessageService) GetConversations(userID uuid.UUID) ([]map[string]interface{}, error) {
	return s.Repo.GetConversations(userID)
}