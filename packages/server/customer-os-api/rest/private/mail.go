package private

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type MailHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewMailHandler(services *cosapi_services.Services, responseHandler *response.Response) *MailHandler {
	return &MailHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

func (h *MailHandler) SendEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "mail/send", c.Request.Header)
		defer span.Finish()

		tenant := c.GetString(security.KEY_TENANT_NAME)

		customCtx := &common.CustomContext{}
		if c.Keys[security.KEY_TENANT_NAME] != nil {
			customCtx.Tenant = c.Keys[security.KEY_TENANT_NAME].(string)
		}
		if c.Keys[security.KEY_USER_ROLES] != nil {
			customCtx.Roles = c.Keys[security.KEY_USER_ROLES].([]string)
		}
		if c.Keys[security.KEY_USER_ID] != nil {
			customCtx.UserId = c.Keys[security.KEY_USER_ID].(string)
		}
		if c.Keys[security.KEY_USER_EMAIL] != nil {
			customCtx.UserEmail = c.Keys[security.KEY_USER_EMAIL].(string)
		}

		ctx = common.WithCustomContext(ctx, customCtx)

		var request *postgres_entity.EmailMessage

		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		request.Tenant = tenant
		request.ProducerId = "N/A"
		request.ProducerType = "N/A"

		span.LogFields(log.Object("request", request))

		err := h.services.CommonServices.MailService.SendMail(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		h.responseHandler.HandleSuccess(c, gin.H{})

	}
}

func (h *MailHandler) TrackEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerOSInternalIdentifier := c.Param("customerOSInternalIdentifier")
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(context.Background(), "/mail/"+customerOSInternalIdentifier+"/track", c.Request.Header)
		defer span.Finish()

		// Preload 1px transparent image
		px := image.NewRGBA(image.Rect(0, 0, 1, 1))
		px.Set(0, 0, color.Transparent)
		var spyPixel bytes.Buffer
		err := png.Encode(&spyPixel, px)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to encode image"))
		}
		spyPixelBytes := spyPixel.Bytes()

		if customerOSInternalIdentifier == "" {
			message := "Missing customerOSInternalIdentifier"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		span.LogFields(log.String("customerOSInternalIdentifier", customerOSInternalIdentifier))

		// log all headers
		for name, values := range c.Request.Header {
			for _, value := range values {
				span.LogFields(log.String("Header: "+name, value))
			}
		}

		interactionEventNode, err := h.services.Repositories.Neo4jRepositories.InteractionEventReadRepository.GetInteractionEventByCustomerOSIdentifier(ctx, customerOSInternalIdentifier)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		if interactionEventNode == nil {
			span.LogFields(log.String("interactionEventId", "not found"))
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		interactionEvent := mapper.MapDbNodeToInteractionEventEntity(interactionEventNode)

		span.LogFields(log.String("interactionEventId", interactionEvent.Id))

		tenant := model.GetTenantFromLabels(interactionEventNode.Labels, model.NodeLabelInteractionEvent)
		if tenant == "" {
			span.LogFields(log.String("tenant", "not identified"))
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		span.SetTag(tracing.SpanTagTenant, tenant)

		metadata, err := utils.ToJson(map[string]interface{}{
			"User-Agent":       c.GetHeader("User-Agent"),
			"Cf-Connecting-Ip": c.GetHeader("Cf-Connecting-Ip"),
		})
		if err != nil {
			message := "Error while converting metadata to json"
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		_, err = h.services.Repositories.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, interactionEvent.Id, model.INTERACTION_EVENT, enum.ActionInteractionEventRead, "", metadata, utils.Now(), "user-admin-api")
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		h.responseHandler.HandleDataStream(c, "image/png", spyPixelBytes)
	}
}
