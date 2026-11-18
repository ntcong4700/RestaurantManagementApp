package models

import "github.com/google/uuid"

type MenuItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ItemNumber uint      `gorm:"not null;uniqueIndex" json:"item_number"`
	Name       string    `gorm:"not null" json:"name"`
	Category   string    `gorm:"not null" json:"category"`
	Price      float32   `gorm:"not null; check:price > 0" json:"price"`
}
