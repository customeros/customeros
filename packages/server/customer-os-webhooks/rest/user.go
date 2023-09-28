package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"io/ioutil"
	"net/http"
)

// Define a struct to represent the JSON structure in the request body
type PocRequest struct {
	Message string `json:"message"`
}

func PocHandler(serviceContainer *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse the JSON request body
		var req PocRequest
		if err := json.Unmarshal(requestBody, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Print the parsed message
		fmt.Println("Received Message:", req.Message)

		c.JSON(http.StatusOK, gin.H{"message": "Message received successfully"})

	}
}
