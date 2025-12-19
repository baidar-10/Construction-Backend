package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"github.com/google/uuid"
)

type CustomerService struct {
	Repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{Repo: repo}
}

func (s *CustomerService) GetProfile(userID uuid.UUID) (*models.Customer, error) {
	return s.Repo.FindByUserID(userID)
}

func (s *CustomerService) ToggleFavorite(customerID, workerID uuid.UUID) error {
	return s.Repo.AddFavorite(customerID, workerID)
}

func (s *CustomerService) RemoveFavorite(customerID, workerID uuid.UUID) error {
	return s.Repo.RemoveFavorite(customerID, workerID)
}

func (s *CustomerService) GetFavorites(customerID uuid.UUID) ([]models.Worker, error) {
	return s.Repo.GetFavorites(customerID)
}