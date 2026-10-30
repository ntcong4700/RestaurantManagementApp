package models

import (
	"time"
	"uuid"
)

const (
	OrderStatusOpen  = "open"
	OrderStatusClose = "close"
	OrderStatusPaid  = "paid"
)

type Order struct {
	ID uuid.UUID `gorm:"primary_key;AUTO_INCREMENT"`

	TableID uint32 `gorm:"not null"`
	Table   Tables `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Status   string     `gorm:"not null;default:open"`
	CreateAt time.Time  `gorm:"not null"`
	ClosedAt *time.Time `gorm:"not null"`

	Bills      []Bill      `gorm:"foreignkey:OrderID"`
	OrderItems []OrderItem `gorm:"foreignkey:OrderID"`
}
