package config

import (
	"fmt"
	"log"
	"os"

	"github.com/moxhyusuf/ai-resume-analyzer/internal/resume"
	"github.com/moxhyusuf/ai-resume-analyzer/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	GroqAPIKey  string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: mustEnv("DATABASE_URL"),
		JWTSecret:   mustEnv("JWT_SECRET"),
		GroqAPIKey:  mustEnv("GROQ_API_KEY"),
	}
}

func ConnectDB(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Database connected")
	return db
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&user.User{},
		&resume.Resume{},
		&resume.AnalysisResult{},
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return val
}
