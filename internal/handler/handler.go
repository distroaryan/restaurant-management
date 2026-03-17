package handler

type Handler struct {
	FoodHandler FoodHandler
	MenuHandler MenuHandler
	OrderHandler OrderHandler
	TableHandler TableHandler
}