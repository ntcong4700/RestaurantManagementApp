package repository

import (
	"restaurant-management-app/backend/internal/models"

	"gorm.io/gorm"
)

type TableRepository interface {
	Create(table *models.DiningTable) error
}

type DiningTableRepository struct {
	db *gorm.DB
}

func NewDiningTableRepository(db *gorm.DB) *DiningTableRepository {
	return &DiningTableRepository{
		db: db,
	}
}

func (dtr *DiningTableRepository) Create(table *models.DiningTable) error {
	return dtr.db.Create(table).Error
}
