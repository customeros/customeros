package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"net/http"
)

func CacheHandler(serviceContainer *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {

		tenant := c.Keys[security.KEY_TENANT_NAME].(string)

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.WriteHeader(http.StatusOK)

		// Get the http.Flusher interface to flush data to the client
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.String(http.StatusInternalServerError, "Streaming unsupported!")
			return
		}

		organizations := serviceContainer.Caches.GetOrganizations(tenant)

		// Stream data in chunks
		for i := 0; i < len(organizations); i++ {

			data, err := json.Marshal(organizations[i])
			if err != nil {
				//tracing.TraceErr(span, err)
				//s.log.Errorf("Failed to marshal orgs to json: %v", err)
				//return err
			}
			jsonStr := string(data)

			fmt.Fprintf(c.Writer, "%s\n", jsonStr)
			flusher.Flush()
		}

		c.JSON(200, organizations)
	}
}
