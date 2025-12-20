package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"

	"github.com/google/uuid"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (s *ReviewService) CreateReview(review *models.Review) error {
	return s.reviewRepo.Create(review)
}

func (s *ReviewService) GetWorkerReviews(workerID uuid.UUID) ([]models.Review, error) {
	return s.reviewRepo.FindByWorkerID(workerID)
}

func (s *ReviewService) GetCustomerReviews(customerID uuid.UUID) ([]models.Review, error) {
	return s.reviewRepo.FindByCustomerID(customerID)
}

func (s *ReviewService) GetReviewByID(id uuid.UUID) (*models.Review, error) {
	return s.reviewRepo.FindByID(id)
}

func (s *ReviewService) UpdateReview(id uuid.UUID, rating int, comment string) error {
	review, err := s.reviewRepo.FindByID(id)
	if err != nil {
		return err
	}

	review.Rating = rating
	review.Comment = comment

	return s.reviewRepo.Update(review)
}

func (s *ReviewService) DeleteReview(id uuid.UUID) error {
	return s.reviewRepo.Delete(id)
}