package models

import "github.com/google/uuid"

type MenuItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	ItemNumber uint      `gorm:"not null: uniqueIndex"`
	Name       string    `gorm:"not null"`
	Category   string    `gorm:"not null"`
	Price      float32   `gorm:"not null; check:price > 0"`
}
