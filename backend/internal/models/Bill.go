package models

import "time"

type Bill struct {
	ID uint32 `gorm:"primary_key;AUTO_INCREMENT"`

	OrderID uint32 `gorm:"not null"`
	Order   Order  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	TotalAmount float32 `gorm:"not null; check:total_amount>0"`
	CreatedAt   time.Time

	Payments []Payment `gorm:"foreignkey:BillID"`
}
