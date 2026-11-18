package controller

import (
	"net/http"
	"restaurant-management-app/internal/models"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
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
	ItemNumber uint    `json:"item_number" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Category   string  `json:"category" binding:"required"`
	Price      float32 `json:"price" binding:"required,gt=0"`
}

func (mc *MenuItemController) CreateItem(c *gin.Context) {
	var req CreateMenuItemRequest
	if !bindJSON(c, &req) {
		return
	}

	createdItem, err := mc.MenuItemService.CreateMenuItem(models.MenuItem{
		ItemNumber: req.ItemNumber,
		Name:       req.Name,
		Category:   req.Category,
		Price:      req.Price,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, createdItem)
}

func (mc *MenuItemController) GetMenuItemById(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := mc.MenuItemService.GetMenuItemByID(id.String())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (mc *MenuItemController) GetMenuItemList(c *gin.Context) {
	items, err := mc.MenuItemService.GetAllItems()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// UpdateMenuItemRequest is a partial update: only the fields sent are changed.
type UpdateMenuItemRequest struct {
	ItemNumber *uint    `json:"item_number" binding:"omitempty,gt=0"`
	Name       *string  `json:"name" binding:"omitempty,min=1"`
	Category   *string  `json:"category" binding:"omitempty,min=1"`
	Price      *float32 `json:"price" binding:"omitempty,gt=0"`
}

func (mc *MenuItemController) UpdateItem(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req UpdateMenuItemRequest
	if !bindJSON(c, &req) {
		return
	}
	if req.ItemNumber == nil && req.Name == nil && req.Category == nil && req.Price == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided for update"})
		return
	}

	item, err := mc.MenuItemService.GetMenuItemByID(id.String())
	if err != nil {
		respondError(c, err)
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
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedItem)
}

func (mc *MenuItemController) DeleteItem(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := mc.MenuItemService.DeleteItems(id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
