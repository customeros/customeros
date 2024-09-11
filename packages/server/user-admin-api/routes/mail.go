package routes

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
)

func addMailRoutes(rg *gin.RouterGroup, conf *config.Config, services *service.Services) {

	//Preload 1px transparent image
	px := image.NewRGBA(image.Rect(0, 0, 1, 1))
	px.Set(0, 0, color.Transparent)

	var spyPixel bytes.Buffer
	err := png.Encode(&spyPixel, px)
	if err != nil {
		log.Printf("unable to encode image: %v", err)
	}
	var spyPixelBytes = spyPixel.Bytes()

	rg.POST("/mail/send",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.USER_ADMIN_API, security.WithCache(caches.NewCommonCache())),
		func(c *gin.Context) {
			ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "mail/send", c.Request.Header)
			defer span.Finish()

			username := c.GetString(security.KEY_USER_EMAIL)
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
			if c.Keys[security.KEY_IDENTITY_ID] != nil {
				customCtx.IdentityId = c.Keys[security.KEY_IDENTITY_ID].(string)
			}

			ctx = common.WithCustomContext(ctx, customCtx)

			var request dto.MailRequest

			if err := c.BindJSON(&request); err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
				return
			}

			span.LogFields(tracingLog.Object("request", request))

			uniqueInternalIdentifier := utils.GenerateRandomString(64)
			request.UniqueInternalIdentifier = &uniqueInternalIdentifier

			footer := `
					<div>
						<div style="font-size: 12px; font-weight: normal; font-family: Barlow, sans-serif; color: rgb(102, 112, 133); line-height: 32px;">
							<img width="16px" src="https://customer-os.imgix.net/website/favicon.png" alt="CustomerOS" style="vertical-align: middle; margin-right: 5px; margin-bottom: 2px;" />
							Sent from <a href="https://customeros.ai/?utm_content=timeline_email&utm_medium=email" style="text-decoration: underline; color: rgb(102, 112, 133);">CustomerOS</a>
						</div>
					</div>
					`
			request.Content += footer

			// Append an image tag pointing to the spy endpoint to the request content
			imgTag := "<img id=\"customer-os-email-track-open\" height=1 width=1 src=\"" + conf.Service.PublicPath + "/mail/" + uniqueInternalIdentifier + "/track\" />"
			request.Content += imgTag

			replyMail, err := services.CommonServices.MailService.SendMail(ctx, request, &username)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
				return
			}

			interactionEventId, err := services.CommonServices.MailService.SaveMail(ctx, request, replyMail, tenant, username, uniqueInternalIdentifier)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			span.LogFields(tracingLog.String("result - interactionEventId", *interactionEventId))
			c.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprintf("interaction event created with id: %s", *interactionEventId),
			})

		})

	rg.GET("/mail/:customerOSInternalIdentifier/track", func(c *gin.Context) {
		customerOSInternalIdentifier := c.Param("customerOSInternalIdentifier")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(context.Background(), "/mail/"+customerOSInternalIdentifier+"/track", c.Request.Header)
		defer span.Finish()

		if customerOSInternalIdentifier == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing customerOSInternalIdentifier"})
			return
		}

		span.LogFields(tracingLog.String("customerOSInternalIdentifier", customerOSInternalIdentifier))

		//log all headers
		for name, values := range c.Request.Header {
			for _, value := range values {
				span.LogFields(tracingLog.String("Header: "+name, value))
			}
		}

		interactionEventNode, err := services.CommonServices.Neo4jRepositories.InteractionEventReadRepository.GetInteractionEventByCustomerOSIdentifier(ctx, customerOSInternalIdentifier)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if interactionEventNode == nil {
			span.LogFields(tracingLog.String("interactionEventId", "not found"))
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		interactionEvent := neo4jmapper.MapDbNodeToInteractionEventEntity(interactionEventNode)

		span.LogFields(tracingLog.String("interactionEventId", interactionEvent.Id))

		tenant := model.GetTenantFromLabels(interactionEventNode.Labels, model.NodeLabelInteractionEvent)
		if tenant == "" {
			span.LogFields(tracingLog.String("tenant", "not identified"))
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

		_, err = services.CommonServices.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, interactionEvent.Id, model.INTERACTION_EVENT, neo4jenum.ActionInteractionEventRead, "", metadata, utils.Now(), "user-admin-api")
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		c.Data(http.StatusOK, "image/png", spyPixelBytes)
	})
}
