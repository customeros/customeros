package private

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func SendEmail(s *cosapi_services.Services) gin.HandlerFunc {
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

		var request *entity.EmailMessage

		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		request.Tenant = tenant
		request.ProducerId = "N/A"
		request.ProducerType = "N/A"

		span.LogFields(log.Object("request", request))

		err := s.CommonServices.MailService.SendMail(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func TrackEmail(s *cosapi_services.Services) gin.HandlerFunc {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing customerOSInternalIdentifier"})
			return
		}

		span.LogFields(log.String("customerOSInternalIdentifier", customerOSInternalIdentifier))

		// log all headers
		for name, values := range c.Request.Header {
			for _, value := range values {
				span.LogFields(log.String("Header: "+name, value))
			}
		}

		interactionEventNode, err := s.Repositories.Neo4jRepositories.InteractionEventReadRepository.GetInteractionEventByCustomerOSIdentifier(ctx, customerOSInternalIdentifier)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if interactionEventNode == nil {
			span.LogFields(log.String("interactionEventId", "not found"))
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		interactionEvent := mapper.MapDbNodeToInteractionEventEntity(interactionEventNode)

		span.LogFields(log.String("interactionEventId", interactionEvent.Id))

		tenant := model.GetTenantFromLabels(interactionEventNode.Labels, model.NodeLabelInteractionEvent)
		if tenant == "" {
			span.LogFields(log.String("tenant", "not identified"))
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		span.SetTag(tracing.SpanTagTenant, tenant)

		metadata, err := utils.ToJson(map[string]interface{}{
			"User-Agent":       c.GetHeader("User-Agent"),
			"Cf-Connecting-Ip": c.GetHeader("Cf-Connecting-Ip"),
		})
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while converting metadata to json"})
			return
		}

		_, err = s.Repositories.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, interactionEvent.Id, model.INTERACTION_EVENT, enum.ActionInteractionEventRead, "", metadata, utils.Now(), "user-admin-api")
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		c.Data(http.StatusOK, "image/png", spyPixelBytes)
	}
}
