package service

import (
	"errors"
	"restaurant-management-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrTableNotFound     = errors.New("table not found")
	ErrTableHasOpenOrder = errors.New("table already has an open order")
	ErrOrderHasItems     = errors.New("order has items and cannot be deleted")
)

type OrderService struct {
	DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		DB: db,
	}
}

// CreateOrder opens a new order on a table. The table must exist and must not
// already have an open order; it is marked as occupied.
func (os *OrderService) CreateOrder(order models.Order) (*models.Order, error) {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	order.Status = models.OrderStatusOpen
	order.ClosedAt = nil

	err := os.DB.Transaction(func(tx *gorm.DB) error {
		var table models.Tables
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", order.TableID).First(&table).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTableNotFound
			}
			return err
		}

		var openOrders int64
		if err := tx.Model(&models.Order{}).
			Where("table_id = ? AND status = ?", order.TableID, models.OrderStatusOpen).
			Count(&openOrders).Error; err != nil {
			return err
		}
		if openOrders > 0 {
			return ErrTableHasOpenOrder
		}

		if err := tx.Omit(clause.Associations).Create(&order).Error; err != nil {
			return err
		}

		return tx.Model(&table).Update("status", models.TableStatusOccupied).Error
	})
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (os *OrderService) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := os.DB.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// DeleteOrder removes an order that has no items yet (e.g. opened by mistake)
// and releases its table if the order was still open.
func (os *OrderService) DeleteOrder(orderID uuid.UUID) error {
	return os.DB.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", orderID).First(&order).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrderNotFound
			}
			return err
		}

		var itemCount int64
		if err := tx.Model(&models.OrderItem{}).Where("order_id = ?", orderID).Count(&itemCount).Error; err != nil {
			return err
		}
		if itemCount > 0 {
			return ErrOrderHasItems
		}

		if err := tx.Delete(&order).Error; err != nil {
			return err
		}

		if order.Status == models.OrderStatusOpen {
			return tx.Model(&models.Tables{}).Where("id = ?", order.TableID).
				Update("status", models.TableStatusAvailable).Error
		}
		return nil
	})
}

func (os *OrderService) GetOrderById(orderID string) (*models.Order, error) {
	var order models.Order
	if err := os.DB.Preload("OrderItems.MenuItem").Where("id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}
