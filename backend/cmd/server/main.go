package main

import (
	"log"
	config2 "restaurant-management-app/backend/cmd/config"
	router2 "restaurant-management-app/backend/internal/router"
)

func main() {

	config, err := config2.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config2.DBConnection(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	router := router2.NewRouter(db)
	router.Run()
}
