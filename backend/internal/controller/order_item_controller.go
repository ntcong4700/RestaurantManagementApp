package controller

import (
	"net/http"
	"restaurant-management-app/internal/models"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderItemController struct {
	OrderItemService *service.OrderItemService
}

func NewOrderItemController(orderItemService *service.OrderItemService) *OrderItemController {
	return &OrderItemController{
		OrderItemService: orderItemService,
	}
}

type CreateOrderItemRequest struct {
	MenuItemID uuid.UUID `json:"menu_item_id" binding:"required"`
	Quantity   uint32    `json:"quantity" binding:"required,gt=0"`
}

// CreateOrderItem handles POST /orders/:id/items.
func (oc *OrderItemController) CreateOrderItem(c *gin.Context) {
	orderID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req CreateOrderItemRequest
	if !bindJSON(c, &req) {
		return
	}

	item, err := oc.OrderItemService.CreateOrderItem(models.OrderItem{
		OrderID:    orderID,
		MenuItemID: req.MenuItemID,
		Quantity:   req.Quantity,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// GetOrderItemsByOrder handles GET /orders/:id/items.
func (oc *OrderItemController) GetOrderItemsByOrder(c *gin.Context) {
	orderID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	items, err := oc.OrderItemService.GetOrderItemsByOrderID(orderID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetOrderItemByID handles GET /order-items/:id.
func (oc *OrderItemController) GetOrderItemByID(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := oc.OrderItemService.GetOrderItemByID(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

type UpdateOrderItemRequest struct {
	Quantity uint32 `json:"quantity" binding:"required,gt=0"`
}

// UpdateOrderItem handles PATCH /order-items/:id.
func (oc *OrderItemController) UpdateOrderItem(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req UpdateOrderItemRequest
	if !bindJSON(c, &req) {
		return
	}

	item, err := oc.OrderItemService.UpdateOrderItemQuantity(id, req.Quantity)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteOrderItem handles DELETE /order-items/:id.
func (oc *OrderItemController) DeleteOrderItem(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := oc.OrderItemService.DeleteOrderItem(id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
