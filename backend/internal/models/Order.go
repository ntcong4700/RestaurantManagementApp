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
	ID uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`

	TableID uuid.UUID `gorm:"not null" json:"table_id"`
	Table   Tables    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	Status    string     `gorm:"not null;default:open" json:"status"`
	CreatedAt time.Time  `gorm:"not null;autoCreateTime" json:"created_at"`
	ClosedAt  *time.Time `json:"closed_at"`

	Bills      []Bill      `gorm:"foreignkey:OrderID" json:"bills,omitempty"`
	OrderItems []OrderItem `gorm:"foreignkey:OrderID" json:"order_items,omitempty"`
}
