package controller

import (
	"net/http"
	"restaurant-management-app/internal/models"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderController struct {
	OrderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) *OrderController {
	return &OrderController{
		OrderService: orderService,
	}
}

type CreateOrderRequest struct {
	TableID uuid.UUID `json:"table_id" binding:"required"`
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if !bindJSON(c, &req) {
		return
	}

	createdOrder, err := oc.OrderService.CreateOrder(models.Order{TableID: req.TableID})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, createdOrder)
}

func (oc *OrderController) GetAllOrders(c *gin.Context) {
	orders, err := oc.OrderService.GetAllOrders()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (oc *OrderController) GetOrderByID(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	order, err := oc.OrderService.GetOrderById(id.String())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, order)
}

// DeleteOrder removes an order that has no items yet (opened by mistake).
func (oc *OrderController) DeleteOrder(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := oc.OrderService.DeleteOrder(id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
