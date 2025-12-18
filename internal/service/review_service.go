package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
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

func (s *ReviewService) GetWorkerReviews(workerID uint) ([]models.Review, error) {
	return s.Repo.GetByWorkerID(workerID)
}