package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Create a mock for the healthCheckHandler
type MockHealthCheckHandler struct {
	mock.Mock
}

func (m *MockHealthCheckHandler) Handle(c *gin.Context) {
	m.Called(c)
}

func TestHealthCheckHandler2(t *testing.T) {
	// Create a new Gin router and replace the actual healthCheckHandler with the mock
	router := gin.Default()
	mockHandler := new(MockHealthCheckHandler)
	router.GET("/health", mockHandler.Handle)

	// Create a request and recorder for testing
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	// Serve the request to the router
	router.ServeHTTP(w, req)

	// Assert that the mock handler's Handle function was called
	mockHandler.AssertExpectations(t)

	// Assert the response status code is 200 (OK)
	assert.Equal(t, http.StatusOK, w.Code)

	// You can further assert the response body or headers if needed
	// For example, assert the response body contains JSON with status "OK"
	assert.Contains(t, w.Body.String(), `"status": "OK"`)
}
