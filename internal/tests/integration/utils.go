package integration

import (
	"context"
	"testing"

	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/models"
	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// -------------------------
// Collection Helpers
// -------------------------

func clearCollection(t *testing.T, db *database.Database, name string) {
	t.Helper()
	err := db.Client.Database("restaurant").Collection(name).Drop(context.Background())
	require.NoError(t, err)
}

func clearMenus(t *testing.T, db *database.Database) {
	clearCollection(t, db, "menus")
}

func clearFoods(t *testing.T, db *database.Database) {
	clearCollection(t, db, "foods")
}

func clearTables(t *testing.T, db *database.Database) {
	clearCollection(t, db, "tables")
}

func clearOrders(t *testing.T, db *database.Database) {
	clearCollection(t, db, "orders")
}

// ----------------------------------------------------------------------
// Seed helpers - used to insert some data in database
// ----------------------------------------------------------------------
func seedMenu(t *testing.T, menuRepo *repository.MenuRepository, name, description string) *models.Menu {
	t.Helper()
	menu := &models.Menu{
		Name:        name,
		Description: description,
	}

	err := menuRepo.CreateMenu(context.Background(), menu)
	require.NoError(t, err)
	return menu
}

func seedFood(t *testing.T, foodRepo *repository.FoodRepository, name string, price int, menuId bson.ObjectID) *models.Food {
	t.Helper()
	food := &models.Food{
		Name:   name,
		Price:  float64(price),
		MenuID: menuId,
	}
	err := foodRepo.CreateFood(context.Background(), food)
	require.NoError(t, err)
	return food
}

func seedOrder(t *testing.T, orderRepo *repository.OrderRepository, tableId bson.ObjectID, items []models.OrderItem) *models.Order {
	t.Helper()
	order := &models.Order{
		TableID: tableId,
		Items:   items,
	}
	err := orderRepo.CreateOrder(context.Background(), order)
	require.NoError(t, err)
	return order
}

func seedTable(t *testing.T, tableRepo *repository.TableRepository, name string, capacity int, reservedSeats int) *models.Table {
	t.Helper()
	table := &models.Table{
		Name:          name,
		Capacity:      capacity,
		ReservedSeats: reservedSeats,
		Status:        models.TableStatusAvailable,
	}
	err := tableRepo.CreateTable(context.Background(), table)
	require.NoError(t, err)
	return table
}
