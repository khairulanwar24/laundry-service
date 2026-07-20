// Package routes (expense) mendaftarkan endpoint HTTP untuk domain pengeluaran.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// ExpenseRoute mewadahi controller registry & router Fiber.
type ExpenseRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IExpenseRoute adalah kontrak pendaftaran rute domain pengeluaran.
type IExpenseRoute interface {
	Run()
}

// NewExpenseRoute membuat instance baru.
func NewExpenseRoute(controller controllers.IControllerRegistry, router fiber.Router) IExpenseRoute {
	return &ExpenseRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain pengeluaran.
func (r *ExpenseRoute) Run() {
	api := r.router.Group("/outlets/:id_outlet/expenses", middleware.LaundryJWTMiddleware, middleware.RequireOutletMember)

	api.Get("/", r.controller.GetExpense().List)
	api.Post("/", middleware.ValidateForm(&dto.ExpenseForm{}), r.controller.GetExpense().Create)
	api.Get("/:id_expense", r.controller.GetExpense().Get)
	api.Put("/:id_expense", middleware.ValidateForm(&dto.ExpenseForm{}), r.controller.GetExpense().Update)
	api.Delete("/:id_expense", r.controller.GetExpense().Delete)
}
