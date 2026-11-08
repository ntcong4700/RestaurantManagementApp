package router

import (
	"restaurant-management-app/internal/controller"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	tableService := service.NewTableService(db)
	tableController := controller.NewTableController(tableService)

	public := router.Group("/api")
	{
		public.POST("/tables", tableController.CreateTable)
		public.GET("/tables", tableController.GetTables)
		public.PUT("/tables/:id", tableController.UpdateTable)
		public.DELETE("/tables/:id", tableController.DeleteTable)
	}

	return router
}
