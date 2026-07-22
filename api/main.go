package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables. ENV_FILE selects an alternative file (for
	// example .env.supabase) so the same binary can target the local database
	// or the cloud stack without editing .env.
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: no %s file found, using system environment variables", envFile)
	}

	// Initialize Database
	config.ConnectDB()

	// Demo data for local development only (see SEED_DEV_DATA)
	config.SeedDevData()

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName:   "iSpark Authentication API",
		BodyLimit: 10 * 1024 * 1024, // 10MB limit (frontend will restrict to 5MB)
	})

	// Add Middlewares
	app.Use(recover.New()) // Recovers from panics to keep app running
	app.Use(logger.New())  // Logs HTTP request details

	// Setup CORS. ALLOWED_ORIGINS extends the local defaults with deployed
	// web origins (comma-separated), e.g. https://ispark-test.vercel.app
	allowOrigins := "http://localhost:5173, http://localhost:3000, http://127.0.0.1:5173, http://127.0.0.1:3000"
	if extra := os.Getenv("ALLOWED_ORIGINS"); extra != "" {
		allowOrigins = allowOrigins + ", " + extra
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Liveness probe for Render health checks and uptime monitors.
	// Deliberately DB-free so a paused free-tier database doesn't flap it.
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Setup Routes
	routes.SetupRoutes(app)

	// Fallback/Not Found route
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint not found",
		})
	})

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
