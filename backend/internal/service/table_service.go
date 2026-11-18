package service

import (
	"errors"
	"restaurant-management-app/internal/models"

	"gorm.io/gorm"
)

var ErrInvalidTableStatus = errors.New("invalid table status")

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

func (t *TableService) GetAllTable() ([]*models.Tables, error) {
	var tables []*models.Tables
	if err := t.DB.Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (t *TableService) UpdateTables(table models.Tables) (*models.Tables, error) {
	if !isValidTableStatus(table.Status) {
		return nil, ErrInvalidTableStatus
	}
	result := t.DB.Model(&table).Select("table_number", "status").Updates(&table)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrTableNotFound
	}
	return &table, nil
}

func isValidTableStatus(status string) bool {
	switch status {
	case models.TableStatusAvailable, models.TableStatusOccupied, models.TableStatusReserved:
		return true
	}
	return false
}

func (t *TableService) DeleteTable(tableID string) error {
	result := t.DB.Where("id = ?", tableID).Delete(&models.Tables{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTableNotFound
	}
	return nil
}
