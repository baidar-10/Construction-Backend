package repository

import (
	"construction-backend/internal/database"
	"construction-backend/internal/models"

	"github.com/google/uuid"
)

type ReviewRepository struct {
	db *database.Database
}

func NewReviewRepository(db *database.Database) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) FindByWorkerID(workerID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("Customer.User").Where("worker_id = ?", workerID).
		Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) FindByCustomerID(customerID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("Worker.User").Where("customer_id = ?", customerID).
		Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) FindByID(id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("Customer.User").Preload("Worker.User").
		First(&review, "id = ?", id).Error
	return &review, err
}

func (r *ReviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *ReviewRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Review{}, "id = ?", id).Error
}