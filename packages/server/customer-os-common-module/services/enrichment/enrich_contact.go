package enrichment

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/opentracing/opentracing-go/log"
	"strconv"
	"strings"
	"time"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neoenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

const (
	globalContactTTL = 90 * 24 * time.Hour
)

func (s *enrichmentService) EnrichContact(ctx context.Context, contactId, socialId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "enrichmentService.EnrichContact")
	defer spans.Finish()

	spans.TagEntity(contactId)
	spans.LogKV("socialId", socialId)

	tenant := common.GetTenantFromContext(ctx)

	// skip enrichment if disabled in tenant settings
	tenantSettingsDbNode, err := s.neo4jRepository.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "TenantReadRepository.GetTenantSettings"))
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)
	if !tenantSettingsEntity.EnrichContacts {
		spans.LogKV("result", "enrichment disabled")
		return nil
	}

	// skip enrichment if contact is already enriched
	contactEntity, err := s.contactService.GetContactById(ctx, contactId)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "ContactService.GetContactById"))
		return nil
	}

	if contactEntity.EnrichDetails.EnrichedAt != nil {
		spans.LogKV("result", "contact already enriched")
		return nil
	}

	// load social entity if socialId is provided
	var socialEntity *neo4j_entity.SocialEntity
	if socialId != "" {
		socialEntity, err = s.socialService.GetById(ctx, socialId)
		if err != nil {
			spans.TraceError(err)
			return err
		}
		if !socialEntity.IsLinkedin() {
			spans.LogKV("result", "social not linkedin")
			return nil
		}
	}

	enrichedWithGlobalContacts, err := s.enrichContactWithGlobalContacts(ctx, contactEntity, socialEntity)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "enrichContactWithGlobalContacts"))
	}

	if !enrichedWithGlobalContacts {
		return s.enrichContactWithScrapin(ctx, contactEntity, socialEntity)
	}
	return nil
}

func (s *enrichmentService) enrichContactWithGlobalContacts(ctx context.Context, contactEntity *neo4j_entity.ContactEntity, socialEntity *neo4j_entity.SocialEntity) (bool, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.enrichContactWithGlobalContacts")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	if socialEntity == nil || contactEntity == nil {
		return false, nil
	}

	// find global contacts by linkedin url, alias, or person identifier
	var globalContacts []*postgres_entity.GlobalContact
	globalContacts, err := s.globalContactService.GetGlobalContactsByLinkedIn(ctx, socialEntity.ExternalId)
	if err != nil {
		spans.TraceError(err)
		return false, err
	}
	if len(globalContacts) == 0 {
		globalContacts, err = s.globalContactService.GetGlobalContactsByLinkedIn(ctx, socialEntity.Alias)
		if err != nil {
			spans.TraceError(err)
			return false, err
		}
	}
	if len(globalContacts) == 0 {
		globalContacts, err = s.globalContactService.GetGlobalContactsByLinkedIn(ctx, socialEntity.ExtractLinkedinPersonIdentifierFromUrl())
		if err != nil {
			spans.TraceError(err)
			return false, err
		}
	}
	if len(globalContacts) == 0 {
		return false, nil
	}

	contactEnriched := false
	contactFieldsUpdated := false
	for _, globalContact := range globalContacts {
		// skip if global contact data was fetched more than 90 days ago
		if globalContact.DataFetchedAt == nil || globalContact.DataFetchedAt.Add(globalContactTTL).Before(time.Now()) {
			continue
		}

		// skip if global contact record job already ended
		if globalContact.JobEndedAt != nil && globalContact.JobEndedAt.Before(utils.Now()) {
			continue
		}

		// update common contact fields only once
		if !contactFieldsUpdated {
			contactFieldsUpdated = true

			contactFields := data_fields.ContactFields{}
			if globalContact.FirstName != "" {
				contactFields.FirstName = utils.StringPtr(globalContact.FirstName)
			}
			if globalContact.LastName != "" {
				contactFields.LastName = utils.StringPtr(globalContact.LastName)
			}
			if globalContact.ProfilePhotoPath != "" {
				contactFields.ProfilePhotoUrl = utils.StringPtr(constants.S3ImagesCDN + globalContact.ProfilePhotoPath)
			} else if globalContact.ProfilePhotoExternalUrl != "" {
				contactFields.ProfilePhotoUrl = utils.StringPtr(globalContact.ProfilePhotoExternalUrl)
			}
			if globalContact.LocationText != "" {
				contactLocation, err := s.createLocationForContact(ctx, contactEntity, globalContact.LocationText)
				if err != nil {
					spans.TraceError(err)
					return false, err
				}
				if contactLocation != nil && contactEntity.Timezone != "" && contactLocation.TimeZone != "" {
					contactFields.Timezone = utils.StringPtr(contactLocation.TimeZone)
				}
			}
			_, err := s.contactService.Save(ctx, nil, &contactEntity.Id, contactFields, false)
			if err != nil {
				spans.TraceError(errors.Wrap(err, "ContactService.Save"))
				return false, err
			}
			err = s.neo4jRepository.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonModel.NodeLabelContact, contactEntity.Id, string(neo4j_entity.ContactPropertyEnrichedGlobalContactId), strconv.FormatUint(globalContact.ID, 10))
			if err != nil {
				spans.TraceError(errors.Wrap(err, "CommonWriteRepository.UpdateStringProperty"))
			}
			err = s.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonModel.NodeLabelContact, contactEntity.Id, string(neo4j_entity.ContactPropertyEnrichedAt), utils.NowPtr())
			if err != nil {
				spans.TraceError(errors.Wrap(err, "CommonWriteRepository.UpdateTimeProperty"))
				s.log.Errorf("Error updating enriched at property: %s", err.Error())
			}
			contactEnriched = true
		}

		// find, create or show organization by domain
		organizationEntity, err := s.organizationService.GetOrganizationByDomain(ctx, globalContact.PrimaryDomain, false)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "OrganizationService.GetOrganizationByDomain"))
			continue
		}
		if organizationEntity != nil {
			if organizationEntity.IsHidden() {
				err = s.organizationService.Show(ctx, nil, organizationEntity.ID)
				if err != nil {
					spans.TraceError(errors.Wrap(err, "OrganizationService.Show"))
					continue
				}
			}
		} else {
			orgId, err := s.organizationService.CreateFromGlobalOrganizationByDomain(ctx, nil, globalContact.LastName, data_fields.OrganizationFields{})
			if err != nil {
				continue
			}
			organizationEntity, err = s.organizationService.GetById(ctx, tenant, orgId)
			if err != nil {
				spans.TraceError(errors.Wrap(err, "OrganizationService.GetById"))
				continue
			}
		}

		// link contact to organization
		if organizationEntity != nil {
			err = s.contactService.LinkContactWithOrganization(ctx, nil, contactEntity.Id, organizationEntity.ID,
				globalContact.JobTitle, "", neo4j_entity.DataSourceOpenline.String(), globalContact.JobEndedAt == nil,
				globalContact.JobStartedAt, globalContact.JobEndedAt)
			if err != nil {
				spans.TraceError(errors.Wrap(err, "ContactService.LinkContactWithOrganization"))
				continue
			}
		}
	}

	return contactEnriched, nil
}

func (s *enrichmentService) enrichContactWithScrapin(ctx context.Context, contactEntity *neo4j_entity.ContactEntity, socialEntity *neo4j_entity.SocialEntity) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.enrichContactWithScrapin")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	linkedInUrl, emailAddress, firstName, lastName, domain, companyName := "", "", "", "", "", ""
	// if linkedInUrl is empty fetch all data for searching person
	if socialEntity == nil {
		// prepare linked in for searching person
		socialDbNodes, err := s.neo4jRepository.SocialReadRepository.GetAllForEntities(ctx, tenant, commonModel.CONTACT, []string{contactEntity.Id})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "SocialReadRepository.GetAllForEntities"))
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
		emailAddress, err = s.getContactEmailAddress(ctx, contactEntity.Id)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "getContactEmail"))
			return err
		}

		// prepare organization name for searching person
		result, err := s.neo4jRepository.OrganizationReadRepository.GetPrimaryOrganizationsWithJobRoleForContacts(ctx, tenant, []string{contactEntity.Id})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "OrganizationReadRepository.GetPrimaryOrganizationsWithJobRoleForContacts"))
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
				spans.TraceError(errors.Wrap(err, "DomainReadRepository.GetForOrganizations"))
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
			if !s.cache.IsPersonalEmailProvider(emailDomain) {
				domain = emailDomain
			}
		}
		firstName = contactEntity.FirstName
		lastName = contactEntity.LastName
	} else {
		linkedInUrl = socialEntity.Url
	}

	spans.LogFields(
		log.String("emailAddress", emailAddress),
		log.String("firstName", firstName),
		log.String("lastName", lastName),
		log.String("domain", domain),
		log.String("companyName", companyName))
	if linkedInUrl != "" || emailAddress != "" || (firstName != "" && lastName != "" && domain != "" && companyName != "") {
		err := s.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonModel.NodeLabelContact, contactEntity.Id, string(neo4j_entity.ContactPropertyEnrichRequestedAt), utils.NowPtr())
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to update enrich requested at"))
		}
		s.events.Publisher.PublishNotification(ctx, tenant, contactEntity.Id, commonModel.CONTACT, utils.NewEventCompletedDetails().WithUpdate())

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
			spans.TraceError(errors.Wrap(err, "callApiEnrichPerson"))
			err = s.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonModel.NodeLabelContact, contactEntity.Id, string(neo4j_entity.ContactPropertyEnrichFailedAt), utils.NowPtr())
			if err != nil {
				spans.TraceError(errors.Wrap(err, "failed to update enrich failed at"))
			}
		}
		if scrapinResponseBody != nil {
			err = s.updateContactWithScrapInEnrichDetails(ctx, tenant, contactEntity, scrapinResponseBody, utils.IfNotNilUint64(recordID))
			if err != nil {
				spans.TraceError(err)
			}
		}
	} else {
		spans.LogKV("result", "no linkedInUrl, email or name")
	}

	return nil
}

func (s *enrichmentService) getContactEmailAddress(ctx context.Context, contactId string) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "enrichmentService.getContactEmailAddress")
	defer spans.Finish()
	spans.LogKV("contactId", contactId)

	tenant := common.GetTenantFromContext(ctx)

	records, err := s.neo4jRepository.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, commonModel.CONTACT, []string{contactId})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "EmailReadRepository.GetAllEmailNodesForLinkedEntityIds"))
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

func (s *enrichmentService) updateContactWithScrapInEnrichDetails(ctx context.Context, tenant string, contact *neo4j_entity.ContactEntity, scrapinContactResponse *postgres_entity.ScrapInResponseBody, recordID uint64) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "enrichmentService.updateContactWithScrapInEnrichDetails")
	defer spans.Finish()

	if scrapinContactResponse == nil || scrapinContactResponse.Person == nil {
		return nil
	}

	if !scrapinContactResponse.Success || scrapinContactResponse.Person == nil {
		spans.LogKV("result", "person not found")

		// mark contact as failed to enrich
		err := s.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonModel.NodeLabelContact, contact.Id, string(neo4j_entity.ContactPropertyEnrichFailedAt), utils.NowPtr())
		if err != nil {
			spans.TraceError(errors.Wrap(err, "CommonWriteRepository.UpdateTimeProperty"))
		}

		err = s.neo4jRepository.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonModel.NodeLabelContact, contact.Id, string(neo4j_entity.ContactPropertyEnrichedScrapinRecordId), strconv.FormatUint(recordID, 10))
		if err != nil {
			spans.TraceError(errors.Wrap(err, "CommonWriteRepository.UpdateStringProperty"))
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
		contactLocation, err := s.createLocationForContact(ctx, contact, scrapinContactResponse.Person.Location)
		if err != nil {
			return err
		}

		if contactLocation != nil && contact.Timezone != "" && contactLocation.TimeZone != "" {
			updateContact = true
			contactFields.Timezone = utils.StringPtr(contactLocation.TimeZone)
		}
	}

	if updateContact {
		_, err := s.contactService.Save(ctx, nil, &contact.Id, contactFields, false)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "ContactService.Save"))
			s.log.Errorf("Error updating contact: %s", err.Error())
		}
	}

	err := s.neo4jRepository.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonModel.NodeLabelContact, contact.Id, string(neo4j_entity.ContactPropertyEnrichedScrapinRecordId), strconv.FormatUint(recordID, 10))
	if err != nil {
		spans.TraceError(errors.Wrap(err, "CommonWriteRepository.UpdateStringProperty"))
		s.log.Errorf("Error updating enriched scrap in person search param property: %s", err.Error())
	}

	// mark contact as enriched
	err = s.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonModel.NodeLabelContact, contact.Id, string(neo4j_entity.ContactPropertyEnrichedAt), utils.NowPtr())
	if err != nil {
		spans.TraceError(errors.Wrap(err, "CommonWriteRepository.UpdateTimeProperty"))
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
			spans.TraceError(errors.Wrap(err, "SocialReadRepository.GetAllForEntities"))
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
			spans.TraceError(errors.Wrap(err, "SocialService.AddSocialToEntity"))
			s.log.Errorf("Error adding social profile: %s", err.Error())
		}
	}

	// Create main company of the linked in response if missing
	if scrapinContactResponse.Company != nil {
		var organizationDbNode *dbtype.Node

		// step1 - check org exists by linkedin url
		organizationDbNodes, err := s.neo4jRepository.OrganizationReadRepository.GetOrganizationsByLinkedIn(ctx, tenant, scrapinContactResponse.Company.LinkedInUrl, scrapinContactResponse.Company.UniversalName, scrapinContactResponse.Company.LinkedInId)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "OrganizationReadRepository.GetOrganizationsByLinkedIn"))
			s.log.Errorf("Error getting organization by social url: %s", err.Error())
		}
		if len(organizationDbNodes) > 0 {
			organizationDbNode = organizationDbNodes[0]
		}

		// step 2 - check org exists by domain
		if organizationDbNodes == nil {
			domain := s.domainService.GetPrimaryDomainForOrganizationWebsite(ctx, scrapinContactResponse.Company.WebsiteUrl)
			spans.LogKV("extractedDomainFromWebsite", domain)
			if domain != "" {
				organizationDbNode, err = s.neo4jRepository.OrganizationReadRepository.GetOrganizationByDomain(ctx, nil, tenant, domain)
				if err != nil {
					spans.TraceError(errors.Wrap(err, "OrganizationReadRepository.GetOrganizationByDomain"))
					s.log.Errorf("Error getting organization by domain: %s", err.Error())
					return err
				}
				if organizationDbNode != nil {
					organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
					if organizationEntity.IsHidden() {
						err = s.organizationService.Show(ctx, nil, organizationEntity.ID)
						if err != nil {
							spans.TraceError(errors.Wrap(err, "OrganizationService.Show"))
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
						spans.TraceError(errors.Wrap(err, "SocialService.AddSocialToEntity"))
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
				Stage:        utils.ToPtr(enum.Target),
				AppSource:    utils.StringPtr(string(enum.SourceScrapin)),
			})
			if err != nil {
				spans.TraceError(errors.Wrap(err, "OrganizationClient.UpsertOrganization"))
				s.log.Errorf("Error creating organization: %s", err.Error())
			} else if orgId == "" {
				spans.TraceError(errors.New("organization id is missing"))
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
					spans.TraceError(errors.Wrap(err, "SocialService.AddSocialToEntity"))
					s.log.Errorf("Error adding social profile: %s", err.Error())
				}
			}
		}
	}

	// minimize the impact on the batch processing
	time.Sleep(1 * time.Second)

	if len(scrapinContactResponse.Person.Positions.PositionHistory) > 0 {
		// process positions in reverse order
		for i := len(scrapinContactResponse.Person.Positions.PositionHistory) - 1; i >= 0; i-- {
			position := scrapinContactResponse.Person.Positions.PositionHistory[i]
			positionName := ""
			var positionStartedAt, positionEndedAt *time.Time

			// find organization by linkedin url
			organizationDbNodes, err := s.neo4jRepository.OrganizationReadRepository.GetOrganizationsByLinkedIn(ctx, tenant, position.LinkedInUrl, "", position.LinkedInId)
			if err != nil {
				spans.TraceError(errors.Wrap(err, "OrganizationReadRepository.GetOrganizationsByLinkedIn"))
				s.log.Errorf("Error getting organization by social url: %s", err.Error())
			}
			if len(organizationDbNodes) > 0 {
				organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNodes[0])
				if organizationEntity.IsHidden() {
					err = s.organizationService.Show(ctx, nil, organizationEntity.ID)
					if err != nil {
						spans.TraceError(errors.Wrap(err, "OrganizationService.Show"))
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
					spans.TraceError(errors.Wrap(err, "ContactClient.LinkWithOrganization"))
					s.log.Errorf("Error linking contact with organization: %s", err.Error())
					return err
				}
			}
		}
	}

	return nil
}

func (s *enrichmentService) createLocationForContact(ctx context.Context, contact *neo4j_entity.ContactEntity, rawLocation string) (*data_fields.LocationFields, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.createLocationForContact")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	locationFields, err := s.locationService.ExtractAndEnrichLocation(ctx, tenant, rawLocation)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if locationFields == nil {
		return nil, nil
	}

	locationFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
	locationFields.RawAddress = rawLocation

	_, err = s.locationService.Create(ctx, nil, *locationFields, &common_srv.LinkWith{
		Id:   contact.Id,
		Type: commonModel.CONTACT,
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return locationFields, nil
}
