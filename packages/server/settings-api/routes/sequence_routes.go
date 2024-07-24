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
	//r.GET("/sequences/:id/senders",
	//	security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
	//	security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
	//	getSequenceSendersHandler(services))
	//
	//r.POST("/sequences/:id/senders",
	//	security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
	//	security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
	//	postSequenceSenderHandler(services))
	//
	//r.DELETE("/sequences/:id/senders/:senderId",
	//	security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
	//	security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
	//	deleteSequenceContactHandler(services))

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
// @Param   sequence  body    postgresEntity.Sequence  true  "Sequence entity to be created / updated"
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

		var request postgresEntity.Sequence

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		if request.ID == "" {
			request.T = postgresEntity.Tenant{
				Name: tenant,
			}
			request.CreatedAt = now
		}

		request.UpdatedAt = now

		sequence, err := services.CommonServices.SequenceService.StoreSequence(ctx, &request)
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
// @Param   sequence  body    postgresEntity.SequenceStep  true  "Sequence step entity to be created / updated"
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

		var request postgresEntity.SequenceStep

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		var sequenceStep *postgresEntity.SequenceStep

		if request.ID == "" {
			sequenceStep = &postgresEntity.SequenceStep{}
			sequenceStep.SequenceId = sequence.ID
			sequenceStep.CreatedAt = now
		} else {
			sequenceStep, err = services.CommonServices.SequenceService.GetSequenceStepById(ctx, tenant, request.ID)
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
// @Param   sequence  body    postgresEntity.SequenceContact  true  "Sequence contact entity to be created / updated"
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

		sequence, err := services.CommonServices.SequenceService.GetSequenceById(ctx, tenant, sequenceId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(404)
			return
		}

		var request postgresEntity.SequenceContact

		if err := c.BindJSON(&request); err != nil {
			span.LogFields(tracingLog.Object("request", request))
			tracing.TraceErr(span, err)
			c.Status(400)
			return
		}

		now := utils.Now()

		var sequenceContact *postgresEntity.SequenceContact

		if request.ID == "" {
			sequenceContact = &postgresEntity.SequenceContact{}
			sequenceContact.SequenceId = sequence.ID
			sequenceContact.CreatedAt = now
		} else {
			sequenceContact, err = services.CommonServices.SequenceService.GetSequenceContactById(ctx, tenant, request.ID)
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

		err = services.CommonServices.SequenceService.DeleteSequenceContact(ctx, tenant, sequenceContactId)
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
