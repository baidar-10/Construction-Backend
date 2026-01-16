package repository

import (
"construction-backend/internal/database"
"construction-backend/internal/models"

"github.com/google/uuid"
)

type BookingRepository struct {
db *database.Database
}

func NewBookingRepository(db *database.Database) *BookingRepository {
return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *models.Booking) error {
return r.db.Create(booking).Error
}

func (r *BookingRepository) FindByID(id uuid.UUID) (*models.Booking, error) {
var booking models.Booking
err := r.db.Preload("Customer.User").Preload("Worker.User").
First(&booking, "id = ?", id).Error
return &booking, err
}

func (r *BookingRepository) FindByCustomerID(customerID uuid.UUID) ([]models.Booking, error) {
var bookings []models.Booking
err := r.db.Preload("Worker.User").Preload("Review").Where("customer_id = ?", customerID).
Order("scheduled_date DESC").Find(&bookings).Error
return bookings, err
}

func (r *BookingRepository) FindByWorkerID(workerID uuid.UUID) ([]models.Booking, error) {
var bookings []models.Booking
err := r.db.Preload("Customer.User").Where("worker_id = ?", workerID).
Order("scheduled_date DESC").Find(&bookings).Error
return bookings, err
}

func (r *BookingRepository) Update(booking *models.Booking) error {
return r.db.Save(booking).Error
}

func (r *BookingRepository) Delete(id uuid.UUID) error {
return r.db.Delete(&models.Booking{}, "id = ?", id).Error
}

func (r *BookingRepository) UpdateStatus(id uuid.UUID, status string) error {
return r.db.Model(&models.Booking{}).Where("id = ?", id).
Update("status", status).Error
}

func (r *BookingRepository) FindOpenBookings() ([]models.Booking, error) {
var bookings []models.Booking
err := r.db.Preload("Customer.User").
Where("is_open = ? AND status = ? AND worker_id IS NULL", true, "pending").
Order("created_at DESC").Find(&bookings).Error
return bookings, err
}

func (r *BookingRepository) ClaimOpenBooking(bookingID uuid.UUID, workerID uuid.UUID) error {
return r.db.Model(&models.Booking{}).Where("id = ? AND is_open = ? AND worker_id IS NULL", bookingID, true).
Updates(map[string]interface{}{
"worker_id": workerID,
"is_open":   false,
}).Error
}
