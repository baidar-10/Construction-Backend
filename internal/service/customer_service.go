package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
)

type CustomerService struct {
	Repo *repository.UserRepository
}

func NewCustomerService(repo *repository.UserRepository) *CustomerService {
	return &CustomerService{Repo: repo}
}

func (s *CustomerService) GetProfile(userID uint) (*models.User, error) {
	return s.Repo.GetByID(userID)
}

func (s *CustomerService) ToggleFavorite(userID, workerID uint) error {
	// In a real app, this would call a specific Favorites repository
	// For now, we'll assume the repository handles the many-to-many relationship
	return s.Repo.AddFavorite(userID, workerID)
}

func (s *CustomerService) GetFavorites(userID uint) ([]models.WorkerProfile, error) {
	return s.Repo.GetFavorites(userID)
}