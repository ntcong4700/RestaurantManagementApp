package models

type OrderItem struct {
	ID uint32 `gorm:"primary_key;AUTO_INCREMENT"`

	OrderID uint32 `gorm:"not null; index"`
	Order   Order  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	MenuItemID uint32   `gorm:"not null"`
	MenuItem   MenuItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Quantity uint32 `gorm:"not null; check:quantity > 0"`

	BillID *uint32 `gorm:"not null"`
	Bill   *Bill   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
