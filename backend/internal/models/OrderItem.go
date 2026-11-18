package models

import "github.com/google/uuid"

type OrderItem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`

	OrderID uuid.UUID `gorm:"not null; index" json:"order_id"`
	Order   Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	MenuItemID uuid.UUID `gorm:"not null" json:"menu_item_id"`
	MenuItem   MenuItem  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"menu_item"`

	Quantity uint32 `gorm:"not null; check:quantity > 0" json:"quantity"`
	// UnitPrice is the menu price at the time the item was ordered, so later
	// menu price changes do not affect this order or its bills.
	UnitPrice float32 `gorm:"not null;default:0" json:"unit_price"`

	BillID *uuid.UUID `json:"bill_id"`
	Bill   *Bill      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}
