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
	return &MenuItemService{
		DB: db,
	}
}

func (ms *MenuItemService) CreateMenuItem(item models.MenuItem) (*models.MenuItem, error) {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	if err := ms.DB.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (ms *MenuItemService) GetAllItems() ([]*models.MenuItem, error) {
	var menuItems []*models.MenuItem
	if err := ms.DB.Order("item_number").Find(&menuItems).Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

func (ms *MenuItemService) UpdateItems(item models.MenuItem) (*models.MenuItem, error) {
	result := ms.DB.Model(&item).Select("item_number", "name", "category", "price").Updates(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrMenuItemNotFound
	}
	return &item, nil
}

func (ms *MenuItemService) DeleteItems(itemID uuid.UUID) error {
	result := ms.DB.Where("id = ?", itemID).Delete(&models.MenuItem{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrMenuItemNotFound
	}
	return nil
}

func (ms *MenuItemService) GetMenuItemByID(itemID string) (*models.MenuItem, error) {
	var item models.MenuItem
	if err := ms.DB.Where("id = ?", itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMenuItemNotFound
		}
		return nil, err
	}
	return &item, nil
}
