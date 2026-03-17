package handler

import (
	"net/http"

	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/gin-gonic/gin"
)

type FoodHandler struct {
	foodRepo *repository.FoodRepository
}

func NewFoodHandler(foodRepo *repository.FoodRepository) *FoodHandler {
	return &FoodHandler{
		foodRepo: foodRepo,
	}
}

func (h *FoodHandler) GetAllFoods(c *gin.Context) {
	foods, err := h.foodRepo.GetAllFoods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch foods",
		})
		return
	}
	c.JSON(http.StatusOK, foods)
}

func (h *FoodHandler) GetFoodById(c *gin.Context) {
	foodId := c.Param("foodId")

	food, err := h.foodRepo.GetFoodById(c.Request.Context(), foodId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch foods",
		})
		return
	}
	c.JSON(http.StatusOK, food)
}

func (h *FoodHandler) GetFoodByMenu(c *gin.Context) {
	menuId := c.Param("menuId")

	food, err := h.foodRepo.GetFoodByMenu(c.Request.Context(), menuId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch foods",
		})
		return 
	}
	c.JSON(http.StatusOK, food)
}
