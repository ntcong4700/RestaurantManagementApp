package models

import "uuid"

type MenuItem struct {
	ID       uuid.UUID `gorm:"primary_key;AUTO_INCREMENT"`
	Name     string    `gorm:"not null"`
	Category string    `gorm:"not null"`
	Price    float32   `gorm:"not null; check:price > 0"`
}
