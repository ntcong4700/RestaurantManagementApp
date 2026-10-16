package models

type OrderItem struct {
	ID      uint64 `gorm:"primary_key" json:"id"`
	OrderID uint64 `gorm:"not null" json:"order_id"`

	Order    Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	MenuItem MenuItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}
