package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/distroaryan/restaurant-management/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompleteRestaurantWorkflow(t *testing.T) {
	app, cleanup := SetupApp(t)
	defer cleanup()

	// 1. Seed database with mock data to reflect realistic operating environment
	menu1 := seedMenu(t, app.MenuRepo, "Breakfast", "Morning meals")
	food1 := seedFood(t, app.FoodRepo, "Pancakes", 10, menu1.ID)
	food2 := seedFood(t, app.FoodRepo, "Omelette", 15, menu1.ID)
	
	table1 := seedTable(t, app.TableRepo, "Table 1", 4, 0)
    
	// 2. Client fetches menus
	resp, err := http.Get(app.Server.URL + "/api/v1/menus")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var menus []models.Menu
	err = json.NewDecoder(resp.Body).Decode(&menus)
	require.NoError(t, err)
	assert.Len(t, menus, 1)
	assert.Equal(t, menu1.ID, menus[0].ID)

	// 3. Client views foods for a menu
	resp, err = http.Get(app.Server.URL + "/api/v1/foods/menu/" + menu1.ID.Hex())
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var foods []models.Food
	err = json.NewDecoder(resp.Body).Decode(&foods)
	require.NoError(t, err)
	assert.Len(t, foods, 2)

	// 4. Client fetches all tables
	resp, err = http.Get(app.Server.URL + "/api/v1/tables")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tables []models.Table
	err = json.NewDecoder(resp.Body).Decode(&tables)
	require.NoError(t, err)
	assert.Len(t, tables, 1)

	// 5. Client books table
	reqBody := []byte(`{"seats": 2}`)
	resp, err = http.Post(app.Server.URL+"/api/v1/tables/book-seats/"+table1.ID.Hex(), "application/json", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 6. Client creates order
	orderReqBody := []byte(`{
		"table_id": "` + table1.ID.Hex() + `",
		"items": [
			{"food_id": "` + food1.ID.Hex() + `", "quantity": 1},
			{"food_id": "` + food2.ID.Hex() + `", "quantity": 2}
		]
	}`)
	resp, err = http.Post(app.Server.URL+"/api/v1/orders/create-order", "application/json", bytes.NewBuffer(orderReqBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var createdOrder models.Order
	err = json.NewDecoder(resp.Body).Decode(&createdOrder)
	require.NoError(t, err)
	assert.Equal(t, table1.ID, createdOrder.TableID)
	assert.Len(t, createdOrder.Items, 2)
	assert.Equal(t, float64(10), createdOrder.Items[0].UnitPrice)
	assert.Equal(t, float64(30), createdOrder.Items[1].UnitPrice) // 15 * 2
	
	// 7. Client fetches order status
	resp, err = http.Get(app.Server.URL + "/api/v1/orders/" + createdOrder.ID.Hex())
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedOrder models.Order
	err = json.NewDecoder(resp.Body).Decode(&fetchedOrder)
	require.NoError(t, err)
	assert.Equal(t, createdOrder.ID, fetchedOrder.ID)

	// 8. Client releases the table seats
	releaseReqBody := []byte(`{"seats": 2}`)
	resp, err = http.Post(app.Server.URL+"/api/v1/tables/releaseSeats/"+table1.ID.Hex(), "application/json", bytes.NewBuffer(releaseReqBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
