package contact

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/location"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	socialpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/social"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact"
	contactevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ContactEventHandler struct {
	services             *service.Services
	log                  logger.Logger
	cfg                  *config.Config
	caches               caches.Cache
	grpcClients          *grpc_client.Clients
	locationEventHandler *location.LocationEventHandler
}

func NewContactEventHandler(services *service.Services, log logger.Logger, cfg *config.Config, caches caches.Cache, grpcClients *grpc_client.Clients) *ContactEventHandler {
	contactEventHandler := ContactEventHandler{
		services:             services,
		log:                  log,
		cfg:                  cfg,
		caches:               caches,
		grpcClients:          grpcClients,
		locationEventHandler: location.NewLocationEventHandler(services, log, cfg, grpcClients),
	}

	return &contactEventHandler
}

func (h *ContactEventHandler) OnEnrichContactRequested(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnEnrichContactRequested")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData contactevent.ContactRequestEnrich
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	return h.enrichContact(ctx, eventData.Tenant, contactId, "")
}

func (h *ContactEventHandler) OnSocialAddedToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnSocialAddedToContact")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData contactevent.ContactAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	if strings.Contains(eventData.Url, "linkedin.com") {
		return h.enrichContact(ctx, eventData.Tenant, contactId, eventData.Url)
	}

	return nil
}

func (h *ContactEventHandler) enrichContact(ctx context.Context, tenant, contactId, linkedInUrl string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.enrichContact")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contactId", contactId), log.String("linkedInUrl", linkedInUrl))

	// skip enrichment if disabled in tenant settings
	tenantSettings, err := h.services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
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
	contactDbNode, err := h.services.CommonServices.Neo4jRepositories.ContactReadRepository.GetContact(ctx, tenant, contactId)
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

	emailAddress, firstName, lastName, domain := "", "", "", ""
	if linkedInUrl == "" {
		socialDbNodes, err := h.services.CommonServices.Neo4jRepositories.SocialReadRepository.GetAllForEntities(ctx, tenant, model.CONTACT, []string{contactId})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "SocialReadRepository.GetAllForEntities"))
			h.log.Errorf("Error getting social profiles for contact: %s", err.Error())
		} else {
			for _, socialDbNode := range socialDbNodes {
				socialEntity := neo4jmapper.MapDbNodeToSocialEntity(socialDbNode.Node)
				if strings.Contains(socialEntity.Url, "linkedin.com") {
					linkedInUrl = socialEntity.Url
					break
				}
			}
		}

		// get email from contact
		emailAddress, err = h.getContactEmail(ctx, tenant, contactId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "getContactEmail"))
			h.log.Errorf("Error getting contact email: %s", err.Error())
			return err
		}

		domains, _ := h.services.CommonServices.Neo4jRepositories.ContactReadRepository.GetLinkedOrgDomains(ctx, tenant, contactEntity.Id)
		emailDomain := utils.ExtractDomainFromEmail(emailAddress)
		if utils.Contains(domains, emailDomain) {
			domain = emailDomain
		} else if len(domains) > 0 {
			domain = domains[0]
		}
		if domain == "" {
			domain = emailDomain
		}
		firstName, lastName = contactEntity.DeriveFirstAndLastNames()
	}

	span.LogFields(log.String("emailAddress", emailAddress), log.String("firstName", firstName), log.String("lastName", lastName), log.String("domain", domain))
	if linkedInUrl != "" || emailAddress != "" || (firstName != "" && lastName != "" && domain != "") {
		err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contactEntity.Id, string(neo4jentity.ContactPropertyEnrichRequestedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich requested at"))
		}
		utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

		apiResponse, err := h.callApiEnrichPerson(ctx, tenant, linkedInUrl, emailAddress, firstName, lastName, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "callApiEnrichPerson"))
			err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contactEntity.Id, string(neo4jentity.ContactPropertyEnrichFailedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
			}
		} else {
			err = h.enrichContactWithScrapInEnrichDetails(ctx, tenant, contactEntity, apiResponse)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "enrichContactWithScrapInEnrichDetails"))
			}
		}
	} else {
		span.LogFields(log.String("result", "no linkedInUrl, email or name"))
	}

	return nil
}

func (h *ContactEventHandler) getContactEmail(ctx context.Context, tenant, contactId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.getContactEmail")
	defer span.Finish()
	span.LogFields(log.String("contactId", contactId))

	records, err := h.services.CommonServices.Neo4jRepositories.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, model.CONTACT, []string{contactId})
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

func (h *ContactEventHandler) enrichContactWithScrapInEnrichDetails(ctx context.Context, tenant string, contact *neo4jentity.ContactEntity, enrichPersonResponse *enrichmentmodel.EnrichPersonScrapinResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.enrichContactWithScrapInEnrichDetails")
	defer span.Finish()

	if enrichPersonResponse == nil || enrichPersonResponse.Data == nil || enrichPersonResponse.Data.PersonProfile == nil {
		return nil
	}

	scrapinContactResponse := enrichPersonResponse.Data.PersonProfile

	if !scrapinContactResponse.Success || scrapinContactResponse.Person == nil {
		span.LogFields(log.String("result", "person not found"))

		// mark contact as failed to enrich
		err := h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contact.Id, string(neo4jentity.ContactPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateTimeProperty"))
			h.log.Errorf("Error updating enriched at scrap in person search property: %s", err.Error())
		}

		err = h.services.CommonServices.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contact.Id, neo4jentity.ContactPropertyEnrichedScrapinRecordId, strconv.FormatUint(enrichPersonResponse.RecordId, 10))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateAnyProperty"))
			h.log.Errorf("Error updating enriched scrap in person search param property: %s", err.Error())
		}

		return nil
	}

	updateContact := false
	contactFields := neo4jrepository.ContactFields{}
	if strings.TrimSpace(contact.FirstName) == "" && scrapinContactResponse.Person.FirstName != "" {
		updateContact = true
		contactFields.FirstName = scrapinContactResponse.Person.FirstName
		contactFields.UpdateFirstName = true
	}
	if strings.TrimSpace(contact.LastName) == "" && scrapinContactResponse.Person.LastName != "" {
		updateContact = true
		contactFields.LastName = scrapinContactResponse.Person.LastName
		contactFields.UpdateLastName = true
	}
	if strings.TrimSpace(contact.ProfilePhotoUrl) == "" && scrapinContactResponse.Person.PhotoUrl != "" {
		updateContact = true
		contactFields.ProfilePhotoUrl = scrapinContactResponse.Person.PhotoUrl
		contactFields.UpdateProfilePhotoUrl = true
	}
	if strings.TrimSpace(contact.Description) == "" && scrapinContactResponse.Person.Summary != "" {
		updateContact = true
		contactFields.Description = scrapinContactResponse.Person.Summary
		contactFields.UpdateDescription = true
	}

	// add location
	if scrapinContactResponse.Person.Location != "" {
		contactLocation, err := h.locationEventHandler.ExtractAndEnrichLocation(ctx, tenant, scrapinContactResponse.Person.Location)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ExtractAndEnrichLocation"))
		}
		if contactLocation != nil {
			tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
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
						Country:       contactLocation.Country,
						CountryCodeA2: contactLocation.CountryCodeA2,
						CountryCodeA3: contactLocation.CountryCodeA3,
						Region:        contactLocation.Region,
						Locality:      contactLocation.Locality,
						AddressLine1:  contactLocation.Address,
						AddressLine2:  contactLocation.Address2,
						ZipCode:       contactLocation.Zip,
						AddressType:   contactLocation.AddressType,
						HouseNumber:   contactLocation.HouseNumber,
						PostalCode:    contactLocation.PostalCode,
						Commercial:    contactLocation.Commercial,
						Predirection:  contactLocation.Predirection,
						District:      contactLocation.District,
						Street:        contactLocation.Street,
						Latitude:      utils.FloatToString(contactLocation.Latitude),
						Longitude:     utils.FloatToString(contactLocation.Longitude),
						TimeZone:      contactLocation.TimeZone,
						UtcOffset:     contactLocation.UtcOffset,
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ContactClient.AddLocationToContact"))
				h.log.Errorf("Error adding location to contact: %s", err.Error())
			}
			// update timezone on contact
			if contact.Timezone != "" && contactLocation.TimeZone != "" {
				updateContact = true
				contactFields.Timezone = contactLocation.TimeZone
				contactFields.UpdateTimezone = true
			}
		}
	}

	if updateContact {
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppScrapin,
		})
		_, err := h.services.CommonServices.ContactService.SaveContact(innerCtx, &contact.Id, contactFields, "", neo4jmodel.ExternalSystem{})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactService.SaveContact"))
			h.log.Errorf("Error updating contact: %s", err.Error())
		}
	}

	err := h.services.CommonServices.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contact.Id, neo4jentity.ContactPropertyEnrichedScrapinRecordId, strconv.FormatUint(enrichPersonResponse.RecordId, 10))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateAnyProperty"))
		h.log.Errorf("Error updating enriched scrap in person search param property: %s", err.Error())
	}

	// mark contact as enriched
	err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contact.Id, string(neo4jentity.ContactPropertyEnrichedAt), utils.NowPtr())
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
		socialDbNodes, err := h.services.CommonServices.Neo4jRepositories.SocialReadRepository.GetAllForEntities(ctx, tenant, model.CONTACT, []string{contact.Id})
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

		_, err = h.services.CommonServices.SocialService.MergeSocialWithEntity(ctx, tenant, contact.Id, model.CONTACT,
			neo4jentity.SocialEntity{
				Id:             socialId,
				Url:            url,
				Alias:          scrapinContactResponse.Person.PublicIdentifier,
				ExternalId:     scrapinContactResponse.Person.LinkedInIdentifier,
				FollowersCount: int64(scrapinContactResponse.Person.FollowerCount),
				Source:         neo4jentity.DataSourceOpenline,
				AppSource:      constants.AppScrapin,
			})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactClient.AddSocial"))
			h.log.Errorf("Error adding social profile: %s", err.Error())
		}
	}

	if scrapinContactResponse.Company != nil {
		var organizationDbNode *dbtype.Node

		// step1 - check org exists by linkedin url
		organizationDbNode, err = h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationBySocialUrl(ctx, tenant, scrapinContactResponse.Company.LinkedInUrl)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationBySocialUrl"))
			h.log.Errorf("Error getting organization by social url: %s", err.Error())
		}

		// step 2 - check org exists by social url
		if organizationDbNode == nil {
			// step 2 - find by domain
			domain, _ := h.services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, scrapinContactResponse.Company.WebsiteUrl)
			span.LogFields(log.String("extractedDomainFromWebsite", domain))
			if domain != "" {
				organizationDbNode, err = h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, tenant, domain)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationByDomain"))
					h.log.Errorf("Error getting organization by domain: %s", err.Error())
					return err
				}
				if organizationDbNode != nil {
					orgId := utils.GetStringPropOrEmpty(organizationDbNode.Props, "id")
					_, err = utils.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
						return h.grpcClients.OrganizationClient.AddSocial(ctx, &organizationpb.AddSocialGrpcRequest{
							Tenant:         tenant,
							OrganizationId: orgId,
							Url:            scrapinContactResponse.Company.LinkedInUrl,
							FollowersCount: int64(scrapinContactResponse.Company.FollowerCount),
						})
					})
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "OrganizationClient.AddSocial"))
						h.log.Errorf("Error adding social profile: %s", err.Error())
					}
				}
			}
		}

		// step 3 if not found - create organization
		if organizationDbNode == nil {
			orgId, err := h.services.CommonServices.OrganizationService.Save(ctx, nil, tenant, nil, &neo4jrepository.OrganizationSaveFields{
				Name:               scrapinContactResponse.Company.Name,
				Website:            scrapinContactResponse.Company.WebsiteUrl,
				Relationship:       neo4jenum.Prospect,
				Stage:              neo4jenum.Lead,
				UpdateName:         true,
				UpdateWebsite:      true,
				UpdateRelationship: true,
				UpdateStage:        true,
				SourceFields: neo4jmodel.SourceFields{
					Source:    constants.SourceOpenline,
					AppSource: constants.AppScrapin,
				},
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "OrganizationClient.UpsertOrganization"))
				h.log.Errorf("Error creating organization: %s", err.Error())
			} else if orgId == nil {
				tracing.TraceErr(span, errors.New("organization id is nil"))
				return errors.New("organization id is nil")
			} else {
				_, err = utils.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
					return h.grpcClients.OrganizationClient.AddSocial(ctx, &organizationpb.AddSocialGrpcRequest{
						Tenant:         tenant,
						OrganizationId: *orgId,
						Url:            scrapinContactResponse.Company.LinkedInUrl,
						FollowersCount: int64(scrapinContactResponse.Company.FollowerCount),
					})
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "OrganizationClient.AddSocial"))
					h.log.Errorf("Error adding social profile: %s", err.Error())
				}
			}
		}
	}

	//minimize the impact on the batch processing
	time.Sleep(3 * time.Second)

	if len(scrapinContactResponse.Person.Positions.PositionHistory) > 0 {
		positionName := ""
		var positionStartedAt, positionEndedAt *time.Time
		for _, position := range scrapinContactResponse.Person.Positions.PositionHistory {
			// find organization by linkedin url
			orgByLinkedinUrlNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationBySocialUrl(ctx, tenant, position.LinkedInUrl)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationBySocialUrl"))
				h.log.Errorf("Error getting organization by social url: %s", err.Error())
				continue
			}
			if orgByLinkedinUrlNode != nil {
				positionName = position.Title
				if position.StartEndDate.Start != nil {
					positionStartedAt = utils.TimePtr(utils.FirstTimeOfMonth(position.StartEndDate.Start.Year, position.StartEndDate.Start.Month))
				}
				if position.StartEndDate.End != nil {
					positionEndedAt = utils.TimePtr(utils.FirstTimeOfMonth(position.StartEndDate.End.Year, position.StartEndDate.End.Month))
				}
				organizationId := utils.GetStringPropOrEmpty(orgByLinkedinUrlNode.Props, "id")
				// link contact with organization
				_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
					return h.grpcClients.ContactClient.LinkWithOrganization(ctx, &contactpb.LinkWithOrganizationGrpcRequest{
						Tenant:         tenant,
						ContactId:      contact.Id,
						OrganizationId: organizationId,
						JobTitle:       positionName,
						StartedAt:      utils.ConvertTimeToTimestampPtr(positionStartedAt),
						EndedAt:        utils.ConvertTimeToTimestampPtr(positionEndedAt),
					})
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "ContactClient.LinkWithOrganization"))
					h.log.Errorf("Error upserting organization: %s", err.Error())
					return err
				}
			}
		}
	}

	return nil
}

func (h *ContactEventHandler) callApiEnrichPerson(ctx context.Context, tenant, linkedinUrl, email, firstName, lastName, domain string) (*enrichmentmodel.EnrichPersonScrapinResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.callApiEnrichPerson")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("linkedinUrl", linkedinUrl), log.String("email", email), log.String("firstName", firstName), log.String("lastName", lastName), log.String("domain", domain))

	requestJSON, err := json.Marshal(enrichmentmodel.EnrichPersonRequest{
		Email:       email,
		LinkedinUrl: linkedinUrl,
		FirstName:   firstName,
		LastName:    lastName,
		Domain:      domain,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("GET", h.cfg.Services.EnrichmentApi.Url+"/enrichPerson", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, h.cfg.Services.EnrichmentApi.ApiKey)
	req.Header.Set(security.TenantHeader, tenant)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.statusCode", response.StatusCode))

	var enrichPersonApiResponse enrichmentmodel.EnrichPersonScrapinResponse
	err = json.NewDecoder(response.Body).Decode(&enrichPersonApiResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode enrich person response"))
		return nil, err
	}
	return &enrichPersonApiResponse, nil
}
