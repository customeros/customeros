package organization

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
)

type WebscrapeRequest struct {
	Domain string `json:"scrape"`
}

type WebscrapeResponseV1 struct {
	CompanyName      string `json:"companyName,omitempty"`
	Market           string `json:"market,omitempty"`
	Industry         string `json:"industry,omitempty"`
	IndustryGroup    string `json:"industryGroup,omitempty"`
	SubIndustry      string `json:"subIndustry,omitempty"`
	TargetAudience   string `json:"targetAudience,omitempty"`
	ValueProposition string `json:"valueProposition,omitempty"`
	Github           string `json:"github,omitempty"`
	Linkedin         string `json:"linkedin,omitempty"`
	Twitter          string `json:"twitter,omitempty"`
	Youtube          string `json:"youtube,omitempty"`
	Instagram        string `json:"instagram,omitempty"`
	Facebook         string `json:"facebook,omitempty"`
}

type organizationEventHandler struct {
	repositories         *repository.Repositories
	organizationCommands *commands.OrganizationCommands
	log                  logger.Logger
	cfg                  *config.Config
}

func (h *organizationEventHandler) WebscrapeOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.WebscrapeOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkDomainEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	if eventData.Domain == "" {
		tracing.TraceErr(span, errors.New("domain is empty"))
		h.log.Errorf("Missing domain in event data: %v", eventData)
		return nil
	}

	alreadyWebscraped, err := h.repositories.OrganizationRepository.OrganizationWebscrapedForDomain(ctx, eventData.Tenant, organizationId, eventData.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error checking if organization %s already webscraped: %v", organizationId, err)
		return nil
	}
	if alreadyWebscraped {
		h.log.Infof("Organization %s already webscraped for domain %s", organizationId, eventData.Domain)
		return nil
	}

	webscrapeRequest := WebscrapeRequest{
		Domain: eventData.Domain,
	}

	requestBodyJson, err := json.Marshal(webscrapeRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error marshalling webscrape request: %v", webscrapeRequest)
		return nil
	}
	req, err := http.NewRequest("POST", h.cfg.Services.WebscrapeApi, bytes.NewReader(requestBodyJson))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error creating webscrape request: %v", webscrapeRequest)
		return nil
	}
	// Set the request headers
	req.Header.Set(constants.TenantKeyHeader, h.cfg.Services.WebscrapeApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error making webscrape request: %v", webscrapeRequest)
		return nil
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		tracing.TraceErr(span, errors.New("response status code is not 200"))
		h.log.Errorf("Response code %v received while making webscrape request for %v", response.StatusCode, webscrapeRequest)
		return nil
	}
	var result WebscrapeResponseV1
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error decoding webscrape response: %v", webscrapeRequest)
		return nil
	}

	org, err := h.repositories.OrganizationRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization with id %s: %v", organizationId, err)
		return nil
	}
	currentOrgName := utils.GetStringPropOrEmpty(org.Props, "name")
	currentOrgMarket := utils.GetStringPropOrEmpty(org.Props, "market")

	err = h.organizationCommands.UpdateOrganization.Handle(ctx,
		commands.NewUpdateOrganizationCommand(organizationId, eventData.Tenant, constants.SourceWebscrape,
			models.OrganizationDataFields{
				Name:             h.prepareOrgName(result.CompanyName, currentOrgName),
				Market:           data.AdjustOrganizationMarket(result.Market, currentOrgMarket),
				Industry:         result.Industry,
				IndustryGroup:    result.IndustryGroup,
				SubIndustry:      result.SubIndustry,
				TargetAudience:   result.TargetAudience,
				ValueProposition: result.ValueProposition,
			},
			utils.TimePtr(utils.Now()), true))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error updating organization: %v", err)
		return nil
	}
	if result.Youtube != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "youtube", result.Youtube)
	}
	if result.Twitter != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "twitter", result.Twitter)
	}
	if result.Linkedin != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "linkedin", result.Linkedin)
	}
	if result.Github != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "github", result.Github)
	}
	if result.Instagram != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "instagram", result.Instagram)
	}
	if result.Facebook != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "facebook", result.Facebook)
	}

	return nil
}

func (h *organizationEventHandler) addSocial(ctx context.Context, organizationId, tenant, platform, url string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.addSocial")
	defer span.Finish()
	span.LogFields(log.String("organizationId", organizationId), log.String("tenant", tenant), log.String("platform", platform), log.String("url", url))

	err := h.organizationCommands.AddSocialCommand.Handle(ctx,
		commands.NewAddSocialCommand(organizationId, tenant, uuid.New().String(),
			platform, url, constants.SourceWebscrape, constants.SourceWebscrape, constants.AppSourceEventProcessingPlatform, nil, nil))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error adding %s social: %v", platform, err)
	}
}

func (h *organizationEventHandler) prepareOrgName(webscrabedOrgName, currentOrgName string) string {
	if currentOrgName == "" {
		return webscrabedOrgName
	} else {
		return currentOrgName
	}
}
