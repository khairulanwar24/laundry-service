package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	SetupAuthRoutes(app)
	// SetupUserRoutes dipindah ke pola berlapis (routes/user) & didaftarkan via NewRouteRegistry di main.go
	// SetupMasterRoutes dipindah ke pola berlapis (routes/masterapp) & didaftarkan via NewRouteRegistry di main.go
	SetupMstGrupAksesRoutes(app)
	SetupMstMenuRoutes(app)
	// SetupRefRoutes dipindah ke pola berlapis (routes/ref) & didaftarkan via NewRouteRegistry di main.go

	// api.Post("/users", controllers.CreateUser)

	// Authorization := app.Group("/apps", middleware.JWTMiddleware)
	// Authorization.Get("/users/getusers/:id_user", controllers.GetUsers)
	// Authorization.Get("/users/getusers/all", controllers.GetAllUsers)
	// allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	// Authorization.Put("/users/updateusers/:id_user", middleware.FileUploadMiddleware("avatar", "avatar", allowedTypes), controllers.UpdateUsers)
	// Authorization.Put("/users/updatepassword/:id_user", controllers.UpdatePassword)

	// users.Get("/:id_user", controllers.GetUsers)

	// users.Delete("/:id_user", controllers.DeleteUsers)

	// users.Put("/:id_user", controllers.UpdateUsers)

}
