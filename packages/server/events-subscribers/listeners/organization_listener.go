package listeners

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/constants"
)

type OrganizationListener interface {
	enrichOrganization(ctx context.Context, tenant, organizationId, domain string) error
}

type organizationListenerImpl struct {
	services *service.Services
	log      logger.Logger
}

func NewOrganizationListener(services *service.Services, log logger.Logger) OrganizationListener {
	return &organizationListenerImpl{services: services, log: log}
}

func OnRequestedEnrichOrganization(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnRequestedEnrichOrganization")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)
	organizationId := message.Event.EntityId
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		tracing.TraceErr(span, err)
		return nil
	}
	messageData := message.Event.Data.(*dto.RequestEnrichOrganization)

	span.SetTag(tracing.SpanTagEntityId, organizationId)

	if services.GlobalConfig.InternalServices.EnrichmentApiConfig.Url == "" || services.GlobalConfig.InternalServices.EnrichmentApiConfig.ApiKey == "" {
		err := errors.New("enrichment api url or api key is not set")
		tracing.TraceErr(span, err)
		return err
	}

	if services.GlobalConfig.InternalServices.AiApiConfig.Url == "" || services.GlobalConfig.InternalServices.AiApiConfig.ApiKey == "" {
		err := errors.New("ai api url or api key is not set")
		tracing.TraceErr(span, err)
		return err
	}

	l := NewOrganizationListener(services, services.Logger)

	domain, _ := services.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, messageData.Url)
	if domain == "" {
		return nil
	}

	return l.enrichOrganization(ctx, common.GetTenantFromContext(ctx), organizationId, domain)
}

func (l *organizationListenerImpl) enrichOrganization(ctx context.Context, tenant, organizationId, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationListener.enrichOrganization")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, organizationId)

	if domain == "" {
		tracing.TraceErr(span, errors.New("domain is empty"))
		return nil
	}

	// check if domain is primary
	domainEntity, err := l.services.DomainService.GetDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get domain"))
		l.log.Errorf("Error getting domain %s: %s", domain, err.Error())
		return nil
	}
	if domainEntity.IsPrimary == nil || !*domainEntity.IsPrimary {
		l.log.Infof("Domain %s is not primary", domain)
		return nil
	}

	organizationDbNode, err := l.services.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		l.log.Errorf("Error getting organization with id %s: %v", organizationId, err)
		return nil
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	if organizationEntity.EnrichDetails.EnrichedAt != nil {
		l.log.Infof("Organization %s already enriched", organizationId)
		return nil
	}

	err = l.services.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4jentity.OrganizationPropertyEnrichRequestedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich requested at"))
	}

	l.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

	enrichOrganizationResponse, err := l.callApiEnrichOrganization(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call enrich organization API"))
		l.log.Errorf("Error calling enrich organization API: %s", err.Error())
		err = l.services.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4jentity.OrganizationPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
		}
		return nil
	}
	if enrichOrganizationResponse != nil && enrichOrganizationResponse.Success == true {
		l.updateOrganizationWithEnrichData(ctx, tenant, domain, enrichOrganizationResponse.PrimaryEnrichSource, *organizationEntity, enrichOrganizationResponse.Data)
	} else {
		err = l.services.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4jentity.OrganizationPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
		}
	}

	return nil
}

func (l *organizationListenerImpl) callApiEnrichOrganization(ctx context.Context, tenant, domain string) (*enrichmentmodel.EnrichOrganizationResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationListener.callApiEnrichOrganization")
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
	req, err := http.NewRequestWithContext(ctx, "GET", l.services.GlobalConfig.InternalServices.EnrichmentApiConfig.Url+"/enrichOrganization", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, l.services.GlobalConfig.InternalServices.EnrichmentApiConfig.ApiKey)
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
		l.log.Errorf("Enrich organization API response status is : %d", response.StatusCode)
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

func (l *organizationListenerImpl) updateOrganizationWithEnrichData(ctx context.Context, tenant, domain, enrichSource string, organizationEntity neo4jentity.OrganizationEntity, data *enrichmentmodel.EnrichOrganizationResponseData) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationListener.updateOrganizationWithEnrichData")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "data", data)
	tracing.TagTenant(span, tenant)

	orgFields := data_fields.OrganizationFields{
		Source:       utils.StringPtr(neo4jentity.DataSourceOpenline.String()),
		EnrichDomain: utils.StringPtr(domain),
		EnrichSource: utils.StringPtr(enrichSource),
	}

	if organizationEntity.Employees == 0 && data.Employees > 0 {
		orgFields.Employees = utils.Int64Ptr(data.Employees)
	}
	if (organizationEntity.YearFounded == nil || *organizationEntity.YearFounded < 1000) && data.FoundedYear > 0 {
		orgFields.YearFounded = utils.Int64Ptr(data.FoundedYear)
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

	_, err := l.services.OrganizationService.Save(ctx, nil, &organizationEntity.ID, orgFields)
	if err != nil {
		tracing.TraceErr(span, err)
		l.log.Errorf("Error updaing organization with enrich data: %s", err.Error())
	}

	// add location
	if !data.Location.IsEmpty() {
		_, err := l.services.LocationService.Create(ctx, nil, data_fields.LocationFields{
			Country:       data.Location.Country,
			CountryCodeA2: data.Location.CountryCodeA2,
			CountryCodeA3: data.Location.CountryCodeA3,
			Region:        data.Location.Region,
			Locality:      data.Location.Locality,
			PostalCode:    data.Location.PostalCode,
			Address:       data.Location.AddressLine1,
			Address2:      data.Location.AddressLine2,
			AppSource:     utils.StringPtr(constants.AppEnrichment),
			Source:        utils.StringPtr(constants.SourceOpenline),
		},
			&service.LinkWith{
				Id:   organizationEntity.ID,
				Type: commonmodel.ORGANIZATION,
			})
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	// add socials
	for _, social := range data.Socials {
		l.addSocial(ctx, organizationEntity.ID, tenant, social.Url, social.Alias, social.Id, constants.AppEnrichment)
	}
}

func (l *organizationListenerImpl) addSocial(ctx context.Context, organizationId, tenant, url, alias, externalId, appSource string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationListener.addSocial")
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

	_, err := l.services.SocialService.AddSocialToEntity(ctx, nil, service.LinkWith{
		Id:   organizationId,
		Type: commonmodel.ORGANIZATION,
	}, socialEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		l.log.Errorf("Error adding %s social: %s", url, err.Error())
	}
}
