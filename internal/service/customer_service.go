package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"

	"github.com/google/uuid"
)

type CustomerService struct {
	customerRepo *repository.CustomerRepository
}

func NewCustomerService(customerRepo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) GetCustomerProfile(id uuid.UUID) (*models.Customer, error) {
	return s.customerRepo.FindByID(id)
}

func (s *CustomerService) UpdateCustomerProfile(id uuid.UUID, updates map[string]interface{}) (*models.Customer, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if address, ok := updates["address"].(string); ok {
		customer.Address = address
	}
	if city, ok := updates["city"].(string); ok {
		customer.City = city
	}
	if state, ok := updates["state"].(string); ok {
		customer.State = state
	}
	if postalCode, ok := updates["postalCode"].(string); ok {
		customer.PostalCode = postalCode
	}

	if err := s.customerRepo.Update(customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) GetFavoriteWorkers(customerID uuid.UUID) ([]models.Worker, error) {
	return s.customerRepo.GetFavorites(customerID)
}

func (s *CustomerService) AddFavoriteWorker(customerID, workerID uuid.UUID) error {
	return s.customerRepo.AddFavorite(customerID, workerID)
}

func (s *CustomerService) RemoveFavoriteWorker(customerID, workerID uuid.UUID) error {
	return s.customerRepo.RemoveFavorite(customerID, workerID)
}