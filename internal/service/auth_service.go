package service

import (
	"construction-backend/internal/models"
	"construction-backend/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	workerRepo   *repository.WorkerRepository
	customerRepo *repository.CustomerRepository
	jwtSecret    string
	db           *gorm.DB
}

func NewAuthService(userRepo *repository.UserRepository, workerRepo *repository.WorkerRepository, customerRepo *repository.CustomerRepository, jwtSecret string, db *gorm.DB) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		workerRepo:   workerRepo,
		customerRepo: customerRepo,
		jwtSecret:    jwtSecret,
		db:           db,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if user exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	// If the error is something else than record not found, return it
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if req.FirstName == "" {
		return nil, errors.New("first name is required")
	}
	if req.LastName == "" {
		return nil, errors.New("last name is required")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		UserType:     req.UserType,
		IsActive:     true,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create profile based on user type
	if req.UserType == "worker" {
		worker := &models.Worker{
			UserID:             user.ID,
			Specialty:          req.Specialty,
			HourlyRate:         req.HourlyRate,
			ExperienceYears:    req.ExperienceYears,
			Bio:                req.Bio,
			Location:           req.Location,
			AvailabilityStatus: "available",
			Rating:             0.0,
			TotalReviews:       0,
			TotalJobs:          0,
		}

		if err := tx.Create(worker).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Add skills if provided
		if len(req.Skills) > 0 {
			for _, skill := range req.Skills {
				workerSkill := &models.WorkerSkill{
					WorkerID: worker.ID,
					Skill:    skill,
				}
				if err := tx.Create(workerSkill).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	} else if req.UserType == "customer" {
		customer := &models.Customer{
			UserID:     user.ID,
			Address:    req.Address,
			City:       req.City,
			State:      req.State,
			PostalCode: req.PostalCode,
		}

		if err := tx.Create(customer).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"userId":   user.ID.String(),
		"email":    user.Email,
		"userType": user.UserType,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) GetCurrentUser(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) ChangePassword(userID uuid.UUID, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}