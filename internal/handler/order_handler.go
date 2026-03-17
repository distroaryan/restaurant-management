package handler

import (
	"net/http"
	"sync"

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
	TableID string             `json:"table_id" binding:"required"`
	Items   []orderItemRequest `json:"items"    binding:"required,min=1"`
}

func NewOrderHandler(orderRepo *repository.OrderRepository, tableRepo *repository.TableRepository, foodRepo *repository.FoodRepository) *OrderHandler {
	return &OrderHandler{
		orderRepo: orderRepo,
		tableRepo: tableRepo,
		foodRepo:  foodRepo,
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

	// 2. Verify the tableID
	table, err := h.tableRepo.GetTableById(c.Request.Context(), req.TableID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong tableId recieved"})
	}

	food_items := len(req.Items)

	type result struct {
		err  error
		item models.OrderItem
	}

	var wg sync.WaitGroup
	results := make([]result, len(req.Items))
	var mu sync.Mutex

	for idx := range food_items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			food, err := h.foodRepo.GetFoodById(c.Request.Context(), req.Items[idx].FoodID)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				results[idx].err = err
				return
			}

			results[idx] = result{
				item: models.OrderItem{
					FoodID:    food.ID,
					Quantity:  req.Items[idx].Quantity,
					UnitPrice: food.Price * float64(req.Items[idx].Quantity),
				},
			}

		}()
	}

	wg.Wait()

	totalAmount := 0
	var orderItems []models.OrderItem

	for _, r := range results {
		if r.err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "one or more food items not found"})
			return
		}
		totalAmount += int(r.item.UnitPrice)
		orderItems = append(orderItems, r.item)
	}

	order := &models.Order{
		TableID: table.ID,
		Status:  models.OrderStatusPending,
		Items:   orderItems,
	}

	err = h.orderRepo.CreateOrder(c.Request.Context(), order)
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch order")
		return
	}
	c.JSON(http.StatusOK, order)
}
