package models

import "github.com/google/uuid"

type OrderItem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;"`

	OrderID uuid.UUID `gorm:"not null; index"`
	Order   Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	MenuItemID uuid.UUID `gorm:"not null"`
	MenuItem   MenuItem  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Quantity uint32 `gorm:"not null; check:quantity > 0"`

	BillID *uuid.UUID
	Bill   *Bill `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
