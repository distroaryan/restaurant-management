package e2e

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/distroaryan/restaurant-management/internal/config"
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
	DB        *database.Database
	MenuRepo  *repository.MenuRepository
	FoodRepo  *repository.FoodRepository
	TableRepo *repository.TableRepository
	OrderRepo *repository.OrderRepository
}

func SetupApp(t *testing.T) (*TestApp, func()) {
	ctx := context.Background()

	container, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		require.NoError(t, err)
	}

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	db := database.Connect(uri, "restaurant_test")
	r := repository.NewRepository(db)
	h := handler.NewHandler(r)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := &config.Config{JwtSecret: "secret"}
	routes.RegisterRoutes(router, h, cfg)

	testServer := httptest.NewServer(router)

	cleanup := func() {
		testServer.Close()
		db.Close(ctx)
		container.Terminate(ctx)
	}

	return &TestApp{
		Server:    testServer,
		DB:        db,
		MenuRepo:  r.Menu,
		FoodRepo:  r.Food,
		TableRepo: r.Table,
		OrderRepo: r.Order,
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
func seedTable(t *testing.T, tableRepo *repository.TableRepository, name string) *models.Table {
	t.Helper()
	table := &models.Table{
		Name:          name,
		Status:        models.TableStatusAvailable,
	}
	err := tableRepo.CreateTable(context.Background(), table)
	require.NoError(t, err)
	return table
}
