package private

import (
	"bytes"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/constants"
	"image"
	"image/color"
	"image/png"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
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
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

		spans, ctx := telemetry.StartRestSpan(ctx, "MailHandler.SendEmail")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)

		var request *postgres_entity.EmailMessage

		if err := c.BindJSON(&request); err != nil {
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		request.Tenant = tenant
		request.ProducerId = "N/A"
		request.ProducerType = "N/A"

		spans.LogObjectAsJson("request", request)

		err := h.services.CommonServices.MailService.SendMail(ctx, request)
		if err != nil {
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		h.responseHandler.HandleSuccess(c, gin.H{})

	}
}

func (h *MailHandler) TrackEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerOSInternalIdentifier := c.Param("customerOSInternalIdentifier")
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "MailHandler.TrackEmail")
		defer spans.Finish()

		// Preload 1px transparent image
		px := image.NewRGBA(image.Rect(0, 0, 1, 1))
		px.Set(0, 0, color.Transparent)
		var spyPixel bytes.Buffer
		err := png.Encode(&spyPixel, px)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "unable to encode image"))
		}
		spyPixelBytes := spyPixel.Bytes()

		if customerOSInternalIdentifier == "" {
			message := "Missing customerOSInternalIdentifier"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		spans.LogKV("customerOSInternalIdentifier", customerOSInternalIdentifier)

		// log all headers
		for name, values := range c.Request.Header {
			for _, value := range values {
				spans.LogKV("Header: "+name, value)
			}
		}

		interactionEventNode, err := h.services.Repositories.Neo4jRepositories.InteractionEventReadRepository.GetInteractionEventByCustomerOSIdentifier(ctx, customerOSInternalIdentifier)
		if err != nil {
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		if interactionEventNode == nil {
			spans.LogKV("interactionEventId", "not found")
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		interactionEvent := mapper.MapDbNodeToInteractionEventEntity(interactionEventNode)

		spans.LogKV("interactionEventId", interactionEvent.Id)

		tenant := model.GetTenantFromLabels(interactionEventNode.Labels, model.NodeLabelInteractionEvent)
		if tenant == "" {
			spans.LogKV("tenant", "not identified")
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		spans.LogKV("tenant", tenant)

		metadata, err := utils.ToJson(map[string]interface{}{
			"User-Agent":       c.GetHeader("User-Agent"),
			"Cf-Connecting-Ip": c.GetHeader("Cf-Connecting-Ip"),
		})
		if err != nil {
			message := "Error while converting metadata to json"
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		_, err = h.services.Repositories.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, interactionEvent.Id, model.INTERACTION_EVENT, enum.ActionInteractionEventRead, "", metadata, utils.Now(), constants.AppSourceCustomerOsApi)
		if err != nil {
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		h.responseHandler.HandleDataStream(c, "image/png", spyPixelBytes)
	}
}
