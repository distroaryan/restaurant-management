package integration

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/models"
	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var testDb *database.Database

// ----------------------------------------------------------------
// TestMain - one container for entire test suites
// ----------------------------------------------------------------

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		panic("Failed to start mongodb container: " + err.Error())
	}

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		panic("Failed to get connection string: " + err.Error())
	}

	testDb = database.Connect(uri, "restaurant_test")

	code := m.Run()

	testDb.Close(ctx)
	container.Terminate(ctx)

	os.Exit(code)
}

// ============================================================
// Menu Repository Tests
// ============================================================

func TestMenuRepository(t *testing.T) {
	clearMenus(t, testDb)

	menuRepo := repository.NewMenuRepository(testDb)
	ctx := context.Background()

	expected_breakfast := seedMenu(t, menuRepo, "Breakfast", "Morning Meals")

	t.Run("Testing GetMenuById method", func(t *testing.T) {
		actual_breakfast, err := menuRepo.GetMenuById(ctx, expected_breakfast.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, expected_breakfast, actual_breakfast)

		_, err = menuRepo.GetMenuById(ctx, bson.NewObjectID().Hex())
		assert.Error(t, err)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)

		_, err = menuRepo.GetMenuById(ctx, "not a valid mongodb id")
		assert.Error(t, err)
	})

	t.Run("Testing GetAllMenu method", func(t *testing.T) {
		expected_lunch := seedMenu(t, menuRepo, "Lunch", "Lunch meals")
		expected_dinner := seedMenu(t, menuRepo, "Dinner", "Dinner meals")

		expected_menus := []*models.Menu{expected_breakfast, expected_lunch, expected_dinner}

		actual_menus, err := menuRepo.GetAllMenu(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected_menus, actual_menus)

		clearMenus(t, testDb)

		menus, err := menuRepo.GetAllMenu(ctx)
		assert.NoError(t, err)
		assert.Empty(t, menus)
	})
}

// ============================================================
// Food Repository Tests
// ============================================================

func TestFoodRepository(t *testing.T) {
	clearFoods(t, testDb)

	ctx := context.Background()
	foodRepo := repository.NewFoodRepository(testDb)

	pizzaMenuId := bson.NewObjectID()
	margheritaPizza := seedFood(t, foodRepo, "Margerrita pizza", 1, pizzaMenuId)
	garlicBread := seedFood(t, foodRepo, "Garlic Bread", 2, pizzaMenuId)

	t.Run("Testing GetFoodById and GetFoodByName method", func(t *testing.T) {
		fetchedPizza, err := foodRepo.GetFoodById(ctx, margheritaPizza.ID.Hex())
		assert.NoError(t, err, "should not return error for valid food ID")
		assert.Equal(t, margheritaPizza, fetchedPizza, "fetched food should match inserted food")

		_, err = foodRepo.GetFoodById(ctx, bson.NewObjectID().Hex())
		assert.Error(t, err, "should return error for non-existent food ID")
		assert.ErrorIs(t, err, mongo.ErrNoDocuments, "error should be ErrNoDocuments")

		_, err = foodRepo.GetFoodById(ctx, "invalid_id")
		assert.Error(t, err, "should return error for invalid ObjectID string")
	})

	t.Run("Testing GetFoodByMenu method", func(t *testing.T) {
		fetchedFoodsByMenu, err := foodRepo.GetFoodByMenu(ctx, pizzaMenuId.Hex())
		expectedFoods := []*models.Food{margheritaPizza, garlicBread}
		assert.NoError(t, err, "should not return error for valid menu ID")
		assert.Equal(t, expectedFoods, fetchedFoodsByMenu, "should return all foods associated with the menu")

		_, err = foodRepo.GetFoodByMenu(ctx, bson.NewObjectID().Hex())
		assert.NoError(t, err, "should not return error for random menu ID but return empty array")

		_, err = foodRepo.GetFoodByMenu(ctx, "invalid_id")
		assert.Error(t, err, "should return error for invalid menu ObjectID string")
	})

	t.Run("Testing GetAllFoods method", func(t *testing.T) {
		allFoods, err := foodRepo.GetAllFoods(ctx)
		expectedFoods := []*models.Food{margheritaPizza, garlicBread}
		assert.NoError(t, err, "should not return error when fetching all foods")
		assert.Equal(t, expectedFoods, allFoods, "should return all inserted foods")

		// Test Empty Foods List
		clearFoods(t, testDb)
		emptyFoods, err := foodRepo.GetAllFoods(ctx)
		assert.NoError(t, err, "should not return error when fetching all foods on empty collection")
		assert.Empty(t, emptyFoods, "should return empty slice when no foods exist")
	})
}

// ============================================================
// Order Repository Tests
// ============================================================

func TestOrderRepository(t *testing.T) {
	clearOrders(t, testDb)

	ctx := context.Background()
	orderRepo := repository.NewOrderRepository(testDb)

	tableOneID := bson.NewObjectID()
	foodItems := []models.OrderItem{
		{FoodID: bson.NewObjectID(), Quantity: 1, UnitPrice: 15.0},
		{FoodID: bson.NewObjectID(), Quantity: 2, UnitPrice: 10.0},
	}

	tableOneOrder := seedOrder(t, orderRepo, tableOneID, foodItems)

	t.Run("Testing GetOrderById method", func(t *testing.T) {
		fetchedOrder, err := orderRepo.GetOrderById(ctx, tableOneOrder.ID.Hex())
		assert.NoError(t, err, "should not return error for valid order ID")
		assert.Equal(t, tableOneOrder, fetchedOrder, "fetched order should match inserted order")

		_, err = orderRepo.GetOrderById(ctx, bson.NewObjectID().Hex())
		assert.Error(t, err, "should return error for non-existent order ID")
		assert.ErrorIs(t, err, mongo.ErrNoDocuments, "error should be ErrNoDocuments")

		_, err = orderRepo.GetOrderById(ctx, "invalid_id")
		assert.Error(t, err, "should return error for invalid ObjectID string")
	})

	t.Run("Testing GetOrdersByUserID method", func(t *testing.T) {
		emptyTableOrders, err := orderRepo.GetOrdersByUserID(ctx, bson.NewObjectID().Hex())
		assert.NoError(t, err, "should not return error for non-existent table ID")
		assert.Empty(t, emptyTableOrders, "should return empty list for table with no orders")
	})

	t.Run("Testing UpdateOrderStatus method", func(t *testing.T) {
		err := orderRepo.UpdateOrderStatus(ctx, tableOneOrder.ID.Hex(), models.OrderStatusCompleted)
		assert.NoError(t, err, "should successfully update order status for valid order")

		updatedOrder, _ := orderRepo.GetOrderById(ctx, tableOneOrder.ID.Hex())
		assert.Equal(t, models.OrderStatusCompleted, updatedOrder.Status, "order status should be updated to completed")

		err = orderRepo.UpdateOrderStatus(ctx, bson.NewObjectID().Hex(), models.OrderStatusProcessing)
		assert.NoError(t, err, "updating non-existent order may not return error depending on mongo driver behavior")

		err = orderRepo.UpdateOrderStatus(ctx, "invalid_id", models.OrderStatusCancelled)
		assert.Error(t, err, "should return error when updating order with invalid ID")
	})
}

// ============================================================
// Table Repository Tests
// ============================================================
func TestTableRepository(t *testing.T) {
	clearTables(t, testDb)

	ctx := context.Background()
	tableRepo := repository.NewTableRepositroy(testDb)

	tableOne := seedTable(t, tableRepo, "Table1")
	tableTwo := seedTable(t, tableRepo, "Table2")

	t.Run("testing GetTableById method", func(t *testing.T) {
		actualTable, err := tableRepo.GetTableById(ctx, tableOne.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, tableOne, actualTable)
	
		_, err = tableRepo.GetTableById(ctx, bson.NewObjectID().Hex())
		assert.Error(t, err)
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
	
		_, err = tableRepo.GetTableById(ctx, "invalid-mongo-id")
		assert.Error(t, err)
	})

	t.Run("Testing BookTable method", func(t *testing.T) {
		err := tableRepo.BookTable(ctx, tableOne.ID.Hex(), "user1")
		assert.NoError(t, err)
	
		bookTableOne, err := tableRepo.GetTableById(ctx, tableOne.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, "user1", bookTableOne.UserID)
		assert.Equal(t, models.TableStatusFull, bookTableOne.Status)
	
		err = tableRepo.BookTable(ctx, tableOne.ID.Hex(), "user2")
		assert.Error(t, err)
	
		err = tableRepo.BookTable(ctx, bson.NewObjectID().Hex(), "user1")
		assert.Error(t, err)
	
		err = tableRepo.BookTable(ctx, "invalid-mongo-id", "user1")
		assert.Error(t, err)

	})

	t.Run("Testing ReleaseTable method", func(t *testing.T) {
		err := tableRepo.ReleaseTable(ctx, tableOne.ID.Hex())
		assert.NoError(t, err)
	
		releaseTableOne, err := tableRepo.GetTableById(ctx, tableOne.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, "", releaseTableOne.UserID)
		assert.Equal(t, models.TableStatusAvailable, releaseTableOne.Status)
	
		err = tableRepo.ReleaseTable(ctx, bson.NewObjectID().Hex())
		assert.Error(t, err)
	
		err = tableRepo.ReleaseTable(ctx, "invalid-mongo-id")
		assert.Error(t, err)

	})

	t.Run("Testing GetAllTables method", func(t *testing.T) {
		expectedTables := []*models.Table{tableOne, tableTwo}
		tables, err := tableRepo.GetAllTables(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedTables, tables)
	
		clearTables(t, testDb)
		tables, err = tableRepo.GetAllTables(ctx)
		assert.NoError(t, err)
		assert.Empty(t, tables)
	})
}

func TestConcurrentTableBooking(t *testing.T) {
	clearTables(t, testDb)

	ctx := context.Background()
	tableRepo := repository.NewTableRepositroy(testDb)

	table := seedTable(t, tableRepo, "TableOne")

	// Testing Booking concurrency by simulatenously attempting to book
	users := 100
	var wg sync.WaitGroup
	var successCount atomic.Int32

	for range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := tableRepo.BookTable(ctx, table.ID.Hex(), "concurrentUser")
			if err == nil {
				successCount.Add(1)
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, int32(1), successCount.Load())

	finalState, err := tableRepo.GetTableById(ctx, table.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, models.TableStatusFull, finalState.Status)
	assert.Equal(t, "concurrentUser", finalState.UserID)
}
