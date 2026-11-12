package service

import (
	"errors"
	"restaurant-management-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuItemService struct {
	DB *gorm.DB
}

func NewMenuItemsService(db *gorm.DB) *MenuItemService {
	return &MenuItemService{}
}

func (ms *MenuItemService) CreateMenuItem(item models.MenuItem) (*models.MenuItem, error) {
	if err := ms.DB.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (ms MenuItemService) GetAllItems() ([]*models.MenuItem, error) {
	var menuItems []*models.MenuItem
	if err := ms.DB.Find(&menuItems).Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

func (ms *MenuItemService) UpdateItems(item models.MenuItem) (*models.MenuItem, error) {
	if err := ms.DB.Save(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (ms *MenuItemService) DeleteItems(itemID uuid.UUID) error {
	result := ms.DB.Where("id = ?", itemID).Delete(&models.Tables{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Table not found")
	}
	return nil
}

func (ms *MenuItemService) GetMenuItemByID(itemID uuid.UUID) (*models.MenuItem, error) {
	var item *models.MenuItem
	result := ms.DB.Where("id = ?", itemID).Find(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("Item not found")
	}
	return item, nil
}
