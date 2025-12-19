package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"errors"
	"time"
	"github.com/google/uuid"
)

type BookingService struct {
	Repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{Repo: repo}
}

func (s *BookingService) CreateBooking(booking *models.Booking) error {
	if booking.ScheduledDate.Before(time.Now()) {
		return errors.New("cannot book dates in the past")
	}
	booking.Status = "pending"
	return s.Repo.Create(booking)
}

func (s *BookingService) GetUserBookings(userID uuid.UUID, isWorker bool) ([]models.Booking, error) {
	if isWorker {
		return s.Repo.FindByWorkerID(userID)
	}
	return s.Repo.FindByCustomerID(userID)
}

func (s *BookingService) UpdateStatus(bookingID uuid.UUID, status string) error {
	return s.Repo.UpdateStatus(bookingID, status)
}