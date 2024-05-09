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

		if trackingInformation.Identity.SessionId == "" || trackingInformation.Identity.Identifier.Email == "" {
			return
		}

		domain := commonUtils.ExtractDomain(trackingInformation.Identity.Identifier.Email)

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

			contactExists, err := services.CommonServices.Neo4jRepositories.ContactReadRepository.ContactExistsInOrganizationByEmail(ctx, *tenant, organizationId, trackingInformation.Identity.Identifier.Email)
			if err != nil {
				return
			}

			if !contactExists {
				contactId, err = services.CustomerOsClient.CreateContact(*tenant, "", "", "", trackingInformation.Identity.Identifier.Email, nil)
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
			}

			//TODO insert timeline items
		}

		ginContext.JSON(http.StatusOK, gin.H{"tenant": tenant})
	})

}
