package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminService struct {
	db          *gorm.DB
	userRepo    *repository.UserRepository
	workerRepo  *repository.WorkerRepository
	bookingRepo *repository.BookingRepository
	reviewRepo  *repository.ReviewRepository
}

func NewAdminService(
	db *gorm.DB,
	userRepo *repository.UserRepository,
	workerRepo *repository.WorkerRepository,
	bookingRepo *repository.BookingRepository,
	reviewRepo *repository.ReviewRepository,
) *AdminService {
	return &AdminService{
		db:          db,
		userRepo:    userRepo,
		workerRepo:  workerRepo,
		bookingRepo: bookingRepo,
		reviewRepo:  reviewRepo,
	}
}

type DashboardStats struct {
	TotalUsers     int64 `json:"totalUsers"`
	TotalWorkers   int64 `json:"totalWorkers"`
	TotalCustomers int64 `json:"totalCustomers"`
	TotalBookings  int64 `json:"totalBookings"`
	ActiveBookings int64 `json:"activeBookings"`
	TotalReviews   int64 `json:"totalReviews"`
}

func (s *AdminService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Total users
	if err := s.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Total workers
	if err := s.db.Model(&models.Worker{}).Count(&stats.TotalWorkers).Error; err != nil {
		return nil, err
	}

	// Total customers
	if err := s.db.Model(&models.Customer{}).Count(&stats.TotalCustomers).Error; err != nil {
		return nil, err
	}

	// Total bookings
	if err := s.db.Model(&models.Booking{}).Count(&stats.TotalBookings).Error; err != nil {
		return nil, err
	}

	// Active bookings (pending, accepted, in_progress)
	if err := s.db.Model(&models.Booking{}).
		Where("status IN ?", []string{"pending", "accepted", "in_progress"}).
		Count(&stats.ActiveBookings).Error; err != nil {
		return nil, err
	}

	// Total reviews
	if err := s.db.Model(&models.Review{}).Count(&stats.TotalReviews).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *AdminService) GetAllUsers(page, limit int, userType string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit
	query := s.db.Model(&models.User{})

	if userType != "" {
		query = query.Where("user_type = ?", userType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *AdminService) ToggleUserVerification(userID uuid.UUID) error {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	user.IsVerified = !user.IsVerified
	return s.db.Save(&user).Error
}

func (s *AdminService) ToggleUserStatus(userID uuid.UUID) error {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	user.IsActive = !user.IsActive
	return s.db.Save(&user).Error
}

func (s *AdminService) DeleteUser(userID uuid.UUID) error {
	// Delete user and all related data (cascade)
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete worker or customer profile
		if err := tx.Where("user_id = ?", userID).Delete(&models.Worker{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&models.Customer{}).Error; err != nil {
			return err
		}

		// Delete user
		if err := tx.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *AdminService) GetAllBookings(page, limit int, status string) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	offset := (page - 1) * limit
	query := s.db.Model(&models.Booking{}).
		Preload("Customer").
		Preload("Customer.User").
		Preload("Worker").
		Preload("Worker.User")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&bookings).Error; err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

func (s *AdminService) VerifyWorker(workerID uuid.UUID) error {
	// Add a verified field to worker model in the future
	// For now, just check if worker exists
	var worker models.Worker
	if err := s.db.First(&worker, "id = ?", workerID).Error; err != nil {
		return fmt.Errorf("worker not found")
	}

	// You can add a "verified" boolean field to the workers table
	// and update it here
	return nil
}
