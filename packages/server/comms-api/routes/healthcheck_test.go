package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var healthRouter *gin.Engine

func init() {
	healthRouter = gin.Default()
	route := healthRouter.Group("/")

	addHealthRoutes(route)
}

type HealthResponse struct {
	Status string `json:"status,required"`
}

func TestHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	healthRouter.ServeHTTP(w, req)
	if !assert.Equal(t, w.Code, 200) {
		return
	}
	var resp HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("TestHealthCheck, unable to decode json: %v\n", err.Error())
		return
	}
	assert.Equal(t, resp.Status, "ok")
}
