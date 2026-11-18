package service

import (
	"errors"
	"restaurant-management-app/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrBillNotFound        = errors.New("bill not found")
	ErrNoItemsToBill       = errors.New("no unbilled items to bill")
	ErrItemNotInOrder      = errors.New("order item does not belong to this order")
	ErrBillHasPayments     = errors.New("bill already has payments")
	ErrDuplicateOrderItems = errors.New("duplicate order item ids")
)

type BillService struct {
	DB *gorm.DB
}

func NewBillService(db *gorm.DB) *BillService {
	return &BillService{
		DB: db,
	}
}

// CreateBill creates a bill for an open order.
// If orderItemIDs is empty, every unbilled item of the order goes on the bill;
// otherwise only the given items do (split bill). The total is computed from
// menu prices × quantities and the items are linked to the new bill.
func (s *BillService) CreateBill(orderID uuid.UUID, orderItemIDs []uuid.UUID) (*models.Bill, error) {
	var bill models.Bill

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := ensureOrderOpen(tx, orderID); err != nil {
			return err
		}

		var items []models.OrderItem
		// The order row is locked by ensureOrderOpen, so concurrent bills for
		// the same order are serialized and cannot grab the same items.
		query := tx.Preload("MenuItem")

		if len(orderItemIDs) == 0 {
			if err := query.Where("order_id = ? AND bill_id IS NULL", orderID).Find(&items).Error; err != nil {
				return err
			}
		} else {
			unique := make(map[uuid.UUID]struct{}, len(orderItemIDs))
			for _, id := range orderItemIDs {
				if _, dup := unique[id]; dup {
					return ErrDuplicateOrderItems
				}
				unique[id] = struct{}{}
			}

			if err := query.Where("id IN ?", orderItemIDs).Find(&items).Error; err != nil {
				return err
			}
			if len(items) != len(orderItemIDs) {
				return ErrOrderItemNotFound
			}
			for _, item := range items {
				if item.OrderID != orderID {
					return ErrItemNotInOrder
				}
				if item.BillID != nil {
					return ErrOrderItemBilled
				}
			}
		}

		if len(items) == 0 {
			return ErrNoItemsToBill
		}

		var total float32
		ids := make([]uuid.UUID, 0, len(items))
		for _, item := range items {
			total += item.MenuItem.Price * float32(item.Quantity)
			ids = append(ids, item.ID)
		}

		bill = models.Bill{
			ID:          uuid.New(),
			OrderID:     orderID,
			TotalAmount: total,
			CreatedAt:   time.Now(),
		}
		if err := tx.Omit(clause.Associations).Create(&bill).Error; err != nil {
			return err
		}

		return tx.Model(&models.OrderItem{}).Where("id IN ?", ids).Update("bill_id", bill.ID).Error
	})
	if err != nil {
		return nil, err
	}

	return s.GetBillByID(bill.ID)
}

// GetBillByID returns a bill with its items (and their menu items) and payments.
func (s *BillService) GetBillByID(id uuid.UUID) (*models.Bill, error) {
	var bill models.Bill
	err := s.DB.Preload("OrderItems.MenuItem").Preload("Payments").
		Where("id = ?", id).First(&bill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBillNotFound
		}
		return nil, err
	}
	return &bill, nil
}

// GetBillsByOrderID returns all bills of an order.
func (s *BillService) GetBillsByOrderID(orderID uuid.UUID) ([]models.Bill, error) {
	var bills []models.Bill
	err := s.DB.Preload("OrderItems.MenuItem").Preload("Payments").
		Where("order_id = ?", orderID).Order("created_at").Find(&bills).Error
	if err != nil {
		return nil, err
	}
	return bills, nil
}

// GetRemainingAmount returns how much of the bill is still unpaid.
func (s *BillService) GetRemainingAmount(billID uuid.UUID) (float32, error) {
	var bill models.Bill
	if err := s.DB.Select("id", "total_amount").Where("id = ?", billID).First(&bill).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrBillNotFound
		}
		return 0, err
	}
	paid, err := paidAmount(s.DB, billID)
	if err != nil {
		return 0, err
	}
	return bill.TotalAmount - paid, nil
}

// DeleteBill cancels a bill that has no payments yet; its items become unbilled again.
func (s *BillService) DeleteBill(id uuid.UUID) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var bill models.Bill
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&bill).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBillNotFound
			}
			return err
		}

		if err := ensureOrderOpen(tx, bill.OrderID); err != nil {
			return err
		}

		var paymentCount int64
		if err := tx.Model(&models.Payment{}).Where("bill_id = ?", id).Count(&paymentCount).Error; err != nil {
			return err
		}
		if paymentCount > 0 {
			return ErrBillHasPayments
		}

		if err := tx.Model(&models.OrderItem{}).Where("bill_id = ?", id).Update("bill_id", nil).Error; err != nil {
			return err
		}
		return tx.Delete(&bill).Error
	})
}

func paidAmount(tx *gorm.DB, billID uuid.UUID) (float32, error) {
	var paid float64
	err := tx.Model(&models.Payment{}).Where("bill_id = ?", billID).
		Select("COALESCE(SUM(amount), 0)").Scan(&paid).Error
	return float32(paid), err
}
