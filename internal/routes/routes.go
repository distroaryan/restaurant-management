package routes

import (
	"github.com/distroaryan/restaurant-management/internal/config"
	"github.com/distroaryan/restaurant-management/internal/handler"
	"github.com/distroaryan/restaurant-management/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *handler.Handler, cfg *config.Config) {
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
	tables.Use(middleware.Auth(cfg))
	{
		tables.GET("", handler.Table.GetAllTables)
		tables.GET("/:tableId", handler.Table.GetTableById)
		tables.POST("book-table/:tableId", handler.Table.BookTable)
		tables.POST("release-table/:tableId", handler.Table.ReleaseTable)
	}

	// Order
	orders := v1.Group("/orders")
	orders.Use(middleware.Auth(cfg))
	{
		orders.GET("/:orderId", handler.Order.GetOrderById)
		orders.POST("/create-order", handler.Order.CreateOrder)
	}
}