package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"construction-backend/internal/middleware"
	"construction-backend/internal/models"
	"construction-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user, "message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, _ := middleware.GetUserIDFromContext(c)
	
	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "message": "Invalid user ID"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
		return
	}

	user, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "message": "Profile updated successfully"})
}

// UploadAvatar allows authenticated users to upload a profile avatar image
func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "message": "Invalid user ID"})
		return
	}

	// Ensure requester is authenticated and the same user
	authUserID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Authentication required"})
		return
	}
	if authUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "You can only upload avatar for your own profile"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file is required", "message": "avatar file is required"})
		return
	}

	// Validate extension
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type", "message": "Only jpg and png allowed"})
		return
	}

	// Ensure uploads folder exists
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Use userID to form filename
	filename := fmt.Sprintf("avatar_%s%s", userID.String(), ext)
	destination := filepath.Join("./uploads", filename)

	if err := c.SaveUploadedFile(file, destination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file", "message": err.Error()})
		return
	}

	// Update user profile with avatar URL
	avatarURL := fmt.Sprintf("/uploads/%s", filename)
	_, err = h.authService.UpdateProfile(userID, &models.UpdateProfileRequest{AvatarURL: avatarURL})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatarUrl": avatarURL, "message": "Avatar uploaded"})
} 