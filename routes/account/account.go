// Package routes (account) mendaftarkan endpoint HTTP untuk domain akun laundry.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// AccountRoute mewadahi controller registry & router Fiber.
type AccountRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IAccountRoute adalah kontrak pendaftaran rute domain akun.
type IAccountRoute interface {
	Run()
}

// NewAccountRoute membuat instance baru.
func NewAccountRoute(controller controllers.IControllerRegistry, router fiber.Router) IAccountRoute {
	return &AccountRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain akun.
func (r *AccountRoute) Run() {
	api := r.router.Group("/account")

	api.Post("/register", middleware.ValidateForm(&dto.RegisterForm{}), r.controller.GetAccount().Register)
	api.Post("/login", middleware.ValidateForm(&dto.LoginForm{}), r.controller.GetAccount().Login)
	api.Post("/forgot-password", middleware.ValidateForm(&dto.ForgotPasswordForm{}), r.controller.GetAccount().ForgotPassword)
	api.Post("/reset-password", middleware.ValidateForm(&dto.ResetPasswordForm{}), r.controller.GetAccount().ResetPassword)

	api.Get("/me", middleware.LaundryJWTMiddleware, r.controller.GetAccount().Me)
	api.Post("/logout", middleware.LaundryJWTMiddleware, r.controller.GetAccount().Logout)
}
