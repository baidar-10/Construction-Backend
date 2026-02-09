package service

import (
	"construction-backend/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PromotionService interface {
	GetPromotionPricing() ([]models.PromotionPricing, error)
	PromoteWorker(workerID uuid.UUID, promotionType string, durationDays int) error
	ExpirePromotions() error
	GetTopWorkers(limit int) ([]models.Worker, error)
	CancelPromotion(workerID uuid.UUID) error
	GetPromotionHistory(workerID uuid.UUID) ([]models.PromotionHistory, error)
	// Promotion requests
	CreatePromotionRequest(workerID uuid.UUID, promotionType string, durationDays int, message string) error
	GetPromotionRequests(status string) ([]models.PromotionRequest, error)
	GetWorkerPromotionRequests(workerID uuid.UUID) ([]models.PromotionRequest, error)
	ApprovePromotionRequest(requestID uuid.UUID, adminID uuid.UUID, notes string) error
	RejectPromotionRequest(requestID uuid.UUID, adminID uuid.UUID, notes string) error
}

type promotionService struct {
	db *gorm.DB
}

func NewPromotionService(db *gorm.DB) PromotionService {
	return &promotionService{db: db}
}

// GetPromotionPricing retrieves all available promotion pricing options
func (s *promotionService) GetPromotionPricing() ([]models.PromotionPricing, error) {
	var pricing []models.PromotionPricing
	result := s.db.Where("is_active = ?", true).Order("display_order ASC").Find(&pricing)
	return pricing, result.Error
}

// PromoteWorker marks a worker as promoted
func (s *promotionService) PromoteWorker(workerID uuid.UUID, promotionType string, durationDays int) error {
	// Start a transaction to ensure atomicity
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get promotion pricing
		var pricing models.PromotionPricing
		result := tx.Where("promotion_type = ? AND is_active = ?", promotionType, true).First(&pricing)
		if result.Error != nil {
			return fmt.Errorf("promotion type not found: %w", result.Error)
		}

		// Calculate expiration date
		expiresAt := time.Now().AddDate(0, 0, durationDays)
		promotionPrice := pricing.PricePerDay * float64(durationDays)

		// Update worker with promotion
		updates := map[string]interface{}{
			"is_promoted":            true,
			"promotion_type":         promotionType,
			"promotion_expires_at":   expiresAt,
			"promotion_payment_date": time.Now(),
			"promotion_price":        promotionPrice,
		}

		if err := tx.Model(&models.Worker{}).Where("id = ?", workerID).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update worker promotion: %w", err)
		}

		// Create history record
		history := models.PromotionHistory{
			WorkerID:      workerID.String(),
			PromotionType: promotionType,
			DurationDays:  durationDays,
			PaymentAmount: &promotionPrice,
			ExpiresAt:     &expiresAt,
			Status:        "active",
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("failed to create promotion history: %w", err)
		}

		return nil
	})
}

// ExpirePromotions checks and expires outdated promotions
func (s *promotionService) ExpirePromotions() error {
	now := time.Now()

	// Start a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Find workers with expired promotions
		var expiredWorkers []models.Worker
		if err := tx.Where("is_promoted = ? AND promotion_expires_at < ?", true, now).Find(&expiredWorkers).Error; err != nil {
			return err
		}

		// Update workers
		if err := tx.Model(&models.Worker{}).
			Where("is_promoted = ? AND promotion_expires_at < ?", true, now).
			Updates(map[string]interface{}{
				"is_promoted":    false,
				"promotion_type": "none",
			}).Error; err != nil {
			return err
		}

		// Update history status for expired promotions
		for _, worker := range expiredWorkers {
			workerID, err := uuid.Parse(worker.ID.String())
			if err != nil {
				continue
			}

			tx.Model(&models.PromotionHistory{}).
				Where("worker_id = ? AND status = ?", workerID.String(), "active").
				Update("status", "expired")
		}

		return nil
	})
}

// GetTopWorkers retrieves promoted workers sorted by promotion type and rating
func (s *promotionService) GetTopWorkers(limit int) ([]models.Worker, error) {
	var workers []models.Worker
	result := s.db.
		Preload("User").
		Where("is_promoted = ? AND promotion_expires_at > ?", true, time.Now()).
		Order("CASE WHEN promotion_type = 'premium' THEN 1 WHEN promotion_type = 'top' THEN 2 WHEN promotion_type = 'featured' THEN 3 ELSE 4 END, rating DESC").
		Limit(limit).
		Find(&workers)
	return workers, result.Error
}

// CancelPromotion removes promotion from a worker
func (s *promotionService) CancelPromotion(workerID uuid.UUID) error {
	// Start a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update worker
		if err := tx.Model(&models.Worker{}).Where("id = ?", workerID).Updates(map[string]interface{}{
			"is_promoted":    false,
			"promotion_type": "none",
		}).Error; err != nil {
			return err
		}

		// Update promotion history status
		if err := tx.Model(&models.PromotionHistory{}).
			Where("worker_id = ? AND status = ?", workerID.String(), "active").
			Update("status", "cancelled").Error; err != nil {
			return err
		}

		return nil
	})
}

// GetPromotionHistory retrieves promotion history for a worker
func (s *promotionService) GetPromotionHistory(workerID uuid.UUID) ([]models.PromotionHistory, error) {
	var history []models.PromotionHistory
	result := s.db.Where("worker_id = ?", workerID.String()).Order("created_at DESC").Find(&history)
	return history, result.Error
}

// CreatePromotionRequest creates a new promotion request from a worker
func (s *promotionService) CreatePromotionRequest(workerID uuid.UUID, promotionType string, durationDays int, message string) error {
	request := models.PromotionRequest{
		WorkerID:      workerID,
		PromotionType: promotionType,
		DurationDays:  durationDays,
		Message:       message,
		Status:        "pending",
	}
	return s.db.Create(&request).Error
}

// GetPromotionRequests retrieves promotion requests, optionally filtered by status
func (s *promotionService) GetPromotionRequests(status string) ([]models.PromotionRequest, error) {
	var requests []models.PromotionRequest
	query := s.db.Preload("Worker.User").Order("created_at DESC")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	result := query.Find(&requests)
	return requests, result.Error
}

// GetWorkerPromotionRequests retrieves all promotion requests for a specific worker
func (s *promotionService) GetWorkerPromotionRequests(workerID uuid.UUID) ([]models.PromotionRequest, error) {
	var requests []models.PromotionRequest
	result := s.db.Where("worker_id = ?", workerID).Order("created_at DESC").Find(&requests)
	return requests, result.Error
}

// ApprovePromotionRequest approves a promotion request and promotes the worker
func (s *promotionService) ApprovePromotionRequest(requestID uuid.UUID, adminID uuid.UUID, notes string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get the request
		var request models.PromotionRequest
		if err := tx.First(&request, requestID).Error; err != nil {
			return fmt.Errorf("request not found: %w", err)
		}

		if request.Status != "pending" {
			return fmt.Errorf("request already processed")
		}

		// Update request status
		now := time.Now()
		if err := tx.Model(&request).Updates(map[string]interface{}{
			"status":      "approved",
			"admin_notes": notes,
			"reviewed_by": adminID,
			"reviewed_at": now,
		}).Error; err != nil {
			return err
		}

		// Promote the worker
		if err := s.PromoteWorker(request.WorkerID, request.PromotionType, request.DurationDays); err != nil {
			return fmt.Errorf("failed to promote worker: %w", err)
		}

		return nil
	})
}

// RejectPromotionRequest rejects a promotion request
func (s *promotionService) RejectPromotionRequest(requestID uuid.UUID, adminID uuid.UUID, notes string) error {
	var request models.PromotionRequest
	if err := s.db.First(&request, requestID).Error; err != nil {
		return fmt.Errorf("request not found: %w", err)
	}

	if request.Status != "pending" {
		return fmt.Errorf("request already processed")
	}

	now := time.Now()
	return s.db.Model(&request).Updates(map[string]interface{}{
		"status":      "rejected",
		"admin_notes": notes,
		"reviewed_by": adminID,
		"reviewed_at": now,
	}).Error
}
