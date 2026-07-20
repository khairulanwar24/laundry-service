// Package controllers (expense) adalah lapisan HTTP handler (Fiber) untuk domain pengeluaran.
package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// ExpenseController membungkus service registry.
type ExpenseController struct {
	service services.IServiceRegistry
}

// IExpenseController adalah kontrak handler HTTP domain pengeluaran.
type IExpenseController interface {
	List(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	Get(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
}

// NewExpenseController membuat instance baru.
func NewExpenseController(service services.IServiceRegistry) IExpenseController {
	return &ExpenseController{service: service}
}

func (ctrl *ExpenseController) List(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	kategori := c.Query("kategori", "")
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 15)
	data := ctrl.service.GetExpense().List(c.UserContext(), idOutlet, kategori, startDate, endDate, page, perPage)
	return c.JSON(data)
}

func (ctrl *ExpenseController) Create(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	form := c.Locals("validatedForm").(*dto.ExpenseForm)
	data := ctrl.service.GetExpense().Create(c.UserContext(), idOutlet, *form)
	return c.JSON(data)
}

func (ctrl *ExpenseController) Get(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idExpense := c.Params("id_expense")
	data := ctrl.service.GetExpense().Get(c.UserContext(), idOutlet, idExpense)
	return c.JSON(data)
}

func (ctrl *ExpenseController) Update(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idExpense := c.Params("id_expense")
	form := c.Locals("validatedForm").(*dto.ExpenseForm)
	data := ctrl.service.GetExpense().Update(c.UserContext(), idOutlet, idExpense, *form)
	return c.JSON(data)
}

func (ctrl *ExpenseController) Delete(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idExpense := c.Params("id_expense")
	data := ctrl.service.GetExpense().Delete(c.UserContext(), idOutlet, idExpense)
	return c.JSON(data)
}
