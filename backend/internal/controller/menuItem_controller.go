package controller

import (
	"net/http"
	"restaurant-management-app/internal/models"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MenuItemController struct {
	MenuItemService *service.MenuItemService
}

func NewMenuItemController(menuItemService *service.MenuItemService) *MenuItemController {
	return &MenuItemController{
		MenuItemService: menuItemService,
	}
}

type CreateMenuItemRequest struct {
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Price    float32 `json:"price" binding:"required"`
}

func (mc *MenuItemController) Create(c *gin.Context) {
	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	item := models.MenuItem{
		ID:       uuid.New(),
		Name:     req.Name,
		Category: req.Category,
		Price:    req.Price,
	}
	createdItem, err := mc.MenuItemService.CreateMenuItem(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, createdItem)
}
