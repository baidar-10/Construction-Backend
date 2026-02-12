package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	Port           string
	AllowedOrigins []string
	// MinIO
	MinioEndpoint       string
	MinioPublicEndpoint string
	MinioAccessKey      string
	MinioSecretKey      string
	MinioUseSSL         bool
	MinioBucket         string
}

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	allowed := getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:5174,http://localhost:3000")

	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "admin"),
		DBPassword:     getEnv("DB_PASSWORD", "admin123"),
		DBName:         getEnv("DB_NAME", "construction_db"),
		JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		Port:           getEnv("PORT", "8080"),
		AllowedOrigins: strings.Split(allowed, ","),
		// MinIO
		MinioEndpoint:       getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioPublicEndpoint: getEnv("MINIO_PUBLIC_ENDPOINT", "http://localhost:9000"),
		MinioAccessKey:      getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:      getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinioUseSSL:         getEnv("MINIO_USE_SSL", "false") == "true",
		MinioBucket:         getEnv("MINIO_BUCKET", "verification-documents"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 