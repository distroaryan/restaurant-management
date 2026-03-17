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
		menus.GET("", handler.MenuHandler.GetAllMenus)
	}

	// Foods
	foods := v1.Group("/foods")
	{
		foods.GET("", handler.FoodHandler.GetAllFoods)
		foods.GET("/:foodId", handler.FoodHandler.GetFoodById)
		foods.GET("/:menuId", handler.FoodHandler.GetFoodByMenu)
	}

	// Tables
	tables := v1.Group("/tables")
	{
		tables.GET("", handler.TableHandler.GetAllTables)
		tables.GET("/:tableId", handler.TableHandler.GetTableById)
		tables.POST("book-seats/:tableId", handler.TableHandler.BookSeats)
		tables.POST("releaseSeats/:tableId", handler.TableHandler.ReleaseSeats)
	}

	// Order
	orders := v1.Group("/orders")
	{
		orders.GET("/:orderId", handler.OrderHandler.GetOrderById)
		orders.POST("/create-order", handler.OrderHandler.CreateOrder)
	}
}