package service

import (
	"errors"
	"restaurant-management-app/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInvalidTableStatus = errors.New("invalid table status")
	ErrTableStatusManaged = errors.New("table status is managed by orders: a table with an open order must stay occupied, and occupied is only set by opening an order")
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
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		var current models.Tables
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", table.ID).First(&current).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTableNotFound
			}
			return err
		}

		var openOrders int64
		if err := tx.Model(&models.Order{}).
			Where("table_id = ? AND status = ?", table.ID, models.OrderStatusOpen).
			Count(&openOrders).Error; err != nil {
			return err
		}
		hasOpenOrder := openOrders > 0
		if hasOpenOrder != (table.Status == models.TableStatusOccupied) {
			return ErrTableStatusManaged
		}

		return tx.Model(&table).Select("table_number", "status").Updates(&table).Error
	})
	if err != nil {
		return nil, err
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
