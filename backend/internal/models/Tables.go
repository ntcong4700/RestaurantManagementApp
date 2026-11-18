package models

import "github.com/google/uuid"

const (
	TableStatusAvailable = "available"
	TableStatusOccupied  = "occupied"
	TableStatusReserved  = "reserved"
)

type Tables struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	TableNumber int       `gorm:"not null; uniqueIndex" json:"table_number"`
	Status      string    `gorm:"not null;default:available" json:"status"`
}
