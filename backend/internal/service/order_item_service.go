package service

import (
	"errors"
	"restaurant-management-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrOrderItemNotFound = errors.New("order item not found")
	ErrOrderNotFound     = errors.New("order not found")
	ErrOrderNotOpen      = errors.New("order is not open")
	ErrMenuItemNotFound  = errors.New("menu item not found")
	ErrInvalidQuantity   = errors.New("quantity must be greater than 0")
	ErrOrderItemBilled   = errors.New("order item has already been billed")
)

type OrderItemService struct {
	DB *gorm.DB
}

func NewOrderItemService(db *gorm.DB) *OrderItemService {
	return &OrderItemService{
		DB: db,
	}
}

// CreateOrderItem adds an item to an open order. The order must exist and be open,
// and the menu item must exist.
func (s *OrderItemService) CreateOrderItem(orderItem models.OrderItem) (*models.OrderItem, error) {
	if orderItem.Quantity == 0 {
		return nil, ErrInvalidQuantity
	}
	if orderItem.ID == uuid.Nil {
		orderItem.ID = uuid.New()
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := ensureOrderOpen(tx, orderItem.OrderID); err != nil {
			return err
		}

		var count int64
		if err := tx.Model(&models.MenuItem{}).Where("id = ?", orderItem.MenuItemID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return ErrMenuItemNotFound
		}

		return tx.Omit(clause.Associations).Create(&orderItem).Error
	})
	if err != nil {
		return nil, err
	}

	return s.GetOrderItemByID(orderItem.ID)
}

// GetOrderItemByID returns one order item with its menu item preloaded.
func (s *OrderItemService) GetOrderItemByID(id uuid.UUID) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	err := s.DB.Preload("MenuItem").Where("id = ?", id).First(&orderItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderItemNotFound
		}
		return nil, err
	}
	return &orderItem, nil
}

// GetOrderItemsByOrderID returns all items of an order with their menu items preloaded.
func (s *OrderItemService) GetOrderItemsByOrderID(orderID uuid.UUID) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	err := s.DB.Preload("MenuItem").Where("order_id = ?", orderID).Find(&orderItems).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}

// UpdateOrderItemQuantity changes the quantity of an item. Only allowed while the
// order is open and the item has not been billed yet.
func (s *OrderItemService) UpdateOrderItemQuantity(id uuid.UUID, quantity uint32) (*models.OrderItem, error) {
	if quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		orderItem, err := findEditableOrderItem(tx, id)
		if err != nil {
			return err
		}
		return tx.Model(orderItem).Update("quantity", quantity).Error
	})
	if err != nil {
		return nil, err
	}

	return s.GetOrderItemByID(id)
}

// DeleteOrderItem removes an item from an order. Only allowed while the order is
// open and the item has not been billed yet.
func (s *OrderItemService) DeleteOrderItem(id uuid.UUID) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		orderItem, err := findEditableOrderItem(tx, id)
		if err != nil {
			return err
		}
		return tx.Delete(orderItem).Error
	})
}

func ensureOrderOpen(tx *gorm.DB, orderID uuid.UUID) error {
	var order models.Order
	// Lock the order row so adding items, billing and the final payment
	// that closes the order cannot interleave.
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id", "status").Where("id = ?", orderID).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}
	if order.Status != models.OrderStatusOpen {
		return ErrOrderNotOpen
	}
	return nil
}

func findEditableOrderItem(tx *gorm.DB, id uuid.UUID) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	err := tx.Where("id = ?", id).First(&orderItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderItemNotFound
		}
		return nil, err
	}
	if orderItem.BillID != nil {
		return nil, ErrOrderItemBilled
	}
	if err := ensureOrderOpen(tx, orderItem.OrderID); err != nil {
		return nil, err
	}
	return &orderItem, nil
}
