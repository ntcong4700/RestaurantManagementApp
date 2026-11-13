package controller

import (
	"errors"
	"net/http"
	"restaurant-management-app/internal/models"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	ItemNumber uint    `json:"item_number" validate:"required"`
	Name       string  `json:"name" binding:"required"`
	Category   string  `json:"category" binding:"required"`
	Price      float32 `json:"price" binding:"required"`
}

func (mc *MenuItemController) CreateItem(c *gin.Context) {
	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	item := models.MenuItem{
		ID:         uuid.New(),
		Name:       req.Name,
		Category:   req.Category,
		Price:      req.Price,
		ItemNumber: req.ItemNumber,
	}
	createdItem, err := mc.MenuItemService.CreateMenuItem(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, createdItem)
}

func (mc *MenuItemController) GetMenuItemById(c *gin.Context) {
	id := c.Param("id")
	item, err := mc.MenuItemService.GetMenuItemByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Item not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve order",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (mc *MenuItemController) GetMenuItemList(c *gin.Context) {
	item, err := mc.MenuItemService.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, item)
}

type UpdateMenuItemRequest struct {
	ItemNumber *uint    `json:"item_number" validate:"required"`
	Name       *string  `json:"name" binding:"required"`
	Category   *string  `json:"category" binding:"required"`
	Price      *float32 `json:"price" binding:"required"`
}

func (mc *MenuItemController) UpdateItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid menu item ID",
		})
		return
	}

	var req UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.ItemNumber == nil &&
		req.Name == nil &&
		req.Category == nil &&
		req.Price == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No fields provided for update",
		})
		return
	}

	item, err := mc.MenuItemService.GetMenuItemByID(id.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Menu item not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve menu item",
		})
		return
	}

	if req.ItemNumber != nil {
		item.ItemNumber = *req.ItemNumber
	}

	if req.Name != nil {
		item.Name = *req.Name
	}

	if req.Category != nil {
		item.Category = *req.Category
	}

	if req.Price != nil {
		item.Price = *req.Price
	}

	updatedItem, err := mc.MenuItemService.UpdateItems(*item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update menu item",
		})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}