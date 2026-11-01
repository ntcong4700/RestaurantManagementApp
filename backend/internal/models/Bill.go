package models

import (
	"time"

	"github.com/google/uuid"
)

type Bill struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	OrderID uuid.UUID `gorm:"not null"`
	Order   Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	TotalAmount float32 `gorm:"not null; check:total_amount>0"`
	CreatedAt   time.Time

	Payments []Payment `gorm:"foreignkey:BillID"`
}
