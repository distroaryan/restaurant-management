package handler

import (
	"github.com/distroaryan/restaurant-management/internal/repository"
)

type Handler struct {
	Food  *FoodHandler
	Menu  *MenuHandler
	Order *OrderHandler
	Table *TableHandler
}

func NewHandler(repository *repository.Repository) *Handler {
	return &Handler{
		Food: NewFoodHandler(repository),
		Menu: NewMenuHandler(repository),
		Order: NewOrderHandler(repository),
		Table: NewTableRepositroy(repository),
	}
}