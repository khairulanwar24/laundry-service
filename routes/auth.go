package routes

import (
	"sso-service/controllers"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {

	api := app.Group("/auth")

	api.Post("/login", middleware.ValidateForm(&controllers.LoginForm{}), controllers.LoginController)
	api.Post("/logout", controllers.LogoutController)
	api.Post("/resetpassword", middleware.ValidateForm(&controllers.ResetForm{}), controllers.ResetPassword)
	api.Post("/changepassword", middleware.ValidateForm(&controllers.ChangePasswordForm{}), controllers.ChangePassword)
	api.Post("/cekotp", middleware.ValidateForm(&controllers.OtpForm{}), controllers.CekOtp)
	api.Get("/gettoken", controllers.RefreshToken)

}
