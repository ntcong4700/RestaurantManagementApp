package models

import "time"

type Bill struct {
	ID      uint64 `gorm:"primary_key" json:"id"`
	OrderID uint64 `gorm:"not null" json:"order_id"`
	Order   Order  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	Status    string    `gorm:"not null,default:'Not paid'" json:"status"`
	CreatedAt time.Time `json:"createdAt"`

	Payments   []Payment   `gorm:"foreignkey:BillID" json:"payments,omitempty"`
	OrderItems []OrderItem `gorm:"foreignkey:BillID" json:"order_items,omitempty"`
}
