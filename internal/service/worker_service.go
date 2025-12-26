package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkerService struct {
	workerRepo *repository.WorkerRepository
	userRepo   *repository.UserRepository
}

func NewWorkerService(workerRepo *repository.WorkerRepository, userRepo *repository.UserRepository) *WorkerService {
	return &WorkerService{workerRepo: workerRepo, userRepo: userRepo}
}

// GetOrCreateWorkerByUserID returns the worker profile for the given user id or creates a minimal profile if not found
func (s *WorkerService) GetOrCreateWorkerByUserID(userID uuid.UUID) (*models.Worker, error) {
	// Try to find existing worker
	worker, err := s.workerRepo.FindByUserID(userID)
	if err == nil && worker != nil {
		return worker, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Ensure user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Create minimal worker profile
	newWorker := &models.Worker{
		UserID:             user.ID,
		Specialty:          "General",
		HourlyRate:         0.0,
		ExperienceYears:    0,
		Bio:                "",
		Location:           "",
		AvailabilityStatus: "available",
		Rating:             0.0,
		TotalReviews:       0,
		TotalJobs:          0,
	}

	if err := s.workerRepo.Create(newWorker); err != nil {
		return nil, err
	}

	return newWorker, nil
}

func (s *WorkerService) GetAllWorkers(filters map[string]interface{}) ([]models.Worker, error) {
	return s.workerRepo.FindAll(filters)
}

func (s *WorkerService) GetWorkerByID(id uuid.UUID) (*models.Worker, error) {
	return s.workerRepo.FindByID(id)
}

func (s *WorkerService) GetWorkerByUserID(userID uuid.UUID) (*models.Worker, error) {
	return s.workerRepo.FindByUserID(userID)
}

func (s *WorkerService) SearchWorkers(query string) ([]models.Worker, error) {
	return s.workerRepo.Search(query)
}

func (s *WorkerService) FilterBySkill(skill string) ([]models.Worker, error) {
	return s.workerRepo.FilterBySkill(skill)
}

func (s *WorkerService) UpdateWorker(id uuid.UUID, updates map[string]interface{}) (*models.Worker, error) {
	worker, err := s.workerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if specialty, ok := updates["specialty"].(string); ok {
		worker.Specialty = specialty
	}
	if hourlyRate, ok := updates["hourlyRate"].(float64); ok {
		worker.HourlyRate = hourlyRate
	}
	if experienceYears, ok := updates["experienceYears"].(float64); ok {
		worker.ExperienceYears = int(experienceYears)
	}
	if bio, ok := updates["bio"].(string); ok {
		worker.Bio = bio
	}
	if location, ok := updates["location"].(string); ok {
		worker.Location = location
	}
	if availabilityStatus, ok := updates["availabilityStatus"].(string); ok {
		worker.AvailabilityStatus = availabilityStatus
	}

	if err := s.workerRepo.Update(worker); err != nil {
		return nil, err
	}

	// Handle skills update
	if skills, ok := updates["skills"].([]interface{}); ok {
		for _, skill := range skills {
			if skillStr, ok := skill.(string); ok {
				s.workerRepo.AddSkill(id, skillStr)
			}
		}
	}

	return s.workerRepo.FindByID(id)
}

func (s *WorkerService) AddPortfolio(workerID uuid.UUID, title, description, imageURL string) error {
	portfolio := &models.Portfolio{
		WorkerID:    workerID,
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
	}
	return s.workerRepo.AddPortfolio(portfolio)
}