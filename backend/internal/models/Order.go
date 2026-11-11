package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrderStatusOpen = "open"
	OrderStatusPaid = "paid"
)

type Order struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;"`

	TableID uuid.UUID `gorm:"not null"`
	Table   Tables    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Status    string    `gorm:"not null;default:open"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	ClosedAt  *time.Time

	Bills      []Bill      `gorm:"foreignkey:OrderID"`
	OrderItems []OrderItem `gorm:"foreignkey:OrderID"`
}
