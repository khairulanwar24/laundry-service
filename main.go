package main

import (
	"log"
	"os"

	"laundry-service/controllers"
	"laundry-service/database"
	"laundry-service/repositories"
	"laundry-service/routes"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize the database connection
	database.Connect()

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "3001"
	}

	// Initialize Fiber
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		ExposeHeaders:    "Content-Length",
		MaxAge:           3600,
		AllowOrigins:     getEnvOrDefault("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://localhost:5173"),
	}))

	// Dependency injection untuk seluruh domain laundry (pola registry berlapis).
	repository := repositories.NewRepositoryRegistry(database.DB)
	service := services.NewServiceRegistry(repository)
	controller := controllers.NewControllerRegistry(service)
	routes.NewRouteRegistry(controller, app).Serve()

	log.Println("HTTP server running on port", httpPort)
	if err := app.Listen(":" + httpPort); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func init() {
	err := godotenv.Load(("app/.env"))
	if err != nil {
		log.Println("⚠️  .env tidak ditemukan, lanjutkan dengan env system")
	}
}
