// Package controllers (customer) adalah lapisan HTTP handler (Fiber) untuk domain pelanggan.
package controllers

import (
	"strconv"

	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// CustomerController membungkus service registry.
type CustomerController struct {
	service services.IServiceRegistry
}

// ICustomerController adalah kontrak handler HTTP domain pelanggan.
type ICustomerController interface {
	List(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	Get(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
	ListOrders(*fiber.Ctx) error
}

// NewCustomerController membuat instance baru.
func NewCustomerController(service services.IServiceRegistry) ICustomerController {
	return &CustomerController{service: service}
}

func parseIsActiveQuery(c *fiber.Ctx) *bool {
	raw := c.Query("is_active", "")
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return nil
	}
	return &v
}

func (ctrl *CustomerController) List(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	search := c.Query("search", "")
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 15)
	data := ctrl.service.GetCustomer().List(c.UserContext(), idOutlet, search, parseIsActiveQuery(c), page, perPage)
	return c.JSON(data)
}

func (ctrl *CustomerController) Create(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	form := c.Locals("validatedForm").(*dto.CustomerForm)
	data := ctrl.service.GetCustomer().Create(c.UserContext(), idOutlet, *form)
	return c.JSON(data)
}

func (ctrl *CustomerController) Get(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idCustomer := c.Params("id_customer")
	data := ctrl.service.GetCustomer().Get(c.UserContext(), idOutlet, idCustomer)
	return c.JSON(data)
}

func (ctrl *CustomerController) Update(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idCustomer := c.Params("id_customer")
	form := c.Locals("validatedForm").(*dto.CustomerForm)
	data := ctrl.service.GetCustomer().Update(c.UserContext(), idOutlet, idCustomer, *form)
	return c.JSON(data)
}

func (ctrl *CustomerController) Delete(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idCustomer := c.Params("id_customer")
	data := ctrl.service.GetCustomer().Delete(c.UserContext(), idOutlet, idCustomer)
	return c.JSON(data)
}

func (ctrl *CustomerController) ListOrders(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idCustomer := c.Params("id_customer")
	status := c.Query("status", "")
	fromDate := c.Query("from_date", "")
	toDate := c.Query("to_date", "")
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 15)
	data := ctrl.service.GetCustomer().ListOrders(c.UserContext(), idOutlet, idCustomer, status, fromDate, toDate, page, perPage)
	return c.JSON(data)
}
