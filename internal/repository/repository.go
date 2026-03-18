package repository

import "github.com/distroaryan/restaurant-management/internal/database"

type Repository struct {
	Food  *FoodRepository
	Table *TableRepository
	Menu  *MenuRepository
	Order *OrderRepository
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{
		Food:  NewFoodRepository(db),
		Table: NewTableRepositroy(db),
		Menu:  NewMenuRepository(db),
		Order: NewOrderRepository(db),
	}
}
