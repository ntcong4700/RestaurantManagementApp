package models

import "time"

type Payment struct {
	ID uint32 `gorm:"primary_key;AUTO_INCREMENT"`

	BillID uint32 `gorm:"not null"`
	Bill   Bill   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Amount float32 `gorm:"not null; check:amount>0"`
	PaidAt time.Time
}
