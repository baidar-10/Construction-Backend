package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"errors"
	"github.com/google/uuid"
)

type WorkerService struct {
	Repo *repository.WorkerRepository
}

func NewWorkerService(repo *repository.WorkerRepository) *WorkerService {
	return &WorkerService{Repo: repo}
}

func (s *WorkerService) GetWorkerByID(id uuid.UUID) (*models.Worker, error) {
	return s.Repo.FindByID(id)
}

func (s *WorkerService) SearchWorkers(query string, skill string) ([]models.Worker, error) {
	if skill != "" {
		return s.Repo.FilterBySkill(skill)
	}
	return s.Repo.Search(query)
}

func (s *WorkerService) ListWorkers(filters map[string]interface{}) ([]models.Worker, error) {
	return s.Repo.FindAll(filters)
}

func (s *WorkerService) UpdateProfile(profile *models.Worker) error {
	if profile.ID == uuid.Nil {
		return errors.New("invalid worker profile ID")
	}
	return s.Repo.Update(profile)
}