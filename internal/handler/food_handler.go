package handler

import (
	"net/http"

	"github.com/distroaryan/restaurant-management/internal/errs"
	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/gin-gonic/gin"
)

type FoodHandler struct {
	foodRepo *repository.FoodRepository
}

func NewFoodHandler(repository *repository.Repository) *FoodHandler {
	return &FoodHandler{
		foodRepo: repository.Food,
	}
}

func (h *FoodHandler) GetAllFoods(c *gin.Context) {
	foods, err := h.foodRepo.GetAllFoods(c.Request.Context())
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch foods")
		return
	}
	c.JSON(http.StatusOK, foods)
}

func (h *FoodHandler) GetFoodById(c *gin.Context) {
	foodId := c.Param("foodId")

	food, err := h.foodRepo.GetFoodById(c.Request.Context(), foodId)
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch foods")
		return
	}
	c.JSON(http.StatusOK, food)
}

func (h *FoodHandler) GetFoodByMenu(c *gin.Context) {
	menuId := c.Param("menuId")

	food, err := h.foodRepo.GetFoodByMenu(c.Request.Context(), menuId)
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch foods")
		return 
	}
	c.JSON(http.StatusOK, food)
}
