package models

import (
	"time"
	"uuid"
)

type Payment struct {
	ID uuid.UUID `gorm:"primary_key;AUTO_INCREMENT"`

	BillID uint32 `gorm:"not null"`
	Bill   Bill   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Amount float32 `gorm:"not null; check:amount>0"`
	PaidAt time.Time
}
