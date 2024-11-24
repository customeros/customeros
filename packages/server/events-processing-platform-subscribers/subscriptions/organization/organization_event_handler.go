package organization

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	ai "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	Unknown = "Unknown"
)

type Socials struct {
	Github    string `json:"github,omitempty"`
	Linkedin  string `json:"linkedin,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Youtube   string `json:"youtube,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
}

type organizationEventHandler struct {
	log         logger.Logger
	cfg         *config.Config
	caches      caches.Cache
	aiModel     ai.AiModel
	grpcClients *grpc_client.Clients
	services    *service.Services
}

func NewOrganizationEventHandler(services *service.Services, log logger.Logger, cfg *config.Config, caches caches.Cache, aiModel ai.AiModel, grpcClients *grpc_client.Clients) *organizationEventHandler {
	return &organizationEventHandler{
		log:         log,
		cfg:         cfg,
		caches:      caches,
		aiModel:     aiModel,
		grpcClients: grpcClients,
		services:    services,
	}
}

func (h *organizationEventHandler) EnrichOrganizationByRequest(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.EnrichOrganizationByRequest")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationRequestEnrich
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    eventData.Tenant,
		AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
	})

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	domain, _ := h.services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(innerCtx, eventData.Website)
	if domain == "" {
		return nil
	}

	return h.enrichOrganization(innerCtx, eventData.Tenant, organizationId, domain)
}

func (h *organizationEventHandler) enrichOrganization(ctx context.Context, tenant, organizationId, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.enrichOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	if domain == "" {
		tracing.TraceErr(span, errors.New("domain is empty"))
		return nil
	}

	// check if domain is primary
	domainEntity, err := h.services.CommonServices.DomainService.GetDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get domain"))
		h.log.Errorf("Error getting domain %s: %s", domain, err.Error())
		return nil
	}
	if domainEntity.IsPrimary == nil || !*domainEntity.IsPrimary {
		h.log.Infof("Domain %s is not primary", domain)
		return nil
	}

	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization with id %s: %v", organizationId, err)
		return nil
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	if organizationEntity.EnrichDetails.EnrichedAt != nil {
		h.log.Infof("Organization %s already enriched", organizationId)
		return nil
	}

	err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4jentity.OrganizationPropertyEnrichRequestedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich requested at"))
	}
	utils.EventCompleted(ctx, tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	enrichOrganizationResponse, err := h.callApiEnrichOrganization(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call enrich organization API"))
		h.log.Errorf("Error calling enrich organization API: %s", err.Error())
		err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4jentity.OrganizationPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
		}
		return nil
	}
	if enrichOrganizationResponse != nil && enrichOrganizationResponse.Success == true {
		h.updateOrganizationWithEnrichData(ctx, tenant, domain, enrichOrganizationResponse.PrimaryEnrichSource, *organizationEntity, &enrichOrganizationResponse.Data)
	} else {
		err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4jentity.OrganizationPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
		}
	}

	return nil
}

func (h *organizationEventHandler) callApiEnrichOrganization(ctx context.Context, tenant, domain string) (*enrichmentmodel.EnrichOrganizationResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.callApiEnrichOrganization")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain)

	requestJSON, err := json.Marshal(enrichmentmodel.EnrichOrganizationRequest{
		Domain: domain,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequestWithContext(ctx, "GET", h.cfg.Services.EnrichmentApi.Url+"/enrichOrganization", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, h.cfg.Services.EnrichmentApi.ApiKey)
	req.Header.Set(security.TenantHeader, tenant)

	// Make the HTTP request, retry once if response status is 502
	var response *http.Response
	client := &http.Client{}

	for attempt := 1; attempt <= 2; attempt++ {
		// Make the HTTP request
		response, err = client.Do(req)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
			return nil, err
		}
		defer response.Body.Close() // Ensures the body is closed only once

		// Retry on 502 and 400
		if response.StatusCode == http.StatusBadGateway || response.StatusCode == http.StatusBadRequest {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}

	if response == nil {
		tracing.TraceErr(span, errors.New("Enrich organization response is nil"))
		return nil, errors.New("Enrich organization response is nil")
	}

	span.LogFields(log.Int("response.statusCode", response.StatusCode))

	if response.StatusCode != http.StatusOK {
		h.log.Errorf("Enrich organization API response status is : %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		span.LogFields(log.String("response.body", string(body)))
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return nil, err
	}

	var enrichOrganizationApiResponse enrichmentmodel.EnrichOrganizationResponse
	// read the response body
	err = json.Unmarshal(body, &enrichOrganizationApiResponse)
	if err != nil {
		span.LogFields(log.String("response.body", string(body)))
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal enrich organization response"))
		return nil, err
	}
	return &enrichOrganizationApiResponse, nil
}

func (h *organizationEventHandler) updateOrganizationWithEnrichData(ctx context.Context, tenant, domain, enrichSource string, organizationEntity neo4jentity.OrganizationEntity, data *enrichmentmodel.EnrichOrganizationResponseData) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.updateOrganizationWithEnrichData")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "data", data)
	tracing.TagTenant(span, tenant)

	orgFields := data_fields.OrganizationFields{
		Source:       utils.StringPtr(constants.SourceOpenline),
		EnrichDomain: utils.StringPtr(domain),
		EnrichSource: utils.StringPtr(enrichSource),
	}

	if organizationEntity.Employees == 0 && data.Employees > 0 {
		orgFields.Employees = utils.Int64Ptr(data.Employees)
	}
	if (organizationEntity.YearFounded == nil || *organizationEntity.YearFounded < 1000) && data.FoundedYear > 0 {
		orgFields.YearFounded = utils.Int64Ptr(data.FoundedYear)
	}
	if organizationEntity.ValueProposition == "" && data.ShortDescription != "" {
		orgFields.ValueProposition = utils.StringPtr(data.ShortDescription)
	}
	if organizationEntity.Description == "" && data.LongDescription != "" {
		orgFields.Description = utils.StringPtr(data.LongDescription)
	}
	if data.Public != nil {
		orgFields.IsPublic = data.Public
	}

	// Set organization name
	if data.Name != "" {
		orgFields.Name = utils.StringPtr(data.Name)
	} else if organizationEntity.Name == "" {
		if data.Domain != "" {
			domainPrefixCapitalized := utils.CapitalizeAllParts(utils.GetDomainWithoutTLD(data.Domain), []string{"-", "_", "."})
			orgFields.Name = utils.StringPtr(domainPrefixCapitalized)
		}
	}

	// Set company website
	if organizationEntity.Website == "" {
		if data.Website != "" {
			orgFields.Website = utils.StringPtr(data.Website)
		} else if data.Domain != "" {
			orgFields.Website = utils.StringPtr(data.Domain)
		}
	}

	// Set company logo and icon urls
	if organizationEntity.LogoUrl == "" && len(data.Logos) > 0 {
		orgFields.LogoUrl = utils.StringPtr(data.Logos[0])
	}
	if organizationEntity.IconUrl == "" && len(data.Icons) > 0 {
		orgFields.IconUrl = utils.StringPtr(data.Icons[0])
	}

	// set industry
	if organizationEntity.Industry == "" && data.Industry != "" {
		orgFields.Industry = utils.StringPtr(data.Industry)
	}

	_, err := h.services.CommonServices.OrganizationService.Save(ctx, nil, &organizationEntity.ID, orgFields)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error updaing organization with enrich data: %s", err.Error())
	}

	//add location
	if !data.Location.IsEmpty() {
		_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*locationpb.LocationIdGrpcResponse](func() (*locationpb.LocationIdGrpcResponse, error) {
			return h.grpcClients.OrganizationClient.AddLocation(ctx, &organizationpb.OrganizationAddLocationGrpcRequest{
				Tenant:         tenant,
				OrganizationId: organizationEntity.ID,
				LocationDetails: &locationpb.LocationDetails{
					Country:       data.Location.Country,
					CountryCodeA2: data.Location.CountryCodeA2,
					CountryCodeA3: data.Location.CountryCodeA3,
					Locality:      data.Location.Locality,
					Region:        data.Location.Region,
					PostalCode:    data.Location.PostalCode,
					AddressLine1:  data.Location.AddressLine1,
					AddressLine2:  data.Location.AddressLine2,
				},
				SourceFields: &commonpb.SourceFields{
					AppSource: constants.AppEnrichment,
					Source:    constants.SourceOpenline,
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	//add socials
	for _, social := range data.Socials {
		h.addSocial(ctx, organizationEntity.ID, tenant, social.Url, social.Alias, social.Id, constants.AppEnrichment)
	}
}

func (h *organizationEventHandler) addSocial(ctx context.Context, organizationId, tenant, url, alias, externalId, appSource string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.addSocial")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, organizationId)
	span.LogKV("url", url, "alias", alias, "externalId", externalId, "appSource", appSource)

	socialEntity := neo4jentity.SocialEntity{
		Url:        url,
		Alias:      alias,
		ExternalId: externalId,
		AppSource:  appSource,
		Source:     neo4jentity.DataSourceOpenline,
	}

	_, err := h.services.CommonServices.SocialService.AddSocialToEntity(ctx, nil, commonservice.LinkWith{
		Id:   organizationId,
		Type: commonmodel.ORGANIZATION,
	}, socialEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error adding %s social: %s", url, err.Error())
	}
}

func (h *organizationEventHandler) saveOrganizationIndustryAndMarket(ctx context.Context, tenant, organizationId, market, industry string, updateMarket, updateIndustry bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.saveOrganizationIndustryAndMarket")
	defer span.Finish()

	if !updateMarket && !updateIndustry {
		h.log.Infof("No need to update organization %s", organizationId)
		return nil
	}

	//delay to avoid updating organization before main event
	time.Sleep(250 * time.Millisecond)

	_, err := h.services.CommonServices.OrganizationService.Save(ctx, nil, &organizationId, data_fields.OrganizationFields{
		Market:   utils.StringPtr(market),
		Industry: utils.StringPtr(industry),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error updating organization %s: %s", organizationId, err.Error())
		return err
	}
	return nil
}

func (h *organizationEventHandler) mapMarketValue(inputMarket string) string {
	return data.AdjustOrganizationMarket(inputMarket)
}

func (h *organizationEventHandler) mapIndustryToGICS(ctx context.Context, tenant, orgId, inputIndustry string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.mapIndustryToGICS")
	defer span.Finish()
	span.LogFields(log.String("inputIndustry", inputIndustry))
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, orgId)

	trimmedInputIndustry := strings.TrimSpace(inputIndustry)

	if trimmedInputIndustry == "" {
		return ""
	}

	var industry = trimmedInputIndustry
	if industryValue, ok := h.caches.GetIndustry(trimmedInputIndustry); ok {
		span.LogFields(log.Bool("result.industryFoundInCache", true))
		span.LogFields(log.String("result.cacheMapping", industryValue))
		industry = industryValue
	} else {
		span.LogFields(log.Bool("result.industryFoundInCache", false))
		h.log.Infof("Industry %s not found in cache, asking AI", trimmedInputIndustry)
		industry = h.mapIndustryToGICSWithAI(ctx, tenant, orgId, trimmedInputIndustry)
		if utils.Contains(data.GICSIndustryValues, industry) {
			h.caches.SetIndustry(trimmedInputIndustry, industry)
			span.LogFields(log.String("result.newMapping", industry))
		} else {
			industry = trimmedInputIndustry
		}
	}

	return strings.TrimSpace(industry)
}

func (h *organizationEventHandler) mapIndustryToGICSWithAI(ctx context.Context, tenant, orgId, inputIndustry string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.mapIndustryToGICSWithAI")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("inputIndustry", inputIndustry))

	firstPrompt := fmt.Sprintf(h.cfg.Services.Anthropic.IndustryLookupPrompt1, inputIndustry)

	promptLog1 := postgresEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		Provider:       constants.Anthropic,
		Model:          "claude-2",
		PromptType:     constants.PromptType_MapIndustry,
		Tenant:         &tenant,
		NodeId:         &orgId,
		NodeLabel:      utils.StringPtr(commonmodel.NodeLabelOrganization),
		PromptTemplate: &h.cfg.Services.Anthropic.IndustryLookupPrompt1,
		Prompt:         firstPrompt,
	}
	promptStoreLogId1, err := h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.Store(promptLog1)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store prompt log"))
		h.log.Errorf("Error storing prompt log: %v", err)
	} else {
		span.LogFields(log.String("promptStoreLogId1", promptStoreLogId1))
	}

	firstResult, err := h.aiModel.Inference(ctx, firstPrompt) // ai.InvokeAnthropic(ctx, h.cfg, h.log, firstPrompt)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to invoke AI for first prompt"))
		h.log.Errorf("Error invoking AI: %v", err)
		storeErr := h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId1, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with error"))
			h.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return ""
	} else {
		storeErr := h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId1, firstResult)
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with ai response"))
			h.log.Errorf("Error updating prompt log with ai response: %v", storeErr)
		}
	}
	if firstResult == "" || firstResult == Unknown {
		return firstResult
	}
	secondPrompt := fmt.Sprintf(h.cfg.Services.Anthropic.IndustryLookupPrompt2, firstResult)

	promptLog2 := postgresEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		Provider:       constants.Anthropic,
		Model:          "claude-2",
		PromptType:     constants.PromptType_ExtractIndustryValue,
		Tenant:         &tenant,
		NodeId:         &orgId,
		NodeLabel:      utils.StringPtr(commonmodel.NodeLabelOrganization),
		PromptTemplate: &h.cfg.Services.Anthropic.IndustryLookupPrompt2,
		Prompt:         secondPrompt,
	}
	promptStoreLogId2, err := h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.Store(promptLog2)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store prompt log"))
		h.log.Errorf("Error storing prompt log with error: %v", err)
	}
	secondResult, err := h.aiModel.Inference(ctx, secondPrompt) // ai.InvokeAnthropic(ctx, h.cfg, h.log, secondPrompt)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to invoke AI for second prompt"))
		h.log.Errorf("Error invoking AI: %v", err)
		err = h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId2, err.Error())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update prompt log with error"))
			h.log.Errorf("Error updating prompt log with error: %v", err)
		}
		return ""
	} else {
		err = h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId2, secondResult)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update prompt log with ai response"))
			h.log.Errorf("Error updating prompt log with ai response: %v", err)
		}
	}
	return secondResult
}

func (h *organizationEventHandler) OnAdjustIndustry(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnAdjustIndustry")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationAdjustIndustryEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    eventData.Tenant,
		AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
	})

	orgDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(innerCtx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization with id %s: %v", organizationId, err)
		return err
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)

	industry := h.mapIndustryToGICS(innerCtx, eventData.Tenant, organizationId, organizationEntity.Industry)

	if industry != "" && organizationEntity.Industry != industry {
		_, err = h.services.CommonServices.OrganizationService.Save(innerCtx, nil, &organizationId, data_fields.OrganizationFields{
			Industry: utils.StringPtr(industry),
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error updating organization %s: %s", organizationId, err.Error())
			return err
		}
	}
	return nil
}
