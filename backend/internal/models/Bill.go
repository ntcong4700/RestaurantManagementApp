package models

import (
	"time"

	"github.com/google/uuid"
)

type Bill struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	OrderID uuid.UUID `gorm:"not null" json:"order_id"`
	Order   Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	TotalAmount float32   `gorm:"not null; check:total_amount>0" json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`

	Payments   []Payment   `gorm:"foreignkey:BillID" json:"payments"`
	OrderItems []OrderItem `gorm:"foreignkey:BillID" json:"order_items"`
}
