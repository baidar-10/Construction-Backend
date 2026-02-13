package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"

	"github.com/google/uuid"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	workerRepo *repository.WorkerRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository, workerRepo *repository.WorkerRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		workerRepo: workerRepo,
	}
}

func (s *ReviewService) CreateReview(review *models.Review) error {
	// Create the review
	if err := s.reviewRepo.Create(review); err != nil {
		return err
	}

	// Update worker's rating and totalReviews
	return s.updateWorkerRating(review.WorkerID)
}

// updateWorkerRating recalculates and updates the worker's average rating and total reviews
func (s *ReviewService) updateWorkerRating(workerID uuid.UUID) error {
	// Get all reviews for the worker
	reviews, err := s.reviewRepo.FindByWorkerID(workerID)
	if err != nil {
		return err
	}

	// Calculate average rating
	var totalRating float64
	reviewCount := len(reviews)

	for _, review := range reviews {
		totalRating += float64(review.Rating)
	}

	var avgRating float64
	if reviewCount > 0 {
		avgRating = totalRating / float64(reviewCount)
	}

	// Get the worker and update their rating
	worker, err := s.workerRepo.FindByID(workerID)
	if err != nil {
		return err
	}

	worker.Rating = avgRating
	worker.TotalReviews = reviewCount

	return s.workerRepo.Update(worker)
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
	// Get the review first to know which worker to update
	review, err := s.reviewRepo.FindByID(id)
	if err != nil {
		return err
	}

	workerID := review.WorkerID

	// Delete the review
	if err := s.reviewRepo.Delete(id); err != nil {
		return err
	}

	// Update worker's rating after deletion
	return s.updateWorkerRating(workerID)
}