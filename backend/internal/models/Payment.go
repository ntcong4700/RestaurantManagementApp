package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`

	BillID uuid.UUID `gorm:"not null" json:"bill_id"`
	Bill   Bill      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	Amount float32   `gorm:"not null; check:amount>0" json:"amount"`
	PaidAt time.Time `json:"paid_at"`
}
