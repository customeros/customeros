package response

import (
	"encoding/json"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/gin-gonic/gin"
)

type APIStatus string

const (
	APIStatusError          APIStatus = "error"
	APIStatusPartialSuccess APIStatus = "partial_success"
	APIStatusProcessing     APIStatus = "processing"
	APIStatusSuccess        APIStatus = "success"
	APIStatusWarning        APIStatus = "warning"
)

type Response struct{}

func NewRestResponseHandler() *Response {
	return &Response{}
}

// HandleSuccess handles successful responses with optional data
func (h *Response) HandleSuccess(c *gin.Context, data any) {
	c.Header("Content-Type", "application/json")

	response := gin.H{
		"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
		"status":    string(APIStatusSuccess),
	}

	if data != nil {
		// Convert data struct to map
		b, _ := json.Marshal(data)
		var dataMap map[string]interface{}
		json.Unmarshal(b, &dataMap)

		// Append all fields from data to response
		for k, v := range dataMap {
			response[k] = v
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *Response) HandleSuccessWithCustomRequstID(c *gin.Context, requestId string, data any) {
	c.Header("Content-Type", "application/json")

	response := gin.H{
		"requestId": requestId,
		"status":    string(APIStatusSuccess),
	}

	if data != nil {
		// Convert data struct to map
		b, _ := json.Marshal(data)
		var dataMap map[string]interface{}
		json.Unmarshal(b, &dataMap)

		// Append all fields from data to response
		for k, v := range dataMap {
			response[k] = v
		}
	}

	c.JSON(http.StatusOK, response)
}

// HandleCreated handles creation responses with the new resource ID
func (h *Response) HandleCreated(c *gin.Context, data any) {
	response := gin.H{
		"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
		"status":    string(APIStatusSuccess),
	}

	if data != nil {
		// Convert data struct to map
		b, _ := json.Marshal(data)
		var dataMap map[string]interface{}
		json.Unmarshal(b, &dataMap)

		// Append all fields from data to response
		for k, v := range dataMap {
			response[k] = v
		}
	}

	c.JSON(http.StatusCreated, response)
}

// HandleAccepted handles accepted but not yet processed responses
func (h *Response) HandleAccepted(c *gin.Context) {
	response := gin.H{
		"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
		"status":    string(APIStatusProcessing),
	}

	c.JSON(http.StatusAccepted, response)
}

// HandleError handles error responses with optional message
func (h *Response) HandleError(c *gin.Context, statusCode int, message *string) {
	response := gin.H{
		"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
		"status":    string(APIStatusError),
		"code":      http.StatusText(statusCode),
	}

	if message != nil {
		response["message"] = *message
	}

	c.JSON(statusCode, response)
}

func (h *Response) AbortAndHandleError(c *gin.Context, statusCode int, message *string) {
	response := gin.H{
		"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
		"status":    string(APIStatusError),
		"code":      http.StatusText(statusCode),
	}

	if message != nil {
		response["message"] = *message
	}

	c.JSON(500, response)
}

// method for partial success responses
func (h *Response) HandlePartialSuccess(c *gin.Context, responseObjectName string, data any, warnings []string) {
	response := gin.H{
		"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
		"status":    string(APIStatusPartialSuccess),
	}

	if data != nil {
		response[responseObjectName] = data
	}

	if len(warnings) > 0 {
		response["warnings"] = warnings
	}

	c.JSON(http.StatusOK, response)
}

func (h *Response) HandleDataStream(c *gin.Context, contentType string, data []byte) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Type", contentType)

	c.Data(http.StatusOK, contentType, data)
}
