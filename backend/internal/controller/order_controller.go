package controller

import (
	"errors"
	"net/http"
	"restaurant-management-app/internal/models"
	"restaurant-management-app/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order := models.Order{
		TableID:   req.TableID,
		Status:    models.OrderStatusOpen,
		CreatedAt: time.Now(),
	}
	createdOrder, err := oc.OrderService.CreateOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdOrder)
}

func (oc *OrderController) GetAllOrders(c *gin.Context) {
	orders, err := oc.OrderService.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (oc *OrderController) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")

	order, err := oc.OrderService.GetOrderById(orderID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve order",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}
