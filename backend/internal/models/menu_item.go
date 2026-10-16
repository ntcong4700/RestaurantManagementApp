package models

type MenuItem struct {
	ID       uint64  `gorm:"primary_key" json:"id"`
	Name     string  `gorm:"not null" json:"name"`
	Category string  `gorm:"not null" json:"category"`
	Price    float64 `gorm:"not null" json:"price"`
}
