// Package controllers (dashboard) adalah lapisan HTTP handler (Fiber) untuk domain dashboard.
package controllers

import (
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// DashboardController membungkus service registry.
type DashboardController struct {
	service services.IServiceRegistry
}

// IDashboardController adalah kontrak handler HTTP domain dashboard.
type IDashboardController interface {
	Summary(*fiber.Ctx) error
	DailyReport(*fiber.Ctx) error
	FinancialSummary(*fiber.Ctx) error
	RevenueSummary(*fiber.Ctx) error
	ProfitLossSummary(*fiber.Ctx) error
	CustomerReport(*fiber.Ctx) error
	CustomerSearch(*fiber.Ctx) error
}

// NewDashboardController membuat instance baru.
func NewDashboardController(service services.IServiceRegistry) IDashboardController {
	return &DashboardController{service: service}
}

func (ctrl *DashboardController) Summary(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetDashboard().Summary(c.UserContext(), idOutlet, c.Query("date", ""))
	return c.JSON(data)
}

func (ctrl *DashboardController) DailyReport(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetDashboard().DailyReport(c.UserContext(), idOutlet, c.Query("date", ""))
	return c.JSON(data)
}

func (ctrl *DashboardController) FinancialSummary(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetDashboard().FinancialSummary(c.UserContext(), idOutlet, c.Query("date", ""))
	return c.JSON(data)
}

func (ctrl *DashboardController) RevenueSummary(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetDashboard().RevenueSummary(c.UserContext(), idOutlet)
	return c.JSON(data)
}

func (ctrl *DashboardController) ProfitLossSummary(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetDashboard().ProfitLossSummary(c.UserContext(), idOutlet)
	return c.JSON(data)
}

func (ctrl *DashboardController) CustomerReport(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetDashboard().CustomerReport(c.UserContext(), idOutlet,
		c.Query("search", ""), c.Query("status", "all"), c.Query("date_from", ""), c.Query("date_to", ""))
	return c.JSON(data)
}

func (ctrl *DashboardController) CustomerSearch(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetDashboard().CustomerSearch(c.UserContext(), idOutlet,
		c.Query("search", ""), c.Query("status", "all"), c.Query("date_from", ""), c.Query("date_to", ""),
		c.QueryInt("page", 1), c.QueryInt("per_page", 15))
	return c.JSON(data)
}
