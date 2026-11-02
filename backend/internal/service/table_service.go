package service

import (
	"restaurant-management-app/backend/internal/models"

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
	if err := t.DB.Create(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func (t *TableService) GetTable() ([]*models.Tables, error) {
	var tables []*models.Tables
	if err := t.DB.Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (t *TableService) UpdateTables(table models.Tables) (*models.Tables, error) {
	if err := t.DB.Save(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func (t *TableService) DeleteTable(tableID string) error {
	if err := t.DB.Where("id = ?", tableID).Delete(&models.Tables{}).Error; err != nil {
		return err
	}
	return nil
}
