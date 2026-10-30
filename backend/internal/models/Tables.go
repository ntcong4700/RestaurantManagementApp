package models

import "uuid"

const (
	TableStatusAvailable = "available"
	TableStatusOccupied  = "occupied"
	TableStatusReserved  = "reserved"
)

type Tables struct {
	ID          uuid.UUID `gorm:"primary_key,AUTO_INCREMENT"`
	TableNumber string    `gorm:"not null; uniqueIndex"`
	Status      bool      `gorm:"not null;default:available"`
}
