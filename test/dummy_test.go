package test

import (
	"bytes"
{{- if index .Modules "postgres"}}
	"context"
{{- end}}
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func truncateDummy(t *testing.T) {
{{- if index .Modules "postgres"}}
	t.Helper()
	_, err := testDB.ExecContext(context.Background(), "TRUNCATE TABLE dummy RESTART IDENTITY")
	require.NoError(t, err)
{{- end}}
}

func createDummy(t *testing.T, name string) map[string]interface{} {
	t.Helper()
	body := fmt.Sprintf(`{"name": %q}`, name)
	resp, err := http.Post(ts.URL+"/api/v1/dummy", "application/json", bytes.NewBufferString(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	return result
}

// --- POST /api/v1/dummy ---

func TestCreateDummy_Success(t *testing.T) {
	truncateDummy(t)

	body := `{"name": "integration_test"}`
	resp, err := http.Post(ts.URL+"/api/v1/dummy", "application/json", bytes.NewBufferString(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "integration_test", result["name"])
	assert.NotZero(t, result["id"])
}

func TestCreateDummy_EmptyName(t *testing.T) {
	body := `{"name": ""}`
	resp, err := http.Post(ts.URL+"/api/v1/dummy", "application/json", bytes.NewBufferString(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "bad_request", result["code"])
	assert.Equal(t, "name is required", result["message"])
}

func TestCreateDummy_InvalidJSON(t *testing.T) {
	body := `not valid json`
	resp, err := http.Post(ts.URL+"/api/v1/dummy", "application/json", bytes.NewBufferString(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "bad_request", result["code"])
}

// --- GET /api/v1/dummy ---

func TestListDummies_ReturnsCreated(t *testing.T) {
	truncateDummy(t)

	createDummy(t, "first")
	createDummy(t, "second")

	resp, err := http.Get(ts.URL + "/api/v1/dummy")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

{{- if index .Modules "postgres"}}
func TestListDummies_Empty(t *testing.T) {
	truncateDummy(t)

	resp, err := http.Get(ts.URL + "/api/v1/dummy")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestListDummies_FilterByID(t *testing.T) {
	truncateDummy(t)

	created := createDummy(t, "find_me")
	createDummy(t, "ignore_me")

	createdID := int(created["id"].(float64))

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/dummy?id=%d", ts.URL, createdID))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "find_me", result[0]["name"])
	assert.Equal(t, float64(createdID), result[0]["id"])
}

func TestListDummies_FilterByName(t *testing.T) {
	truncateDummy(t)

	createDummy(t, "alpha_one")
	createDummy(t, "beta_two")
	createDummy(t, "alpha_three")

	resp, err := http.Get(ts.URL + "/api/v1/dummy?name=alpha")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	for _, item := range result {
		assert.Contains(t, item["name"], "alpha")
	}
}

func TestListDummies_Pagination(t *testing.T) {
	truncateDummy(t)

	for i := 1; i <= 5; i++ {
		createDummy(t, fmt.Sprintf("item_%d", i))
	}

	// First page: limit=2
	resp, err := http.Get(ts.URL + "/api/v1/dummy?limit=2")
	require.NoError(t, err)
	defer resp.Body.Close()

	var page1 []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&page1)
	require.NoError(t, err)
	assert.Len(t, page1, 2)

	// Second page: limit=2, offset=2
	resp2, err := http.Get(ts.URL + "/api/v1/dummy?limit=2&offset=2")
	require.NoError(t, err)
	defer resp2.Body.Close()

	var page2 []map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&page2)
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	// Pages should contain different items
	assert.NotEqual(t, page1[0]["id"], page2[0]["id"])
	assert.NotEqual(t, page1[1]["id"], page2[1]["id"])

	// Third page: limit=2, offset=4 — should return 1 remaining item
	resp3, err := http.Get(ts.URL + "/api/v1/dummy?limit=2&offset=4")
	require.NoError(t, err)
	defer resp3.Body.Close()

	var page3 []map[string]interface{}
	err = json.NewDecoder(resp3.Body).Decode(&page3)
	require.NoError(t, err)
	assert.Len(t, page3, 1)
}
{{- end}}
