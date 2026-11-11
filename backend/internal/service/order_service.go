package service

import (
	"errors"
	"restaurant-management-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{}
}

func (os *OrderService) CreateOrder(order models.Order) (*models.Order, error) {
	if err := os.DB.Create(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (os *OrderService) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := os.DB.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (os *OrderService) Update(order models.Order) (*models.Order, error) {
	if err := os.DB.Save(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (os *OrderService) DeleteOrder(orderID string) error {
	result := os.DB.Where("id = ?", orderID).Delete(&models.Tables{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Order not found")
	}
	return nil
}

func (os *OrderService) GetOrderById(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := os.DB.Where("id = ?", orderID).Find(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
