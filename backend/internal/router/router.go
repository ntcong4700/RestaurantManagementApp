package router

import (
	"restaurant-management-app/internal/controller"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, corsAllowedOrigins []string) *gin.Engine {
	router := gin.New()
	router.Use(
		// Health checks run every few seconds; keep them out of the access log.
		gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/health", "/health/ready"}}),
		gin.Recovery(),
	)
	// Not behind a reverse proxy: use the real remote address as client IP
	// instead of trusting X-Forwarded-For from anyone.
	if err := router.SetTrustedProxies(nil); err != nil {
		panic(err)
	}
	if len(corsAllowedOrigins) > 0 {
		router.Use(corsMiddleware(corsAllowedOrigins))
	}

	tableController := controller.NewTableController(service.NewTableService(db))
	menuItemController := controller.NewMenuItemController(service.NewMenuItemsService(db))
	orderController := controller.NewOrderController(service.NewOrderService(db))
	orderItemController := controller.NewOrderItemController(service.NewOrderItemService(db))
	billController := controller.NewBillController(service.NewBillService(db))
	paymentController := controller.NewPaymentController(service.NewPaymentService(db))
	healthController := controller.NewHealthController(db)

	router.GET("/health", healthController.Live)
	router.GET("/health/ready", healthController.Ready)

	api := router.Group("/api")
	{
		// Tables
		api.POST("/tables", tableController.CreateTable)
		api.GET("/tables", tableController.GetAllTables)
		api.GET("/getalltables", tableController.GetAllTables) // kept for backward compatibility
		api.PUT("/tables/:id", tableController.UpdateTable)
		api.DELETE("/tables/:id", tableController.DeleteTable)

		// Menu items
		api.POST("/menu-items", menuItemController.CreateItem)
		api.GET("/menu-items", menuItemController.GetMenuItemList)
		api.GET("/menu-items/:id", menuItemController.GetMenuItemById)
		api.PATCH("/menu-items/:id", menuItemController.UpdateItem)
		api.DELETE("/menu-items/:id", menuItemController.DeleteItem)

		// Orders
		api.POST("/orders", orderController.CreateOrder)
		api.GET("/orders", orderController.GetAllOrders)
		api.GET("/orders/:id", orderController.GetOrderByID)
		api.DELETE("/orders/:id", orderController.DeleteOrder)

		// Order items
		api.POST("/orders/:id/items", orderItemController.CreateOrderItem)
		api.GET("/orders/:id/items", orderItemController.GetOrderItemsByOrder)
		api.GET("/order-items/:id", orderItemController.GetOrderItemByID)
		api.PATCH("/order-items/:id", orderItemController.UpdateOrderItem)
		api.DELETE("/order-items/:id", orderItemController.DeleteOrderItem)

		// Bills
		api.POST("/orders/:id/bills", billController.CreateBill)
		api.GET("/orders/:id/bills", billController.GetBillsByOrder)
		api.GET("/bills/:id", billController.GetBillByID)
		api.DELETE("/bills/:id", billController.DeleteBill)

		// Payments
		api.POST("/bills/:id/payments", paymentController.CreatePayment)
		api.GET("/bills/:id/payments", paymentController.GetPaymentsByBill)
		api.GET("/payments/:id", paymentController.GetPaymentByID)
	}

	return router
}
