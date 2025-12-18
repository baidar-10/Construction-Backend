package repository

import (
	"construction-backend/internal/database"
	"construction-backend/internal/models"

	"github.com/google/uuid"
)

type CustomerRepository struct {
	db *database.Database
}

func NewCustomerRepository(db *database.Database) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) FindByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Preload("User").First(&customer, "id = ?", id).Error
	return &customer, err
}

func (r *CustomerRepository) FindByUserID(userID uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Preload("User").Where("user_id = ?", userID).First(&customer).Error
	return &customer, err
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) GetFavorites(customerID uuid.UUID) ([]models.Worker, error) {
	var favorites []models.FavoriteWorker
	err := r.db.Where("customer_id = ?", customerID).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	var workerIDs []uuid.UUID
	for _, fav := range favorites {
		workerIDs = append(workerIDs, fav.WorkerID)
	}

	var workers []models.Worker
	if len(workerIDs) > 0 {
		err = r.db.Preload("User").Where("id IN ?", workerIDs).Find(&workers).Error
	}

	return workers, err
}

func (r *CustomerRepository) AddFavorite(customerID, workerID uuid.UUID) error {
	favorite := &models.FavoriteWorker{
		CustomerID: customerID,
		WorkerID:   workerID,
	}
	return r.db.Create(favorite).Error
}

func (r *CustomerRepository) RemoveFavorite(customerID, workerID uuid.UUID) error {
	return r.db.Where("customer_id = ? AND worker_id = ?", customerID, workerID).
		Delete(&models.FavoriteWorker{}).Error
}