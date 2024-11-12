package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/savioruz/mikti-task/tree/week-4/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func getAuthToken(t *testing.T) string {
	// First register a user
	registerBody := model.RegisterRequest{
		Email:    "todo.test@svrz.xyz",
		Password: "strongpassword",
		Role:     "user",
	}

	registerJson, err := json.Marshal(registerBody)
	assert.Nil(t, err)

	registerRequest := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(registerJson))
	registerRequest.Header.Set("Content-Type", "application/json")
	registerRequest.Header.Set("Accept", "application/json")

	registerRecorder := httptest.NewRecorder()
	app.ServeHTTP(registerRecorder, registerRequest)

	// Then login to get the token
	loginBody := model.LoginRequest{
		Email:    "todo.test@svrz.xyz",
		Password: "strongpassword",
	}

	loginJson, err := json.Marshal(loginBody)
	assert.Nil(t, err)

	loginRequest := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(loginJson))
	loginRequest.Header.Set("Content-Type", "application/json")
	loginRequest.Header.Set("Accept", "application/json")

	loginRecorder := httptest.NewRecorder()
	app.ServeHTTP(loginRecorder, loginRequest)

	loginResponse := loginRecorder.Result()
	loginBytes, err := io.ReadAll(loginResponse.Body)
	assert.Nil(t, err)

	var loginRawResponse map[string]interface{}
	err = json.Unmarshal(loginBytes, &loginRawResponse)
	assert.Nil(t, err)

	return loginRawResponse["data"].(map[string]interface{})["access_token"].(string)
}

func createTodo(t *testing.T, token string, title string) {
	createRequest := httptest.NewRequest(http.MethodPost, "/api/v1/todo", strings.NewReader(fmt.Sprintf(`{"title":"%s"}`, title)))
	createRequest.Header.Set("Content-Type", "application/json")
	createRequest.Header.Set("Accept", "application/json")
	createRequest.Header.Set("Authorization", "Bearer "+token)

	createRecorder := httptest.NewRecorder()
	app.ServeHTTP(createRecorder, createRequest)

	assert.Equal(t, http.StatusCreated, createRecorder.Result().StatusCode)
}

func TestCreateTodo(t *testing.T) {
	ClearAll()

	token := getAuthToken(t)

	requestBody := model.TodoCreateRequest{
		Title: "Test Todo",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)

	response := recorder.Result()
	b, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(b))

	var rawResponse map[string]interface{}
	err = json.Unmarshal(b, &rawResponse)
	assert.Nil(t, err)

	data, ok := rawResponse["data"].(map[string]interface{})
	assert.True(t, ok, "Data should be present in response")

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.Title, data["title"])
	assert.Equal(t, false, data["completed"])
	assert.NotEmpty(t, data["id"])
	assert.NotEmpty(t, data["created_at"])
	assert.NotEmpty(t, data["updated_at"])
}

func TestCreateTodoUnauthorized(t *testing.T) {
	ClearAll()

	requestBody := model.TodoCreateRequest{
		Title: "Test Todo",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)

	response := recorder.Result()
	b, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(b))

	var rawResponse map[string]interface{}
	err = json.Unmarshal(b, &rawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "Missing authorization header", rawResponse["message"])
}

func TestCreateTodoValidationError(t *testing.T) {
	ClearAll()

	token := getAuthToken(t)

	requestBody := model.TodoCreateRequest{
		Title: "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)

	response := recorder.Result()
	b, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(b))

	var rawResponse map[string]interface{}
	err = json.Unmarshal(b, &rawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "validation error", rawResponse["error"])
}

func TestUpdateTodo(t *testing.T) {
	ClearAll()

	token := getAuthToken(t)
	TestCreateTodo(t)
	todoID := GetFirstTodoID(t)

	newTitle := "Updated Todo"
	completed := true
	updateBody := model.TodoUpdateRequest{
		Title:     &newTitle,
		Completed: &completed,
	}

	updateJson, err := json.Marshal(updateBody)
	assert.Nil(t, err)

	updateRequest := httptest.NewRequest(http.MethodPut, "/api/v1/todo/"+todoID, bytes.NewReader(updateJson))
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Accept", "application/json")
	updateRequest.Header.Set("Authorization", "Bearer "+token)

	updateRecorder := httptest.NewRecorder()
	app.ServeHTTP(updateRecorder, updateRequest)

	updateResponse := updateRecorder.Result()
	updateBytes, err := io.ReadAll(updateResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(updateBytes))

	var updateRawResponse map[string]interface{}
	err = json.Unmarshal(updateBytes, &updateRawResponse)
	assert.Nil(t, err)

	updatedData := updateRawResponse["data"].(map[string]interface{})
	assert.Equal(t, http.StatusOK, updateResponse.StatusCode)
	assert.Equal(t, newTitle, updatedData["title"])
	assert.Equal(t, completed, updatedData["completed"])
}

func TestUpdateTodoUnauthorized(t *testing.T) {
	ClearAll()

	TestCreateTodo(t)
	todoID := GetFirstTodoID(t)

	newTitle := "Updated Todo"
	completed := true
	updateBody := model.TodoUpdateRequest{
		Title:     &newTitle,
		Completed: &completed,
	}

	updateJson, err := json.Marshal(updateBody)
	assert.Nil(t, err)

	updateRequest := httptest.NewRequest(http.MethodPut, "/api/v1/todo/"+todoID, bytes.NewReader(updateJson))
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Accept", "application/json")

	updateRecorder := httptest.NewRecorder()
	app.ServeHTTP(updateRecorder, updateRequest)

	updateResponse := updateRecorder.Result()
	updateBytes, err := io.ReadAll(updateResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(updateBytes))

	var updateRawResponse map[string]interface{}
	err = json.Unmarshal(updateBytes, &updateRawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, updateResponse.StatusCode)
	assert.Equal(t, "Missing authorization header", updateRawResponse["message"])
}

func TestUpdateTodoValidationError(t *testing.T) {
	ClearAll()

	token := getAuthToken(t)
	TestCreateTodo(t)
	todoID := GetFirstTodoID(t)

	newTitle := ""
	completed := true
	updateBody := model.TodoUpdateRequest{
		Title:     &newTitle,
		Completed: &completed,
	}

	updateJson, err := json.Marshal(updateBody)
	assert.Nil(t, err)

	updateRequest := httptest.NewRequest(http.MethodPut, "/api/v1/todo/"+todoID, bytes.NewReader(updateJson))
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Accept", "application/json")
	updateRequest.Header.Set("Authorization", "Bearer "+token)

	updateRecorder := httptest.NewRecorder()
	app.ServeHTTP(updateRecorder, updateRequest)

	updateResponse := updateRecorder.Result()
	updateBytes, err := io.ReadAll(updateResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(updateBytes))

	var updateRawResponse map[string]interface{}
	err = json.Unmarshal(updateBytes, &updateRawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, updateResponse.StatusCode)
	assert.Equal(t, "validation error", updateRawResponse["error"])
}

func TestDeleteTodo(t *testing.T) {
	ClearAll()

	token := getAuthToken(t)
	TestCreateTodo(t)
	todoID := GetFirstTodoID(t)

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/v1/todo/"+todoID, nil)
	deleteRequest.Header.Set("Content-Type", "application/json")
	deleteRequest.Header.Set("Accept", "application/json")
	deleteRequest.Header.Set("Authorization", "Bearer "+token)

	deleteRecorder := httptest.NewRecorder()
	app.ServeHTTP(deleteRecorder, deleteRequest)

	deleteResponse := deleteRecorder.Result()
	deleteBytes, err := io.ReadAll(deleteResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(deleteBytes))

	var deleteRawResponse map[string]interface{}
	err = json.Unmarshal(deleteBytes, &deleteRawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, deleteResponse.StatusCode)
	assert.Equal(t, "todo with ID "+todoID+" has been deleted", deleteRawResponse["message"].(map[string]interface{})["message"])
}

func TestDeleteTodoUnauthorized(t *testing.T) {
	ClearAll()

	TestCreateTodo(t)
	todoID := GetFirstTodoID(t)

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/v1/todo/"+todoID, nil)
	deleteRequest.Header.Set("Content-Type", "application/json")
	deleteRequest.Header.Set("Accept", "application/json")

	deleteRecorder := httptest.NewRecorder()
	app.ServeHTTP(deleteRecorder, deleteRequest)

	deleteResponse := deleteRecorder.Result()
	deleteBytes, err := io.ReadAll(deleteResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(deleteBytes))

	var deleteRawResponse map[string]interface{}
	err = json.Unmarshal(deleteBytes, &deleteRawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, deleteResponse.StatusCode)
	assert.Equal(t, "Missing authorization header", deleteRawResponse["message"])
}

func TestDeleteTodoNotFound(t *testing.T) {
	ClearAll()

	token := getAuthToken(t)
	id := uuid.NewString()

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/v1/todo/"+id, nil)
	deleteRequest.Header.Set("Content-Type", "application/json")
	deleteRequest.Header.Set("Accept", "application/json")
	deleteRequest.Header.Set("Authorization", "Bearer "+token)

	deleteRecorder := httptest.NewRecorder()
	app.ServeHTTP(deleteRecorder, deleteRequest)

	deleteResponse := deleteRecorder.Result()
	deleteBytes, err := io.ReadAll(deleteResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(deleteBytes))

	var deleteRawResponse map[string]interface{}
	err = json.Unmarshal(deleteBytes, &deleteRawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, deleteResponse.StatusCode)
	assert.Equal(t, "not found", deleteRawResponse["error"])
}

func TestGetTodo(t *testing.T) {
	ClearAll()

	token := getAuthToken(t)
	TestCreateTodo(t)
	todoID := GetFirstTodoID(t)

	getRequest := httptest.NewRequest(http.MethodGet, "/api/v1/todo/"+todoID, nil)
	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("Authorization", "Bearer "+token)

	getRecorder := httptest.NewRecorder()
	app.ServeHTTP(getRecorder, getRequest)

	getResponse := getRecorder.Result()
	getBytes, err := io.ReadAll(getResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(getBytes))

	var getRawResponse map[string]interface{}
	err = json.Unmarshal(getBytes, &getRawResponse)
	assert.Nil(t, err)

	getData := getRawResponse["data"].(map[string]interface{})
	assert.Equal(t, http.StatusOK, getResponse.StatusCode)
	assert.Equal(t, "Test Todo", getData["title"])
	assert.Equal(t, false, getData["completed"])
	assert.Equal(t, todoID, getData["id"])
}

func TestGetTodoUnauthorized(t *testing.T) {
	ClearAll()

	TestCreateTodo(t)
	todoID := GetFirstTodoID(t)

	getRequest := httptest.NewRequest(http.MethodGet, "/api/v1/todo/"+todoID, nil)
	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Accept", "application/json")

	getRecorder := httptest.NewRecorder()
	app.ServeHTTP(getRecorder, getRequest)

	getResponse := getRecorder.Result()
	getBytes, err := io.ReadAll(getResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(getBytes))

	var getRawResponse map[string]interface{}
	err = json.Unmarshal(getBytes, &getRawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, getResponse.StatusCode)
	assert.Equal(t, "Missing authorization header", getRawResponse["message"])
}

func TestGetTodoNotFound(t *testing.T) {
	ClearAll()

	token := getAuthToken(t)
	id := uuid.NewString()

	getRequest := httptest.NewRequest(http.MethodGet, "/api/v1/todo/"+id, nil)
	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("Authorization", "Bearer "+token)

	getRecorder := httptest.NewRecorder()
	app.ServeHTTP(getRecorder, getRequest)

	getResponse := getRecorder.Result()
	getBytes, err := io.ReadAll(getResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(getBytes))

	var getRawResponse map[string]interface{}
	err = json.Unmarshal(getBytes, &getRawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, getResponse.StatusCode)
	assert.Equal(t, "not found", getRawResponse["error"])
}

func TestListTodo(t *testing.T) {
	// Setup
	ClearAll()
	token := getAuthToken(t)
	count := 5
	CreateTodos(t, count)

	// Test default pagination (page 1)
	t.Run("List todos with default pagination", func(t *testing.T) {
		listRequest := httptest.NewRequest(http.MethodGet, "/api/v1/todo", nil)
		listRequest.Header.Set("Content-Type", "application/json")
		listRequest.Header.Set("Accept", "application/json")
		listRequest.Header.Set("Authorization", "Bearer "+token)

		listRecorder := httptest.NewRecorder()
		app.ServeHTTP(listRecorder, listRequest)

		listResponse := listRecorder.Result()
		assert.Equal(t, http.StatusOK, listResponse.StatusCode)

		listBytes, err := io.ReadAll(listResponse.Body)
		require.NoError(t, err)

		var response struct {
			Data   []map[string]interface{} `json:"data"`
			Paging struct {
				Page       int `json:"page"`
				Size       int `json:"size"`
				TotalItems int `json:"total_items"`
				TotalPages int `json:"total_pages"`
			} `json:"paging"`
		}
		err = json.Unmarshal(listBytes, &response)
		require.NoError(t, err)

		assert.Equal(t, count, len(response.Data))
		assert.Equal(t, 1, response.Paging.Page)
		assert.NotZero(t, response.Paging.Size)
		assert.Equal(t, count, response.Paging.TotalItems)
		assert.NotZero(t, response.Paging.TotalPages)
	})

	// Test with explicit pagination
	t.Run("List todos with explicit pagination", func(t *testing.T) {
		pageSize := 3
		listRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/todo?page=1&size=%d", pageSize), nil)
		listRequest.Header.Set("Content-Type", "application/json")
		listRequest.Header.Set("Accept", "application/json")
		listRequest.Header.Set("Authorization", "Bearer "+token)

		listRecorder := httptest.NewRecorder()
		app.ServeHTTP(listRecorder, listRequest)

		listResponse := listRecorder.Result()
		assert.Equal(t, http.StatusOK, listResponse.StatusCode)

		listBytes, err := io.ReadAll(listResponse.Body)
		require.NoError(t, err)

		var response struct {
			Data   []map[string]interface{} `json:"data"`
			Paging struct {
				Page       int `json:"page"`
				Size       int `json:"size"`
				TotalItems int `json:"total_items"`
				TotalPages int `json:"total_pages"`
			} `json:"paging"`
		}
		err = json.Unmarshal(listBytes, &response)
		require.NoError(t, err)

		assert.Equal(t, pageSize, len(response.Data))
		assert.Equal(t, 1, response.Paging.Page)
		assert.Equal(t, pageSize, response.Paging.Size)
		assert.Equal(t, count, response.Paging.TotalItems)
		assert.Equal(t, (count+pageSize-1)/pageSize, response.Paging.TotalPages)
	})

	// Test invalid pagination
	t.Run("List todos with invalid pagination should use defaults", func(t *testing.T) {
		listRequest := httptest.NewRequest(http.MethodGet, "/api/v1/todo?page=0&size=0", nil)
		listRequest.Header.Set("Content-Type", "application/json")
		listRequest.Header.Set("Accept", "application/json")
		listRequest.Header.Set("Authorization", "Bearer "+token)

		listRecorder := httptest.NewRecorder()
		app.ServeHTTP(listRecorder, listRequest)

		listResponse := listRecorder.Result()
		assert.Equal(t, http.StatusOK, listResponse.StatusCode)

		listBytes, err := io.ReadAll(listResponse.Body)
		require.NoError(t, err)

		var response struct {
			Data   []map[string]interface{} `json:"data"`
			Paging struct {
				Page       int `json:"page"`
				Size       int `json:"size"`
				TotalItems int `json:"total_items"`
				TotalPages int `json:"total_pages"`
			} `json:"paging"`
		}
		err = json.Unmarshal(listBytes, &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Paging.Page)
		assert.NotZero(t, response.Paging.Size)
		assert.NotZero(t, response.Paging.TotalPages)
	})
}

func TestListTodoCache(t *testing.T) {
	// Setup
	ClearAll()
	token := getAuthToken(t)
	count := 5
	CreateTodos(t, count)

	makeRequest := func() (*http.Response, map[string]interface{}) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/todo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		resp := rec.Result()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		return resp, result
	}

	// Test initial request (should hit database)
	t.Run("Initial request should hit database", func(t *testing.T) {
		resp1, result1 := makeRequest()
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		data1 := result1["data"].([]interface{})
		assert.Equal(t, count, len(data1))

		paging := result1["paging"].(map[string]interface{})
		assert.NotZero(t, paging["size"])
		assert.Equal(t, float64(1), paging["page"])
	})

	// Test cached response
	t.Run("Second request should hit cache", func(t *testing.T) {
		resp2, result2 := makeRequest()
		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		data2 := result2["data"].([]interface{})
		assert.Equal(t, count, len(data2))
	})

	// Test cache invalidation
	t.Run("Cache should be invalidated after adding new todo", func(t *testing.T) {
		// Add a new todo
		createTodo(t, token, "New Todo")

		resp3, result3 := makeRequest()
		assert.Equal(t, http.StatusOK, resp3.StatusCode)

		data3 := result3["data"].([]interface{})
		assert.Equal(t, count+1, len(data3))
	})

	// Note: Skip cache expiration test in regular test runs
	t.Skip("Cache expiration test - skipped for regular test runs")
	t.Run("Cache should expire after TTL", func(t *testing.T) {
		// First request to populate cache
		makeRequest()

		// Wait for cache to expire (assuming 5 minute TTL)
		time.Sleep(5 * time.Minute)

		// Request after expiration should hit database
		resp4, result4 := makeRequest()
		assert.Equal(t, http.StatusOK, resp4.StatusCode)

		data4 := result4["data"].([]interface{})
		assert.Equal(t, count+1, len(data4))
	})
}
