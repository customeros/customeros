package routes

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
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
		ctx, cancel := commonUtils.GetMediumLivedContext(context.Background())
		defer cancel()

		origin := ginContext.GetHeader("Origin")
		referer := ginContext.GetHeader("Referer")
		userAgent := ginContext.GetHeader("User-Agent")
		trackPayload := ginContext.GetHeader("X-Tracker-Payload")

		if origin == "" || referer == "" || userAgent == "" || trackPayload == "" {
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		tenant, err := services.CommonServices.PostgresRepositories.TrackingAllowedOriginRepository.GetTenantForOrigin(origin)
		if err != nil {
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		if tenant == nil || *tenant == "" {
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		trackingInformation := TrackingInformation{}
		err = json.Unmarshal([]byte(trackPayload), &trackingInformation)
		if err != nil {
			return
		}

		email := trackingInformation.Identity.Identifier.Email
		if trackingInformation.Identity.SessionId == "" || email == "" {
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
				return
			}

			if organizationByDomain == nil {
				prospect := model.OrganizationRelationshipProspect
				lead := model.OrganizationStageLead
				leadSource := "tracking"
				organizationId, err = services.CustomerOsClient.CreateOrganization(*tenant, "", model.OrganizationInput{Relationship: &prospect, Stage: &lead, Domains: []string{domain}, LeadSource: &leadSource})
				if err != nil {
					return
				}
			} else {
				organizationId = mapper.MapDbNodeToOrganizationEntity(organizationByDomain).ID
			}

			contactNode, err := services.CommonServices.Neo4jRepositories.ContactReadRepository.GetContactInOrganizationByEmail(ctx, *tenant, organizationId, email)
			if err != nil {
				return
			}

			if contactNode == nil {
				contactId, err = services.CustomerOsClient.CreateContact(*tenant, "", "", "", email, nil)
				if err != nil {
					return
				}

				if organizationId == "" || contactId == "" {
					return
				}

				_, err = services.CustomerOsClient.LinkContactToOrganization(*tenant, contactId, organizationId)
				if err != nil {
					return
				}
			} else {
				contactId = mapper.MapDbNodeToContactEntity(contactNode).Id
			}

			if trackingInformation.Activity != "" {

				channelValue := "TRACKING"

				interactionSessionNode, err := services.CommonServices.Neo4jRepositories.InteractionSessionReadRepository.GetByIdentifierAndChannel(ctx, *tenant, trackingInformation.Identity.SessionId, channelValue)
				if err != nil {
					return
				}

				sessionId := ""
				sessionIdentifier := trackingInformation.Identity.SessionId
				appSource := APP_SOURCE

				if interactionSessionNode == nil {
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

					sessionIdResponse, err := services.CustomerOsClient.CreateInteractionSession(*tenant, email, sessionOpts...)
					if err != nil || sessionIdResponse == nil {
						return
					}
					sessionId = *sessionIdResponse
				} else {
					sessionId = mapper.MapDbNodeToInteractionSessionEntity(interactionSessionNode).Id
				}

				activityParts := strings.Split(trackingInformation.Activity, ",")

				contentType := "text/plain"
				now := time.Now()

				if len(activityParts)%2 != 0 {
					return
				}

				for i := 0; i < len(activityParts); i += 2 {

					if activityParts[i+1] == "" || activityParts[i] == "" || err != nil {
						return
					}

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

					_, err = services.CustomerOsClient.CreateInteractionEvent(*tenant, email, eventOpts...)
					if err != nil {
						return
					}

				}

			}
		}

		ginContext.JSON(http.StatusOK, gin.H{"tenant": tenant})
	})

}
