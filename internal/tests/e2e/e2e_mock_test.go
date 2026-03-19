package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/distroaryan/restaurant-management/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type orderItemRequest struct {
	FoodID   string `json:"food_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type createOrderRequest struct {
	TableID string             `json:"table_id,omitempty"`
	Items   []orderItemRequest `json:"items"    binding:"required,min=1"`
}

func TestRestaurantEndpoints(t *testing.T) {
	t.Log("🚀 Starting E2E Mock Test Suite...")
	srv, cleanup := SetUpMockServer(t)
	defer func() {
		cleanup()
		t.Log("🎉 Completed E2E Mock Test Suite!")
	}()

	testUserID := "test-user-id"
	token := GenerateTestToken(testUserID)
	reqAuthHeaderKey := "Authorization"
	reqAuthHeaderValue := "Bearer " + token

	serverURL := srv.Server.URL
	t.Run("Test Menu Routes", func(t *testing.T) {
		resp, err := http.Get(srv.Server.URL + "/api/v1/menus")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		var fetchedMenus []models.Menu
		err = json.NewDecoder(resp.Body).Decode(&fetchedMenus)
		assert.NoError(t, err)
		assert.Equal(t, srv.TestData.Menus, fetchedMenus)
	})

	t.Run("Test Food Routes", func(t *testing.T) {
		baseFoodRoute := "/api/v1/foods"
		allFoodRoute := serverURL + baseFoodRoute
		resp, err := http.Get(allFoodRoute)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var fetchedFoods []models.Food
		err = json.NewDecoder(resp.Body).Decode(&fetchedFoods)
		assert.NoError(t, err)
		assert.Equal(t, srv.TestData.Foods, fetchedFoods)
		resp.Body.Close()

		menuFoodRoute := serverURL + baseFoodRoute + fmt.Sprintf("/menu/%s", srv.TestData.Menus[0].ID.Hex())
		resp, err = http.Get(menuFoodRoute)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&fetchedFoods)
		assert.NoError(t, err)
		assert.Equal(t, srv.TestData.Foods, fetchedFoods)
		resp.Body.Close()

		specificFoodRoute := serverURL + baseFoodRoute + fmt.Sprintf("/%s", srv.TestData.Foods[0].ID.Hex())
		resp, err = http.Get(specificFoodRoute)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var fetchedSingleFood models.Food
		err = json.NewDecoder(resp.Body).Decode(&fetchedSingleFood)
		assert.NoError(t, err)
		assert.Equal(t, srv.TestData.Foods[0], fetchedSingleFood)
	})

	t.Run("Testing Table GET routes", func(t *testing.T) {
		baseTableRoute := "/api/v1/tables"

		allTableRoutes := serverURL + baseTableRoute
		req, err := http.NewRequest("GET", allTableRoutes, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var fetchedTables []models.Table
		err = json.NewDecoder(resp.Body).Decode(&fetchedTables)
		assert.NoError(t, err)
		assert.Equal(t, srv.TestData.Tables, fetchedTables)

		specificTableRoute := serverURL + baseTableRoute + fmt.Sprintf("/%s", srv.TestData.Tables[0].ID.Hex())
		req, err = http.NewRequest("GET", specificTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var singleTable models.Table
		err = json.NewDecoder(resp.Body).Decode(&singleTable)
		assert.NoError(t, err)
		assert.Equal(t, srv.TestData.Tables[0], singleTable)
	})

	t.Run("Test Table Booking and Releasing", func(t *testing.T) {
		baseTableRoute := "/api/v1/tables"
		tableID := srv.TestData.Tables[0].ID.Hex()

		// 1. Book the table
		bookTableRoute := serverURL + baseTableRoute + fmt.Sprintf("/book-table/%s", tableID)
		req, err := http.NewRequest("POST", bookTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// 2. Retrieve details and verify
		specificTableRoute := serverURL + baseTableRoute + fmt.Sprintf("/%s", tableID)
		req, err = http.NewRequest("GET", specificTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var singleTable models.Table
		err = json.NewDecoder(resp.Body).Decode(&singleTable)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Equal(t, "test-user-id", singleTable.UserID, "UserID should be set to test-user-id")
		assert.Equal(t, models.TableStatusFull, singleTable.Status, "Table status should be FULL")

		// 3. Try booking the table again
		req, err = http.NewRequest("POST", bookTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		resp.Body.Close()

		// 4. Release the table and verify its new status
		releaseTableRoute := serverURL + baseTableRoute + fmt.Sprintf("/release-table/%s", tableID)
		req, err = http.NewRequest("POST", releaseTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Retrieve details again to verify it is released
		req, err = http.NewRequest("GET", specificTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var releasedTable models.Table
		err = json.NewDecoder(resp.Body).Decode(&releasedTable)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Empty(t, releasedTable.UserID, "UserID should be empty after releasing")
		assert.Equal(t, models.TableStatusAvailable, releasedTable.Status, "Table status should be AVAILABLE")

		// 5. Try releasing the table again (should fail)
		req, err = http.NewRequest("POST", releaseTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		resp.Body.Close()

		// 6. Try releasing the table of any other table (should fail)
		releaseTableRoute = serverURL + baseTableRoute + fmt.Sprintf("/release-table/%s", srv.TestData.Tables[2].ID.Hex())
		req, err = http.NewRequest("POST", releaseTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Test Creating Order Endpoint", func(t *testing.T) {
		baseOrderRoute := "/api/v1/orders"
		createOrderRoute := serverURL + baseOrderRoute + "/create-order"
		reqContentHeaderKey := "Content-Type"
		reqContentHeaderValue := "application/json"

		tableId := srv.TestData.Tables[0].ID.Hex()
		foodID1 := srv.TestData.Foods[0].ID.Hex()
		foodID2 := srv.TestData.Foods[1].ID.Hex()
		reqBody := createOrderRequest{
			TableID: tableId,
			Items: []orderItemRequest{
				{
					FoodID:   foodID1,
					Quantity: 10,
				},
				{
					FoodID:   foodID2,
					Quantity: 15,
				},
			},
		}

		body, err := json.Marshal(reqBody)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", createOrderRoute, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)
		req.Header.Set(reqContentHeaderKey, reqContentHeaderValue)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		fetchOrderRoute := serverURL + baseOrderRoute + fmt.Sprintf("/user/%s", testUserID)
		req, err = http.NewRequest("GET", fetchOrderRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var fetchdOrders []models.Order
		err = json.NewDecoder(resp.Body).Decode(&fetchdOrders)
		require.NoError(t, err)
		assert.NotEmpty(t, fetchdOrders)

		fetchOrderRoute = serverURL + baseOrderRoute + fmt.Sprintf("/%s", "invalid-test-user-id")
		req, err = http.NewRequest("GET", fetchOrderRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestConcurrentBooking_SameTable(t *testing.T) {
	t.Log("🚀 Starting E2E Mock Test Suite...")
	srv, cleanup := SetUpMockServer(t)
	defer func() {
		cleanup()
		t.Log("🎉 Completed E2E Mock Test Suite!")
	}()

	users := 100

	tableID := srv.TestData.Tables[0].ID.Hex()
	bookTableRoute := srv.Server.URL + "/api/v1/tables" + fmt.Sprintf("/book-table/%s", tableID)

	var requests []*http.Request
	for i := range users {
		testUserID := fmt.Sprintf("test-user-id-%d", i)
		token := GenerateTestToken(testUserID)
		reqAuthHeaderKey := "Authorization"
		reqAuthHeaderValue := "Bearer " + token

		req, err := http.NewRequest("POST", bookTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)
		requests = append(requests, req)
	}

	var successCount atomic.Int32
	var wg sync.WaitGroup

	for i := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.DefaultClient.Do(requests[i])
			if err == nil && resp.StatusCode == http.StatusOK {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(1), successCount.Load())
}

func TestConcurrentBooking_DifferentTable(t *testing.T) {
	t.Log("🚀 Starting E2E Mock Test Suite...")
	srv, cleanup := SetUpMockServer(t)
	defer func() {
		cleanup()
		t.Log("🎉 Completed E2E Mock Test Suite!")
	}()

	users := 100
	
	var requests []*http.Request
	for i := range users {
		tableID := srv.TestData.Tables[i].ID.Hex()
		bookTableRoute := srv.Server.URL + "/api/v1/tables" + fmt.Sprintf("/book-table/%s", tableID)
		testUserID := fmt.Sprintf("test-user-id-%d", i)
		token := GenerateTestToken(testUserID)
		reqAuthHeaderKey := "Authorization"
		reqAuthHeaderValue := "Bearer " + token

		req, err := http.NewRequest("POST", bookTableRoute, nil)
		require.NoError(t, err)
		req.Header.Set(reqAuthHeaderKey, reqAuthHeaderValue)
		requests = append(requests, req)
	}

	var successCount atomic.Int32
	var wg sync.WaitGroup

	for i := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.DefaultClient.Do(requests[i])
			if err == nil && resp.StatusCode == http.StatusOK {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(100), successCount.Load())
}
