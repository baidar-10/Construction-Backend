package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"errors"
)

type WorkerService struct {
	Repo *repository.WorkerRepository
}

func NewWorkerService(repo *repository.WorkerRepository) *WorkerService {
	return &WorkerService{Repo: repo}
}

func (s *WorkerService) GetWorkerByID(id uint) (*models.WorkerProfile, error) {
	return s.Repo.GetByID(id)
}

func (s *WorkerService) SearchWorkers(query string, skill string) ([]models.WorkerProfile, error) {
	return s.Repo.Search(query, skill)
}

func (s *WorkerService) UpdateProfile(profile *models.WorkerProfile) error {
	if profile.ID == 0 {
		return errors.New("invalid worker profile ID")
	}
	return s.Repo.Update(profile)
}