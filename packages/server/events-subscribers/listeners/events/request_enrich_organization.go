package events_listeners

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type RequestEnrichOrganizationListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewRequestEnrichOrganizationListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &HideContactListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.RequestEnrichOrganization](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *RequestEnrichOrganizationListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestEnrichOrganizationListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return err
	}

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	message, err := l.validateMessage(ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	organizationId := event.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	domain, _ := l.dependencies.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, message.Url)
	if domain == "" {
		return nil
	}

	return l.enrichOrganization(ctx, tenant, organizationId, domain)
}

func (l *RequestEnrichOrganizationListener) enrichOrganization(ctx context.Context, tenant, organizationId, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestEnrichOrganizationListener.enrichOrganization")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, organizationId)

	if domain == "" {
		tracing.TraceErr(span, errors.New("domain is empty"))
		return nil
	}

	// check if domain is primary
	domainEntity, err := l.dependencies.CommonServices.DomainService.GetDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get domain"))
		l.dependencies.Logger.Errorf("Error getting domain %s: %s", domain, err.Error())
		return nil
	}
	if domainEntity.IsPrimary == nil || !*domainEntity.IsPrimary {
		l.dependencies.Logger.Infof("Domain %s is not primary", domain)
		return nil
	}

	organizationDbNode, err := l.dependencies.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		l.dependencies.Logger.Errorf("Error getting organization with id %s: %v", organizationId, err)
		return nil
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	if organizationEntity.EnrichDetails.EnrichedAt != nil {
		l.dependencies.Logger.Infof("Organization %s already enriched", organizationId)
		return nil
	}

	err = l.dependencies.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4j_entity.OrganizationPropertyEnrichRequestedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich requested at"))
	}

	l.dependencies.CommonServices.Events.Publisher.PublishNotification(ctx, tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

	enrichOrganizationResponse, err := l.dependencies.CommonServices.EnrichmentService.FetchEnrichOrganizationData(ctx, &domain, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call enrich organization"))
		l.dependencies.Logger.Errorf("Error calling enrich organization: %s", err.Error())
		err = l.dependencies.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4j_entity.OrganizationPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
		}
		return nil
	}
	if enrichOrganizationResponse != nil {
		l.updateOrganizationWithEnrichData(ctx, tenant, domain, *organizationEntity, enrichOrganizationResponse)
	} else {
		err = l.dependencies.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationId, string(neo4j_entity.OrganizationPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
		}
	}

	return nil
}

func (l *RequestEnrichOrganizationListener) updateOrganizationWithEnrichData(ctx context.Context, tenant, domain string, organizationEntity neo4j_entity.OrganizationEntity, data *interfaces.OrganizationData) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestEnrichOrganizationListener.updateOrganizationWithEnrichData")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "data", data)
	tracing.TagTenant(span, tenant)

	orgFields := data_fields.OrganizationFields{
		Source:       utils.StringPtr(neo4j_entity.DataSourceOpenline.String()),
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

	_, err := l.dependencies.CommonServices.OrganizationService.Save(ctx, nil, &organizationEntity.ID, orgFields)
	if err != nil {
		tracing.TraceErr(span, err)
		l.dependencies.Logger.Errorf("Error updaing organization with enrich data: %s", err.Error())
	}

	// add location
	if !data.Location.IsEmpty() {
		_, err := l.dependencies.CommonServices.LocationService.Create(ctx, nil, data_fields.LocationFields{
			Country:       data.Location.Country,
			CountryCodeA2: data.Location.CountryCodeA2,
			CountryCodeA3: data.Location.CountryCodeA3,
			Region:        data.Location.Region,
			Locality:      data.Location.Locality,
			PostalCode:    data.Location.PostalCode,
			Address:       data.Location.AddressLine1,
			Address2:      data.Location.AddressLine2,
			AppSource:     utils.StringPtr("Enrichment"),
			Source:        utils.StringPtr(string(enum.SourceCustomerOS)),
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
		l.addSocial(ctx, organizationEntity.ID, tenant, social.Url, social.Alias, social.Id, "Enrichment")
	}

	err = l.dependencies.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.NodeLabelOrganization, organizationEntity.ID, string(neo4j_entity.OrganizationPropertyEnrichedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to update enriched at"))
	}
	l.dependencies.CommonServices.Events.Publisher.PublishNotification(ctx, tenant, organizationEntity.ID, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
}

func (l *RequestEnrichOrganizationListener) addSocial(ctx context.Context, organizationId, tenant, url, alias, externalId, appSource string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestEnrichOrganizationListener.addSocial")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, organizationId)
	span.LogKV("url", url, "alias", alias, "externalId", externalId, "appSource", appSource)

	socialEntity := neo4j_entity.SocialEntity{
		Url:        url,
		Alias:      alias,
		ExternalId: externalId,
		AppSource:  appSource,
		Source:     neo4j_entity.DataSourceOpenline,
	}

	_, err := l.dependencies.CommonServices.SocialService.AddSocialToEntity(ctx, nil, common_srv.LinkWith{
		Id:   organizationId,
		Type: commonmodel.ORGANIZATION,
	}, socialEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		l.dependencies.Logger.Errorf("Error adding %s social: %s", url, err.Error())
	}
}

func (l *RequestEnrichOrganizationListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.RequestEnrichOrganization, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestEnrichOrganizationListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.RequestEnrichOrganization)
	if !ok {
		err := fmt.Errorf("expected RequestEnrichOrganization, got %T", event.Event.Data)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if event.Event.EntityId == "" {
		err := errors.New("EntityId not set on event")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return message, nil
}
