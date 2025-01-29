package enrichment

import (
	"context"
	"strconv"
	"strings"
	"time"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neoenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

func (s *enrichmentService) EnrichContact(ctx context.Context, contactId, linkedInUrl string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "enrichmentService.enrichContact")
	defer span.Finish()
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.TagEntity(span, contactId)
	span.LogFields(log.String("linkedInUrl", linkedInUrl))

	tenant := common.GetTenantFromContext(ctx)

	// skip enrichment if disabled in tenant settings
	tenantSettings, err := s.neo4jRepository.TenantReadRepository.GetTenantSettings(ctx, tenant)
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
	contactEntity, err := s.contactService.GetContactById(ctx, contactId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactService.GetContactById"))
		return nil
	}

	if contactEntity.EnrichDetails.EnrichedAt != nil {
		span.LogFields(log.String("result", "contact already enriched"))
		return nil
	}

	emailAddress, firstName, lastName, domain, companyName := "", "", "", "", ""
	// if linkedInUrl is empty fetch all data for searching person
	if linkedInUrl == "" {
		// prepare linked in for searching person
		socialDbNodes, err := s.neo4jRepository.SocialReadRepository.GetAllForEntities(ctx, tenant, commonModel.CONTACT, []string{contactId})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "SocialReadRepository.GetAllForEntities"))
		} else {
			for _, socialDbNode := range socialDbNodes {
				socialEntity := neo4jmapper.MapDbNodeToSocialEntity(socialDbNode.Node)
				if socialEntity.IsLinkedin() {
					linkedInUrl = socialEntity.Url
					break
				}
			}
		}

		// prepare email address for searching person
		emailAddress, err = s.getContactEmailAddress(ctx, contactId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "getContactEmail"))
			return err
		}

		// prepare organization name for searching person
		result, err := s.neo4jRepository.OrganizationReadRepository.GetPrimaryOrganizationsWithJobRoleForContacts(ctx, tenant, []string{contactId})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetPrimaryOrganizationsWithJobRoleForContacts"))
		}
		var organizationEntity *neo4j_entity.OrganizationEntity
		if len(result) > 0 && result[0].Pair.First != nil {
			organizationEntity = neo4jmapper.MapDbNodeToOrganizationEntity(result[0].Pair.First)
			companyName = organizationEntity.Name
		}

		// prepare domain for searching person
		if organizationEntity != nil {
			organizationDomainDbNodes, err := s.neo4jRepository.DomainReadRepository.GetForOrganizations(ctx, tenant, []string{organizationEntity.ID})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "DomainReadRepository.GetForOrganizations"))
			}
			for _, domainDbNode := range organizationDomainDbNodes {
				domainEntity := neo4jmapper.MapDbNodeToDomainEntity(domainDbNode.Node)
				if utils.IfNotNilBool(domainEntity.IsPrimary) {
					domain = domainEntity.Domain
				} else if domain == "" {
					domain = domainEntity.Domain
				}
			}
		}

		// if not found domain from organization, get one from email, check it's not a personal one
		if domain == "" && emailAddress != "" {
			emailDomain := utils.ExtractDomainFromEmail(emailAddress)
			if !caches.IsPersonalEmailProvider(emailDomain) {
				domain = emailDomain
			}
		}
		firstName = contactEntity.FirstName
		lastName = contactEntity.LastName
	}

	span.LogFields(
		log.String("emailAddress", emailAddress),
		log.String("firstName", firstName),
		log.String("lastName", lastName),
		log.String("domain", domain),
		log.String("companyName", companyName))
	if linkedInUrl != "" || emailAddress != "" || (firstName != "" && lastName != "" && domain != "" && companyName != "") {
		err = s.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonModel.NodeLabelContact, contactEntity.Id, string(neo4j_entity.ContactPropertyEnrichRequestedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich requested at"))
		}
		s.events.Publisher.PublishEventCompleted(ctx, tenant, contactId, commonModel.CONTACT, utils.NewEventCompletedDetails().WithUpdate())

		query := interfaces.PersonSearch{
			LinkedinURL: linkedInUrl,
			FirstName:   firstName,
			LastName:    lastName,
			Email:       emailAddress,
			Domain:      domain,
			CompanyName: companyName,
		}

		recordID, scrapinResponseBody, err := s.EnrichPerson(ctx, query)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "callApiEnrichPerson"))
			err = s.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonModel.NodeLabelContact, contactEntity.Id, string(neo4j_entity.ContactPropertyEnrichFailedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to update enrich failed at"))
			}
		}
		if scrapinResponseBody != nil {
			err = s.enrichContactWithScrapInEnrichDetails(ctx, tenant, contactEntity, scrapinResponseBody, utils.IfNotNilUint64(recordID))
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "enrichContactWithScrapInEnrichDetails"))
			}
		}
	} else {
		span.LogFields(log.String("result", "no linkedInUrl, email or name"))
	}

	return nil
}

func (s *enrichmentService) getContactEmailAddress(ctx context.Context, contactId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "enrichmentService.getContactEmailAddress")
	defer span.Finish()
	span.LogFields(log.String("contactId", contactId))

	tenant := common.GetTenantFromContext(ctx)

	records, err := s.neo4jRepository.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, commonModel.CONTACT, []string{contactId})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "EmailReadRepository.GetAllEmailNodesForLinkedEntityIds"))
		return "", err
	}
	foundEmailAddress := ""
	for _, record := range records {
		emailEntity := neo4jmapper.MapDbNodeToEmailEntity(record.Node)
		if emailEntity.Email != "" {
			// Choose email address with primary domain
			if utils.IfNotNilBool(emailEntity.IsPrimaryDomain) {
				foundEmailAddress = emailEntity.Email
			} else if foundEmailAddress == "" {
				foundEmailAddress = emailEntity.Email
			}
		} else if foundEmailAddress == "" && emailEntity.RawEmail != "" && strings.Contains(emailEntity.RawEmail, "@") {
			foundEmailAddress = emailEntity.RawEmail
		}
	}
	return foundEmailAddress, nil
}

func (s *enrichmentService) enrichContactWithScrapInEnrichDetails(ctx context.Context, tenant string, contact *neo4j_entity.ContactEntity, enrichPersonResponse *postgres_entity.ScrapInResponseBody, recordID uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "enrichmentService.enrichContactWithScrapInEnrichDetails")
	defer span.Finish()
	tracing.TagComponentListener(span)
	tracing.TagTenant(span, tenant)

	if enrichPersonResponse == nil || enrichPersonResponse.Person == nil {
		return nil
	}

	scrapinContactResponse := enrichPersonResponse

	if !scrapinContactResponse.Success || scrapinContactResponse.Person == nil {
		span.LogFields(log.String("result", "person not found"))

		// mark contact as failed to enrich
		err := s.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonModel.NodeLabelContact, contact.Id, string(neo4j_entity.ContactPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "CommonWriteRepository.UpdateTimeProperty"))
		}

		err = s.neo4jRepository.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonModel.NodeLabelContact, contact.Id, string(neo4j_entity.ContactPropertyEnrichedScrapinRecordId), strconv.FormatUint(recordID, 10))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "CommonWriteRepository.UpdateStringProperty"))
		}

		return nil
	}

	updateContact := false
	contactFields := data_fields.ContactFields{}
	if scrapinContactResponse.Person.FirstName != "" {
		updateContact = true
		contactFields.FirstName = utils.StringPtr(scrapinContactResponse.Person.FirstName)
	}
	if scrapinContactResponse.Person.LastName != "" {
		updateContact = true
		contactFields.LastName = utils.StringPtr(scrapinContactResponse.Person.LastName)
	}
	if strings.TrimSpace(contact.ProfilePhotoUrl) == "" && scrapinContactResponse.Person.PhotoUrl != "" {
		updateContact = true
		contactFields.ProfilePhotoUrl = utils.StringPtr(scrapinContactResponse.Person.PhotoUrl)
	}
	if strings.TrimSpace(contact.Description) == "" && scrapinContactResponse.Person.Summary != "" {
		updateContact = true
		contactFields.Description = utils.StringPtr(scrapinContactResponse.Person.Summary)
	}

	// add location
	if scrapinContactResponse.Person.Location != "" {
		contactLocation, err := s.locationService.ExtractAndEnrichLocation(ctx, tenant, scrapinContactResponse.Person.Location)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ExtractAndEnrichLocation"))
		}
		if contactLocation != nil {
			contactLocation.RawAddress = scrapinContactResponse.Person.Location
			contactLocation.AppSource = utils.StringPtr(string(enum.SourceScrapin))
			_, err := s.locationService.Create(ctx, nil, *contactLocation, &common_srv.LinkWith{
				Id:   contact.Id,
				Type: commonModel.CONTACT,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			// update timezone on contact
			if contact.Timezone != "" && contactLocation.TimeZone != "" {
				updateContact = true
				contactFields.Timezone = utils.StringPtr(contactLocation.TimeZone)
			}
		}
	}

	if updateContact {
		_, err := s.contactService.Save(ctx, nil, &contact.Id, contactFields, false)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactService.Save"))
			s.log.Errorf("Error updating contact: %s", err.Error())
		}
	}

	err := s.neo4jRepository.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonModel.NodeLabelContact, contact.Id, string(neo4j_entity.ContactPropertyEnrichedScrapinRecordId), strconv.FormatUint(recordID, 10))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "CommonWriteRepository.UpdateStringProperty"))
		s.log.Errorf("Error updating enriched scrap in person search param property: %s", err.Error())
	}

	// mark contact as enriched
	err = s.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonModel.NodeLabelContact, contact.Id, string(neo4j_entity.ContactPropertyEnrichedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "CommonWriteRepository.UpdateTimeProperty"))
		s.log.Errorf("Error updating enriched at property: %s", err.Error())
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
		socialDbNodes, err := s.neo4jRepository.SocialReadRepository.GetAllForEntities(ctx, tenant, commonModel.CONTACT, []string{contact.Id})
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

		_, err = s.socialService.AddSocialToEntity(ctx,
			nil,
			common_srv.LinkWith{
				Id:   contact.Id,
				Type: commonModel.CONTACT,
			},
			neo4j_entity.SocialEntity{
				Id:             socialId,
				Url:            url,
				Alias:          scrapinContactResponse.Person.PublicIdentifier,
				ExternalId:     scrapinContactResponse.Person.LinkedInIdentifier,
				FollowersCount: int64(scrapinContactResponse.Person.FollowerCount),
				Source:         neo4j_entity.DataSourceOpenline,
				AppSource:      string(enum.SourceScrapin),
			})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "SocialService.AddSocialToEntity"))
			s.log.Errorf("Error adding social profile: %s", err.Error())
		}
	}

	// Create main company of the linked in response if missing
	if scrapinContactResponse.Company != nil {
		var organizationDbNode *dbtype.Node

		// step1 - check org exists by linkedin url
		organizationDbNodes, err := s.neo4jRepository.OrganizationReadRepository.GetOrganizationsByLinkedIn(ctx, tenant, scrapinContactResponse.Company.LinkedInUrl, scrapinContactResponse.Company.UniversalName, scrapinContactResponse.Company.LinkedInId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationsByLinkedIn"))
			s.log.Errorf("Error getting organization by social url: %s", err.Error())
		}
		if len(organizationDbNodes) > 0 {
			organizationDbNode = organizationDbNodes[0]
		}

		// step 2 - check org exists by domain
		if organizationDbNodes == nil {
			domain, _ := s.domainService.GetPrimaryDomainForOrganizationWebsite(ctx, scrapinContactResponse.Company.WebsiteUrl)
			span.LogFields(log.String("extractedDomainFromWebsite", domain))
			if domain != "" {
				organizationDbNode, err = s.neo4jRepository.OrganizationReadRepository.GetOrganizationByDomain(ctx, nil, tenant, domain)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationByDomain"))
					s.log.Errorf("Error getting organization by domain: %s", err.Error())
					return err
				}
				if organizationDbNode != nil {
					organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
					if organizationEntity.IsHidden() {
						err = s.organizationService.Show(ctx, nil, organizationEntity.ID)
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "OrganizationService.Show"))
							return err
						}
					}
					_, err = s.socialService.AddSocialToEntity(ctx, nil, common_srv.LinkWith{
						Id:   organizationEntity.ID,
						Type: commonModel.ORGANIZATION,
					}, neo4j_entity.SocialEntity{
						Url:            scrapinContactResponse.Company.LinkedInUrl,
						Alias:          scrapinContactResponse.Company.UniversalName,
						ExternalId:     scrapinContactResponse.Company.LinkedInId,
						FollowersCount: int64(scrapinContactResponse.Company.FollowerCount),
					})
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "SocialService.AddSocialToEntity"))
						s.log.Errorf("Error adding social profile: %s", err.Error())
					}
				}
			}
		}

		// step 3 if not found - create organization
		if organizationDbNode == nil {
			orgId, err := s.organizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
				Name:         utils.StringPtr(scrapinContactResponse.Company.Name),
				Website:      utils.StringPtr(scrapinContactResponse.Company.WebsiteUrl),
				Relationship: utils.ToPtr(neoenum.OrganizationRelationshipProspect),
				Stage:        utils.ToPtr(neoenum.Lead),
				AppSource:    utils.StringPtr(string(enum.SourceScrapin)),
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "OrganizationClient.UpsertOrganization"))
				s.log.Errorf("Error creating organization: %s", err.Error())
			} else if orgId == "" {
				tracing.TraceErr(span, errors.New("organization id is missing"))
				return errors.New("organization id is missing")
			} else {
				_, err = s.socialService.AddSocialToEntity(ctx, nil, common_srv.LinkWith{
					Id:   orgId,
					Type: commonModel.ORGANIZATION,
				}, neo4j_entity.SocialEntity{
					Url:            scrapinContactResponse.Company.LinkedInUrl,
					Alias:          scrapinContactResponse.Company.UniversalName,
					ExternalId:     scrapinContactResponse.Company.LinkedInId,
					FollowersCount: int64(scrapinContactResponse.Company.FollowerCount),
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "SocialService.AddSocialToEntity"))
					s.log.Errorf("Error adding social profile: %s", err.Error())
				}
			}
		}
	}

	// minimize the impact on the batch processing
	time.Sleep(1 * time.Second)

	if len(scrapinContactResponse.Person.Positions.PositionHistory) > 0 {
		positionName := ""
		var positionStartedAt, positionEndedAt *time.Time

		// process positions in reverse order
		for i := len(scrapinContactResponse.Person.Positions.PositionHistory) - 1; i >= 0; i-- {
			position := scrapinContactResponse.Person.Positions.PositionHistory[i]

			// find organization by linkedin url
			organizationDbNodes, err := s.neo4jRepository.OrganizationReadRepository.GetOrganizationsByLinkedIn(ctx, tenant, position.LinkedInUrl, "", position.LinkedInId)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationsByLinkedIn"))
				s.log.Errorf("Error getting organization by social url: %s", err.Error())
			}
			if len(organizationDbNodes) > 0 {
				organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNodes[0])
				if organizationEntity.IsHidden() {
					err = s.organizationService.Show(ctx, nil, organizationEntity.ID)
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "OrganizationService.Show"))
						return err
					}
				}

				positionName = position.Title
				if position.StartEndDate.Start != nil {
					positionStartedAt = utils.TimePtr(utils.FirstTimeOfMonth(position.StartEndDate.Start.Year, position.StartEndDate.Start.Month))
				}
				if position.StartEndDate.End != nil {
					positionEndedAt = utils.TimePtr(utils.FirstTimeOfMonth(position.StartEndDate.End.Year, position.StartEndDate.End.Month))
				}

				// link contact with organization
				err = s.contactService.LinkContactWithOrganization(ctx, nil, contact.Id, organizationEntity.ID, positionName, "",
					neo4j_entity.DataSourceOpenline.String(), false, positionStartedAt, positionEndedAt)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "ContactClient.LinkWithOrganization"))
					s.log.Errorf("Error linking contact with organization: %s", err.Error())
					return err
				}
			}
		}
	}

	return nil
}
