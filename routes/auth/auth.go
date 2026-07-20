// Package routes (auth) mendaftarkan endpoint HTTP untuk domain autentikasi.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// AuthRoute mewadahi controller registry & router Fiber.
type AuthRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IAuthRoute adalah kontrak pendaftaran rute autentikasi.
type IAuthRoute interface {
	Run()
}

// NewAuthRoute membuat instance AuthRoute baru.
func NewAuthRoute(controller controllers.IControllerRegistry, router fiber.Router) IAuthRoute {
	return &AuthRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain autentikasi.
func (r *AuthRoute) Run() {
	api := r.router.Group("/auth")

	api.Post("/login", middleware.ValidateForm(&dto.LoginForm{}), r.controller.GetAuth().LoginController)
	api.Post("/logout", r.controller.GetAuth().LogoutController)
	api.Post("/resetpassword", middleware.ValidateForm(&dto.ResetForm{}), r.controller.GetAuth().ResetPassword)
	api.Post("/changepassword", middleware.ValidateForm(&dto.ChangePasswordForm{}), r.controller.GetAuth().ChangePassword)
	api.Post("/cekotp", middleware.ValidateForm(&dto.OtpForm{}), r.controller.GetAuth().CekOtp)
	api.Get("/gettoken", r.controller.GetAuth().RefreshToken)
}
