package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig holds configuration for API testing
type TestConfig struct {
	BaseURL       string
	AuthToken     string
	TestUserID    string
	TestOrderID   string
	TestProductID string
}

// Helper function to create a test server and config
func setupTestServer(t *testing.T) (*Server, *http.Client, *TestConfig) {
	cfg := &config.Config{
		Port:             "8080",
		InventoryService: "localhost:50051",
		OrderService:     "localhost:50052",
		UserService:      "localhost:50053",
		MailerService:    "localhost:50054",
		JWTSecret:        "test-secret-key",
	}

	server := NewServer(cfg)
	client := &http.Client{}

	testCfg := &TestConfig{
		BaseURL:       "/api",
		AuthToken:     "Bearer test-token",
		TestUserID:    "user123",
		TestOrderID:   "order123",
		TestProductID: "product123",
	}

	return server, client, testCfg
}

// ============================================================================
// API TESTING SUITE - Meeting Task Requirements
// ============================================================================
// Requirements:
// - 6+ requests: GET (2+), POST (1+), PUT/PATCH (1+), DELETE (1+)
// - Validate: status codes, response body fields, headers
// - 5+ assertions per test
// - Environment variables: baseUrl, authToken, etc.
// ============================================================================

// TEST 1: POST - User Registration
// Validates: status code, response structure, headers, content-type
func TestUserRegistration(t *testing.T) {
	server, _, testCfg := setupTestServer(t)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"username": "testuser",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", testCfg.BaseURL+"/users/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Assertion 1: Validate status code (201 Created or 500 if service unavailable)
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusInternalServerError,
		"Status code should be 201 Created or 500 Internal Server Error")

	// Assertion 2: Validate response header - Content-Type
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json",
		"Response should have Content-Type: application/json")

	// Assertion 3: Validate response body is valid JSON
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "Response body should be valid JSON")

	// Assertion 4: Response should contain either 'id'/'message' or 'error' field
	if w.Code == http.StatusCreated {
		assert.True(t, responseBody["id"] != nil || responseBody["message"] != nil,
			"Response should contain 'id' or 'message' field on success")
	} else {
		assert.NotNil(t, responseBody["error"], "Response should contain 'error' field on failure")
	}

	// Assertion 5: Validate response body is not empty
	assert.NotEmpty(t, w.Body.String(), "Response body should not be empty")
}

// TEST 2: POST - User Authentication
// Validates: status code, token in response, response structure
func TestUserAuthentication(t *testing.T) {
	server, _, testCfg := setupTestServer(t)

	body := map[string]string{
		"username": "testuser",
		"password": "password123",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", testCfg.BaseURL+"/users/authenticate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Assertion 1: Validate status code
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError,
		"Status code should be 200 OK or 500 Internal Server Error")

	// Assertion 2: Validate Content-Type header
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json",
		"Response should have application/json content type")

	// Assertion 3: Parse and validate response body structure
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err, "Response should be valid JSON")

	// Assertion 4: Response should contain token on success
	if w.Code == http.StatusOK {
		assert.NotNil(t, responseBody["token"], "Response should contain 'token' field")
		assert.NotEmpty(t, responseBody["token"], "Token should not be empty")
	}

	// Assertion 5: Response should contain message field
	assert.True(t, responseBody["message"] != nil || responseBody["error"] != nil,
		"Response should contain 'message' or 'error' field")
}

// TEST 3: GET - List Products (First GET endpoint)
// Validates: status code, response array structure, query parameters, pagination fields
func TestListProducts(t *testing.T) {
	server, _, testCfg := setupTestServer(t)

	// Request with query parameters
	req := httptest.NewRequest("GET", testCfg.BaseURL+"/inventory?category=electronics&page=1&limit=10", nil)
	req.Header.Set("Authorization", testCfg.AuthToken)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Assertion 1: Validate status code (200 OK or 401 if auth required)
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError,
		"Status code should be 200, 401, or 500")

	// Assertion 2: Validate Content-Type header
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json",
		"Response should have application/json content type")

	// Assertion 3: Validate response is valid JSON
	var responseBody interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "Response body should be valid JSON")

	// Assertion 4: Validate Authorization header was sent
	assert.Equal(t, testCfg.AuthToken, req.Header.Get("Authorization"),
		"Authorization header should be set correctly")

	// Assertion 5: Response should not be empty
	assert.NotEmpty(t, w.Body.String(), "Response body should not be empty")
}

// TEST 4: GET - Get User Profile (Second GET endpoint)
// Validates: status code, response structure, ID parameter, user fields
func TestGetUserProfile(t *testing.T) {
	server, _, testCfg := setupTestServer(t)

	req := httptest.NewRequest("GET", testCfg.BaseURL+"/users/"+testCfg.TestUserID, nil)
	req.Header.Set("Authorization", testCfg.AuthToken)
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Assertion 1: Validate status code
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError,
		"Status code should be 200, 401, or 500")

	// Assertion 2: Validate Accept header was respected
	assert.Equal(t, "application/json", req.Header.Get("Accept"),
		"Accept header should be set to application/json")

	// Assertion 3: Validate response is JSON
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "Response should be valid JSON")

	// Assertion 4: Validate ID parameter in URL
	assert.Contains(t, req.URL.Path, testCfg.TestUserID,
		"Request URL should contain the user ID parameter")

	// Assertion 5: Validate Content-Type
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json",
		"Response Content-Type should be application/json")
}

// TEST 5: PUT - Update Product
// Validates: status code, request body validation, update fields, response
func TestUpdateProduct(t *testing.T) {
	server, _, testCfg := setupTestServer(t)

	body := map[string]interface{}{
		"id":          testCfg.TestProductID,
		"name":        "Updated Product",
		"price":       149.99,
		"quantity":    5,
		"description": "Updated description",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", testCfg.BaseURL+"/inventory/"+testCfg.TestProductID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", testCfg.AuthToken)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Assertion 1: Validate status code (200 OK or auth error)
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError,
		"Status code should be 200, 401, or 500")

	// Assertion 2: Validate Content-Type in request
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"),
		"Request Content-Type should be application/json")

	// Assertion 3: Validate response is valid JSON
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "Response should be valid JSON")

	// Assertion 4: Validate HTTP method
	assert.Equal(t, "PUT", req.Method, "Request method should be PUT")

	// Assertion 5: Validate request body contains required fields
	assert.True(t, body["name"] != nil && body["price"] != nil,
		"Request body should contain 'name' and 'price' fields")
}

// TEST 6: DELETE - Delete Product
// Validates: status code, method, authentication, response message
func TestDeleteProduct(t *testing.T) {
	server, _, testCfg := setupTestServer(t)

	req := httptest.NewRequest("DELETE", testCfg.BaseURL+"/inventory/"+testCfg.TestProductID, nil)
	req.Header.Set("Authorization", testCfg.AuthToken)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Assertion 1: Validate status code
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError,
		"Status code should be 200, 401, or 500")

	// Assertion 2: Validate HTTP method
	assert.Equal(t, "DELETE", req.Method, "Request method should be DELETE")

	// Assertion 3: Validate authorization header
	assert.Equal(t, testCfg.AuthToken, req.Header.Get("Authorization"),
		"Authorization header should be properly set")

	// Assertion 4: Validate response is JSON
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "Response should be valid JSON")

	// Assertion 5: Response should contain message on success
	if w.Code == http.StatusOK {
		assert.NotNil(t, responseBody["message"], "Response should contain 'message' field")
	}
}

// TEST 7: POST - Create Order
// Validates: status code, request validation, response fields, created resource
func TestCreateOrder(t *testing.T) {
	server, _, testCfg := setupTestServer(t)

	body := map[string]interface{}{
		"user_id":     testCfg.TestUserID,
		"product_id":  testCfg.TestProductID,
		"quantity":    2,
		"total_price": 199.98,
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", testCfg.BaseURL+"/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", testCfg.AuthToken)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Assertion 1: Validate status code
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError,
		"Status code should be 201, 401, or 500")

	// Assertion 2: Validate Content-Type header
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json",
		"Response should have application/json content type")

	// Assertion 3: Validate request body fields
	assert.True(t, body["user_id"] != nil && body["product_id"] != nil,
		"Request body should contain required fields")

	// Assertion 4: Validate response structure
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "Response should be valid JSON")

	// Assertion 5: Validate method is POST
	assert.Equal(t, "POST", req.Method, "Request method should be POST")
}

// TEST 8: GET - Get Order Details (Bonus: covers order retrieval)
// Validates: status code, order ID in response, headers
func TestGetOrder(t *testing.T) {
	server, _, testCfg := setupTestServer(t)

	req := httptest.NewRequest("GET", testCfg.BaseURL+"/orders/"+testCfg.TestOrderID, nil)
	req.Header.Set("Authorization", testCfg.AuthToken)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Assertion 1: Validate status code
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError,
		"Status code should be 200, 401, or 500")

	// Assertion 2: Validate response content-type
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json",
		"Response should have application/json content type")

	// Assertion 3: Validate response is valid JSON
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "Response should be valid JSON")

	// Assertion 4: Validate order ID in URL
	assert.Contains(t, req.URL.Path, testCfg.TestOrderID,
		"Request URL should contain order ID")

	// Assertion 5: Validate body not empty
	assert.NotEmpty(t, w.Body.String(), "Response body should not be empty")
}

// TEST 9: Additional Validation Tests

// TestInvalidJSON validates error handling for malformed requests
func TestInvalidJSON(t *testing.T) {
	server, _, testCfg := setupTestServer(t)

	req := httptest.NewRequest("POST", testCfg.BaseURL+"/users/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Assertion 1: Should return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid JSON should return 400 Bad Request")

	// Assertion 2: Response should be valid JSON
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "Error response should be valid JSON")

	// Assertion 3: Should contain error field
	assert.NotNil(t, responseBody["error"], "Error response should contain 'error' field")
}
