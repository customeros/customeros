package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

func OrganizationsCacheHandler(serviceContainer *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/stream/organizations-cache", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys[security.KEY_TENANT_NAME].(string)

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.WriteHeader(http.StatusOK)

		// Get the http.Flusher interface to flush data to the client
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			tracing.TraceErr(span, fmt.Errorf("Streaming unsupported!"))
			c.String(http.StatusInternalServerError, "Streaming unsupported!")
			return
		}

		organizations := serviceContainer.Caches.GetOrganizations(tenant)

		span.LogFields(log.Int("count", len(organizations)))

		// Stream data in chunks
		for i := 0; i < len(organizations); i++ {

			data, err := json.Marshal(organizations[i])
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Bool("streamed", false))
				return
			}
			jsonStr := string(data)

			_, err = fmt.Fprintf(c.Writer, "%s\n", jsonStr)
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Bool("streamed", false))
				return
			}
			flusher.Flush()
		}

		span.LogFields(log.Bool("streamed", true))
	}
}

func OrganizationsPatchesCacheHandler(serviceContainer *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/stream/organizations-cache", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys[security.KEY_TENANT_NAME].(string)

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.WriteHeader(http.StatusOK)

		// Get the http.Flusher interface to flush data to the client
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			tracing.TraceErr(span, fmt.Errorf("Streaming unsupported!"))
			c.String(http.StatusInternalServerError, "Streaming unsupported!")
			return
		}

		patches, err := serviceContainer.CommonServices.ApiCacheService.GetPatchesForApiCache(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			span.LogFields(log.Bool("streamed", false))
			return
		}

		span.LogFields(log.Int("count", len(patches)))

		// Stream data in chunks
		for i := 0; i < len(patches); i++ {

			data, err := json.Marshal(patches[i])
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Bool("streamed", false))
				return
			}
			jsonStr := string(data)

			_, err = fmt.Fprintf(c.Writer, "%s\n", jsonStr)
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Bool("streamed", false))
				return
			}
			flusher.Flush()
		}

		span.LogFields(log.Bool("streamed", true))
	}
}
