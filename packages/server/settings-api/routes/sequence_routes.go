package routes

import (
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

	//sequences
	r.GET("/sequences",
		security.TenantUserContextEnhancer(security.TENANT, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getSequencesHandler(services))

	r.GET("/sequences/:id",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getSequenceHandler(services))

	r.POST("/sequences",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		postSequenceHandler(services))

	r.POST("/sequences/:id/enable",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		enableSequenceHandler(services))

	r.POST("/sequences/:id/disable",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		disableSequenceHandler(services))

	//steps
	r.GET("/sequences/:id/steps",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getSequenceStepsHandler(services))

	r.POST("/sequences/:id/steps",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		postSequenceStepHandler(services))

	r.GET("/sequences/:id/steps/:stepId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getSequenceStepHandler(services))

	r.DELETE("/sequences/:id/steps/:stepId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deleteSequenceStepHandler(services))

	//contacts
	r.GET("/sequences/:id/contacts",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getSequenceContactsHandler(services))

	r.POST("/sequences/:id/contacts",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		postSequenceContactHandler(services))

	r.GET("/sequences/:id/contacts/:contactId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getSequenceContactHandler(services))

	r.DELETE("/sequences/:id/contacts/:contactId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deleteSequenceContactHandler(services))

	//mailboxes
	r.GET("/sequences/:id/senders",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		getSequenceSendersHandler(services))

	r.POST("/sequences/:id/senders",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		postSequenceSenderHandler(services))

	r.DELETE("/sequences/:id/senders/:senderId",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		deleteSequenceSenderHandler(services))

}

// @Accept  json
// @Produce  json
// @Success 200 {array} postgresEntity.Sequence
// @Failure 401
// @Failure 500
// @Router /sequences [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getSequencesHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /sequences", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)

		sequences, err := services.CommonServices.SequenceService.GetSequences(ctx, tenant)
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
// @Param   id     path    string     true  "Sequence ID"
// @Success 200 {object} postgresEntity.Sequence
// @Failure 401
// @Failure 404
// @Router /sequences/{id} [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getSequenceHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /sequences/:id", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) {
			c.Status(404)
			return
		}

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequence == nil {
			c.Status(404)
			return
		}

		c.JSON(200, sequence)
	}
}

// @Accept  json
// @Produce  json
// @Param   sequence  body    SequencePostRequest  true  "Sequence entity to be created / updated"
// @Success 200 {object} postgresEntity.Sequence
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /sequences [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func postSequenceHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /sequences", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)

		var request SequencePostRequest

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		var err error
		var sequence *postgresEntity.Sequence

		if request.Id == nil {
			sequence = &postgresEntity.Sequence{}
			sequence.Tenant = tenant
			sequence.CreatedAt = now
		} else {
			if !isValidId(*request.Id) {
				c.Status(404)
				return
			}

			sequence, err = services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, *request.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				c.Status(500)
				return
			}

			if sequence == nil {
				c.Status(404)
				return
			}
		}

		sequence.UpdatedAt = now
		sequence.Name = request.Name

		sequence, err = services.CommonServices.SequenceService.StoreSequence(ctx, sequence)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		span.SetTag(tracing.SpanTagEntityId, sequence.ID)

		c.JSON(200, sequence)
	}
}

// @Accept  json
// @Produce  json
// @Param   id     path    string     true  "Sequence ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /sequences/{id}/enable [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func enableSequenceHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /sequences/:id/enable", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) {
			c.Status(404)
			return
		}

		err := services.CommonServices.SequenceService.EnableSequence(ctx, tenant, sequenceId)
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
// @Param   id     path    string     true  "Sequence ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /sequences/{id}/disable [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func disableSequenceHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /sequences/:id/disable", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) {
			c.Status(404)
			return
		}

		err := services.CommonServices.SequenceService.DisableSequence(ctx, tenant, sequenceId)
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
// @Param   id     path    string     true  "Sequence ID"
// @Success 200 {array} postgresEntity.SequenceStep
// @Failure 401
// @Failure 500
// @Router /sequences/{id}/steps [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getSequenceStepsHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /sequences/:id/steps", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		steps, err := services.CommonServices.SequenceService.GetSequenceSteps(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, steps)
	}
}

// @Accept  json
// @Produce  json
// @Param   id     path    string     true  "Sequence ID"
// @Param   sequenceStep  body    SequenceStepPostRequest  true  "Sequence step entity to be created / updated"
// @Success 200 {object} postgresEntity.SequenceStep
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /sequences/{id}/steps [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func postSequenceStepHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /sequences/:id/steps", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(404)
			return
		}

		var request SequenceStepPostRequest

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		var sequenceStep *postgresEntity.SequenceStep

		if request.Id == nil {
			sequenceStep = &postgresEntity.SequenceStep{}
			sequenceStep.SequenceId = sequence.ID
			sequenceStep.CreatedAt = now
		} else {
			if !isValidId(sequenceId) {
				c.Status(404)
				return
			}

			sequenceStep, err = services.CommonServices.SequenceService.GetSequenceStepById(ctx, tenant, *request.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				c.Status(500)
				return
			}

			if sequenceStep == nil {
				c.Status(404)
				return
			}

			if sequenceStep.SequenceId != sequence.ID {
				c.Status(404)
				return
			}
		}

		sequenceStep.UpdatedAt = now
		sequenceStep.Type = request.Type
		sequenceStep.Name = request.Name
		sequenceStep.Text = request.Text
		sequenceStep.Template = request.Template

		sequenceStep, err = services.CommonServices.SequenceService.StoreSequenceStep(ctx, tenant, sequenceStep)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, sequenceStep)
	}
}

// @Accept  json
// @Produce  json
// @Param   id     path    string     true  "Sequence ID"
// @Param   stepId     path    string     true  "Sequence step ID"
// @Success 200 {object} postgresEntity.SequenceStep
// @Failure 401
// @Failure 404
// @Router /sequences/{id}/steps/{stepId} [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getSequenceStepHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")
		sequenceStepId := c.Param("stepId")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /sequences/:id/steps/:stepId", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) || !isValidId(sequenceStepId) {
			c.Status(404)
			return
		}

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequence == nil {
			c.Status(404)
			return
		}

		sequenceStep, err := services.CommonServices.SequenceService.GetSequenceStepById(ctx, tenant, sequenceStepId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequenceStep == nil {
			c.Status(404)
			return
		}

		if sequenceStep.SequenceId != sequence.ID {
			c.Status(404)
			return
		}

		c.JSON(200, sequenceStep)
	}
}

// @Accept  json
// @Produce  json
// @Param   id     path    string     true  "Sequence ID"
// @Param   stepId     path    string     true  "Sequence step ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Router /sequences/{id}/steps/{stepId} [delete]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deleteSequenceStepHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")
		sequenceStepId := c.Param("stepId")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "DELETE /sequences/:id/steps/:stepId", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) || !isValidId(sequenceStepId) {
			c.Status(404)
			return
		}

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequence == nil {
			c.Status(404)
			return
		}

		sequenceStep, err := services.CommonServices.SequenceService.GetSequenceStepById(ctx, tenant, sequenceStepId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequenceStep == nil {
			c.Status(404)
			return
		}

		if sequenceStep.SequenceId != sequence.ID {
			c.Status(404)
			return
		}

		err = services.CommonServices.SequenceService.DeleteSequenceStep(ctx, tenant, sequenceStepId)
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
// @Param   id     path    string     true  "Sequence ID"
// @Success 200 {array} postgresEntity.SequenceContact
// @Failure 401
// @Failure 500
// @Router /sequences/{id}/contacts [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getSequenceContactsHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /sequences/:id/contacts", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		contacts, err := services.CommonServices.SequenceService.GetSequenceContacts(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, contacts)
	}
}

// @Accept  json
// @Produce  json
// @Param   id     path    string     true  "Sequence ID"
// @Param   sequenceContact  body    SequenceContactPostRequest  true  "Sequence contact entity to be created / updated"
// @Success 200 {object} postgresEntity.SequenceContact
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /sequences/{id}/contacts [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func postSequenceContactHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /sequences/:id/contacts", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) {
			c.Status(404)
			return
		}

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(404)
			return
		}

		var request SequenceContactPostRequest

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		var sequenceContact *postgresEntity.SequenceContact

		if request.Id == nil {
			sequenceContact = &postgresEntity.SequenceContact{}
			sequenceContact.SequenceId = sequence.ID
			sequenceContact.CreatedAt = now
		} else {
			if !isValidId(*request.Id) {
				c.Status(404)
				return
			}

			sequenceContact, err = services.CommonServices.SequenceService.GetSequenceContactById(ctx, tenant, *request.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				c.Status(500)
				return
			}

			if sequenceContact == nil {
				c.Status(404)
				return
			}

			if sequenceContact.SequenceId != sequence.ID {
				c.Status(404)
				return
			}
		}

		sequenceContact.UpdatedAt = now
		sequenceContact.FirstName = request.FirstName
		sequenceContact.LastName = request.LastName
		sequenceContact.Email = request.Email
		sequenceContact.LinkedinUrl = request.LinkedinUrl

		sequenceContact, err = services.CommonServices.SequenceService.StoreSequenceContact(ctx, tenant, sequenceContact)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, sequenceContact)
	}
}

// @Accept  json
// @Produce  json
// @Param   id     path    string     true  "Sequence ID"
// @Param   contactId     path    string     true  "Sequence contact ID"
// @Success 200 {object} postgresEntity.SequenceContact
// @Failure 401
// @Failure 404
// @Router /sequences/{id}/contacts/{contactId} [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getSequenceContactHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")
		sequenceContactId := c.Param("contactId")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /sequences/:id/contacts/:contactId", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) || !isValidId(sequenceContactId) {
			c.Status(404)
			return
		}

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequence == nil {
			c.Status(404)
			return
		}

		sequenceContact, err := services.CommonServices.SequenceService.GetSequenceContactById(ctx, tenant, sequenceContactId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequenceContact == nil {
			c.Status(404)
			return
		}

		if sequenceContact.SequenceId != sequence.ID {
			c.Status(404)
			return
		}

		c.JSON(200, sequenceContact)
	}
}

// @Accept  json
// @Produce  json
// @Param   id     path    string     true  "Sequence ID"
// @Param   contactId     path    string     true  "Sequence contact ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Router /sequences/{id}/contacts/{contactId} [delete]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deleteSequenceContactHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")
		sequenceContactId := c.Param("contactId")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "DELETE /sequences/:id/contacts/:contactId", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) || !isValidId(sequenceContactId) {
			c.Status(404)
			return
		}

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequence == nil {
			c.Status(404)
			return
		}

		sequenceContact, err := services.CommonServices.SequenceService.GetSequenceContactById(ctx, tenant, sequenceContactId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequenceContact == nil {
			c.Status(404)
			return
		}

		if sequenceContact.SequenceId != sequence.ID {
			c.Status(404)
			return
		}

		err = services.CommonServices.SequenceService.DeleteSequenceContact(ctx, tenant, sequenceContactId)
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
// @Param   id     path    string     true  "Sequence ID"
// @Success 200 {array} postgresEntity.SequenceSender
// @Failure 401
// @Failure 500
// @Router /sequences/{id}/senders [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getSequenceSendersHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /sequences/:id/senders", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) {
			c.Status(404)
			return
		}

		senders, err := services.CommonServices.SequenceService.GetSequenceSenders(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, senders)
	}
}

// @Accept  json
// @Produce  json
// @Param   id     path    string     true  "Sequence ID"
// @Param   sequenceSender  body    SequenceSenderPostRequest  true  "Sequence sender entity to be created / updated"
// @Success 200 {object} postgresEntity.SequenceSender
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /sequences/{id}/senders [post]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func postSequenceSenderHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "POST /sequences/:id/senders", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) {
			c.Status(404)
			return
		}

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(404)
			return
		}

		var request SequenceSenderPostRequest

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

		sequenceSender := &postgresEntity.SequenceSender{}
		sequenceSender.SequenceId = sequence.ID
		sequenceSender.CreatedAt = now
		sequenceSender.UpdatedAt = now
		sequenceSender.MailboxId = request.MailboxId

		sequenceSender, err = services.CommonServices.SequenceService.StoreSequenceSender(ctx, tenant, sequenceSender)
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
// @Param   id           path    string     true  "Sequence ID"
// @Param   senderId     path    string     true  "Sequence sender ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Router /sequences/{id}/senders/{senderId} [delete]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func deleteSequenceSenderHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		sequenceId := c.Param("id")
		sequenceSenderId := c.Param("senderId")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "DELETE /sequences/:id/senders/:senderId", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)
		span.SetTag(tracing.SpanTagEntityId, sequenceId)

		if !isValidId(sequenceId) || !isValidId(sequenceSenderId) {
			c.Status(404)
			return
		}

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequence == nil {
			c.Status(404)
			return
		}

		sequenceSender, err := services.CommonServices.SequenceService.GetSequenceSenderById(ctx, tenant, sequenceSenderId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		if sequenceSender == nil {
			c.Status(404)
			return
		}

		if sequenceSender.SequenceId != sequence.ID {
			c.Status(404)
			return
		}

		err = services.CommonServices.SequenceService.DeleteSequenceSender(ctx, tenant, sequenceSenderId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.Status(200)
	}
}

func isValidId(id string) bool {
	idParsed, err := uuid.Parse(id)
	if err != nil || idParsed.String() != id {
		return false
	}

	return true
}

type SequencePostRequest struct {
	Id   *string
	Name string
}

type SequenceStepPostRequest struct {
	Id *string

	Order int
	Type  string
	Name  string

	Text     *string
	Template *string
}

type SequenceContactPostRequest struct {
	Id *string

	FirstName   *string
	LastName    *string
	Email       string
	LinkedinUrl *string
}

type SequenceSenderPostRequest struct {
	MailboxId string
}
