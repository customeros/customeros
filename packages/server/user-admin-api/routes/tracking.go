package routes

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

type TrackingInformation struct {
	Identity struct {
		SessionId  string `json:"sessionId"`
		Identifier struct {
			Email string `json:"email"`
		} `json:"identifier"`
	} `json:"identity"`
	Activity string `json:"activity"`
}

func addTrackingRoutes(rg *gin.RouterGroup, services *service.Services) {
	rg.GET("/track", func(ginContext *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(context.Background(), "GET /track", ginContext.Request.Header)
		defer span.Finish()

		origin := ginContext.GetHeader("Origin")
		referer := ginContext.GetHeader("Referer")
		userAgent := ginContext.GetHeader("User-Agent")
		trackPayload := ginContext.GetHeader("X-Tracker-Payload")

		span.LogFields(tracingLog.String("origin", origin))
		span.LogFields(tracingLog.String("referer", referer))
		span.LogFields(tracingLog.String("userAgent", userAgent))
		span.LogFields(tracingLog.String("trackPayload", trackPayload))

		if origin == "" || referer == "" || userAgent == "" || trackPayload == "" {
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		tenant, err := services.CommonServices.PostgresRepositories.TrackingAllowedOriginRepository.GetTenantForOrigin(ctx, origin)
		if err != nil {
			tracing.TraceErr(span, err)
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		if tenant == nil || *tenant == "" {
			span.LogFields(tracingLog.String("result.info", "tenant not found for origin"))
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		trackingInformation := TrackingInformation{}
		err = json.Unmarshal([]byte(trackPayload), &trackingInformation)
		if err != nil {
			tracing.TraceErr(span, err)
			span.LogFields(tracingLog.String("error", err.Error()))
			return
		}
		span.LogFields(tracingLog.String("result.trackingInformation", trackPayload))

		email := trackingInformation.Identity.Identifier.Email
		if trackingInformation.Identity.SessionId == "" || email == "" {
			span.LogFields(tracingLog.String("result.info", "session id or email not found"))
			return
		}

		organizationId, contactId, err := services.RegistrationService.CreateOrganizationAndContact(ctx, *tenant, email, false, "tracking")
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if organizationId == nil || contactId == nil {
			tracing.TraceErr(span, errors.New("organization id or contact id empty"))
			return
		}

		span.LogFields(tracingLog.String("result.organizationId", *organizationId))
		span.LogFields(tracingLog.String("result.contactId", *contactId))

		if trackingInformation.Activity != "" {
			activityParts := strings.Split(trackingInformation.Activity, ",")

			if len(activityParts) == 0 {
				tracing.TraceErr(span, errors.New("activity parts empty"))
				return
			}

			if len(activityParts)%2 != 0 {
				tracing.TraceErr(span, errors.New("activity parts not even"))
				return
			}

			for i := 0; i < len(activityParts); i += 2 {
				if activityParts[i+1] == "" || activityParts[i] == "" || err != nil {
					tracing.TraceErr(span, errors.New("activity parts empty"))
					return
				}

				_, err := commonUtils.UnmarshalDateTime(activityParts[i])
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
			}

			channelValue := "TRACKING"

			interactionSessionNode, err := services.CommonServices.Neo4jRepositories.InteractionSessionReadRepository.GetByIdentifierAndChannel(ctx, *tenant, trackingInformation.Identity.SessionId, channelValue)
			if err != nil {
				span.LogFields(tracingLog.String("error", err.Error()))
				return
			}

			sessionId := ""
			sessionIdentifier := trackingInformation.Identity.SessionId
			appSource := APP_SOURCE

			if interactionSessionNode == nil {
				span.LogFields(tracingLog.String("result.info", "creating new interaction session"))

				sessionName := "Tracking activity"
				sessionStatus := "ACTIVE"
				sessionType := "THREAD"
				sessionOpts := []service.InteractionSessionBuilderOption{
					service.WithSessionIdentifier(&sessionIdentifier),
					service.WithSessionChannel(&channelValue),
					service.WithSessionName(&sessionName),
					service.WithSessionAppSource(&appSource),
					service.WithSessionStatus(&sessionStatus),
					service.WithSessionType(&sessionType),
				}

				sessionIdResponse, err := services.CustomerOsClient.CreateInteractionSession(*tenant, "", sessionOpts...)
				if err != nil {
					span.LogFields(tracingLog.String("error", err.Error()))
					return
				}
				sessionId = *sessionIdResponse
			} else {
				sessionId = mapper.MapDbNodeToInteractionSessionEntity(interactionSessionNode).Id
			}

			if sessionId == "" {
				tracing.TraceErr(span, errors.New("session id empty"))
				return
			}
			span.LogFields(tracingLog.String("result.sessionId", sessionId))

			contentType := "text/plain"
			now := time.Now()

			for i := 0; i < len(activityParts); i += 2 {
				//part[0] is timstamp in unix format
				//part[1] is the activity
				activity := activityParts[i+1]
				timestamp, err := commonUtils.UnmarshalDateTime(activityParts[i])
				if timestamp == nil {
					timestamp = &now
				}
				utc := timestamp.UTC()

				eventOpts := []service.InteractionEventBuilderOption{
					service.WithSessionId(&sessionId),
					service.WithChannel(&channelValue),
					service.WithCreatedAt(&utc),
					service.WithContent(&activity),
					service.WithContentType(&contentType),
					service.WithSentBy([]model.InteractionEventParticipantInput{
						{
							ContactID: contactId,
						},
					}),
					service.WithAppSource(&appSource),
				}

				interactionEventId, err := services.CustomerOsClient.CreateInteractionEvent(*tenant, "", eventOpts...)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
				span.LogFields(tracingLog.String("interactionEventId", *interactionEventId))
			}

		} else {
			tracing.TraceErr(span, errors.New("activity not found"))
		}

		ginContext.JSON(http.StatusOK, gin.H{})
	})

}
