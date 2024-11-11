package listeners

import (
	"bytes"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	socialpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/social"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/constants"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ContactListener interface {
	enrichContact(ctx context.Context, contactId, linkedInUrl string) error
	getContactEmail(ctx context.Context, contactId string) (string, error)
	callApiEnrichPerson(ctx context.Context, tenant, linkedinUrl, email, firstName, lastName, domain string) (*enrichmentmodel.EnrichPersonScrapinResponse, error)
	enrichContactWithScrapInEnrichDetails(ctx context.Context, tenant string, contact *neo4jentity.ContactEntity, enrichPersonResponse *enrichmentmodel.EnrichPersonScrapinResponse) error
}

type contactListenerImpl struct {
	services *service.Services
	log      logger.Logger
}

func NewContactListener(services *service.Services, log logger.Logger) ContactListener {
	return &contactListenerImpl{services: services, log: log}
}

func OnSocialAddedToContact(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnSocialAddedToContact")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		tracing.TraceErr(span, err)
		return nil
	}
	messageData := message.Event.Data.(*dto.AddSocialToContact)
	socialUrl := messageData.Social
	contactId := message.Event.EntityId

	span.SetTag(tracing.SpanTagEntityId, contactId)

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

	c := NewContactListener(services, services.Logger)

	if strings.Contains(socialUrl, "linkedin.com") {
		err := c.enrichContact(ctx, contactId, socialUrl)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "enrichContact"))
		}
	}

	return nil
}

func OnRequestedEnrichContact(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnSocialAddedToContact")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)
	contactId := message.Event.EntityId

	span.SetTag(tracing.SpanTagEntityId, contactId)

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

	c := NewContactListener(services, services.Logger)

	err := c.enrichContact(ctx, contactId, "")
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "enrichContact"))
	}

	return nil
}

func (c *contactListenerImpl) enrichContact(ctx context.Context, contactId, linkedInUrl string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactListener.enrichContact")
	defer span.Finish()
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.TagEntity(span, contactId)
	span.LogFields(log.String("linkedInUrl", linkedInUrl))

	tenant := common.GetTenantFromContext(ctx)

	// skip enrichment if disabled in tenant settings
	tenantSettings, err := c.services.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "TenantReadRepository.GetTenantSettings"))
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettings)
	if !tenantSettingsEntity.EnrichContacts {
		span.LogFields(log.String("result", "enrichment disabled"))
		return nil
	}

	// skip enrichment if contact is already enriched
	contactDbNode, err := c.services.Neo4jRepositories.ContactReadRepository.GetContact(ctx, tenant, contactId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactReadRepository.GetContact"))
		return nil
	}
	contactEntity := neo4jmapper.MapDbNodeToContactEntity(contactDbNode)

	if contactEntity.EnrichDetails.EnrichedAt != nil {
		span.LogFields(log.String("result", "contact already enriched"))
		return nil
	}

	emailAddress, firstName, lastName, domain := "", "", "", ""
	if linkedInUrl == "" {
		socialDbNodes, err := c.services.Neo4jRepositories.SocialReadRepository.GetAllForEntities(ctx, tenant, model.CONTACT, []string{contactId})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "SocialReadRepository.GetAllForEntities"))
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
		emailAddress, err = c.getContactEmail(ctx, contactId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "getContactEmail"))
			return err
		}

		domains, _ := c.services.Neo4jRepositories.ContactReadRepository.GetLinkedOrgDomains(ctx, tenant, contactEntity.Id)
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
		err = c.services.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contactEntity.Id, string(neo4jentity.ContactPropertyEnrichRequestedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich requested at"))
		}
		utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, c.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())

		apiResponse, err := c.callApiEnrichPerson(ctx, tenant, linkedInUrl, emailAddress, firstName, lastName, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "callApiEnrichPerson"))
			err = c.services.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contactEntity.Id, string(neo4jentity.ContactPropertyEnrichFailedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
			}
		} else {
			err = c.enrichContactWithScrapInEnrichDetails(ctx, tenant, contactEntity, apiResponse)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "enrichContactWithScrapInEnrichDetails"))
			}
		}
	} else {
		span.LogFields(log.String("result", "no linkedInUrl, email or name"))
	}

	return nil
}

func (c *contactListenerImpl) getContactEmail(ctx context.Context, contactId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactListener.getContactEmail")
	defer span.Finish()
	span.LogFields(log.String("contactId", contactId))

	tenant := common.GetTenantFromContext(ctx)

	records, err := c.services.Neo4jRepositories.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, model.CONTACT, []string{contactId})
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

func (c *contactListenerImpl) callApiEnrichPerson(ctx context.Context, tenant, linkedinUrl, email, firstName, lastName, domain string) (*enrichmentmodel.EnrichPersonScrapinResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactListener.callApiEnrichPerson")
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
	req, err := http.NewRequest("GET", c.services.GlobalConfig.InternalServices.EnrichmentApiConfig.Url+"/enrichPerson", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, c.services.GlobalConfig.InternalServices.EnrichmentApiConfig.ApiKey)
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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		span.LogFields(log.String("response.body", string(body)))
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return nil, err
	}

	var enrichPersonApiResponse enrichmentmodel.EnrichPersonScrapinResponse
	err = json.Unmarshal(body, &enrichPersonApiResponse)
	if err != nil {
		span.LogFields(log.String("response.body", string(body)))
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal enrich person response"))
		return nil, err
	}
	return &enrichPersonApiResponse, nil
}

func (c *contactListenerImpl) enrichContactWithScrapInEnrichDetails(ctx context.Context, tenant string, contact *neo4jentity.ContactEntity, enrichPersonResponse *enrichmentmodel.EnrichPersonScrapinResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactListener.enrichContactWithScrapInEnrichDetails")
	defer span.Finish()

	if enrichPersonResponse == nil || enrichPersonResponse.Data == nil || enrichPersonResponse.Data.PersonProfile == nil {
		return nil
	}

	scrapinContactResponse := enrichPersonResponse.Data.PersonProfile

	if !scrapinContactResponse.Success || scrapinContactResponse.Person == nil {
		span.LogFields(log.String("result", "person not found"))

		// mark contact as failed to enrich
		err := c.services.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contact.Id, string(neo4jentity.ContactPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateTimeProperty"))
		}

		err = c.services.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contact.Id, neo4jentity.ContactPropertyEnrichedScrapinRecordId, strconv.FormatUint(enrichPersonResponse.RecordId, 10))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateAnyProperty"))
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
		contactLocation, err := c.services.LocationService.ExtractAndEnrichLocation(ctx, tenant, scrapinContactResponse.Person.Location)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ExtractAndEnrichLocation"))
		}
		if contactLocation != nil {
			tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
			_, err := utils.CallEventsPlatformGRPCWithRetry[*locationpb.LocationIdGrpcResponse](func() (*locationpb.LocationIdGrpcResponse, error) {
				return c.services.GrpcClients.ContactClient.AddLocation(ctx, &contactpb.ContactAddLocationGrpcRequest{
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
				c.log.Errorf("Error adding location to contact: %s", err.Error())
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
		_, err := c.services.ContactService.SaveContact(ctx, &contact.Id, contactFields, "", neo4jmodel.ExternalSystem{})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactService.SaveContact"))
			c.log.Errorf("Error updating contact: %s", err.Error())
		}
	}

	err := c.services.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contact.Id, neo4jentity.ContactPropertyEnrichedScrapinRecordId, strconv.FormatUint(enrichPersonResponse.RecordId, 10))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateAnyProperty"))
		c.log.Errorf("Error updating enriched scrap in person search param property: %s", err.Error())
	}

	// mark contact as enriched
	err = c.services.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contact.Id, string(neo4jentity.ContactPropertyEnrichedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateTimeProperty"))
		c.log.Errorf("Error updating enriched at property: %s", err.Error())
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
		socialDbNodes, err := c.services.Neo4jRepositories.SocialReadRepository.GetAllForEntities(ctx, tenant, model.CONTACT, []string{contact.Id})
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

		_, err = c.services.SocialService.AddSocialToEntity(ctx,
			service.LinkWith{
				Id:   contact.Id,
				Type: model.CONTACT,
			},
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
			tracing.TraceErr(span, errors.Wrap(err, "SocialService.AddSocialToEntity"))
			c.log.Errorf("Error adding social profile: %s", err.Error())
		}
	}

	if scrapinContactResponse.Company != nil {
		var organizationDbNode *dbtype.Node

		// step1 - check org exists by linkedin url
		organizationDbNode, err = c.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationBySocialUrl(ctx, tenant, scrapinContactResponse.Company.LinkedInUrl)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationBySocialUrl"))
			c.log.Errorf("Error getting organization by social url: %s", err.Error())
		}

		// step 2 - check org exists by social url
		if organizationDbNode == nil {
			// step 2 - find by domain
			domain, _ := c.services.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, scrapinContactResponse.Company.WebsiteUrl)
			span.LogFields(log.String("extractedDomainFromWebsite", domain))
			if domain != "" {
				organizationDbNode, err = c.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, tenant, domain)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationByDomain"))
					c.log.Errorf("Error getting organization by domain: %s", err.Error())
					return err
				}
				if organizationDbNode != nil {
					orgId := utils.GetStringPropOrEmpty(organizationDbNode.Props, "id")
					_, err = utils.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
						return c.services.GrpcClients.OrganizationClient.AddSocial(ctx, &organizationpb.AddSocialGrpcRequest{
							Tenant:         tenant,
							OrganizationId: orgId,
							Url:            scrapinContactResponse.Company.LinkedInUrl,
							FollowersCount: int64(scrapinContactResponse.Company.FollowerCount),
						})
					})
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "OrganizationClient.AddSocial"))
						c.log.Errorf("Error adding social profile: %s", err.Error())
					}
				}
			}
		}

		// step 3 if not found - create organization
		if organizationDbNode == nil {
			orgId, err := c.services.OrganizationService.Save(ctx, nil, tenant, nil, &neo4jrepository.OrganizationSaveFields{
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
				c.log.Errorf("Error creating organization: %s", err.Error())
			} else if orgId == nil {
				tracing.TraceErr(span, errors.New("organization id is nil"))
				return errors.New("organization id is nil")
			} else {
				_, err = utils.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
					return c.services.GrpcClients.OrganizationClient.AddSocial(ctx, &organizationpb.AddSocialGrpcRequest{
						Tenant:         tenant,
						OrganizationId: *orgId,
						Url:            scrapinContactResponse.Company.LinkedInUrl,
						FollowersCount: int64(scrapinContactResponse.Company.FollowerCount),
					})
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "OrganizationClient.AddSocial"))
					c.log.Errorf("Error adding social profile: %s", err.Error())
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
			orgByLinkedinUrlNode, err := c.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationBySocialUrl(ctx, tenant, position.LinkedInUrl)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationBySocialUrl"))
				c.log.Errorf("Error getting organization by social url: %s", err.Error())
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
				err = c.services.ContactService.LinkContactWithOrganization(ctx, contact.Id, organizationId, positionName, "",
					neo4jentity.DataSourceOpenline.String(), false, positionStartedAt, positionEndedAt)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "ContactClient.LinkWithOrganization"))
					c.log.Errorf("Error linking contact with organization: %s", err.Error())
					return err
				}
			}
		}
	}

	return nil
}
