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
			span.LogFields(tracingLog.String("error", err.Error()))
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		if tenant == nil || *tenant == "" {
			span.LogFields(tracingLog.String("info", "tenant not found for origin"))
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		trackingInformation := TrackingInformation{}
		err = json.Unmarshal([]byte(trackPayload), &trackingInformation)
		if err != nil {
			span.LogFields(tracingLog.String("error", err.Error()))
			return
		}
		span.LogFields(tracingLog.String("trackingInformation", trackPayload))

		email := trackingInformation.Identity.Identifier.Email
		if trackingInformation.Identity.SessionId == "" || email == "" {
			span.LogFields(tracingLog.String("info", "session id or email not found"))
			return
		}

		domain := commonUtils.ExtractDomain(email)

		isPersonalEmail := false
		//check if the user is using a personal email provider
		for _, personalEmailProvider := range services.Cache.GetPersonalEmailProviders() {
			if strings.Contains(domain, personalEmailProvider) {
				isPersonalEmail = true
				break
			}
		}

		if !isPersonalEmail {

			organizationId := ""
			contactId := ""

			organizationByDomain, err := services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationWithDomain(ctx, *tenant, domain)
			if err != nil {
				span.LogFields(tracingLog.String("error", err.Error()))
				return
			}

			if organizationByDomain == nil {
				prospect := model.OrganizationRelationshipProspect
				lead := model.OrganizationStageLead
				leadSource := "tracking"
				organizationId, err = services.CustomerOsClient.CreateOrganization(*tenant, "", model.OrganizationInput{Relationship: &prospect, Stage: &lead, Domains: []string{domain}, LeadSource: &leadSource})
				if err != nil {
					span.LogFields(tracingLog.String("error", err.Error()))
					return
				}
			} else {
				organizationId = mapper.MapDbNodeToOrganizationEntity(organizationByDomain).ID
			}

			if organizationId == "" {
				span.LogFields(tracingLog.String("error", "organization id empty"))
				return
			}
			span.LogFields(tracingLog.String("organizationId", organizationId))

			contactNode, err := services.CommonServices.Neo4jRepositories.ContactReadRepository.GetContactInOrganizationByEmail(ctx, *tenant, organizationId, email)
			if err != nil {
				return
			}

			if contactNode == nil {
				contactInput := model.ContactInput{
					Email: &model.EmailInput{
						Email: email,
					},
					ProfilePhotoURL: nil,
				}

				contactId, err = services.CustomerOsClient.CreateContact(*tenant, "", contactInput)
				if err != nil {
					return
				}

				_, err = services.CustomerOsClient.LinkContactToOrganization(*tenant, contactId, organizationId)
				if err != nil {
					return
				}
			} else {
				contactId = mapper.MapDbNodeToContactEntity(contactNode).Id
			}

			if contactId == "" {
				span.LogFields(tracingLog.String("error", "contact id empty"))
				return
			}
			span.LogFields(tracingLog.String("contactId", contactId))

			if trackingInformation.Activity != "" {

				activityParts := strings.Split(trackingInformation.Activity, ",")

				if len(activityParts)%2 != 0 {
					span.LogFields(tracingLog.String("error", "activity parts not even"))
					return
				}

				for i := 0; i < len(activityParts); i += 2 {
					if activityParts[i+1] == "" || activityParts[i] == "" || err != nil {
						span.LogFields(tracingLog.String("error", "activity parts empty"))
						return
					}

					_, err := commonUtils.UnmarshalDateTime(activityParts[i])
					if err != nil {
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
					span.LogFields(tracingLog.String("info", "creating new interaction session"))

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
					span.LogFields(tracingLog.String("error", "session id empty"))
					return
				}
				span.LogFields(tracingLog.String("sessionId", sessionId))

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
								ContactID: &contactId,
							},
						}),
						service.WithAppSource(&appSource),
					}

					interactionEventId, err := services.CustomerOsClient.CreateInteractionEvent(*tenant, "", eventOpts...)
					if err != nil {
						span.LogFields(tracingLog.String("error", err.Error()))
						return
					}
					span.LogFields(tracingLog.String("interactionEventId", *interactionEventId))
				}

			} else {
				span.LogFields(tracingLog.String("error", "activity not found"))
			}
		}

		ginContext.JSON(http.StatusOK, gin.H{})
	})

}
