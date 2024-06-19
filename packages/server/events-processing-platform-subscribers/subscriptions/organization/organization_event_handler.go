package organization

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"

	ai "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

const (
	Unknown = "Unknown"
)

var nonRetryableErrors = []string{"Invalid Domain Name"}
var knownBrandfetchErrors = []string{"Invalid Domain Name", "User is not authorized to access this resource with an explicit deny"}

type Socials struct {
	Github    string `json:"github,omitempty"`
	Linkedin  string `json:"linkedin,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Youtube   string `json:"youtube,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
}

type BrandfetchResponse struct {
	Message         string            `json:"message,omitempty"`
	Id              string            `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	Domain          string            `json:"domain,omitempty"`
	Claimed         bool              `json:"claimed"`
	Description     string            `json:"description,omitempty"`
	LongDescription string            `json:"longDescription,omitempty"`
	Links           []BrandfetchLink  `json:"links,omitempty"`
	Logos           []BranfetchLogo   `json:"logos,omitempty"`
	QualityScore    float64           `json:"qualityScore,omitempty"`
	Company         BrandfetchCompany `json:"company,omitempty"`
	IsNsfw          bool              `json:"isNsfw"`
}

type BrandfetchLink struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type BranfetchLogo struct {
	Theme   string                `json:"theme,omitempty"`
	Type    string                `json:"type,omitempty"`
	Formats []BranfetchLogoFormat `json:"formats,omitempty"`
}

type BranfetchLogoFormat struct {
	Src        string `json:"src,omitempty"`
	Background string `json:"background,omitempty"`
	Format     string `json:"format,omitempty"`
	Height     int64  `json:"height,omitempty"`
	Width      int64  `json:"width,omitempty"`
	Size       int64  `json:"size,omitempty"`
}

type BrandfetchCompany struct {
	Employees   any                  `json:"employees,omitempty"`
	FoundedYear int64                `json:"foundedYear,omitempty"`
	Industries  []BrandfetchIndustry `json:"industries,omitempty"`
	Kind        string               `json:"kind,omitempty"`
	Location    struct {
		City          string `json:"city,omitempty"`
		Country       string `json:"country,omitempty"`
		CountryCodeA2 string `json:"countryCode,omitempty"`
		Region        string `json:"region,omitempty"`
		State         string `json:"state,omitempty"`
		SubRegion     string `json:"subRegion,omitempty"`
	} `json:"location,omitempty"`
}

func (b *BrandfetchCompany) LocationIsEmpty() bool {
	return b.Location.City == "" && b.Location.Country == "" && b.Location.Region == "" && b.Location.State == "" && b.Location.SubRegion == ""
}

type BrandfetchIndustry struct {
	Score  float64 `json:"score,omitempty"`
	Name   string  `json:"name,omitempty"`
	Emoji  string  `json:"emoji,omitempty"`
	Slug   string  `json:"slug,omitempty"`
	Parent struct {
		Emoji string `json:"emoji,omitempty"`
		Name  string `json:"name,omitempty"`
		Slug  string `json:"slug,omitempty"`
	} `json:"parent,omitempty"`
}

type organizationEventHandler struct {
	repositories *repository.Repositories
	log          logger.Logger
	cfg          *config.Config
	caches       caches.Cache
	aiModel      ai.AiModel
	grpcClients  *grpc_client.Clients
}

func NewOrganizationEventHandler(repositories *repository.Repositories, log logger.Logger, cfg *config.Config, caches caches.Cache, aiModel ai.AiModel, grpcClients *grpc_client.Clients) *organizationEventHandler {
	return &organizationEventHandler{
		repositories: repositories,
		log:          log,
		cfg:          cfg,
		caches:       caches,
		aiModel:      aiModel,
		grpcClients:  grpcClients,
	}
}

func (h *organizationEventHandler) EnrichOrganizationByDomain(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.EnrichOrganizationByDomain")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkDomainEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	return h.enrichOrganization(ctx, eventData.Tenant, organizationId, eventData.Domain)
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
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	if eventData.Website == "" {
		return nil
	}

	domain := utils.ExtractDomain(eventData.Website)
	if domain == "" {
		return nil
	}

	return h.enrichOrganization(ctx, eventData.Tenant, organizationId, domain)

	return nil
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

	organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
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

	// create domain node if not exist
	err = h.repositories.Neo4jRepositories.DomainWriteRepository.MergeDomain(ctx, domain, constants.SourceOpenline, constants.AppSourceEventProcessingPlatformSubscribers, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error creating domain node: %v", err)
		return nil
	}
	domainNode, err := h.repositories.Neo4jRepositories.DomainReadRepository.GetDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting domain node: %v", err)
		return nil
	}
	domainEntity := neo4jmapper.MapDbNodeToDomainEntity(domainNode)

	daysAgo10 := utils.Now().Add(-time.Hour * 24 * 10)
	daysAgo365 := utils.Now().Add(-time.Hour * 24 * 365)

	// if domain is not enriched
	// or last enrich attempt was more than 30 days ago,
	// or last enrich was more than 365 days ago
	// enrich it
	justEnriched := false
	if (domainEntity.EnrichDetails.EnrichedAt == nil && (domainEntity.EnrichDetails.EnrichRequestedAt == nil || domainEntity.EnrichDetails.EnrichRequestedAt.Before(daysAgo10))) ||
		(domainEntity.EnrichDetails.EnrichedAt != nil && domainEntity.EnrichDetails.EnrichedAt.Before(daysAgo365)) {

		if !utils.Contains(nonRetryableErrors, domainEntity.EnrichDetails.EnrichError) {
			err = h.enrichDomain(ctx, tenant, domain)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error enriching domain: %v", err)
				return nil
			}
			justEnriched = true
		}
	}

	// re-fetch latest domain node
	if justEnriched {
		domainNode, err = h.repositories.Neo4jRepositories.DomainReadRepository.GetDomain(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error getting domain node: %v", err)
			return nil
		}
		domainEntity = neo4jmapper.MapDbNodeToDomainEntity(domainNode)
	}

	// Convert enrich data to struct
	var brandfetchResponse BrandfetchResponse
	if domainEntity.EnrichDetails.EnrichSource == neo4jenum.Brandfetch && domainEntity.EnrichDetails.EnrichData != "" {
		err = json.Unmarshal([]byte(domainEntity.EnrichDetails.EnrichData), &brandfetchResponse)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error unmarshalling brandfetch enrich data: %v", err)
			return nil
		}
		h.updateOrganizationFromBrandfetch(ctx, tenant, domain, *organizationEntity, brandfetchResponse)
	}

	return nil
}

func (h *organizationEventHandler) enrichDomain(ctx context.Context, tenant, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.enrichDomain")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("domain", domain))

	brandfetchUrl := h.cfg.Services.BrandfetchApi

	if brandfetchUrl == "" {
		err := errors.New("Brandfetch URL not set")
		tracing.TraceErr(span, err)
		h.log.Errorf("Brandfetch URL not set")
		return err
	}

	// get current month in format yyyy-mm
	currentMonth := utils.Now().Format("2006-01")

	queryResult := h.repositories.PostgresRepositories.ExternalAppKeysRepository.GetAppKeys(ctx, constants.AppBrandfetch, currentMonth, h.cfg.Services.BrandfetchLimit)
	if queryResult.Error != nil {
		tracing.TraceErr(span, queryResult.Error)
		h.log.Errorf("Error getting brandfetch app keys: %v", queryResult.Error)
		return queryResult.Error
	}
	branfetchAppKeys := queryResult.Result.([]postgresEntity.ExternalAppKeys)
	if len(branfetchAppKeys) == 0 {
		err := errors.New(fmt.Sprintf("no brandfetch app keys available for %s", currentMonth))
		tracing.TraceErr(span, err)
		h.log.Errorf("No brandfetch app keys available for %s", currentMonth)
		return err
	}
	// pick random app key from list
	appKey := branfetchAppKeys[rand.Intn(len(branfetchAppKeys))]

	body, err := makeBrandfetchHTTPRequest(brandfetchUrl, appKey.AppKey, domain)

	enrichFailed := false
	errMsg := ""
	if err != nil {
		enrichFailed = true
		errMsg = err.Error()
		tracing.TraceErr(span, err)
		h.log.Errorf("Error making Brandfetch HTTP request: %v", err)
	}

	// Increment usage count of the app key
	queryResult = h.repositories.PostgresRepositories.ExternalAppKeysRepository.IncrementUsageCount(ctx, appKey.ID)
	if queryResult.Error != nil {
		tracing.TraceErr(span, queryResult.Error)
		h.log.Errorf("Error incrementing app key usage count: %v", queryResult.Error)
	}

	var brandfetchResponse BrandfetchResponse
	err = json.Unmarshal(body, &brandfetchResponse)
	if err != nil {
		enrichFailed = true
		errMsg = err.Error()
		tracing.TraceErr(span, err)
		h.log.Errorf("Error unmarshalling brandfetch response: %v", err)
	}

	if utils.Contains(knownBrandfetchErrors, brandfetchResponse.Message) {
		enrichFailed = true
		errMsg = brandfetchResponse.Message
	}

	if enrichFailed {
		innerErr := h.repositories.Neo4jRepositories.DomainWriteRepository.EnrichFailed(ctx, domain, errMsg, neo4jenum.Brandfetch, utils.Now())
		if innerErr != nil {
			tracing.TraceErr(span, innerErr)
			h.log.Errorf("Error saving enriching domain results: %v", innerErr.Error())
		}
		return nil
	}

	bodyAsString := string(body)
	err = h.repositories.Neo4jRepositories.DomainWriteRepository.EnrichSuccess(ctx, domain, bodyAsString, neo4jenum.Brandfetch, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error saving enriching domain results: %v", err.Error())
		return err
	}
	return nil
}

func makeBrandfetchHTTPRequest(baseUrl, apiKey, domain string) ([]byte, error) {
	url := baseUrl + "/" + domain

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	return body, err
}

func (h *organizationEventHandler) updateOrganizationFromBrandfetch(ctx context.Context, tenant, domain string, organizationEntity neo4jEntity.OrganizationEntity, brandfetch BrandfetchResponse) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.updateOrganizationFromBrandfetch")
	defer span.Finish()

	organizationFieldsMask := make([]organizationpb.OrganizationMaskField, 0)
	updateGrpcRequest := organizationpb.UpdateOrganizationGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationEntity.ID,
		SourceFields: &commonpb.SourceFields{
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
			Source:    constants.SourceOpenline,
		},
		EnrichDomain: domain,
		EnrichSource: neo4jenum.Brandfetch.String(),
	}
	sEmployees, ok := brandfetch.Company.Employees.(string)
	if ok {
		if brandfetch.Company.Employees != "" {
			employees := int64(0)
			if strings.Contains(sEmployees, "-") {
				// Handle range case
				parts := strings.Split(sEmployees, "-")
				employees, _ = strconv.ParseInt(parts[0], 10, 64)
			} else {
				employees, _ = strconv.ParseInt(sEmployees, 10, 64)
			}
			if employees > 0 {
				organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_EMPLOYEES)
				updateGrpcRequest.Employees = employees
			}
		}
	}
	iEmployees, ok := brandfetch.Company.Employees.(int64)
	if ok {
		if iEmployees > 0 {
			organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_EMPLOYEES)
			updateGrpcRequest.Employees = iEmployees
		}
	}
	fEmployees, ok := brandfetch.Company.Employees.(float64)
	if ok {
		if fEmployees > 0 {
			organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_EMPLOYEES)
			updateGrpcRequest.Employees = int64(fEmployees)
		}
	}

	if brandfetch.Company.FoundedYear > 0 {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_YEAR_FOUNDED)
		updateGrpcRequest.YearFounded = &brandfetch.Company.FoundedYear
	}
	if brandfetch.Description != "" {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_VALUE_PROPOSITION)
		updateGrpcRequest.ValueProposition = brandfetch.Description
	}
	if brandfetch.LongDescription != "" {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_DESCRIPTION)
		updateGrpcRequest.Description = brandfetch.LongDescription
	}

	// Set headquarters
	if !brandfetch.Company.LocationIsEmpty() {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_HEADQUARTERS)
		headquarters := brandfetch.Company.Location.City
		if brandfetch.Company.Location.State != "" {
			headquarters += ", " + brandfetch.Company.Location.State
		}
		if brandfetch.Company.Location.Country != "" {
			headquarters += ", " + brandfetch.Company.Location.Country
		}
		if brandfetch.Company.Location.Region != "" {
			headquarters += ", " + brandfetch.Company.Location.Region
		}
		if brandfetch.Company.Location.SubRegion != "" {
			headquarters += ", " + brandfetch.Company.Location.SubRegion
		}
		updateGrpcRequest.Headquarters = brandfetch.Company.Location.City + ", " + brandfetch.Company.Location.Country
	}

	// Set public indicator
	if brandfetch.Company.Kind != "" {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_IS_PUBLIC)
		if brandfetch.Company.Kind == "PUBLIC_COMPANY" {
			updateGrpcRequest.IsPublic = true
		} else {
			updateGrpcRequest.IsPublic = false
		}
	}

	// Set company name
	if brandfetch.Name != "" {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_NAME)
		updateGrpcRequest.Name = brandfetch.Name
	} else if brandfetch.Domain != "" {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_NAME)
		updateGrpcRequest.Name = brandfetch.Domain
	}

	if brandfetch.Domain != "" && organizationEntity.Website == "" {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_WEBSITE)
		updateGrpcRequest.Website = brandfetch.Domain
	}

	// Set company logo and icon urls
	logoUrl := ""
	iconUrl := ""
	if len(brandfetch.Logos) > 0 {
		for _, logo := range brandfetch.Logos {
			if logo.Type == "icon" {
				iconUrl = logo.Formats[0].Src
			} else if logo.Type == "symbol" && iconUrl == "" {
				iconUrl = logo.Formats[0].Src
			} else if logo.Type == "logo" {
				logoUrl = logo.Formats[0].Src
			} else if logo.Type == "other" && logoUrl == "" {
				logoUrl = logo.Formats[0].Src
			}
		}
	}
	if logoUrl != "" {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_LOGO_URL)
		updateGrpcRequest.LogoUrl = logoUrl
	}
	if iconUrl != "" {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_ICON_URL)
		updateGrpcRequest.IconUrl = iconUrl
	}

	// set industry
	industryName := ""
	industryMaxScore := float64(0)
	if len(brandfetch.Company.Industries) > 0 {
		for _, industry := range brandfetch.Company.Industries {
			if industry.Name != "" && industry.Score > industryMaxScore {
				industryName = industry.Name
				industryMaxScore = industry.Score
			}
		}
	}
	if industryName != "" {
		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_INDUSTRY)
		updateGrpcRequest.Industry = industryName
	}

	updateGrpcRequest.FieldsMask = organizationFieldsMask
	tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.UpdateOrganization(ctx, &updateGrpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error updating organization: %s", err.Error())
	}

	for _, link := range brandfetch.Links {
		h.addSocial(ctx, organizationEntity.ID, tenant, link.Url)
	}
}

func (h *organizationEventHandler) addSocial(ctx context.Context, organizationId, tenant, url string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.addSocial")
	defer span.Finish()
	span.LogFields(log.String("organizationId", organizationId), log.String("tenant", tenant), log.String("url", url))
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.AddSocial(ctx, &organizationpb.AddSocialGrpcRequest{
			Tenant:         tenant,
			OrganizationId: organizationId,
			SourceFields: &commonpb.SourceFields{
				AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
				Source:    constants.SourceOpenline,
			},
			Url: url,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error adding %s social: %s", url, err.Error())
	}
}

func (h *organizationEventHandler) AdjustNewOrganizationFields(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.AdjustNewOrganizationFields")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	market := h.mapMarketValue(eventData.Market)
	industry := h.mapIndustryToGICS(ctx, eventData.Tenant, organizationId, eventData.Industry)

	// wait for organization to be created in neo4j before updating it
	for attempt := 1; attempt <= constants.MaxRetriesCheckDataInNeo4j; attempt++ {
		exists, err := h.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, eventData.Tenant, organizationId, neo4jutil.NodeLabelOrganization)
		if err == nil && exists {
			break
		}
		time.Sleep(utils.BackOffExponentialDelay(attempt))
	}

	if eventData.Market != market || eventData.Industry != industry {
		err := h.callUpdateOrganizationCommand(ctx, eventData.Tenant, organizationId, eventData.SourceOfTruth, market, industry, span)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		h.log.Infof("No need to update organization %s", organizationId)
	}
	return nil
}

func (h *organizationEventHandler) AdjustUpdatedOrganizationFields(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.AdjustUpdatedOrganizationFields")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))
	tracing.LogObjectAsJson(span, "eventData", evt)

	var eventData events.OrganizationUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	market := h.mapMarketValue(eventData.Market)
	industry := h.mapIndustryToGICS(ctx, eventData.Tenant, organizationId, eventData.Industry)

	if eventData.Market != market || eventData.Industry != industry {
		err := h.callUpdateOrganizationCommand(ctx, eventData.Tenant, organizationId, eventData.Source, market, industry, span)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		h.log.Infof("No need to update organization %s", organizationId)
	}
	return nil
}

func (h *organizationEventHandler) callUpdateOrganizationCommand(ctx context.Context, tenant, organizationId, source, market, industry string, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.UpdateOrganization(ctx, &organizationpb.UpdateOrganizationGrpcRequest{
			Tenant:         tenant,
			OrganizationId: organizationId,
			SourceFields: &commonpb.SourceFields{
				AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
				Source:    source,
			},
			Market:   market,
			Industry: industry,
			FieldsMask: []organizationpb.OrganizationMaskField{
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_MARKET,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_INDUSTRY,
			},
		})
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
	trimmedInputIndustry := strings.TrimSpace(inputIndustry)

	if inputIndustry == "" {
		return ""
	}

	var industry string
	if industryValue, ok := h.caches.GetIndustry(trimmedInputIndustry); ok {
		industry = industryValue
	} else {
		h.log.Infof("Industry %s not found in cache, asking AI", trimmedInputIndustry)
		industry = h.mapIndustryToGICSWithAI(ctx, tenant, orgId, trimmedInputIndustry)
		if industry != "" && len(industry) < 45 {
			h.caches.SetIndustry(trimmedInputIndustry, industry)
			h.log.Infof("Industry %s mapped to %s", trimmedInputIndustry, industry)
		} else {
			h.log.Warnf("Industry %s mapped wrongly to (%s) with AI, returning input value", industry, trimmedInputIndustry)
			return trimmedInputIndustry
		}
	}
	if industry == Unknown {
		h.log.Infof("Unknown industry %s, returning as is", trimmedInputIndustry)
		return trimmedInputIndustry
	}

	return strings.TrimSpace(industry)
}

func (h *organizationEventHandler) mapIndustryToGICSWithAI(ctx context.Context, tenant, orgId, inputIndustry string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.mapIndustryToGICSWithAI")
	defer span.Finish()
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
		NodeLabel:      utils.StringPtr(neo4jutil.NodeLabelOrganization),
		PromptTemplate: &h.cfg.Services.Anthropic.IndustryLookupPrompt1,
		Prompt:         firstPrompt,
	}
	promptStoreLogId1, err := h.repositories.PostgresRepositories.AiPromptLogRepository.Store(promptLog1)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error storing prompt log: %v", err)
	} else {
		span.LogFields(log.String("promptStoreLogId1", promptStoreLogId1))
	}

	firstResult, err := h.aiModel.Inference(ctx, firstPrompt) // ai.InvokeAnthropic(ctx, h.cfg, h.log, firstPrompt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error invoking AI: %v", err)
		storeErr := h.repositories.PostgresRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId1, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
			h.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return ""
	} else {
		storeErr := h.repositories.PostgresRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId1, firstResult)
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
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
		NodeLabel:      utils.StringPtr(neo4jutil.NodeLabelOrganization),
		PromptTemplate: &h.cfg.Services.Anthropic.IndustryLookupPrompt2,
		Prompt:         secondPrompt,
	}
	promptStoreLogId2, err := h.repositories.PostgresRepositories.AiPromptLogRepository.Store(promptLog2)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error storing prompt log with error: %v", err)
	}
	secondResult, err := h.aiModel.Inference(ctx, secondPrompt) // ai.InvokeAnthropic(ctx, h.cfg, h.log, secondPrompt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error invoking AI: %v", err)
		err = h.repositories.PostgresRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId2, err.Error())
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error updating prompt log with error: %v", err)
		}
		return ""
	} else {
		err = h.repositories.PostgresRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId2, secondResult)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error updating prompt log with ai response: %v", err)
		}
	}
	return secondResult
}
