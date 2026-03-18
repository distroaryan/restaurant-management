package handler

import (
	"net/http"

	"github.com/distroaryan/restaurant-management/internal/errs"
	"github.com/distroaryan/restaurant-management/internal/models"
	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderRepo *repository.OrderRepository
	tableRepo *repository.TableRepository
	foodRepo  *repository.FoodRepository
}

type orderItemRequest struct {
	FoodID   string `json:"food_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type createOrderRequest struct {
	TableID string             `json:"table_id,omitempty"`
	Items   []orderItemRequest `json:"items"    binding:"required,min=1"`
}

func NewOrderHandler(repository *repository.Repository) *OrderHandler {
	return &OrderHandler{
		orderRepo: repository.Order,
		tableRepo: repository.Table,
		foodRepo:  repository.Food,
	}
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	orderID := c.Param("orderID")

	order, err := h.orderRepo.GetOrderById(c.Request.Context(), orderID)
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch order")
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// 1. Parse the request body
	var req createOrderRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// 2. Setup order and verify tableID if provided
	order := &models.Order{
		UserID:      c.GetString("userId"),
		Status:      models.OrderStatusPending,
	}

	if req.TableID != "" {
		table, err := h.tableRepo.GetTableById(c.Request.Context(), req.TableID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong tableId recieved"})
			return
		}
		order.TableID = table.ID
	}

	var orderItems []models.OrderItem
	var totalAmount float64

	for _, item := range req.Items {
		food, err := h.foodRepo.GetFoodById(c.Request.Context(), item.FoodID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "one or more food items not found"})
			return
		}

		unitPrice := food.Price * float64(item.Quantity)
		totalAmount += unitPrice

		orderItems = append(orderItems, models.OrderItem{
			FoodID:    food.ID,
			Quantity:  item.Quantity,
			UnitPrice: unitPrice,
		})
	}

	order.Items = orderItems
	order.TotalAmount = totalAmount

	err := h.orderRepo.CreateOrder(c.Request.Context(), order)
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch order")
		return
	}
	c.JSON(http.StatusOK, order)
}
