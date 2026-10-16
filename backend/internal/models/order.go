package models

import "time"

type Order struct {
	ID      uint64      `gorm:"primary_key" json:"id"`
	TableID uint64      `gorm:"not null" json:"table_id"`
	Table   DiningTable `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	Status    string     `gorm:"not null;default:available" json:"status"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	ClosedAt  *time.Time `gorm:"not null" json:"closed_at"`

	Bills      []Bill      `gorm:"foreignkey:BillID" json:"bills,omitempty"`
	OrderItems []OrderItem `gorm:"foreignkey:OrderID" json:"order_items,omitempty"`
}
