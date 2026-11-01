package service

import (
	"restaurant-management-app/backend/internal/models"
	"uuid"

	"gorm.io/gorm"
)

type TableService struct {
	DB *gorm.DB
}

func NewTableService(db *gorm.DB) *TableService {
	return &TableService{
		DB: db,
	}
}

func (t *TableService) CreateTable(table models.Tables) (*models.Tables, error) {
	table.ID = uuid.NewV4()
	if err := t.DB.Create(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}
