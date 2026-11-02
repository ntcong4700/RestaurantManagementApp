package models

import "github.com/google/uuid"

const (
	TableStatusAvailable = "available"
	TableStatusOccupied  = "occupied"
	TableStatusReserved  = "reserved"
)

type Tables struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;"`
	TableNumber int       `gorm:"not null; uniqueIndex"`
	Status      string    `gorm:"not null;default:available"`
}
