package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"

	"github.com/google/uuid"
)

type WorkerService struct {
	workerRepo *repository.WorkerRepository
}

func NewWorkerService(workerRepo *repository.WorkerRepository) *WorkerService {
	return &WorkerService{workerRepo: workerRepo}
}

func (s *WorkerService) GetAllWorkers(filters map[string]interface{}) ([]models.Worker, error) {
	return s.workerRepo.FindAll(filters)
}

func (s *WorkerService) GetWorkerByID(id uuid.UUID) (*models.Worker, error) {
	return s.workerRepo.FindByID(id)
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