package listeners

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/events-subscribers/constants"
	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type OrganizationListener interface {
	enrichOrganization(ctx context.Context, tenant, organizationId, domain string) error
}

type organizationListenerImpl struct {
	services          *service.CommonServices
	neo4jRepositories *neo4j_repository.Repositories
	log               logger.Logger
	config            *config.CommonConfig
}

func NewOrganizationListener(
	services *service.CommonServices,
	neo4jRepositories *neo4j_repository.Repositories,
	log logger.Logger,
	config *config.CommonConfig,
) OrganizationListener {
	return &organizationListenerImpl{
		log:               log,
		config:            config,
		neo4jRepositories: neo4jRepositories,
		services:          services,
	}
}

func OnRequestedEnrichOrganization(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnRequestedEnrichOrganization")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message, err := validateEvent(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	organizationId := message.Event.EntityId

	messageData, ok := message.Event.Data.(*dto.RequestEnrichOrganization)
	if !ok {
		err := errors.New("could not cast event data to *dto.RequestEnrichOrganization")
		tracing.TraceErr(span, err)
		return err
	}

	span.SetTag(tracing.SpanTagEntityId, organizationId)

	l := NewOrganizationListener(
		dependencies.CommonServices,
		dependencies.Neo4jRepositories,
		dependencies.Logger,
		dependencies.CommonConfig,
	)

	domain, _ := dependencies.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, messageData.Url)
	if domain == "" {
		return nil
	}

	return l.enrichOrganization(ctx, common.GetTenantFromContext(ctx), organizationId, domain)
}

func OnOrganizationCreated(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnRequestedEnrichOrganization")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message, err := validateEvent(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	organizationId := message.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	_, ok := message.Event.Data.(*dto.CreateOrganization)
	if !ok {
		err := errors.New("could not cast event data to *dto.CreateOrganization")
		tracing.TraceErr(span, err)
		return err
	}

	// lookup active ICP Qualification agents

	// Send to agent for processing

	return nil
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

	organizationDbNode, err := l.neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
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

	err = l.neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4jentity.OrganizationPropertyEnrichRequestedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich requested at"))
	}

	l.services.Events.Publisher.PublishEventCompleted(ctx, tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

	enrichOrganizationResponse, err := l.services.EnrichmentService.FetchEnrichOrganizationData(ctx, &domain, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call enrich organization"))
		l.log.Errorf("Error calling enrich organization: %s", err.Error())
		err = l.neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4jentity.OrganizationPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
		}
		return nil
	}
	if enrichOrganizationResponse != nil {
		l.updateOrganizationWithEnrichData(ctx, tenant, domain, *organizationEntity, enrichOrganizationResponse)
	} else {
		err = l.neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4jentity.OrganizationPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
		}
	}

	return nil
}

func (l *organizationListenerImpl) updateOrganizationWithEnrichData(ctx context.Context, tenant, domain string, organizationEntity neo4jentity.OrganizationEntity, data *interfaces.OrganizationData) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationListener.updateOrganizationWithEnrichData")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "data", data)
	tracing.TagTenant(span, tenant)

	orgFields := data_fields.OrganizationFields{
		Source:       utils.StringPtr(neo4jentity.DataSourceOpenline.String()),
		EnrichDomain: utils.StringPtr(domain),
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
	if organizationEntity.Name == "" {
		if data.Name != "" {
			orgFields.Name = utils.StringPtr(data.Name)
		} else if data.Domain != "" {
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
			&common_srv.LinkWith{
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

	err = l.neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationEntity.ID, string(neo4jentity.OrganizationPropertyEnrichedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to update enriched at"))
	}
	l.services.Events.Publisher.PublishEventCompleted(ctx, tenant, organizationEntity.ID, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
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

	_, err := l.services.SocialService.AddSocialToEntity(ctx, nil, common_srv.LinkWith{
		Id:   organizationId,
		Type: commonmodel.ORGANIZATION,
	}, socialEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		l.log.Errorf("Error adding %s social: %s", url, err.Error())
	}
}
