package e2e

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/handler"
	"github.com/distroaryan/restaurant-management/internal/models"
	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/distroaryan/restaurant-management/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TestApp struct {
	Server    *httptest.Server
	DB        *database.DBEngine
	MenuRepo  *repository.MenuRepository
	FoodRepo  *repository.FoodRepository
	TableRepo *repository.TableRepository
	OrderRepo *repository.OrderRepository
}

func SetupApp(t *testing.T) (*TestApp, func()) {
	ctx := context.Background()

	container, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		t.Fatalf("Failed to start mongodb container: %v", err)
	}

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	db := database.Connect(uri)

	menuRepo := repository.NewMenuRepository(db)
	foodRepo := repository.NewFoodRepository(db)
	tableRepo := repository.NewTableRepositroy(db)
	orderRepo := repository.NewOrderRepository(db)

	h := &handler.Handler{
		MenuHandler:  *handler.NewMenuHandler(menuRepo),
		FoodHandler:  *handler.NewFoodHandler(foodRepo),
		TableHandler: *handler.NewTableRepositroy(tableRepo),
		OrderHandler: *handler.NewOrderHandler(orderRepo, tableRepo, foodRepo),
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.RegisterRoutes(router, h)

	testServer := httptest.NewServer(router)

	cleanup := func() {
		testServer.Close()
		db.Close(ctx)
		container.Terminate(ctx)
	}

	return &TestApp{
		Server:    testServer,
		DB:        db,
		MenuRepo:  menuRepo,
		FoodRepo:  foodRepo,
		TableRepo: tableRepo,
		OrderRepo: orderRepo,
	}, cleanup
}

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
