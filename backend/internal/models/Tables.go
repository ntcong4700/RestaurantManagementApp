package models

const (
	TableStatusAvailable = "available"
	TableStatusOccupied  = "occupied"
	TableStatusReserved  = "reserved"
)

type Tables struct {
	ID          uint64 `gorm:"primary_key,AUTO_INCREMENT"`
	TableNumber string `gorm:"not null; uniqueIndex"`
	Status      bool   `gorm:"not null;default:available"`
}
