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
	// Get promotion pricing
	var pricing models.PromotionPricing
	result := s.db.Where("promotion_type = ? AND is_active = ?", promotionType, true).First(&pricing)
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

	if err := s.db.Model(&models.Worker{}).Where("id = ?", workerID).Updates(updates).Error; err != nil {
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

	if err := s.db.Create(&history).Error; err != nil {
		return fmt.Errorf("failed to create promotion history: %w", err)
	}

	return nil
}

// ExpirePromotions checks and expires outdated promotions
func (s *promotionService) ExpirePromotions() error {
	now := time.Now()
	return s.db.Model(&models.Worker{}).
		Where("is_promoted = ? AND promotion_expires_at < ?", true, now).
		Updates(map[string]interface{}{
			"is_promoted":    false,
			"promotion_type": "none",
		}).Error
}

// GetTopWorkers retrieves promoted workers sorted by promotion type and rating
func (s *promotionService) GetTopWorkers(limit int) ([]models.Worker, error) {
	var workers []models.Worker
	result := s.db.
		Where("is_promoted = ?", true).
		Order("CASE WHEN promotion_type = 'premium' THEN 1 WHEN promotion_type = 'top' THEN 2 WHEN promotion_type = 'featured' THEN 3 END, rating DESC").
		Limit(limit).
		Find(&workers)
	return workers, result.Error
}

// CancelPromotion removes promotion from a worker
func (s *promotionService) CancelPromotion(workerID uuid.UUID) error {
	return s.db.Model(&models.Worker{}).Where("id = ?", workerID).Updates(map[string]interface{}{
		"is_promoted":    false,
		"promotion_type": "none",
	}).Error
}

// GetPromotionHistory retrieves promotion history for a worker
func (s *promotionService) GetPromotionHistory(workerID uuid.UUID) ([]models.PromotionHistory, error) {
	var history []models.PromotionHistory
	result := s.db.Where("worker_id = ?", workerID).Order("created_at DESC").Find(&history)
	return history, result.Error
}
