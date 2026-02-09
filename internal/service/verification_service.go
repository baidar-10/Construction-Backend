package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"construction-backend/internal/storage"

	"github.com/google/uuid"
)

type VerificationService struct {
	repo   *repository.VerificationRepository
	minio  *storage.MinioClient
}

func NewVerificationService(repo *repository.VerificationRepository, minio *storage.MinioClient) *VerificationService {
	return &VerificationService{
		repo:  repo,
		minio: minio,
	}
}

// UploadDocument uploads a verification document for a user
func (s *VerificationService) UploadDocument(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader, documentType string) (*models.VerificationDocument, error) {
	// Check if user already has a pending document
	hasPending, err := s.repo.HasPendingDocument(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check pending documents: %w", err)
	}
	if hasPending {
		return nil, fmt.Errorf("you already have a pending verification request")
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return nil, fmt.Errorf("invalid file type. Only JPEG, PNG and PDF are allowed")
	}

	// Validate file size (max 10MB)
	if header.Size > 10*1024*1024 {
		return nil, fmt.Errorf("file size exceeds 10MB limit")
	}

	// Generate unique file path
	ext := filepath.Ext(header.Filename)
	objectName := fmt.Sprintf("verifications/%s/%s%s", userID.String(), uuid.New().String(), ext)

	// Upload to MinIO
	_, err = s.minio.UploadFile(ctx, objectName, file, header.Size, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Create database record
	doc := &models.VerificationDocument{
		ID:           uuid.New(),
		UserID:       userID,
		DocumentType: documentType,
		FilePath:     objectName,
		FileName:     header.Filename,
		FileSize:     header.Size,
		MimeType:     contentType,
		Status:       "pending",
	}

	if err := s.repo.Create(doc); err != nil {
		// Try to delete uploaded file on failure
		_ = s.minio.DeleteFile(ctx, objectName)
		return nil, fmt.Errorf("failed to save document record: %w", err)
	}

	return doc, nil
}

// GetUserDocuments retrieves all verification documents for a user
func (s *VerificationService) GetUserDocuments(userID uuid.UUID) ([]models.VerificationDocument, error) {
	return s.repo.GetByUserID(userID)
}

// GetVerificationStatus retrieves the verification status for a user
func (s *VerificationService) GetVerificationStatus(userID uuid.UUID) (*models.VerificationStatusResponse, error) {
	docs, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	response := &models.VerificationStatusResponse{
		IsIdentityVerified: false,
		Documents:          docs,
	}

	// Check if any document is approved
	for _, doc := range docs {
		if doc.Status == "approved" {
			response.IsIdentityVerified = true
		}
		if response.LatestStatus == "" {
			response.LatestStatus = doc.Status
		}
	}

	return response, nil
}

// GetPendingDocuments retrieves all pending verification documents (for admin)
func (s *VerificationService) GetPendingDocuments() ([]models.VerificationDocument, error) {
	return s.repo.GetPendingDocuments()
}

// GetAllDocuments retrieves all verification documents with filters (for admin)
func (s *VerificationService) GetAllDocuments(status string, page, limit int) ([]models.VerificationDocument, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetAllDocuments(status, limit, offset)
}

// GetDocumentByID retrieves a verification document by ID
func (s *VerificationService) GetDocumentByID(id uuid.UUID) (*models.VerificationDocument, error) {
	return s.repo.GetByID(id)
}

// GetDocumentURL generates a presigned URL for viewing a document
func (s *VerificationService) GetDocumentURL(ctx context.Context, id uuid.UUID) (string, error) {
	doc, err := s.repo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("document not found: %w", err)
	}

	url, err := s.minio.GetFileURL(ctx, doc.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL: %w", err)
	}

	return url, nil
}

// ApproveDocument approves a verification document
func (s *VerificationService) ApproveDocument(ctx context.Context, docID, adminID uuid.UUID, comment string) error {
	doc, err := s.repo.GetByID(docID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	if doc.Status != "pending" {
		return fmt.Errorf("document has already been reviewed")
	}

	now := time.Now()
	doc.Status = "approved"
	doc.AdminID = &adminID
	doc.AdminComment = comment
	doc.ReviewedAt = &now

	if err := s.repo.Update(doc); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	// Update user's verification status
	if err := s.repo.UpdateUserVerificationStatus(doc.UserID, true); err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}

	return nil
}

// RejectDocument rejects a verification document
func (s *VerificationService) RejectDocument(ctx context.Context, docID, adminID uuid.UUID, comment string) error {
	doc, err := s.repo.GetByID(docID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	if doc.Status != "pending" {
		return fmt.Errorf("document has already been reviewed")
	}

	now := time.Now()
	doc.Status = "rejected"
	doc.AdminID = &adminID
	doc.AdminComment = comment
	doc.ReviewedAt = &now

	if err := s.repo.Update(doc); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// DeleteDocument deletes a verification document (and its file from storage)
func (s *VerificationService) DeleteDocument(ctx context.Context, docID, userID uuid.UUID) error {
	doc, err := s.repo.GetByID(docID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Only allow owner to delete
	if doc.UserID != userID {
		return fmt.Errorf("you don't have permission to delete this document")
	}

	// Only allow deleting pending or rejected documents
	if doc.Status == "approved" {
		return fmt.Errorf("cannot delete an approved document")
	}

	// Delete from MinIO
	if err := s.minio.DeleteFile(ctx, doc.FilePath); err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
	}

	// Delete from database
	return s.repo.Delete(docID)
}

// isValidImageType checks if the content type is a valid image or PDF
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"application/pdf",
	}

	contentType = strings.ToLower(contentType)
	for _, t := range validTypes {
		if t == contentType {
			return true
		}
	}
	return false
}
