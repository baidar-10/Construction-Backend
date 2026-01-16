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
		return errors.New("customer profile not found for this user. Please ensure you have a customer profile")
	}

	if customer == nil {
		return errors.New("customer profile not found. Please complete your customer registration")
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

func (s *BookingService) AcceptBooking(id uuid.UUID) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return err
	}

	if booking.Status != "pending" {
		return errors.New("only pending bookings can be accepted")
	}

	return s.bookingRepo.UpdateStatus(id, "accepted")
}

func (s *BookingService) DeclineBooking(id uuid.UUID) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return err
	}

	if booking.Status != "pending" {
		return errors.New("only pending bookings can be declined")
	}

	return s.bookingRepo.UpdateStatus(id, "declined")
}

func (s *BookingService) CompleteBooking(id uuid.UUID) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return err
	}

	if booking.Status != "accepted" {
		return errors.New("only accepted bookings can be marked as completed")
	}

	return s.bookingRepo.UpdateStatus(id, "completed")
}
func (s *BookingService) GetOpenBookings() ([]models.Booking, error) {
return s.bookingRepo.FindOpenBookings()
}

func (s *BookingService) ClaimOpenBooking(bookingID uuid.UUID, workerID uuid.UUID) error {
// Check if booking exists and is still open
booking, err := s.bookingRepo.FindByID(bookingID)
if err != nil {
return err
}

if !booking.IsOpen {
return errors.New("booking is not open for claims")
}

if booking.WorkerID != nil {
return errors.New("booking has already been claimed")
}

if booking.Status != "pending" {
return errors.New("only pending bookings can be claimed")
}

// Claim the booking
return s.bookingRepo.ClaimOpenBooking(bookingID, workerID)
}
