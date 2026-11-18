package config

import (
	"fmt"
	"restaurant-management-app/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return nil, fmt.Errorf("failed to create UUID extension: %w", err)
	}

	if err := db.AutoMigrate(
		&models.Order{},
		models.OrderItem{},
		models.Bill{},
		models.MenuItem{},
		models.Tables{},
		models.Payment{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}
	// Backfill unit_price for order items created before the column existed.
	if err := db.Exec(`UPDATE order_items oi SET unit_price = mi.price
		FROM menu_items mi WHERE oi.menu_item_id = mi.id AND oi.unit_price = 0`).Error; err != nil {
		return nil, fmt.Errorf("failed to backfill order item prices: %w", err)
	}
	return db, nil
}
