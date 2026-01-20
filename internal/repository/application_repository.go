package repository

import (
	"construction-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(application *models.BookingApplication) error {
	return r.db.Create(application).Error
}

func (r *ApplicationRepository) FindByID(id uuid.UUID) (*models.BookingApplication, error) {
	var application models.BookingApplication
	err := r.db.Preload("Worker.User").Preload("Booking").First(&application, "id = ?", id).Error
	return &application, err
}

func (r *ApplicationRepository) FindByBookingID(bookingID uuid.UUID) ([]models.BookingApplication, error) {
	var applications []models.BookingApplication
	err := r.db.Preload("Worker.User").Where("booking_id = ?", bookingID).Order("created_at DESC").Find(&applications).Error
	return applications, err
}

func (r *ApplicationRepository) FindByWorkerID(workerID uuid.UUID) ([]models.BookingApplication, error) {
	var applications []models.BookingApplication
	err := r.db.Preload("Booking.Customer.User").Where("worker_id = ?", workerID).Order("created_at DESC").Find(&applications).Error
	return applications, err
}

func (r *ApplicationRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.BookingApplication{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ApplicationRepository) CheckExistingApplication(bookingID, workerID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.BookingApplication{}).Where("booking_id = ? AND worker_id = ?", bookingID, workerID).Count(&count).Error
	return count > 0, err
}

func (r *ApplicationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BookingApplication{}, "id = ?", id).Error
}
