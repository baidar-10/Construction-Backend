package service

import (
	"errors"
	"construction-backend/internal/models"
	"construction-backend/internal/repository"

	"github.com/google/uuid"
)

type BookingService struct {
	bookingRepo  *repository.BookingRepository
	customerRepo *repository.CustomerRepository
}

func NewBookingService(bookingRepo *repository.BookingRepository, customerRepo *repository.CustomerRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo, customerRepo: customerRepo}
}

// CreateBooking creates a booking record directly (low-level). Prefer CreateBookingForUser for API use.
func (s *BookingService) CreateBooking(booking *models.Booking) error {
	return s.bookingRepo.Create(booking)
}

// CreateBookingForUser resolves the authenticated user's customer profile and creates a booking associated with that customer
func (s *BookingService) CreateBookingForUser(userID uuid.UUID, booking *models.Booking) error {
	// Find customer profile by user id
	customer, err := s.customerRepo.FindByUserID(userID)
	if err != nil {
		return err
	}

	if customer == nil {
		return errors.New("customer profile not found")
	}

	booking.CustomerID = customer.ID
	return s.bookingRepo.Create(booking)
}

func (s *BookingService) GetBookingByID(id uuid.UUID) (*models.Booking, error) {
	return s.bookingRepo.FindByID(id)
}

func (s *BookingService) GetUserBookings(userID uuid.UUID, userType string) ([]models.Booking, error) {
	// Determine if user is customer or worker and get appropriate bookings
	if userType == "worker" {
		return nil, errors.New("not implemented for worker userID mapping")
	}

	customer, err := s.customerRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return s.bookingRepo.FindByCustomerID(customer.ID)
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