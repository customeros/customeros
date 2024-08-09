package contact

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/additional_services"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/location"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	socialpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/social"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact"
	event2 "github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type ContactEventHandler struct {
	repositories         *repository.Repositories
	log                  logger.Logger
	cfg                  *config.Config
	caches               caches.Cache
	grpcClients          *grpc_client.Clients
	scrapInService       *additional_services.ScrapInService
	locationEventHandler *location.LocationEventHandler
}

func NewContactEventHandler(repositories *repository.Repositories, log logger.Logger, cfg *config.Config, caches caches.Cache, grpcClients *grpc_client.Clients) *ContactEventHandler {
	contactEventHandler := ContactEventHandler{
		repositories:         repositories,
		log:                  log,
		cfg:                  cfg,
		caches:               caches,
		grpcClients:          grpcClients,
		scrapInService:       additional_services.NewScrapInService(log, cfg, repositories.PostgresRepositories),
		locationEventHandler: location.NewLocationEventHandler(repositories, log, cfg, grpcClients),
	}

	return &contactEventHandler
}

func (h *ContactEventHandler) OnEnrichContactRequested(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnEnrichContactRequested")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event2.ContactRequestEnrich
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	return h.enrichContact(ctx, eventData.Tenant, contactId, &additional_services.ScrapInEnrichContactFlow{})
}

func (h *ContactEventHandler) enrichContact(ctx context.Context, tenant, contactId string, flow *additional_services.ScrapInEnrichContactFlow) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.enrichContact")
	defer span.Finish()

	// skip enrichment if disabled in tenant settings
	tenantSettings, err := h.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "TenantReadRepository.GetTenantSettings"))
		h.log.Errorf("Error getting tenant settings: %s", err.Error())
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettings)
	if !tenantSettingsEntity.EnrichContacts {
		span.LogFields(log.String("result", "enrichment disabled"))
		h.log.Infof("Enrichment disabled for tenant %s", tenant)
		return nil
	}

	// skip enrichment if contact is already enriched
	contactDbNode, err := h.repositories.Neo4jRepositories.ContactReadRepository.GetContact(ctx, tenant, contactId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactReadRepository.GetContact"))
		h.log.Errorf("Error getting contact with id %s: %s", contactId, err.Error())
		return nil
	}
	contactEntity := neo4jmapper.MapDbNodeToContactEntity(contactDbNode)

	if contactEntity.EnrichDetails.EnrichedAt != nil {
		span.LogFields(log.String("result", "contact already enriched"))
		h.log.Infof("Contact %s already enriched", contactId)
		return nil
	}

	if flow.Url != "" {
		return h.enrichContactByLinkedInProfile(ctx, tenant, flow, contactEntity)
	} else {
		// get email from contact
		email, err := h.getContactEmail(ctx, tenant, contactId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "getContactEmail"))
			h.log.Errorf("Error getting contact email: %s", err.Error())
			return err
		}
		if email != "" {
			flow.Email = email
			return h.enrichContactByEmail(ctx, tenant, flow, contactEntity)
		}
	}

	return nil
}

func (h *ContactEventHandler) enrichContactByEmail(ctx context.Context, tenant string, flow *additional_services.ScrapInEnrichContactFlow, contactEntity *neo4jentity.ContactEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.enrichContactByEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("email", flow.Email))
	span.LogFields(log.String("contactId", contactEntity.Id))

	// get enrich details for email
	record, err := h.repositories.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, flow.Email, postgresentity.ScrapInFlowPersonSearch)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow"))
		return err
	}

	// if enrich details found and not older than 1 year, use the data, otherwise scrape it
	daysAgo365 := utils.Now().Add(-time.Hour * 24 * 365)
	daysAgo30 := utils.Now().Add(-time.Hour * 24 * 30)
	if record != nil && record.CreatedAt.After(daysAgo365) && record.Data != "" {
		// enrich contact with the data found
		span.LogFields(log.Bool("result.email already enriched", true))

		var scrapinContactResponse postgresentity.ScrapInPersonResponse
		err := json.Unmarshal([]byte(record.Data), &scrapinContactResponse)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
			h.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
			return err
		}
		if scrapinContactResponse.Success || record.CreatedAt.After(daysAgo30) {
			return h.enrichContactWithScrapInEnrichDetails(ctx, tenant, flow, contactEntity, scrapinContactResponse)
		}
	}

	domains, err := h.repositories.Neo4jRepositories.ContactReadRepository.GetLinkedOrgDomains(ctx, tenant, contactEntity.Id)
	domain := ""
	emailDomain := utils.ExtractDomainFromEmail(flow.Email)
	if utils.Contains(domains, emailDomain) {
		domain = emailDomain
	} else if len(domains) > 0 {
		domain = domains[0]
	}
	if domain == "" {
		domain = emailDomain
	}
	firstName, lastName := contactEntity.DeriveFirstAndLastNames()
	scrapinContactResponse, err := h.scrapInService.ScrapInPersonSearch(ctx, tenant, flow.Email, firstName, lastName, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "scrapInPersonSearch"))
		return err
	}
	return h.enrichContactWithScrapInEnrichDetails(ctx, tenant, flow, contactEntity, scrapinContactResponse)
}

func (h *ContactEventHandler) enrichContactByLinkedInProfile(ctx context.Context, tenant string, flow *additional_services.ScrapInEnrichContactFlow, contactEntity *neo4jentity.ContactEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.enrichContactByLinkedInProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("url", flow.Url))
	span.LogFields(log.String("contactId", contactEntity.Id))

	// get enrich details for email
	record, err := h.repositories.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, flow.Url, postgresentity.ScrapInFlowPersonProfile)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow"))
		return err
	}

	// if enrich details found and not older than 1 year, use the data, otherwise scrape it
	daysAgo365 := utils.Now().Add(-time.Hour * 24 * 365)
	daysAgo30 := utils.Now().Add(-time.Hour * 24 * 30)
	if record != nil && record.CreatedAt.After(daysAgo365) && record.Data != "" {
		// enrich contact with the data found
		span.LogFields(log.Bool("result.linkedin profile already enriched", true))

		var scrapinContactResponse postgresentity.ScrapInPersonResponse
		err := json.Unmarshal([]byte(record.Data), &scrapinContactResponse)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
			h.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
			return err
		}
		if scrapinContactResponse.Success || record.CreatedAt.After(daysAgo30) {
			return h.enrichContactWithScrapInEnrichDetails(ctx, tenant, flow, contactEntity, scrapinContactResponse)
		}
	}

	scrapinContactResponse, err := h.scrapInService.ScrapInPersonProfile(ctx, tenant, flow.GetParam1())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "scrapInPersonProfile"))
		return err
	}
	return h.enrichContactWithScrapInEnrichDetails(ctx, tenant, flow, contactEntity, scrapinContactResponse)
}

func (h *ContactEventHandler) getContactEmail(ctx context.Context, tenant, contactId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.getContactEmail")
	defer span.Finish()
	span.LogFields(log.String("contactId", contactId))

	records, err := h.repositories.Neo4jRepositories.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, model.CONTACT, []string{contactId})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "EmailReadRepository.GetAllEmailNodesForLinkedEntityIds"))
		return "", err
	}
	foundEmailAddress := ""
	for _, record := range records {
		emailEntity := neo4jmapper.MapDbNodeToEmailEntity(record.Node)
		if emailEntity.Email != "" && strings.Contains(emailEntity.Email, "@") {
			foundEmailAddress = emailEntity.Email
			break
		}
		if emailEntity.RawEmail != "" && strings.Contains(emailEntity.RawEmail, "@") {
			foundEmailAddress = emailEntity.RawEmail
		}
	}
	return foundEmailAddress, nil
}

func (h *ContactEventHandler) enrichContactWithScrapInEnrichDetails(ctx context.Context, tenant string, flow *additional_services.ScrapInEnrichContactFlow, contact *neo4jentity.ContactEntity, scrapinContactResponse postgresentity.ScrapInPersonResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.enrichContactWithScrapInEnrichDetails")
	defer span.Finish()

	if !scrapinContactResponse.Success || scrapinContactResponse.Person == nil {
		span.LogFields(log.String("result", "person not found"))
		h.log.Infof("Person not found for %s", flow.GetParam1())

		// mark contact as failed to enrich
		err := h.repositories.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contact.Id, string(flow.GetTimeLabel()), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateTimeProperty"))
			h.log.Errorf("Error updating enriched at scrap in person search property: %s", err.Error())
		}

		err = h.repositories.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contact.Id, flow.GetParamLabel(), flow.GetParam1())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateAnyProperty"))
			h.log.Errorf("Error updating enriched scrap in person search param property: %s", err.Error())
		}

		return nil
	}

	// update contact
	tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	upsertContactGrpcRequest := contactpb.UpsertContactGrpcRequest{
		Id:     contact.Id,
		Tenant: tenant,
		SourceFields: &commonpb.SourceFields{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppScrapin,
		},
	}
	fieldsMask := make([]contactpb.ContactFieldMask, 0)
	name := ""
	if scrapinContactResponse.Person.FirstName != "" {
		upsertContactGrpcRequest.FirstName = scrapinContactResponse.Person.FirstName
		name += scrapinContactResponse.Person.FirstName
		fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_FIRST_NAME)
	}
	if scrapinContactResponse.Person.LastName != "" {
		if name != "" {
			name += " "
		}
		name += scrapinContactResponse.Person.LastName
		upsertContactGrpcRequest.LastName = scrapinContactResponse.Person.LastName
		fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_LAST_NAME)
	}
	if name != "" {
		upsertContactGrpcRequest.Name = name
		fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_NAME)
	}
	if scrapinContactResponse.Person.PhotoUrl != "" {
		upsertContactGrpcRequest.ProfilePhotoUrl = scrapinContactResponse.Person.PhotoUrl
		fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_PROFILE_PHOTO_URL)
	}
	if scrapinContactResponse.Person.Summary != "" {
		upsertContactGrpcRequest.Description = scrapinContactResponse.Person.Summary
		fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_DESCRIPTION)
	}

	// add location
	if scrapinContactResponse.Person.Location != "" {
		location, err := h.locationEventHandler.ExtractAndEnrichLocation(ctx, tenant, scrapinContactResponse.Person.Location)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ExtractAndEnrichLocation"))
		}
		if location != nil {
			_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*locationpb.LocationIdGrpcResponse](func() (*locationpb.LocationIdGrpcResponse, error) {
				return h.grpcClients.ContactClient.AddLocation(ctx, &contactpb.ContactAddLocationGrpcRequest{
					ContactId: contact.Id,
					Tenant:    tenant,
					SourceFields: &commonpb.SourceFields{
						Source:    constants.SourceOpenline,
						AppSource: constants.AppScrapin,
					},
					LocationDetails: &locationpb.LocationDetails{
						RawAddress:    scrapinContactResponse.Person.Location,
						Country:       location.Country,
						CountryCodeA2: location.CountryCodeA2,
						CountryCodeA3: location.CountryCodeA3,
						Region:        location.Region,
						Locality:      location.Locality,
						AddressLine1:  location.Address,
						AddressLine2:  location.Address2,
						ZipCode:       location.Zip,
						AddressType:   location.AddressType,
						HouseNumber:   location.HouseNumber,
						PostalCode:    location.PostalCode,
						Commercial:    location.Commercial,
						Predirection:  location.Predirection,
						District:      location.District,
						Street:        location.Street,
						Latitude:      utils.FloatToString(location.Latitude),
						Longitude:     utils.FloatToString(location.Longitude),
						TimeZone:      location.TimeZone,
						UtcOffset:     location.UtcOffset,
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ContactClient.AddLocationToContact"))
				h.log.Errorf("Error adding location to contact: %s", err.Error())
			}
			// update timezone and offset on contact
			if location.TimeZone != "" {
				upsertContactGrpcRequest.Timezone = location.TimeZone
				fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_TIMEZONE)
			}
		}
	}

	// update contact
	if len(fieldsMask) > 0 {
		upsertContactGrpcRequest.FieldsMask = fieldsMask
		_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return h.grpcClients.ContactClient.UpsertContact(ctx, &upsertContactGrpcRequest)
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactClient.UpsertContact"))
			h.log.Errorf("Error updating contact: %s", err.Error())
		}
	}

	err := h.repositories.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contact.Id, string(flow.GetTimeLabel()), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error updating enriched at scrap in person search property: %s", err.Error())
	}

	err = h.repositories.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contact.Id, flow.GetParamLabel(), flow.GetParam1())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateAnyProperty"))
		h.log.Errorf("Error updating enriched scrap in person search param property: %s", err.Error())
	}

	// mark contact as enriched
	err = h.repositories.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contact.Id, string(neo4jentity.ContactPropertyEnrichedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateTimeProperty"))
		h.log.Errorf("Error updating enriched at property: %s", err.Error())
	}

	// add additional enrich details after marking contact as enriched, to avoid re-enriching

	// add social profiles
	if scrapinContactResponse.Person.LinkedInUrl != "" {
		// prepare url, replace LinkedInIdentifier with PublicIdentifier in url
		url := scrapinContactResponse.Person.LinkedInUrl
		if scrapinContactResponse.Person.LinkedInIdentifier != "" {
			url = strings.Replace(url, scrapinContactResponse.Person.LinkedInIdentifier, scrapinContactResponse.Person.PublicIdentifier, 1)
		}
		// add ending / if missing
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}

		// get social id by url if exist for current contact
		socialId := ""
		socialDbNodes, err := h.repositories.Neo4jRepositories.SocialReadRepository.GetAllForEntities(ctx, tenant, model.CONTACT, []string{contact.Id})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "SocialReadRepository.GetAllForEntities"))
		}
		for _, socialDbNode := range socialDbNodes {
			socialEntity := neo4jmapper.MapDbNodeToSocialEntity(socialDbNode.Node)
			if socialEntity.Url == url {
				socialId = socialEntity.Id
				break
			}
		}

		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
			return h.grpcClients.ContactClient.AddSocial(ctx, &contactpb.ContactAddSocialGrpcRequest{
				ContactId:      contact.Id,
				Tenant:         tenant,
				SocialId:       socialId,
				Url:            url,
				Alias:          scrapinContactResponse.Person.PublicIdentifier,
				ExternalId:     scrapinContactResponse.Person.LinkedInIdentifier,
				FollowersCount: int64(scrapinContactResponse.Person.FollowerCount),
				SourceFields: &commonpb.SourceFields{
					Source:    constants.SourceOpenline,
					AppSource: constants.AppScrapin,
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactClient.AddSocial"))
			h.log.Errorf("Error adding social profile: %s", err.Error())
		}
	}

	// add organization
	organizationNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContactId(ctx, tenant, contact.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationByContactId"))
		h.log.Errorf("Error getting organization by contact id: %s", err.Error())
		return err
	}

	if organizationNode == nil && scrapinContactResponse.Company != nil && scrapinContactResponse.Company.WebsiteUrl != "" {

		domain := utils.ExtractDomain(scrapinContactResponse.Company.WebsiteUrl)

		var organizationId string

		organizationByDomainNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationByDomain"))
			h.log.Errorf("Error getting organization by domain: %s", err.Error())
			return err
		}

		if organizationByDomainNode == nil {
			upsertOrganizationRequest := organizationpb.UpsertOrganizationGrpcRequest{
				Tenant:       tenant,
				Name:         scrapinContactResponse.Company.Name,
				Website:      scrapinContactResponse.Company.WebsiteUrl,
				Relationship: neo4jenum.Prospect.String(),
				Stage:        neo4jenum.Lead.String(),
				SourceFields: &commonpb.SourceFields{
					Source:    constants.SourceOpenline,
					AppSource: constants.AppScrapin,
				},
			}

			organizationCreateResponse, err := utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
				return h.grpcClients.OrganizationClient.UpsertOrganization(ctx, &upsertOrganizationRequest)
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "OrganizationClient.UpsertOrganization"))
				h.log.Errorf("Error upserting organization: %s", err.Error())
				return err
			}

			organizationId = organizationCreateResponse.Id
		} else {
			organizationId = utils.GetStringPropOrEmpty(organizationByDomainNode.Props, "id")
		}

		positionName := ""
		if len(scrapinContactResponse.Person.Positions.PositionHistory) > 0 {
			for _, position := range scrapinContactResponse.Person.Positions.PositionHistory {
				if position.Title != "" && position.CompanyName != "" && position.CompanyName == scrapinContactResponse.Company.Name {
					positionName = position.Title
					break
				}
			}
		}

		//minimize the impact on the batch processing
		time.Sleep(3 * time.Second)

		_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return h.grpcClients.ContactClient.LinkWithOrganization(ctx, &contactpb.LinkWithOrganizationGrpcRequest{
				Tenant:         tenant,
				ContactId:      contact.Id,
				OrganizationId: organizationId,
				JobTitle:       positionName,
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactClient.LinkWithOrganization"))
			h.log.Errorf("Error upserting organization: %s", err.Error())
			return err
		}

	}

	return nil
}

func (h *ContactEventHandler) OnSocialAddedToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnSocialAddedToContact")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event2.ContactAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	if strings.Contains(eventData.Url, ".linkedin.com") || strings.Contains(eventData.Url, "/linkedin.com") {
		return h.enrichContact(ctx, eventData.Tenant, contactId, &additional_services.ScrapInEnrichContactFlow{Url: eventData.Url})
	}

	return nil
}
