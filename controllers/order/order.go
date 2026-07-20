// Package controllers (order) adalah lapisan HTTP handler (Fiber) untuk domain pesanan.
package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// OrderController membungkus service registry.
type OrderController struct {
	service services.IServiceRegistry
}

// IOrderController adalah kontrak handler HTTP domain pesanan.
type IOrderController interface {
	Create(*fiber.Ctx) error
	List(*fiber.Ctx) error
	Search(*fiber.Ctx) error
	Get(*fiber.Ctx) error
	GetHistory(*fiber.Ctx) error
	ChangeStatus(*fiber.Ctx) error
	Pay(*fiber.Ctx) error
	Pickup(*fiber.Ctx) error
}

// NewOrderController membuat instance baru.
func NewOrderController(service services.IServiceRegistry) IOrderController {
	return &OrderController{service: service}
}

func (ctrl *OrderController) Create(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idUser, _ := c.Locals("id_user").(string)
	form := c.Locals("validatedForm").(*dto.CreateOrderForm)
	data := ctrl.service.GetOrder().Create(c.UserContext(), idOutlet, idUser, *form)
	return c.JSON(data)
}

func (ctrl *OrderController) List(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	tab := c.Query("tab", "")
	q := c.Query("q", "")
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 15)
	data := ctrl.service.GetOrder().List(c.UserContext(), idOutlet, tab, q, page, perPage)
	return c.JSON(data)
}

func (ctrl *OrderController) Search(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	q := c.Query("q", "")
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 15)
	data := ctrl.service.GetOrder().Search(c.UserContext(), idOutlet, q, page, perPage)
	return c.JSON(data)
}

func (ctrl *OrderController) Get(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idOrder := c.Params("id_order")
	data := ctrl.service.GetOrder().Get(c.UserContext(), idOutlet, idOrder)
	return c.JSON(data)
}

func (ctrl *OrderController) GetHistory(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idOrder := c.Params("id_order")
	data := ctrl.service.GetOrder().GetHistory(c.UserContext(), idOutlet, idOrder)
	return c.JSON(data)
}

func (ctrl *OrderController) ChangeStatus(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idOrder := c.Params("id_order")
	idUser, _ := c.Locals("id_user").(string)
	form := c.Locals("validatedForm").(*dto.ChangeStatusForm)
	data := ctrl.service.GetOrder().ChangeStatus(c.UserContext(), idOutlet, idOrder, idUser, *form)
	return c.JSON(data)
}

func (ctrl *OrderController) Pay(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idOrder := c.Params("id_order")
	idUser, _ := c.Locals("id_user").(string)
	form := c.Locals("validatedForm").(*dto.PayOrderForm)
	data := ctrl.service.GetOrder().Pay(c.UserContext(), idOutlet, idOrder, idUser, *form)
	return c.JSON(data)
}

func (ctrl *OrderController) Pickup(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idOrder := c.Params("id_order")
	idUser, _ := c.Locals("id_user").(string)
	form := c.Locals("validatedForm").(*dto.PickupForm)
	data := ctrl.service.GetOrder().Pickup(c.UserContext(), idOutlet, idOrder, idUser, *form)
	return c.JSON(data)
}
