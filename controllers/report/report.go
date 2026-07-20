// Package controllers (report) adalah lapisan HTTP handler (Fiber) untuk domain laporan.
package controllers

import (
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// ReportController membungkus service registry.
type ReportController struct {
	service services.IServiceRegistry
}

// IReportController adalah kontrak handler HTTP domain laporan.
type IReportController interface {
	Transactions(*fiber.Ctx) error
}

// NewReportController membuat instance baru.
func NewReportController(service services.IServiceRegistry) IReportController {
	return &ReportController{service: service}
}

func (ctrl *ReportController) Transactions(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetReport().Transactions(c.UserContext(), idOutlet,
		c.Query("status_pesanan", ""), c.Query("status_pembayaran", ""),
		c.Query("tanggal_mulai", ""), c.Query("tanggal_akhir", ""),
		c.QueryInt("page", 1), c.QueryInt("per_page", 15))
	return c.JSON(data)
}
