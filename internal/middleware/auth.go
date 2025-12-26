package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"userId"`
	Email    string    `json:"email"`
	UserType string    `json:"userType"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

// Parse claims into a map so we can robustly extract values
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	mapClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	// Extract userId, email, userType
	userIdStr, _ := mapClaims["userId"].(string)
	email, _ := mapClaims["email"].(string)
	userType, _ := mapClaims["userType"].(string)

	if userIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: missing userId"})
		c.Abort()
		return
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: invalid userId"})
		c.Abort()
		return
	}

	// Set user information in context
	c.Set("userId", userID)
	c.Set("email", email)
	c.Set("userType", userType)

		c.Next()
	}
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := userID.(uuid.UUID)
	return id, ok
}

func GetUserTypeFromContext(c *gin.Context) (string, bool) {
	userType, exists := c.Get("userType")
	if !exists {
		return "", false
	}
	ut, ok := userType.(string)
	return ut, ok
}