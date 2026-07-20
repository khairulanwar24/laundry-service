// Package controllers (account) adalah lapisan HTTP handler (Fiber) untuk domain akun laundry.
package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// AccountController membungkus service registry.
type AccountController struct {
	service services.IServiceRegistry
}

// IAccountController adalah kontrak handler HTTP domain akun.
type IAccountController interface {
	Register(*fiber.Ctx) error
	Login(*fiber.Ctx) error
	Me(*fiber.Ctx) error
	Logout(*fiber.Ctx) error
	ForgotPassword(*fiber.Ctx) error
	ResetPassword(*fiber.Ctx) error
}

// NewAccountController membuat instance baru.
func NewAccountController(service services.IServiceRegistry) IAccountController {
	return &AccountController{service: service}
}

func (ctrl *AccountController) Register(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.RegisterForm)
	data := ctrl.service.GetAccount().Register(c.UserContext(), *form)
	return c.JSON(data)
}

func (ctrl *AccountController) Login(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.LoginForm)
	data := ctrl.service.GetAccount().Login(c.UserContext(), *form)
	return c.JSON(data)
}

func (ctrl *AccountController) Me(c *fiber.Ctx) error {
	idUser, _ := c.Locals("id_user").(string)
	data := ctrl.service.GetAccount().Me(c.UserContext(), idUser)
	return c.JSON(data)
}

func (ctrl *AccountController) Logout(c *fiber.Ctx) error {
	data := ctrl.service.GetAccount().Logout(c.UserContext())
	return c.JSON(data)
}

func (ctrl *AccountController) ForgotPassword(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.ForgotPasswordForm)
	data := ctrl.service.GetAccount().ForgotPassword(c.UserContext(), *form)
	return c.JSON(data)
}

func (ctrl *AccountController) ResetPassword(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.ResetPasswordForm)
	data := ctrl.service.GetAccount().ResetPassword(c.UserContext(), *form)
	return c.JSON(data)
}
