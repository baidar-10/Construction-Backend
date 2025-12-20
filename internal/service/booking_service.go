package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"

	"github.com/google/uuid"
)

type BookingService struct {
	bookingRepo *repository.BookingRepository
}

func NewBookingService(bookingRepo *repository.BookingRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo}
}

func (s *BookingService) CreateBooking(booking *models.Booking) error {
	return s.bookingRepo.Create(booking)
}

func (s *BookingService) GetBookingByID(id uuid.UUID) (*models.Booking, error) {
	return s.bookingRepo.FindByID(id)
}

func (s *BookingService) GetUserBookings(userID uuid.UUID, userType string) ([]models.Booking, error) {
	// Determine if user is customer or worker and get appropriate bookings
	// For now, we'll assume customer ID
	return s.bookingRepo.FindByCustomerID(userID)
}

func (s *BookingService) GetCustomerBookings(customerID uuid.UUID) ([]models.Booking, error) {
	return s.bookingRepo.FindByCustomerID(customerID)
}

func (s *BookingService) GetWorkerBookings(workerID uuid.UUID) ([]models.Booking, error) {
	return s.bookingRepo.FindByWorkerID(workerID)
}

func (s *BookingService) UpdateBookingStatus(id uuid.UUID, status string) error {
	return s.bookingRepo.UpdateStatus(id, status)
}

func (s *BookingService) CancelBooking(id uuid.UUID) error {
	return s.bookingRepo.Delete(id)
}