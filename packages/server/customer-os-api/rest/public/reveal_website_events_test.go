package public

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func (h *WebsiteTrackerEventsHandler) TestBuildTrackerDbData(t *testing.T) {
	// 1) Prepare a test JSON body. Include some fields from your sample data.
	sampleJSON := `{
        "ip": "127.0.0.1",
        "userId": "user123",
        "eventType": "click",
        "eventData": "{\"buttonId\":\"submit-btn\",\"customKey\":\"customValue\"}",
        "timestamp": 1672972800000,
        "href": "https://example.com/some/page",
        "origin": "https://example.com",
        "search": "?query=value",
        "hostname": "example.com",
        "pathname": "/some/page",
        "referrer": "https://google.com",
        "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
        "language": "en-US",
        "cookiesEnabled": true,
        "screenResolution": "1920x1080"
    }`

	// 2) Create a new Gin context with the JSON body.
	//    We'll use httptest utilities to simulate an HTTP POST.
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(sampleJSON))
	if err != nil {
		t.Fatalf("failed to create test request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// 3) Call the function under test
	tenant := "testTenant"
	trackerData := h.buildTrackerDbData(c, tenant)
	if trackerData == nil {
		t.Fatalf("expected non-nil trackerData, got nil")
	}

	// 4) Assert the fields were decoded properly
	assert.Equal(t, "127.0.0.1", trackerData.IP)
	assert.Equal(t, "click", trackerData.EventType)
	assert.Equal(t, "testTenant", trackerData.Tenant)
	assert.Equal(t, "https://example.com", trackerData.Origin)
	assert.Equal(t, "Mozilla/5.0 (Windows NT 10.0; Win64; x64)", trackerData.UserAgent)
	assert.Equal(t, "https://google.com", trackerData.Referrer)
	assert.Equal(t, "https://example.com/some/page", trackerData.Href)
	assert.Equal(t, "?query=value", trackerData.Search)
	assert.Equal(t, "example.com", trackerData.Hostname)
	assert.Equal(t, "/some/page", trackerData.Pathname)
	assert.Equal(t, int64(1672972800000), trackerData.Timestamp.UnixNano()/int64(time.Millisecond))
}
