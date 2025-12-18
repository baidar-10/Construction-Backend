package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"errors"
	"time"
)

type BookingService struct {
	Repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{Repo: repo}
}

func (s *BookingService) CreateBooking(booking *models.Booking) error {
	if booking.StartDate.Before(time.Now()) {
		return errors.New("cannot book dates in the past")
	}
	booking.Status = "pending"
	return s.Repo.Create(booking)
}

func (s *BookingService) GetUserBookings(userID uint, isWorker bool) ([]models.Booking, error) {
	if isWorker {
		return s.Repo.GetByWorkerID(userID)
	}
	return s.Repo.GetByCustomerID(userID)
}

func (s *BookingService) UpdateStatus(bookingID uint, status string) error {
	return s.Repo.UpdateStatus(bookingID, status)
}