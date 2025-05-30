// github.com/wronai/media-vault-backend/cmd/main.go
package main

import (
	"log"
	"os"


	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/wronai/media-vault-backend/internal/auth"
	"github.com/wronai/media-vault-backend/internal/database"
	"github.com/wronai/media-vault-backend/internal/handlers"
	"github.com/wronai/media-vault-backend/internal/services"
)

func main() {
	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Keycloak authentication
	if os.Getenv("OAUTH2_ENABLED") == "true" {
		if err := auth.InitKeycloak(); err != nil {
			log.Printf("Warning: Failed to initialize Keycloak: %v", err)
			log.Println("Running without authentication")
		} else {
			log.Println("Keycloak authentication enabled")
		}
	}

	// Initialize services
	vaultService := services.NewVaultService(db)
	photoService := services.NewPhotoService(db)
	descriptionService := services.NewDescriptionService()
	sharingService := services.NewSharingService(db)

	// Initialize auth middleware
	authMiddleware := auth.NewAuthMiddleware(auth.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	})

	// Initialize handlers
	vaultHandler := handlers.NewVaultHandler(vaultService)
	adminHandler := handlers.NewAdminHandler()
	partnerHandler := handlers.NewPartnerHandler(photoService, sharingService)
	uploadHandler := handlers.NewUploadHandler(*vaultService, *photoService, *descriptionService)
	photoHandler := handlers.NewPhotoHandler(photoService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:       "Media Vault API",
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})


	// Metrics
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Media Vault Metrics"}))

	// API v1 routes
	api := app.Group("/api/v1")


	// Auth routes (placeholder - implement these handlers)
	authGroup := api.Group("/auth")
	{
		authGroup.Post("/login", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
				"error": "Login endpoint not implemented",
			})
		})
		authGroup.Post("/register", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
				"error": "Register endpoint not implemented",
			})
		})
		authGroup.Post("/refresh", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
				"error": "Refresh token endpoint not implemented",
			})
		})
		authGroup.Post("/logout", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
				"error": "Logout endpoint not implemented",
			})
		})
	}

	// Protected routes
	protected := api.Use(authMiddleware)
	{
		// Vault routes
		vault := protected.Group("/vault")
		{
			// Vault routes
			vault.Get("", vaultHandler.GetVault)
			vault.Post("/upload", uploadHandler.UploadSingle)
			vault.Post("/upload/bulk", uploadHandler.BulkUpload)
		}

		// Photo routes
		photos := protected.Group("/photos")
		{
			photos.Get("", photoHandler.ListPhotos)
			photos.Get("/:id", photoHandler.GetPhoto)
			photos.Put("/:id", photoHandler.UpdatePhoto)
			photos.Delete("/:id", photoHandler.DeletePhoto)
			photos.Get("/:id/thumbnail", photoHandler.GetThumbnail)
			photos.Post("/:id/description", photoHandler.UpdateDescription)
			photos.Post("/:id/generate-description", photoHandler.GenerateDescription)
			photos.Get("/:id/shared-with", photoHandler.GetSharedWith)
		}

		// Partner routes
		partner := protected.Group("/partner")
		{
			partner.Post("/upload", partnerHandler.BulkUpload)
			partner.Get("/photos", partnerHandler.GetPartnerPhotos)
			partner.Put("/photos/descriptions", partnerHandler.BatchUpdateDescriptions)
			partner.Post("/photos/share", partnerHandler.BatchSharePhotos)
			partner.Get("/photos/:id/analytics", partnerHandler.GetPhotoAnalytics)
			partner.Get("/dashboard", partnerHandler.GetPartnerDashboard)
			partner.Get("/analytics", partnerHandler.GetPartnerAnalytics)
		}

		// Admin routes
		admin := protected.Group("/admin", auth.RequireRole("admin"))
		{
			// Admin routes
			admin.Get("/users", adminHandler.ListUsers)
			admin.Post("/users", adminHandler.CreateUser)
			admin.Get("/users/:id", adminHandler.GetUser)
			admin.Put("/users/:id", adminHandler.UpdateUser)
			admin.Delete("/users/:id", adminHandler.DeleteUser)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
