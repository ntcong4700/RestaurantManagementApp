package main

import (
	router2 "restaurant-management-app/backend/internal/router"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	router := router2.NewRouter()
	router.Run()
}
