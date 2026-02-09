package repository

import (
	"construction-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

// Create creates a new verification document
func (r *VerificationRepository) Create(doc *models.VerificationDocument) error {
	return r.db.Create(doc).Error
}

// GetByID retrieves a verification document by ID
func (r *VerificationRepository) GetByID(id uuid.UUID) (*models.VerificationDocument, error) {
	var doc models.VerificationDocument
	err := r.db.Preload("User").Preload("Admin").First(&doc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetByUserID retrieves all verification documents for a user
func (r *VerificationRepository) GetByUserID(userID uuid.UUID) ([]models.VerificationDocument, error) {
	var docs []models.VerificationDocument
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&docs).Error
	return docs, err
}

// GetLatestByUserID retrieves the latest verification document for a user
func (r *VerificationRepository) GetLatestByUserID(userID uuid.UUID) (*models.VerificationDocument, error) {
	var doc models.VerificationDocument
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetPendingDocuments retrieves all pending verification documents (for admin)
func (r *VerificationRepository) GetPendingDocuments() ([]models.VerificationDocument, error) {
	var docs []models.VerificationDocument
	err := r.db.Preload("User").Where("status = ?", "pending").Order("created_at ASC").Find(&docs).Error
	return docs, err
}

// GetAllDocuments retrieves all verification documents with optional status filter (for admin)
func (r *VerificationRepository) GetAllDocuments(status string, limit, offset int) ([]models.VerificationDocument, int64, error) {
	var docs []models.VerificationDocument
	var total int64

	query := r.db.Model(&models.VerificationDocument{}).Preload("User").Preload("Admin")

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&docs).Error
	return docs, total, err
}

// Update updates a verification document
func (r *VerificationRepository) Update(doc *models.VerificationDocument) error {
	return r.db.Save(doc).Error
}

// Delete deletes a verification document
func (r *VerificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.VerificationDocument{}, "id = ?", id).Error
}

// HasPendingDocument checks if user has a pending verification document
func (r *VerificationRepository) HasPendingDocument(userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.VerificationDocument{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&count).Error
	return count > 0, err
}

// UpdateUserVerificationStatus updates the is_identity_verified field on the user
func (r *VerificationRepository) UpdateUserVerificationStatus(userID uuid.UUID, verified bool) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("is_identity_verified", verified).Error
}
