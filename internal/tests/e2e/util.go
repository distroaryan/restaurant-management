package e2e

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/distroaryan/restaurant-management/internal/config"
	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/handler"
	"github.com/distroaryan/restaurant-management/internal/middleware"
	"github.com/distroaryan/restaurant-management/internal/models"
	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/distroaryan/restaurant-management/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

type TestData struct {
	Menus  []models.Menu
	Foods  []models.Food
	Tables []models.Table
}

type Server struct {
	Server     *httptest.Server
	Repository *repository.Repository
	TestData   TestData
}

func SetUpMockServer(t *testing.T) (*Server, func()) {
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
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	cfg := &config.Config{JwtSecret: "secret"}
	routes.RegisterRoutes(router, h, cfg)

	testServer := httptest.NewServer(router)

	cleanup := func() {
		testServer.Close()
		db.Close(ctx)
		container.Terminate(ctx)
	}

	app := &Server{
		Server:     testServer,
		Repository: r,
	}

	app.TestData = seedDatabase(t, app)

	return app, cleanup
}

func seedDatabase(t *testing.T, s *Server) TestData {
	t.Helper()
	var data TestData

	// Seed 100 Menus
	for i := range 100 {
		menu := &models.Menu{
			Name:        fmt.Sprintf("Menu %d", i),
			Description: fmt.Sprintf("Description %d", i),
		}
		err := s.Repository.Menu.CreateMenu(context.Background(), menu)
		require.NoError(t, err)
		data.Menus = append(data.Menus, *menu)
	}

	// Seed 50 Foods for first Menu
	for i := range 50 {
		food := &models.Food{
			Name:   fmt.Sprintf("Food %d", i),
			Price:  float64(10 + i),
			MenuID: data.Menus[0].ID,
		}
		err := s.Repository.Food.CreateFood(context.Background(), food)
		require.NoError(t, err)
		data.Foods = append(data.Foods, *food)
	}

	// Seed 100 Tables
	for i := range 10 {
		table := &models.Table{
			Name:   fmt.Sprintf("Table %d", i),
			Status: models.TableStatusAvailable,
		}
		err := s.Repository.Table.CreateTable(context.Background(), table)
		require.NoError(t, err)
		data.Tables = append(data.Tables, *table)
	}

	return data
}

func GenerateTestToken(userId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
	})
	tokenString, _ := token.SignedString([]byte("secret"))
	return tokenString
}
