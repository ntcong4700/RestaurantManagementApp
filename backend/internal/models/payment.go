package models

import "time"

type Payment struct {
	ID      uint64 `gorm:"primary_key" json:"id"`
	bill_id uint64 `gorm:"not null" json:"bill_id"`
	Bill    Bill   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	Amount float64   `gorm:"not null" json:"amount"`
	PaidAt time.Time `gorm:"not null" json:"paid_at"`
}
