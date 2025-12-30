package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	FirstName    string    `json:"firstName" gorm:"not null"`
	LastName     string    `json:"lastName" gorm:"not null"`
	Phone        string    `json:"phone"`
	UserType     string    `json:"userType" gorm:"not null"` // 'customer' or 'worker'
	AvatarURL    string    `json:"avatarUrl"`
	IsActive     bool      `json:"isActive" gorm:"default:true"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Worker struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID             uuid.UUID `json:"userId" gorm:"type:uuid;unique;not null"`
	User               User      `json:"user" gorm:"foreignKey:UserID"`
	Specialty          string    `json:"specialty" gorm:"not null"`
	HourlyRate         float64   `json:"hourlyRate"`
	ExperienceYears    int       `json:"experienceYears"`
	Bio                string    `json:"bio"`
	Location           string    `json:"location"`
	AvailabilityStatus string    `json:"availabilityStatus" gorm:"default:'available'"`
	Rating             float64   `json:"rating" gorm:"default:0.0"`
	TotalReviews       int       `json:"totalReviews" gorm:"default:0"`
	TotalJobs          int       `json:"totalJobs" gorm:"default:0"`
	Skills             []string  `json:"skills" gorm:"-"`
	Portfolio          []Portfolio `json:"portfolio,omitempty" gorm:"foreignKey:WorkerID"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type WorkerSkill struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	WorkerID  uuid.UUID `json:"workerId" gorm:"type:uuid;not null"`
	Skill     string    `json:"skill" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
}

type Portfolio struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	WorkerID    uuid.UUID `json:"workerId" gorm:"type:uuid;not null"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
}

// TableName overrides the default table name for GORM
func (Portfolio) TableName() string {
	return "worker_portfolio"
}

type Customer struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID     uuid.UUID `json:"userId" gorm:"type:uuid;unique;not null"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postalCode"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Booking struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CustomerID    uuid.UUID  `json:"customerId" gorm:"type:uuid;not null"`
	Customer      Customer   `json:"customer" gorm:"foreignKey:CustomerID"`
	WorkerID      uuid.UUID  `json:"workerId" gorm:"type:uuid;not null"`
	Worker        Worker     `json:"worker" gorm:"foreignKey:WorkerID"`
	Title         string     `json:"title" gorm:"not null"`
	Description   string     `json:"description"`
	ScheduledDate time.Time  `json:"scheduledDate" gorm:"type:date;not null"`
	ScheduledTime *time.Time `json:"scheduledTime,omitempty" gorm:"type:time"`
	DurationHours int        `json:"durationHours"`
	Location      string     `json:"location" gorm:"not null"`
	Status        string     `json:"status" gorm:"default:'pending'"`
	TotalCost     float64    `json:"totalCost"`
	Notes         string     `json:"notes"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type Review struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	BookingID  uuid.UUID `json:"bookingId" gorm:"type:uuid;not null"`
	WorkerID   uuid.UUID `json:"workerId" gorm:"type:uuid;not null"`
	Worker     Worker    `json:"worker,omitempty" gorm:"foreignKey:WorkerID"`
	CustomerID uuid.UUID `json:"customerId" gorm:"type:uuid;not null"`
	Customer   Customer  `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Rating     int       `json:"rating" gorm:"not null"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type FavoriteWorker struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CustomerID uuid.UUID `json:"customerId" gorm:"type:uuid;not null"`
	WorkerID   uuid.UUID `json:"workerId" gorm:"type:uuid;not null"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Message struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SenderID   uuid.UUID `json:"senderId" gorm:"type:uuid;not null"`
	Sender     User      `json:"sender" gorm:"foreignKey:SenderID"`
	ReceiverID uuid.UUID `json:"receiverId" gorm:"type:uuid;not null"`
	Receiver   User      `json:"receiver" gorm:"foreignKey:ReceiverID"`
	Content    string    `json:"content" gorm:"not null"`
	IsRead     bool      `json:"isRead" gorm:"default:false"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Notification struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"userId" gorm:"type:uuid;not null"`
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"not null"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"isRead" gorm:"default:false"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"createdAt"`
}

// Request/Response DTOs
type RegisterRequest struct {
	Email           string   `json:"email" binding:"required,email"`
	Password        string   `json:"password" binding:"required,min=8"`
	FirstName       string   `json:"firstName" binding:"required"`
	LastName        string   `json:"lastName" binding:"required"`
	Phone           string   `json:"phone"`
	UserType        string   `json:"userType" binding:"required,oneof=customer worker"`
	
	// Worker-specific fields
	Specialty       string   `json:"specialty,omitempty"`
	HourlyRate      float64  `json:"hourlyRate,omitempty"`
	ExperienceYears int      `json:"experienceYears,omitempty"`
	Bio             string   `json:"bio,omitempty"`
	Location        string   `json:"location,omitempty"`
	Skills          []string `json:"skills,omitempty"`
	
	// Customer-specific fields
	Address         string   `json:"address,omitempty"`
	City            string   `json:"city,omitempty"`
	State           string   `json:"state,omitempty"`
	PostalCode      string   `json:"postalCode,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatarUrl"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
}