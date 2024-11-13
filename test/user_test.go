package test

import (
	"bytes"
	"encoding/json"
	"github.com/savioruz/mikti-task/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterRequest{
		Email:    "user@svrz.xyz",
		Password: "strongpassword",
		Role:     "user",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(bodyJson)))
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

	data, ok := rawResponse["data"].(map[string]interface{})
	assert.True(t, ok, "Data should be present in response")

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.Email, data["email"])
	assert.Equal(t, requestBody.Role, data["role"])
	assert.NotEmpty(t, data["id"])
	assert.NotEmpty(t, data["created_at"])
	assert.NotEmpty(t, data["updated_at"])
}

func TestRegisterValidationError(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterRequest{
		Email:    "",
		Password: "",
		Role:     "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(bodyJson)))
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

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "validation error", rawResponse["error"].(map[string]interface{})["message"])
	assert.NotNil(t, rawResponse["error"])
}

func TestRegisterDuplicateEmail(t *testing.T) {
	ClearAll()
	TestRegisterUser(t)

	requestBody := model.RegisterRequest{
		Email:    "user@svrz.xyz",
		Password: "strongpassword",
		Role:     "user",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(bodyJson)))
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

	assert.Equal(t, http.StatusConflict, response.StatusCode)
	assert.NotNil(t, rawResponse["error"])
}

func TestRegisterDuplicateAdminRole(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterRequest{
		Email:    "user@svrz.xyz",
		Password: "strongpassword",
		Role:     "admin",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	// First request
	firstRequest := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyJson))
	firstRequest.Header.Set("Content-Type", "application/json")
	firstRequest.Header.Set("Accept", "application/json")

	firstRecorder := httptest.NewRecorder()
	app.ServeHTTP(firstRecorder, firstRequest)

	// Second request with same role
	secondRequest := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyJson))
	secondRequest.Header.Set("Content-Type", "application/json")
	secondRequest.Header.Set("Accept", "application/json")

	secondRecorder := httptest.NewRecorder()
	app.ServeHTTP(secondRecorder, secondRequest)
	secondResponse := secondRecorder.Result()

	// Read response body
	secondBytes, err := io.ReadAll(secondResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(secondBytes))

	var secondRawResponse map[string]interface{}
	err = json.Unmarshal(secondBytes, &secondRawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, secondResponse.StatusCode)
	assert.NotNil(t, secondRawResponse["error"])
}

func TestLoginUser(t *testing.T) {
	ClearAll()

	// Register user
	registerBody := model.RegisterRequest{
		Email:    "user@svrz.xyz",
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

	// Login user
	loginBody := model.LoginRequest{
		Email:    "user@svrz.xyz",
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

	t.Logf("Response Body: %s", string(loginBytes))

	var loginRawResponse map[string]interface{}
	err = json.Unmarshal(loginBytes, &loginRawResponse)

	assert.Equal(t, http.StatusOK, loginResponse.StatusCode)
	assert.NotNil(t, loginRawResponse["data"])
}

func TestLoginValidationError(t *testing.T) {
	ClearAll()

	requestBody := model.LoginRequest{
		Email:    "",
		Password: "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(bodyJson))
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

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "validation error", rawResponse["error"].(map[string]interface{})["message"])
	assert.NotNil(t, rawResponse["error"])
}

func TestLoginWrongEmail(t *testing.T) {
	ClearAll()
	TestRegisterUser(t)

	requestBody := model.LoginRequest{
		Email:    "wrong@svrz.xyz",
		Password: "strongpassword",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(bodyJson))
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
	assert.Equal(t, "unauthorized", rawResponse["error"].(map[string]interface{})["message"])
	assert.NotNil(t, rawResponse["error"])
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()
	TestRegisterUser(t)

	requestBody := model.LoginRequest{
		Email:    "user@svrz.xyz",
		Password: "wrongpassword",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(bodyJson))
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
	assert.Equal(t, "unauthorized", rawResponse["error"].(map[string]interface{})["message"])
	assert.NotNil(t, rawResponse["error"])
}

func TestRefreshToken(t *testing.T) {
	ClearAll()
	TestRegisterUser(t)

	requestBody := model.LoginRequest{
		Email:    "user@svrz.xyz",
		Password: "strongpassword",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	loginRequest := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(bodyJson))
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

	token := loginRawResponse["data"].(map[string]interface{})["refresh_token"].(string)
	assert.NotEmpty(t, token)

	refreshBody := model.RefreshTokenRequest{
		RefreshToken: token,
	}

	refreshJson, err := json.Marshal(refreshBody)
	assert.Nil(t, err)

	refreshRequest := httptest.NewRequest(http.MethodPost, "/api/v1/users/refresh", bytes.NewReader(refreshJson))
	refreshRequest.Header.Set("Content-Type", "application/json")
	refreshRequest.Header.Set("Accept", "application/json")

	refreshRecorder := httptest.NewRecorder()
	app.ServeHTTP(refreshRecorder, refreshRequest)

	refreshResponse := refreshRecorder.Result()
	refreshBytes, err := io.ReadAll(refreshResponse.Body)
	assert.Nil(t, err)

	t.Logf("Response Body: %s", string(refreshBytes))

	var refreshRawResponse map[string]interface{}
	err = json.Unmarshal(refreshBytes, &refreshRawResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, refreshResponse.StatusCode)
	assert.NotNil(t, refreshRawResponse["data"])
}

func TestRefreshTokenValidationError(t *testing.T) {
	ClearAll()

	requestBody := model.RefreshTokenRequest{
		RefreshToken: "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users/refresh", bytes.NewReader(bodyJson))
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

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "validation error", rawResponse["error"].(map[string]interface{})["message"])
	assert.NotNil(t, rawResponse["error"])
}

func TestRefreshTokenUnauthorized(t *testing.T) {
	ClearAll()

	requestBody := model.RefreshTokenRequest{
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDkxNzcwZmYtYjNjNS00Y2MyLWIzNWItYmJiOTFlYzE2MDMwIiwiZW1haWwiOiJ1c2VyQHN2cnoueHl6Iiwicm9sZSI6InVzZXIiLCJleHAiOjE3MzA5NDU0NDUsImlhdCI6MTczMDM0MDY0NX0.mYM60zSiDbl8JaNR25EZ6nH4F797EwehxnCz9uRwsuk",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users/refresh", bytes.NewReader(bodyJson))
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
	assert.Equal(t, "unauthorized", rawResponse["error"].(map[string]interface{})["message"])
	assert.NotNil(t, rawResponse["error"])
}

func TestRegisterBindError(t *testing.T) {
	ClearAll()

	// Create malformed JSON to trigger Bind() error
	malformedJSON := `{
		"email": "user@svrz.xyz",
		"password": "strongpassword"
		"role": "user" # Missing comma intentionally
	}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(malformedJSON))
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

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "failed to bind request", rawResponse["error"].(map[string]interface{})["message"])
	assert.NotNil(t, rawResponse["error"])
}

func TestLoginBindError(t *testing.T) {
	ClearAll()

	// Create malformed JSON to trigger Bind() error
	malformedJSON := `{
		"email": "user@svrz.xyz"
		"password": "strongpassword"
	}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(malformedJSON))
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

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "failed to bind request", rawResponse["error"].(map[string]interface{})["message"])
}

func TestRefreshTokenBindError(t *testing.T) {
	ClearAll()

	// Create malformed JSON to trigger Bind() error
	malformedJSON := `{
		"refresh_token": aa"
	}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/users/refresh", strings.NewReader(malformedJSON))
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

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "failed to bind request", rawResponse["error"].(map[string]interface{})["message"])
}

// mock db error
/*func TestDbError(t *testing.T) {
	ClearAll()

	registerRequest := model.RegisterRequest{
		Email:    "user@svrz.xyz",
		Password: "strongpassword",
		Role:     "user",
	}

	loginRequest := model.LoginRequest{
		Email:    "user@svrz.xyz",
		Password: "strongpassword",
	}

	_, err := json.Marshal(registerRequest)
	assert.Nil(t, err)

	mockDB := db
	mockDB.Error = errors.New("database error")

	mockUserRepository := &user.UserRepositoryImpl{
		Repository: repositories.Repository[entities.User]{},
		Log:        log,
	}

	u := &usecases.UserUsecase{
		DB:             mockDB,
		Log:            log,
		Validate:       validate,
		UserRepository: mockUserRepository,
		JWTService:     &auth.JWTService{},
	}

	_, err = u.Create(context.Background(), &registerRequest)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, http.StatusInternalServerError)
	assert.Equal(t, "Internal Server Error", err.Error())

	_, err = u.Login(context.Background(), &loginRequest)
	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, http.StatusUnauthorized)
}*/
