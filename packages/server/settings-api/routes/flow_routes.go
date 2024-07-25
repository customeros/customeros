package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
	tracingLog "github.com/opentracing/opentracing-go/log"
)

func InitSequenceRoutes(r *gin.Engine, services *service.Services) {

	//flows
	r.GET("/flows",
		security.TenantUserContextEnhancer(security.TENANT, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getFlowsHandler(services))

	r.GET("/flows/:flowId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getFlowHandler(services))

	r.POST("/flows",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		postFlowHandler(services))

	r.POST("/flows/:flowId/activate",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		activateFlowHandler(services))

	r.POST("/flows/:flowId/deactivate",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deactivateFlowHandler(services))

	r.DELETE("/flows/:flowId/delete",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deleteFlowHandler(services))

	//sequences
	r.GET("/flows/:flowId/sequences",
		security.TenantUserContextEnhancer(security.TENANT, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getFlowSequencesHandler(services))

	r.GET("/flows/:flowId/sequences/:flowSequenceId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getFlowSequenceHandler(services))

	r.POST("/flows/:flowId/sequences",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		postFlowSequenceHandler(services))

	r.POST("/flows/:flowId/sequences/:flowSequenceId/activate",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		activateFlowSequenceHandler(services))

	r.POST("/flows/:flowId/sequences/:flowSequenceId/deactivate",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deactivateFlowSequenceHandler(services))

	r.DELETE("/flows/:flowId/sequences/:flowSequenceId/delete",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deleteFlowSequenceHandler(services))

	//steps
	r.GET("/flows/:flowId/sequences/:flowSequenceId/steps",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getFlowSequenceStepsHandler(services))

	r.POST("/flows/:flowId/sequences/:flowSequenceId/steps",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		postSequenceStepHandler(services))

	r.GET("/flows/:flowId/sequences/:flowSequenceId/steps/:flowSequenceStepId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getSequenceStepHandler(services))

	r.POST("/flows/:flowId/sequences/:flowSequenceId/steps/:flowSequenceStepId/activate",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		activateSequenceStepHandler(services))

	r.POST("/flows/:flowId/sequences/:flowSequenceId/steps/:flowSequenceStepId/deactivate",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deactivateSequenceStepHandler(services))

	r.DELETE("/flows/:flowId/sequences/:flowSequenceId/steps/:flowSequenceStepId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deleteSequenceStepHandler(services))

	//contacts
	r.GET("/flows/:flowId/sequences/:flowSequenceId/contacts",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getFlowSequenceContactsHandler(services))

	r.POST("/flows/:flowId/sequences/:flowSequenceId/contacts",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		postFlowSequenceContactHandler(services))

	r.GET("/flows/:flowId/sequences/:flowSequenceId/contacts/:flowSequenceContactId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getFlowSequenceContactHandler(services))

	r.DELETE("/flows/:flowId/sequences/:flowSequenceId/contacts/:flowSequenceContactId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deleteFlowSequenceContactHandler(services))

	//mailboxes
	r.GET("/flows/:flowId/sequences/:flowSequenceId/senders",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getSequenceSendersHandler(services))

	r.POST("/flows/:flowId/sequences/:flowSequenceId/senders",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		postSequenceSenderHandler(services))

	r.DELETE("/flows/:flowId/sequences/:flowSequenceId/senders/:flowSequenceSenderId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deleteFlowSequenceSenderHandler(services))

}

// @Accept  json
// @Produce  json
// @Success 200 {array} postgresEntity.Flow
// @Failure 401
// @Failure 500
// @Router /flows [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getFlowsHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /flows", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)

		flows, err := services.CommonServices.FlowService.GetFlows(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, flows)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Success 200 {object} postgresEntity.Flow
// @Failure 401
// @Failure 404
// @Router /flows/{flowId} [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getFlowHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /flows/{flowId}", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flow, err := retrieveFlow(services, ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flow == nil {
			c.Status(404)
			return
		}

		c.JSON(200, flow)
	}
}

// @Accept  json
// @Produce  json
// @Param   flow  body    FlowPostRequest  true  "Flow entity to be created / updated"
// @Success 200 {object} postgresEntity.Flow
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /flows [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func postFlowHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)

		var request FlowPostRequest

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		var err error
		var flow *postgresEntity.Flow

		if request.Id == nil {
			flow = &postgresEntity.Flow{}
			flow.Tenant = tenant
			flow.CreatedAt = now
		} else {

			flow, err := retrieveFlow(services, ctx, tenant, *request.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				c.Status(500)
				return
			}

			if flow == nil {
				c.Status(404)
				return
			}
		}

		flow.UpdatedAt = now
		flow.Name = request.Name
		flow.Description = request.Description

		flow, err = services.CommonServices.FlowService.StoreFlow(ctx, flow)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		span.SetTag(tracing.SpanTagEntityId, flow.ID)

		c.JSON(200, flow)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /flows/{flowId}/activate [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func activateFlowHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/activate", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flow, err := retrieveFlow(services, ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flow == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.ActivateFlow(ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /flows/{flowId}/deactivate [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deactivateFlowHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/deactivate", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flow, err := retrieveFlow(services, ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flow == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.DeactivateFlow(ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /flows/{flowId}/delete [delete]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deleteFlowHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "DELETE /flows/{flowId}/delete", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flow, err := retrieveFlow(services, ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flow == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.DeleteFlow(ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce  json
// @Success 200 {array} postgresEntity.FlowSequence
// @Failure 401
// @Failure 500
// @Router /flows/{flowId}/sequences [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getFlowSequencesHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /flows/{flowId}/sequences", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flow, err := retrieveFlow(services, ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flow == nil {
			c.Status(404)
			return
		}

		sequences, err := services.CommonServices.FlowService.GetFlowSequences(ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, sequences)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Success 200 {object} postgresEntity.FlowSequence
// @Failure 401
// @Failure 404
// @Router /flows/{flowId}/sequences/{flowSequenceId} [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getFlowSequenceHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /flows/{flowId}/sequences/{sequenceId}", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		c.JSON(200, flowSequence)
	}
}

// @Accept  json
// @Produce  json
// @Param   sequence  body    FlowSequencePostRequest  true  "FlowSequence entity to be created / updated"
// @Success 200 {object} postgresEntity.FlowSequence
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /flows/{flowId}/sequences [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func postFlowSequenceHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/sequences", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flow, err := retrieveFlow(services, ctx, tenant, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flow == nil {
			c.Status(404)
			return
		}

		var request FlowSequencePostRequest

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		var flowSequence *postgresEntity.FlowSequence

		if request.Id == nil {
			flowSequence = &postgresEntity.FlowSequence{}
			flowSequence.FlowId = flowId
			flowSequence.CreatedAt = now
		} else {
			flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, *request.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				c.Status(500)
				return
			}

			if flowSequence == nil {
				c.Status(404)
				return
			}
		}

		flowSequence.UpdatedAt = now
		flowSequence.Name = request.Name
		flowSequence.Description = request.Description

		flowSequence, err = services.CommonServices.FlowService.StoreFlowSequence(ctx, tenant, flowSequence)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		span.SetTag(tracing.SpanTagEntityId, flowSequence.ID)

		c.JSON(200, flowSequence)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/activate [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func activateFlowSequenceHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/sequences/{flowSequenceId}/activate", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.ActivateFlowSequence(ctx, tenant, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/deactivate [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deactivateFlowSequenceHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/sequences/{flowSequenceId}/deactivate", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.DeactivateFlowSequence(ctx, tenant, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/delete [delete]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deleteFlowSequenceHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "DELETE /flows/{flowId}/sequences/{flowSequenceId}/delete", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.DeleteFlowSequence(ctx, tenant, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Success 200 {array} postgresEntity.FlowSequenceStep
// @Failure 401
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/steps [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getFlowSequenceStepsHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /flows/{flowId}/sequences/{flowSequenceId}/steps", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		flowSequenceSteps, err := services.CommonServices.FlowService.GetFlowSequenceSteps(ctx, tenant, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, flowSequenceSteps)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceStep  body    FlowSequenceStepPostRequest  true  "FlowSequenceStep entity to be created / updated"
// @Success 200 {object} postgresEntity.FlowSequenceStep
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/steps [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func postSequenceStepHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/sequences/{flowSequenceId}/steps", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		var request FlowSequenceStepPostRequest

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		var flowSequenceStep *postgresEntity.FlowSequenceStep

		if request.Id == nil {
			flowSequenceStep = &postgresEntity.FlowSequenceStep{}
			flowSequenceStep.SequenceId = flowSequenceId
			flowSequenceStep.CreatedAt = now
		} else {
			flowSequenceStep, err = retrieveFlowSequenceStep(services, ctx, tenant, flowId, flowSequenceId, *request.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				c.Status(500)
				return
			}

			if flowSequenceStep == nil {
				c.Status(404)
				return
			}
		}

		flowSequenceStep.UpdatedAt = now
		flowSequenceStep.Type = request.Type
		flowSequenceStep.Name = request.Name
		flowSequenceStep.Text = request.Text
		flowSequenceStep.Template = request.Template

		flowSequenceStep, err = services.CommonServices.FlowService.StoreFlowSequenceStep(ctx, tenant, flowSequenceStep)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, flowSequenceStep)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceStepId     path    string     true  "FlowSequenceStep ID"
// @Success 200 {object} postgresEntity.FlowSequenceStep
// @Failure 401
// @Failure 404
// @Router /flows/{flowId}/sequences/{flowSequenceId}/steps/{flowSequenceStepId} [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getSequenceStepHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /flows/{flowId}/sequences/{flowSequenceId}/steps/{flowSequenceStepId}", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		flowSequenceStepId := c.Param("flowSequenceStepId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequenceStep, err := retrieveFlowSequenceStep(services, ctx, tenant, flowId, flowSequenceId, flowSequenceStepId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequenceStep == nil {
			c.Status(404)
			return
		}

		c.JSON(200, flowSequenceStep)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceStepId     path    string     true  "FlowSequenceStep ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/steps/{flowSequenceStepId}/activate [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func activateSequenceStepHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/sequences/{flowSequenceId}/steps/{flowSequenceStepId}/activate", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		flowSequenceStepId := c.Param("flowSequenceStepId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequenceStep, err := retrieveFlowSequenceStep(services, ctx, tenant, flowId, flowSequenceId, flowSequenceStepId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequenceStep == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.ActivateFlowSequenceStep(ctx, tenant, flowSequenceStepId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceStepId     path    string     true  "FlowSequenceStep ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/steps/{flowSequenceStepId}/deactivate [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deactivateSequenceStepHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/sequences/{flowSequenceId}/steps/{flowSequenceStepId}/deactivate", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		flowSequenceStepId := c.Param("flowSequenceStepId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequenceStep, err := retrieveFlowSequenceStep(services, ctx, tenant, flowId, flowSequenceId, flowSequenceStepId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequenceStep == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.DeactivateFlowSequenceStep(ctx, tenant, flowSequenceStepId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceStepId     path    string     true  "FlowSequenceStep ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/steps/{flowSequenceStepId}/deactivate [delete]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deleteSequenceStepHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "DELETE /flows/{flowId}/sequences/{flowSequenceId}/steps/{flowSequenceStepId}/deactivate", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		flowSequenceStepId := c.Param("flowSequenceStepId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequenceStep, err := retrieveFlowSequenceStep(services, ctx, tenant, flowId, flowSequenceId, flowSequenceStepId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequenceStep == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.DeleteFlowSequenceStep(ctx, tenant, flowSequenceStepId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Success 200 {array} postgresEntity.FlowSequenceContact
// @Failure 401
// @Failure 500
// @Router  /flows/{flowId}/sequences/{flowSequenceId}/contacts [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getFlowSequenceContactsHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /flows/{flowId}/sequences/{flowSequenceId}/contacts", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		flowSequenceContacts, err := services.CommonServices.FlowService.GetFlowSequenceContacts(ctx, tenant, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, flowSequenceContacts)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceContact  body    FlowSequenceContactPostRequest  true  "FlowSequence contact entity to be created / updated"
// @Success 200 {object} postgresEntity.FlowSequenceContact
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/contacts [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func postFlowSequenceContactHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/sequences/{flowSequenceId}/contacts", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		var request FlowSequenceContactPostRequest

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		var flowSequenceContact *postgresEntity.FlowSequenceContact

		if request.Id == nil {
			flowSequenceContact = &postgresEntity.FlowSequenceContact{}
			flowSequenceContact.SequenceId = flowSequenceId
			flowSequenceContact.CreatedAt = now
		} else {
			flowSequenceContact, err = retrieveFlowSequenceContact(services, ctx, tenant, flowId, flowSequenceId, *request.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				c.Status(500)
				return
			}

			if flowSequenceContact == nil {
				c.Status(404)
				return
			}
		}

		flowSequenceContact.UpdatedAt = now
		flowSequenceContact.FirstName = request.FirstName
		flowSequenceContact.LastName = request.LastName
		flowSequenceContact.Email = request.Email
		flowSequenceContact.LinkedinUrl = request.LinkedinUrl

		flowSequenceContact, err = services.CommonServices.FlowService.StoreFlowSequenceContact(ctx, tenant, flowSequenceContact)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, flowSequenceContact)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceContactId     path    string     true  "FlowSequenceContact ID"
// @Param   contactId     path    string     true  "FlowSequence contact ID"
// @Success 200 {object} postgresEntity.FlowSequenceContact
// @Failure 401
// @Failure 404
// @Router /flows/{flowId}/sequences/{flowSequenceId}/contacts/{flowSequenceContactId} [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getFlowSequenceContactHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /flows/{flowId}/sequences/{flowSequenceId}/contacts/{flowSequenceContactId}", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		flowSequenceContactId := c.Param("flowSequenceContactId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequenceContact, err := retrieveFlowSequenceContact(services, ctx, tenant, flowId, flowSequenceId, flowSequenceContactId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequenceContact == nil {
			c.Status(404)
			return
		}

		c.JSON(200, flowSequenceContact)
	}
}

// @Accept  json
// @Produce json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceContactId     path    string     true  "FlowSequenceContact ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Router /flows/{flowId}/sequences/{flowSequenceId}/contacts/{flowSequenceContactId} [delete]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deleteFlowSequenceContactHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "DELETE /flows/{flowId}/sequences/{flowSequenceId}/contacts/{flowSequenceContactId}", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		flowSequenceContactId := c.Param("flowSequenceContactId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequenceContact, err := retrieveFlowSequenceContact(services, ctx, tenant, flowId, flowSequenceId, flowSequenceContactId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequenceContact == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.DeleteFlowSequenceContact(ctx, tenant, flowSequenceContactId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Success 200 {array} postgresEntity.FlowSequenceSender
// @Failure 401
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/senders [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getSequenceSendersHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /flows/{flowId}/sequences/{flowSequenceId}/senders", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		flowSequenceSenders, err := services.CommonServices.FlowService.GetFlowSequenceSenders(ctx, tenant, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, flowSequenceSenders)
	}
}

// @Accept  json
// @Produce  json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceSender  body    FlowSequenceSenderPostRequest  true  "FlowSequenceSender sender entity to be created / updated"
// @Success 200 {object} postgresEntity.FlowSequenceSender
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /flows/{flowId}/sequences/{flowSequenceId}/senders [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func postSequenceSenderHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /flows/{flowId}/sequences/{flowSequenceId}/senders", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequence, err := retrieveFlowSequence(services, ctx, tenant, flowId, flowSequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequence == nil {
			c.Status(404)
			return
		}

		var request FlowSequenceSenderPostRequest

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		if !isValidId(request.MailboxId) {
			c.Status(404)
			return
		}

		mailbox, err := services.CommonServices.PostgresRepositories.TenantSettingsMailboxRepository.GetById(ctx, tenant, request.MailboxId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if mailbox == nil {
			c.Status(404)
			return
		}

		now := utils.Now()

		sequenceSender := &postgresEntity.FlowSequenceSender{}
		sequenceSender.SequenceId = flowSequenceId
		sequenceSender.CreatedAt = now
		sequenceSender.UpdatedAt = now
		sequenceSender.MailboxId = request.MailboxId

		sequenceSender, err = services.CommonServices.FlowService.StoreFlowSequenceSender(ctx, tenant, sequenceSender)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, sequenceSender)
	}
}

// @Accept  json
// @Produce json
// @Param   flowId     path    string     true  "Flow ID"
// @Param   flowSequenceId     path    string     true  "FlowSequence ID"
// @Param   flowSequenceSenderId     path    string     true  "FlowSequenceSender ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Router /flows/{flowId}/sequences/{flowSequenceId}/senders/{flowSequenceSenderId} [delete]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deleteFlowSequenceSenderHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "DELETE /flows/{flowId}/sequences/{flowSequenceId}/senders/{flowSequenceSenderId}", c.Request.Header)
		defer span.Finish()

		flowId := c.Param("flowId")
		flowSequenceId := c.Param("flowSequenceId")
		flowSequenceSenderId := c.Param("flowSequenceSenderId")
		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, flowId)

		flowSequenceSender, err := retrieveFlowSequenceSender(services, ctx, tenant, flowId, flowSequenceId, flowSequenceSenderId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if flowSequenceSender == nil {
			c.Status(404)
			return
		}

		err = services.CommonServices.FlowService.DeleteFlowSequenceSender(ctx, tenant, flowSequenceSenderId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

func retrieveFlow(services *service.Services, ctx context.Context, tenant, flowId string) (*postgresEntity.Flow, error) {
	if !isValidId(flowId) {
		return nil, nil
	}

	flow, err := services.CommonServices.FlowService.GetFlowById(ctx, tenant, flowId)
	if err != nil {
		return nil, err
	}

	if flow == nil {
		return nil, nil
	}

	return flow, nil
}

func retrieveFlowSequence(services *service.Services, ctx context.Context, tenant, flowId, sequenceId string) (*postgresEntity.FlowSequence, error) {
	if !isValidId(flowId) {
		return nil, nil
	}
	if !isValidId(sequenceId) {
		return nil, nil
	}

	flow, err := services.CommonServices.FlowService.GetFlowById(ctx, tenant, flowId)
	if err != nil {
		return nil, err
	}

	if flow == nil {
		return nil, nil
	}

	sequence, err := services.CommonServices.FlowService.GetFlowSequenceById(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	if sequence == nil {
		return nil, nil
	}

	if sequence.FlowId != flow.ID {
		return nil, nil
	}

	return sequence, nil
}

func retrieveFlowSequenceStep(services *service.Services, ctx context.Context, tenant, flowId, sequenceId, stepId string) (*postgresEntity.FlowSequenceStep, error) {
	if !isValidId(flowId) {
		return nil, nil
	}
	if !isValidId(sequenceId) {
		return nil, nil
	}
	if !isValidId(stepId) {
		return nil, nil
	}

	flow, err := services.CommonServices.FlowService.GetFlowById(ctx, tenant, flowId)
	if err != nil {
		return nil, err
	}

	if flow == nil {
		return nil, nil
	}

	sequence, err := services.CommonServices.FlowService.GetFlowSequenceById(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	if sequence == nil {
		return nil, nil
	}

	if sequence.FlowId != flow.ID {
		return nil, nil
	}

	step, err := services.CommonServices.FlowService.GetFlowSequenceStepById(ctx, tenant, stepId)
	if err != nil {
		return nil, err
	}

	if step == nil {
		return nil, nil
	}

	if step.SequenceId != sequence.ID {
		return nil, nil
	}

	return step, nil
}

func retrieveFlowSequenceContact(services *service.Services, ctx context.Context, tenant, flowId, sequenceId, contactId string) (*postgresEntity.FlowSequenceContact, error) {
	if !isValidId(flowId) {
		return nil, nil
	}
	if !isValidId(sequenceId) {
		return nil, nil
	}
	if !isValidId(contactId) {
		return nil, nil
	}

	flow, err := services.CommonServices.FlowService.GetFlowById(ctx, tenant, flowId)
	if err != nil {
		return nil, err
	}

	if flow == nil {
		return nil, nil
	}

	sequence, err := services.CommonServices.FlowService.GetFlowSequenceById(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	if sequence == nil {
		return nil, nil
	}

	if sequence.FlowId != flow.ID {
		return nil, nil
	}

	contact, err := services.CommonServices.FlowService.GetFlowSequenceContactById(ctx, tenant, contactId)
	if err != nil {
		return nil, err
	}

	if contact == nil {
		return nil, nil
	}

	if contact.SequenceId != sequence.ID {
		return nil, nil
	}

	return contact, nil
}

func retrieveFlowSequenceSender(services *service.Services, ctx context.Context, tenant, flowId, sequenceId, senderId string) (*postgresEntity.FlowSequenceSender, error) {
	if !isValidId(flowId) {
		return nil, nil
	}
	if !isValidId(sequenceId) {
		return nil, nil
	}
	if !isValidId(senderId) {
		return nil, nil
	}

	flow, err := services.CommonServices.FlowService.GetFlowById(ctx, tenant, flowId)
	if err != nil {
		return nil, err
	}

	if flow == nil {
		return nil, nil
	}

	sequence, err := services.CommonServices.FlowService.GetFlowSequenceById(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	if sequence == nil {
		return nil, nil
	}

	if sequence.FlowId != flow.ID {
		return nil, nil
	}

	sender, err := services.CommonServices.FlowService.GetFlowSequenceSenderById(ctx, tenant, senderId)
	if err != nil {
		return nil, err
	}

	if sender == nil {
		return nil, nil
	}

	if sender.SequenceId != sequence.ID {
		return nil, nil
	}

	return sender, nil
}

func isValidId(id string) bool {
	idParsed, err := uuid.Parse(id)
	if err != nil || idParsed.String() != id {
		return false
	}

	return true
}

type FlowPostRequest struct {
	Id          *string
	Name        string
	Description string
}

type FlowSequencePostRequest struct {
	Id          *string
	Name        string
	Description string
}

type FlowSequenceStepPostRequest struct {
	Id *string

	Order int
	Type  string
	Name  string

	Text     *string
	Template *string
}

type FlowSequenceContactPostRequest struct {
	Id *string

	FirstName   *string
	LastName    *string
	Email       string
	LinkedinUrl *string
}

type FlowSequenceSenderPostRequest struct {
	MailboxId string
}
