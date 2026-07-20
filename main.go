package main

import (
	"log"
	"net"
	"os"
	"laundry-service/controllers"
	"laundry-service/database"
	"laundry-service/repositories"
	"laundry-service/services"

	grpcsso "laundry-service/grpc"
	"laundry-service/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	pb "laundry-service/proto" // Import file proto yang dihasilkan

	"google.golang.org/grpc"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	// Initialize the database connection
	database.Connect()
	database.ConnectDigiclass()

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "3001"
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	// Initialize Fiber
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		// AllowHeaders:     "Content-Type,Authorization",
		ExposeHeaders: "Content-Length",
		MaxAge:        3600,
		// AllowOrigins:  "*",
		AllowOrigins: "https://*.farmasiunissula.com,http://*.farmasiunissula.com",
		// AllowOrigins: "http://localhost:5173,http://sso.test,https://sso.farmasi.local,https://sso.dnglab.id,https://sso.farmasiunissula.com,http://localhost:91",
		// withCredentials: true,

		// AllowOrigins:     "http://example.com, https://example.com", // Ganti dengan origin yang diizinkan
		// AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		// AllowCredentials: true, // Jika ingin mengizinkan cookies
	}))

	// Dependency injection untuk domain yang sudah dimigrasi ke pola berlapis (registry).
	repository := repositories.NewRepositoryRegistry(database.DB, database.DBAkademik, database.DBDigiclass)
	service := services.NewServiceRegistry(repository)
	controller := controllers.NewControllerRegistry(service)
	routes.NewRouteRegistry(controller, app).Serve()

	// Jalankan Fiber di goroutine
	go func() {
		log.Println("HTTP server running on port", httpPort)
		if err := app.Listen(":" + httpPort); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	// Jalankan gRPC server
	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &grpcsso.AuthServiceServer{}) // Daftarkan gRPC service

	log.Println("gRPC server running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}

func init() {
	err := godotenv.Load(("app/.env"))
	if err != nil {
		log.Println("⚠️  .env tidak ditemukan, lanjutkan dengan env system")
	}
}
