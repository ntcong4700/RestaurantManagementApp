package service

import (
	"errors"
	"restaurant-management-app/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// moneyEpsilon absorbs float32 rounding when comparing amounts.
const moneyEpsilon = 0.005

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidAmount   = errors.New("amount must be greater than 0")
	ErrOverpayment     = errors.New("amount exceeds remaining balance")
	ErrBillAlreadyPaid = errors.New("bill is already fully paid")
)

type PaymentService struct {
	DB *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{
		DB: db,
	}
}

// CreatePayment records a (possibly partial) payment for a bill.
// When the payment settles the bill and every item of the order is billed and
// every bill of the order is paid, the order is closed (status "paid") and its
// table is released (status "available").
func (s *PaymentService) CreatePayment(billID uuid.UUID, amount float32) (*models.Payment, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	var payment models.Payment

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var bill models.Bill
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", billID).First(&bill).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBillNotFound
			}
			return err
		}

		if err := ensureOrderOpen(tx, bill.OrderID); err != nil {
			return err
		}

		paid, err := paidAmount(tx, billID)
		if err != nil {
			return err
		}
		remaining := bill.TotalAmount - paid
		if remaining <= moneyEpsilon {
			return ErrBillAlreadyPaid
		}
		if amount > remaining+moneyEpsilon {
			return ErrOverpayment
		}

		payment = models.Payment{
			ID:     uuid.New(),
			BillID: billID,
			Amount: amount,
			PaidAt: time.Now(),
		}
		if err := tx.Omit(clause.Associations).Create(&payment).Error; err != nil {
			return err
		}

		if remaining-amount <= moneyEpsilon {
			return closeOrderIfSettled(tx, bill.OrderID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// GetPaymentByID returns a single payment.
func (s *PaymentService) GetPaymentByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := s.DB.Where("id = ?", id).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// GetPaymentsByBillID returns all payments of a bill, oldest first.
func (s *PaymentService) GetPaymentsByBillID(billID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	if err := s.DB.Where("bill_id = ?", billID).Order("paid_at").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// closeOrderIfSettled marks the order paid and frees its table when nothing is left to pay.
func closeOrderIfSettled(tx *gorm.DB, orderID uuid.UUID) error {
	var unbilled int64
	if err := tx.Model(&models.OrderItem{}).
		Where("order_id = ? AND bill_id IS NULL", orderID).Count(&unbilled).Error; err != nil {
		return err
	}
	if unbilled > 0 {
		return nil
	}

	var unpaidBills int64
	err := tx.Raw(`
		SELECT COUNT(*) FROM bills b
		WHERE b.order_id = ?
		  AND b.total_amount - COALESCE((SELECT SUM(p.amount) FROM payments p WHERE p.bill_id = b.id), 0) > ?`,
		orderID, moneyEpsilon).Scan(&unpaidBills).Error
	if err != nil {
		return err
	}
	if unpaidBills > 0 {
		return nil
	}

	var order models.Order
	if err := tx.Select("id", "table_id").Where("id = ?", orderID).First(&order).Error; err != nil {
		return err
	}

	now := time.Now()
	if err := tx.Model(&models.Order{}).Where("id = ?", orderID).Updates(map[string]any{
		"status":    models.OrderStatusPaid,
		"closed_at": now,
	}).Error; err != nil {
		return err
	}

	return tx.Model(&models.Tables{}).Where("id = ?", order.TableID).
		Update("status", models.TableStatusAvailable).Error
}
