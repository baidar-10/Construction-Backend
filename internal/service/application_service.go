package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type ApplicationService struct {
	applicationRepo *repository.ApplicationRepository
	bookingRepo     *repository.BookingRepository
}

func NewApplicationService(applicationRepo *repository.ApplicationRepository, bookingRepo *repository.BookingRepository) *ApplicationService {
	return &ApplicationService{
		applicationRepo: applicationRepo,
		bookingRepo:     bookingRepo,
	}
}

func (s *ApplicationService) CreateApplication(application *models.BookingApplication) error {
	// Check if booking exists and is open
	booking, err := s.bookingRepo.FindByID(application.BookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if !booking.IsOpen {
		return errors.New("this booking is not accepting applications")
	}

	if booking.WorkerID != nil {
		return errors.New("this booking already has an assigned worker")
	}

	// Check if worker already applied
	exists, err := s.applicationRepo.CheckExistingApplication(application.BookingID, application.WorkerID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("you have already applied to this booking")
	}

	return s.applicationRepo.Create(application)
}

func (s *ApplicationService) GetBookingApplications(bookingID uuid.UUID) ([]models.BookingApplication, error) {
	return s.applicationRepo.FindByBookingID(bookingID)
}

func (s *ApplicationService) GetWorkerApplications(workerID uuid.UUID) ([]models.BookingApplication, error) {
	return s.applicationRepo.FindByWorkerID(workerID)
}

func (s *ApplicationService) AcceptApplication(applicationID uuid.UUID) error {
	application, err := s.applicationRepo.FindByID(applicationID)
	if err != nil {
		return err
	}

	if application.Status != "pending" {
		return errors.New("only pending applications can be accepted")
	}

	// Get booking
	booking, err := s.bookingRepo.FindByID(application.BookingID)
	if err != nil {
		return err
	}

	if booking.WorkerID != nil {
		return errors.New("this booking already has an assigned worker")
	}

	// Assign worker to booking and mark as accepted
	if err := s.bookingRepo.ClaimOpenBooking(application.BookingID, application.WorkerID); err != nil {
		return err
	}

	// Update application status to accepted
	if err := s.applicationRepo.UpdateStatus(applicationID, "accepted"); err != nil {
		return err
	}

	// Reject all other pending applications for this booking
	applications, err := s.applicationRepo.FindByBookingID(application.BookingID)
	if err != nil {
		return err
	}

	for _, app := range applications {
		if app.ID != applicationID && app.Status == "pending" {
			s.applicationRepo.UpdateStatus(app.ID, "rejected")
		}
	}

	return nil
}

func (s *ApplicationService) RejectApplication(applicationID uuid.UUID) error {
	application, err := s.applicationRepo.FindByID(applicationID)
	if err != nil {
		return err
	}

	if application.Status != "pending" {
		return errors.New("only pending applications can be rejected")
	}

	return s.applicationRepo.UpdateStatus(applicationID, "rejected")
}
