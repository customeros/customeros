package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"time"
)

func SyncUsersHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncUsers", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeCommon)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Errorf("(SyncUsers) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var users []model.UserData
		if err = json.Unmarshal(requestBody, &users); err != nil {
			log.Errorf("(SyncUsers) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing users in request"})
			return
		}

		// Context timeout, allocate per user
		timeout := time.Duration(len(users)) * utils.LongDuration
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		err = services.UserService.SyncUsers(ctx, users)
		if err != nil {
			log.Errorf("(SyncUsers) error in sync users: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing users"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Message received successfully"})
		}
	}
}

func SyncUserHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncUsers", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeCommon)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Errorf("(SyncUsers) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var user model.UserData
		if err = json.Unmarshal(requestBody, &user); err != nil {
			log.Errorf("(SyncUsers) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per user
		timeout := utils.LongDuration
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		err = services.UserService.SyncUsers(ctx, []model.UserData{user})
		if err != nil {
			log.Errorf("(SyncUsers) error in sync users: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing users"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Message received successfully"})
		}
	}
}
