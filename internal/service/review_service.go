package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"github.com/google/uuid"
)

type ReviewService struct {
	Repo *repository.ReviewRepository
}

func NewReviewService(repo *repository.ReviewRepository) *ReviewService {
	return &ReviewService{Repo: repo}
}

func (s *ReviewService) CreateReview(review *models.Review) error {
	// Logic to recalculate worker average rating could go here
	return s.Repo.Create(review)
}

func (s *ReviewService) GetWorkerReviews(workerID uuid.UUID) ([]models.Review, error) {
	return s.Repo.FindByWorkerID(workerID)
}