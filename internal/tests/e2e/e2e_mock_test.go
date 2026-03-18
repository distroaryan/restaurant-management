package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/distroaryan/restaurant-management/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestaurantEndpoints(t *testing.T) {
	srv, cleanup := SetUpMockServer(t)
	defer cleanup()

	token := GenerateTestToken("test-user-id")
	requestHeaderKey := "Authorization"
	requestHeaderValue := "Bearer " + token 

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
		req.Header.Set(requestHeaderKey, requestHeaderValue)

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
		req.Header.Set(requestHeaderKey,requestHeaderValue)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var singleTable models.Table
		err = json.NewDecoder(resp.Body).Decode(&singleTable)
		assert.NoError(t, err)
		assert.Equal(t, srv.TestData.Tables[0], singleTable) 
	})

	// t.Run("Test Table Booking and Releasing (Auth Required)", func(t *testing.T) {
	// 	req, err := http.NewRequest("POST", srv.Server.URL+"/api/v1/tables/book-table/"+srv.TestData.Tables[0].ID.Hex(), nil)
	// 	require.NoError(t, err)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	resp, err := http.DefaultClient.Do(req)
	// 	require.NoError(t, err)
	// 	defer resp.Body.Close()
	// 	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 	// Validate that the returned response contains the table id and message
	// 	// e.g: var res map[string]interface{} ...

	// 	// Releasing table
	// 	req, err = http.NewRequest("POST", srv.Server.URL+"/api/v1/tables/release-table/"+srv.TestData.Tables[0].ID.Hex(), nil)
	// 	require.NoError(t, err)
	// 	req.Header.Set("Authorization", "Bearer "+token)

	// 	resp, err = http.DefaultClient.Do(req)
	// 	require.NoError(t, err)
	// 	defer resp.Body.Close()
	// 	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// })
}
