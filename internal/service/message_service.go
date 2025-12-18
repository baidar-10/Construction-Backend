package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"time"
)

type MessageService struct {
	Repo *repository.MessageRepository // You need to ensure this struct exists in your repository package
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

func (s *MessageService) GetConversation(user1, user2 uint) ([]models.Message, error) {
	return s.Repo.GetConversation(user1, user2)
}

func (s *MessageService) MarkAsRead(msgID uint) error {
	return s.Repo.MarkRead(msgID)
}