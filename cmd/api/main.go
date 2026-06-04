// @title           AI Resume Analyzer API
// @version         1.0
// @description     REST API untuk analisis resume menggunakan Groq AI
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

package main

import (
	"log"

	"github.com/moxhyusuf/ai-resume-analyzer/config"
	_ "github.com/moxhyusuf/ai-resume-analyzer/docs"
	"github.com/moxhyusuf/ai-resume-analyzer/internal/auth"
	"github.com/moxhyusuf/ai-resume-analyzer/internal/middleware"
	"github.com/moxhyusuf/ai-resume-analyzer/internal/resume"
	"github.com/moxhyusuf/ai-resume-analyzer/internal/user"
	"github.com/moxhyusuf/ai-resume-analyzer/pkg/groq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	fiberswagger "github.com/swaggo/fiber-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()
	db := config.ConnectDB(cfg)
	config.AutoMigrate(db)

	groqClient := groq.NewClient(cfg.GroqAPIKey)

	userRepo := user.NewUserRepository(db)
	resumeRepo := resume.NewResumeRepository(db)

	authService := auth.NewAuthService(userRepo, cfg.JWTSecret)
	resumeService := resume.NewResumeService(resumeRepo, groqClient)

	authHandler := auth.NewAuthHandler(authService)
	resumeHandler := resume.NewResumeHandler(resumeService)

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
		AppName:   "AI RESUME ANALIZER",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// Swagger ← tambah ini
	app.Get("/swagger/*", fiberswagger.WrapHandler)

	api := app.Group("/api/v1")

	authGroup := api.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	resumeGroup := api.Group("/resumes", middleware.JWTProtected(cfg.JWTSecret))
	resumeGroup.Post("/analyze", resumeHandler.Analyze)
	resumeGroup.Get("/", resumeHandler.GetAll)
	resumeGroup.Get("/:id", resumeHandler.GetByID)
	resumeGroup.Delete("/:id", resumeHandler.Delete)

	log.Printf("Server running on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
