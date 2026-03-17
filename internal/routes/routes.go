package routes

import (
	"github.com/distroaryan/restaurant-management/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *handler.Handler) {
	v1 := r.Group("/api/v1")

	// Menus
	menus := v1.Group("/menus")
	{
		menus.GET("", handler.Menu.GetAllMenus)
	}

	// Foods
	foods := v1.Group("/foods")
	{
		foods.GET("", handler.Food.GetAllFoods)
		foods.GET("/:foodId", handler.Food.GetFoodById)
		foods.GET("/menu/:menuId", handler.Food.GetFoodByMenu)
	}

	// Tables
	tables := v1.Group("/tables")
	{
		tables.GET("", handler.Table.GetAllTables)
		tables.GET("/:tableId", handler.Table.GetTableById)
		tables.POST("book-seats/:tableId", handler.Table.BookSeats)
		tables.POST("releaseSeats/:tableId", handler.Table.ReleaseSeats)
	}

	// Order
	orders := v1.Group("/orders")
	{
		orders.GET("/:orderId", handler.Order.GetOrderById)
		orders.POST("/create-order", handler.Order.CreateOrder)
	}
}