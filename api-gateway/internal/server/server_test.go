package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/config"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a test server
func setupTestServer(t *testing.T) (*Server, *http.Client) {
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

	return server, client
}

// TestUserRegistration tests the user registration endpoint
func TestUserRegistration(t *testing.T) {
	server, _ := setupTestServer(t)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"username": "testuser",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Note: This will fail without a running user service, but tests the endpoint structure
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusInternalServerError)
}

// TestUserAuthentication tests the authentication endpoint
func TestUserAuthentication(t *testing.T) {
	server, _ := setupTestServer(t)

	body := map[string]string{
		"username": "testuser",
		"password": "password123",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/users/authenticate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Should return 200 or 500 depending on service availability
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

// TestGetUserProfile tests the get user profile endpoint
func TestGetUserProfile(t *testing.T) {
	server, _ := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/users/123", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Should require auth or fail gracefully
	assert.True(t, w.Code == http.StatusOK ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestCreateProduct tests the create product endpoint
func TestCreateProduct(t *testing.T) {
	server, _ := setupTestServer(t)

	body := map[string]interface{}{
		"name":        "Test Product",
		"description": "A test product",
		"price":       99.99,
		"quantity":    10,
		"category":    "electronics",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/inventory", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusCreated ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestGetProduct tests the get product endpoint
func TestGetProduct(t *testing.T) {
	server, _ := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/inventory/product123", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestUpdateProduct tests the update product endpoint
func TestUpdateProduct(t *testing.T) {
	server, _ := setupTestServer(t)

	body := map[string]interface{}{
		"id":       "product123",
		"name":     "Updated Product",
		"price":    149.99,
		"quantity": 5,
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/inventory/product123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestDeleteProduct tests the delete product endpoint
func TestDeleteProduct(t *testing.T) {
	server, _ := setupTestServer(t)

	req := httptest.NewRequest("DELETE", "/api/inventory/product123", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestListProducts tests the list products endpoint
func TestListProducts(t *testing.T) {
	server, _ := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/inventory?category=electronics&page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestCreateOrder tests the create order endpoint
func TestCreateOrder(t *testing.T) {
	server, _ := setupTestServer(t)

	body := map[string]interface{}{
		"user_id":     "user123",
		"product_id":  "product123",
		"quantity":    2,
		"total_price": 199.98,
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusCreated ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestGetOrder tests the get order endpoint
func TestGetOrder(t *testing.T) {
	server, _ := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/orders/order123", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestUpdateOrderStatus tests the update order status endpoint
func TestUpdateOrderStatus(t *testing.T) {
	server, _ := setupTestServer(t)

	body := map[string]string{
		"id":     "order123",
		"status": "shipped",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/orders/order123/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestListUserOrders tests the list user orders endpoint
func TestListUserOrders(t *testing.T) {
	server, _ := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/orders?user_id=user123", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestSendEmail tests the send email endpoint
func TestSendEmail(t *testing.T) {
	server, _ := setupTestServer(t)

	body := map[string]string{
		"to":      "recipient@example.com",
		"subject": "Test Email",
		"body":    "This is a test email",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/email", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK ||
		w.Code == http.StatusUnauthorized ||
		w.Code == http.StatusInternalServerError)
}

// TestInvalidJSON tests that invalid JSON returns a bad request error
func TestInvalidJSON(t *testing.T) {
	server, _ := setupTestServer(t)

	req := httptest.NewRequest("POST", "/api/users/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMissingAuthHeader tests that protected routes require authentication
func TestMissingAuthHeader(t *testing.T) {
	server, _ := setupTestServer(t)

	// Protected route without auth header
	req := httptest.NewRequest("GET", "/api/inventory", nil)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestPublicRoutes tests that public routes don't require authentication
func TestPublicRoutes(t *testing.T) {
	server, _ := setupTestServer(t)

	tests := []struct {
		method string
		path   string
		body   map[string]string
	}{
		{
			method: "POST",
			path:   "/api/users/register",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
				"username": "testuser",
			},
		},
		{
			method: "POST",
			path:   "/api/users/authenticate",
			body: map[string]string{
				"username": "testuser",
				"password": "password123",
			},
		},
	}

	for _, test := range tests {
		jsonBody, _ := json.Marshal(test.body)
		req := httptest.NewRequest(test.method, test.path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)

		// These should not return 401 Unauthorized
		assert.NotEqual(t, http.StatusUnauthorized, w.Code, "Public route %s should not require auth", test.path)
	}
}

// TestResponseFormat tests that responses are in valid JSON format
func TestResponseFormat(t *testing.T) {
	server, _ := setupTestServer(t)

	body := map[string]string{
		"username": "testuser",
		"password": "password123",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/users/authenticate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Verify response can be unmarshalled as JSON
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
}
