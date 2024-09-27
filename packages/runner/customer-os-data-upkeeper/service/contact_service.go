package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	cosClient "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/client"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/generic"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

type ContactService interface {
	UpkeepContacts()
	AskForWorkEmailOnBetterContact()
	EnrichWithWorkEmailFromBetterContact()
	CheckBetterContactRequestsWithoutResponse()
	EnrichContacts()
	AskForLinkedInConnections()
	ProcessLinkedInConnections()
	LinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn()
}

type contactService struct {
	cfg                 *config.Config
	log                 logger.Logger
	commonServices      *commonService.Services
	customerOSApiClient cosClient.CustomerOSApiClient
	eventBufferService  *eventbuffer.EventBufferStoreService
}

func NewContactService(cfg *config.Config, log logger.Logger, commonServices *commonService.Services, customerOSApiClient cosClient.CustomerOSApiClient, eventBufferService *eventbuffer.EventBufferStoreService) ContactService {
	return &contactService{
		cfg:                 cfg,
		log:                 log,
		commonServices:      commonServices,
		customerOSApiClient: customerOSApiClient,
		eventBufferService:  eventBufferService,
	}
}

func (s *contactService) UpkeepContacts() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	s.removeEmptySocials(ctx)
	s.removeDuplicatedSocials(ctx)
	s.hideContactsWithGroupEmail(ctx)
	s.updateContactNamesFromEmails(ctx)
}

func (s *contactService) removeEmptySocials(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.removeEmptySocials")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		minutesSinceLastUpdate := 180
		records, err := s.commonServices.Neo4jRepositories.SocialReadRepository.GetEmptySocialsForEntityType(ctx, model.NodeLabelContact, minutesSinceLastUpdate, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting socials: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		//remove socials from contact
		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
				return s.commonServices.GrpcClients.ContactClient.RemoveSocial(ctx, &contactpb.ContactRemoveSocialGrpcRequest{
					Tenant:    record.Tenant,
					ContactId: record.LinkedEntityId,
					SocialId:  record.SocialId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error removing social {%s}: %s", record.SocialId, err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *contactService) removeDuplicatedSocials(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.removeDuplicatedSocials")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		minutesSinceLastUpdate := 180
		records, err := s.commonServices.Neo4jRepositories.SocialReadRepository.GetDuplicatedSocialsForEntityType(ctx, model.NodeLabelContact, minutesSinceLastUpdate, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting socials: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		//remove socials from contact
		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
				return s.commonServices.GrpcClients.ContactClient.RemoveSocial(ctx, &contactpb.ContactRemoveSocialGrpcRequest{
					Tenant:    record.Tenant,
					ContactId: record.LinkedEntityId,
					SocialId:  record.SocialId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error removing social {%s}: %s", record.SocialId, err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *contactService) hideContactsWithGroupEmail(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.hideContactsWithGroupEmail")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 500

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsWithGroupEmail(ctx, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting contacts: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		//hide contact
		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
				return s.commonServices.GrpcClients.ContactClient.HideContact(ctx, &contactpb.ContactIdGrpcRequest{
					Tenant:    record.Tenant,
					ContactId: record.ContactId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error hiding contact {%s}: %s", record.ContactId, err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *contactService) updateContactNamesFromEmails(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.updateContactNamesFromEmails")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 200

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsWithEmailForNameUpdate(ctx, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting contacts: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		// update contact names
		for _, record := range records {
			firstName, lastName := neo4jentity.ContactEntity{}.GetNamesFromString(record.FieldStr1)
			if firstName == "" && lastName == "" {
				err = errors.New("cannot derive names from email")
				tracing.TraceErr(span, err)
				s.log.Errorf("Error updating contact {%s}: %s", record.ContactId, err.Error())
				continue
			}

			upsertRequest := contactpb.UpsertContactGrpcRequest{
				Tenant: record.Tenant,
				Id:     record.ContactId,
				SourceFields: &commonpb.SourceFields{
					AppSource: constants.AppSourceDataUpkeeper,
				},
			}
			var fieldsMask []contactpb.ContactFieldMask
			if firstName != "" {
				upsertRequest.FirstName = firstName
				fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_FIRST_NAME)
			}
			if lastName != "" {
				upsertRequest.LastName = lastName
				fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_LAST_NAME)
			}
			upsertRequest.FieldsMask = fieldsMask
			_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
				return s.commonServices.GrpcClients.ContactClient.UpsertContact(ctx, &upsertRequest)
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ContactClient.UpsertContact"))
				s.log.Errorf("Error updating contact {%s}: %s", record.ContactId, err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *contactService) AskForWorkEmailOnBetterContact() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	s.findEmailsWithBetterContact(ctx)
}

func (s *contactService) EnrichWithWorkEmailFromBetterContact() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	s.enrichWithWorkEmailFromBetterContact(ctx)
}

func (s *contactService) CheckBetterContactRequestsWithoutResponse() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	s.checkBetterContactRequestsWithoutResponse(ctx)
}

func (s *contactService) AskForLinkedInConnections() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.askForLinkedInConnections(ctx)
}

func (s *contactService) ProcessLinkedInConnections() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.processLinkedInConnections(ctx)
}

func (s *contactService) LinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn(ctx)
}

type BetterContactRequestBody struct {
	Data    []BetterContactData `json:"data"`
	Webhook string              `json:"webhook"`
}

type BetterContactData struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	LinkedInUrl   string `json:"linkedin_url"`
	Company       string `json:"company"`
	CompanyDomain string `json:"company_domain"`
}

type BetterContactResponseBody struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (s *contactService) askForLinkedInConnections(c context.Context) {
	span, ctx := tracing.StartTracerSpan(c, "ContactService.askForLinkedInConnections")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	linkedinTokens, err := s.commonServices.PostgresRepositories.BrowserConfigRepository.Get(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	span.LogFields(log.Int("linkedinTokens", len(linkedinTokens)))

	for _, linkedinToken := range linkedinTokens {
		//todo check if there is already a scheduled job for this token today
		err := s.commonServices.PostgresRepositories.BrowserAutomationRunRepository.Add(ctx, entity.BrowserAutomationsRun{
			BrowserConfigId: linkedinToken.Id,
			UserId:          linkedinToken.UserId,
			Tenant:          linkedinToken.Tenant,
			Type:            "FIND_CONNECTIONS",
			Status:          "SCHEDULED",
			Payload:         "\"\"",
		})
		if err != nil {
			tracing.TraceErr(span, err)
			break
		}
	}

}

func (s *contactService) processLinkedInConnections(c context.Context) {
	span, ctx := tracing.StartTracerSpan(c, "ContactService.processLinkedInConnections")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	automationsRuns, err := s.commonServices.PostgresRepositories.BrowserAutomationRunRepository.Get(ctx, "FIND_CONNECTIONS", "COMPLETED")
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	span.LogFields(log.Int("processing", len(automationsRuns)))

	for _, automationRun := range automationsRuns {
		s.processAutomationRunResult(ctx, automationRun)
	}
}

func (s *contactService) processAutomationRunResult(c context.Context, automationRun entity.BrowserAutomationsRun) {
	span, ctx := opentracing.StartSpanFromContext(c, "ContactService.processAutomationRunResult")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	result, err := s.commonServices.PostgresRepositories.BrowserAutomationRunResultRepository.Get(ctx, automationRun.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	if result == nil || result.ResultData == "" {
		span.LogFields(log.String("results", "empty"))
		return
	}

	useByEmailNode, err := s.commonServices.Neo4jRepositories.UserReadRepository.GetUserById(ctx, automationRun.Tenant, automationRun.UserId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "UserReadRepository.GetUserById"))
		return
	}
	if useByEmailNode == nil {
		tracing.TraceErr(span, errors.Wrap(err, "User does not exist"))
		return
	}

	var results []string

	err = json.Unmarshal([]byte(result.ResultData), &results)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	span.LogFields(log.Int("results", len(results)))

	tenant := automationRun.Tenant
	userId := automationRun.UserId

	for _, linkedinUrl := range results {
		err := s.processLinkedInUrl(ctx, tenant, linkedinUrl, userId)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

	}

	err = s.commonServices.PostgresRepositories.BrowserAutomationRunRepository.MarkAsProcessed(ctx, automationRun.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
}

func (s *contactService) processLinkedInUrl(ctx context.Context, tenant, linkedinUrl, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.processLinkedInUrl")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	linkedinProfileUrl := linkedinUrl
	if linkedinProfileUrl != "" && linkedinProfileUrl[len(linkedinProfileUrl)-1] != '/' {
		linkedinProfileUrl = linkedinProfileUrl + "/"
	}

	linkedinProfileUrl = utils.NormalizeString(linkedinProfileUrl)

	contactsWithLinkedin, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsWithSocialUrl(ctx, tenant, linkedinProfileUrl)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactReadRepository.GetContactsWithSocialUrl"))
		return err
	}

	var contactIds []string
	if len(contactsWithLinkedin) == 0 {
		contactInput := cosModel.ContactInput{
			SocialURL: &linkedinProfileUrl,
		}

		contactId, err := s.customerOSApiClient.CreateContact(tenant, "", contactInput)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "CreateContact"))
			return err
		}
		contactIds = append(contactIds, contactId)
	} else {
		for _, contactWithLinkedin := range contactsWithLinkedin {
			contactId := utils.GetStringPropOrEmpty(contactWithLinkedin.Props, "id")
			contactIds = append(contactIds, contactId)
		}
	}

	//link contacts to user
	if userId != "" {
		for _, cid := range contactIds {

			isLinkedWith, err := s.commonServices.Neo4jRepositories.CommonReadRepository.IsLinkedWith(ctx, tenant, cid, model.CONTACT, "CONNECTED_WITH", userId, model.USER)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "CommonReadRepository.IsLinkedWith"))
				return err
			}

			if !isLinkedWith {
				evt := generic.LinkEntityWithEntity{
					BaseEvent: event.BaseEvent{
						EventName:  generic.LinkEntityWithEntityV1,
						Tenant:     tenant,
						CreatedAt:  utils.Now(),
						AppSource:  constants.AppSourceDataUpkeeper,
						Source:     "WECONNECT",
						EntityId:   cid,
						EntityType: model.CONTACT,
					},
					WithEntityId:   userId,
					WithEntityType: model.USER,
					Relationship:   "CONNECTED_WITH",
				}

				err = s.eventBufferService.ParkBaseEvent(ctx, &evt, tenant, utils.Now().Add(time.Minute*1))
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "ParkBaseEvent"))
					return err
				}
			}
		}
	}

	return nil
}

func (s *contactService) linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	orphanContacts, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsEnrichedNotLinkedToOrganization(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	span.LogFields(log.Int("orphanContactsCount", len(orphanContacts)))

	for _, orpanContact := range orphanContacts {
		tenant := orpanContact.Tenant

		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		scrapIn, err := s.commonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, orpanContact.FieldStr1, entity.ScrapInFlowPersonProfile)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if scrapIn != nil && scrapIn.Success && scrapIn.CompanyFound {

			var scrapinContactResponse entity.ScrapInResponseBody
			err := json.Unmarshal([]byte(scrapIn.Data), &scrapinContactResponse)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
				return
			}

			domain := s.commonServices.DomainService.ExtractDomainFromOrganizationWebsite(ctx, scrapinContactResponse.Company.WebsiteUrl)
			if domain == "" {
				continue
			}

			organizationByDomainNode, err := s.commonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, tenant, domain)
			if err != nil {
				//TODO uncomment when data is fixed in DB
				//tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationByDomain"))
				//return
				continue
			}

			if organizationByDomainNode != nil {
				organizationId := utils.GetStringPropOrEmpty(organizationByDomainNode.Props, "id")

				positionName := ""
				if len(scrapinContactResponse.Person.Positions.PositionHistory) > 0 {
					for _, position := range scrapinContactResponse.Person.Positions.PositionHistory {
						if position.Title != "" && position.CompanyName != "" && position.CompanyName == scrapinContactResponse.Company.Name {
							positionName = position.Title
							break
						}
					}
				}

				_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
					return s.commonServices.GrpcClients.ContactClient.LinkWithOrganization(ctx, &contactpb.LinkWithOrganizationGrpcRequest{
						Tenant:         tenant,
						ContactId:      orpanContact.ContactId,
						OrganizationId: organizationId,
						JobTitle:       positionName,
					})
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "ContactClient.LinkWithOrganization"))
					return
				}
			}
		}

	}
}

func (s *contactService) findEmailsWithBetterContact(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.findEmailsWithBetterContact")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// Better contact is limited to 60 requests per minute
	// https://bettercontact.notion.site/Documentation-API-e8e1b352a0d647ee9ff898609bf1a168
	limit := 50

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		minutesFromLastContactUpdate := 2
		records, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsToFindWorkEmailWithBetterContact(ctx, minutesFromLastContactUpdate, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		for _, record := range records {
			enrichmentResponse, err := s.callEnrichmentApiFindWorkEmail(ctx, record)
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Object("record", record))
			} else {
				// mark contact with enrich requested
				err = s.commonServices.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, record.Tenant, record.ContactId, neo4jentity.ContactPropertyFindWorkEmailWithBetterContactRequestedId, enrichmentResponse.BetterContactRequestId)
				if err != nil {
					tracing.TraceErr(span, err)
				}
				err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, record.Tenant, model.NodeLabelContact, record.ContactId, string(neo4jentity.ContactPropertyFindWorkEmailWithBetterContactRequestedAt), utils.NowPtr())
				if err != nil {
					tracing.TraceErr(span, err)
				}
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *contactService) callEnrichmentApiFindWorkEmail(ctx context.Context, details neo4jrepository.ContactsEnrichWorkEmail) (*enrichmentmodel.FindWorkEmailResponse, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.callEnrichmentApiFindWorkEmail")
	defer span.Finish()

	requestJSON, err := json.Marshal(enrichmentmodel.FindWorkEmailRequest{
		LinkedinUrl:   details.LinkedInUrl,
		FirstName:     details.ContactFirstName,
		LastName:      details.ContactLastName,
		CompanyName:   details.OrganizationName,
		CompanyDomain: details.OrganizationDomain,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("GET", s.cfg.EnrichmentApi.Url+"/findWorkEmail", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, s.cfg.EnrichmentApi.ApiKey)
	req.Header.Set(security.TenantHeader, details.Tenant)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.status.findWorkEmail", response.StatusCode))

	var findWorkEmailApiResponse enrichmentmodel.FindWorkEmailResponse
	err = json.NewDecoder(response.Body).Decode(&findWorkEmailApiResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode find work email response"))
		return nil, err
	}
	return &findWorkEmailApiResponse, nil
}

func (s *contactService) enrichWithWorkEmailFromBetterContact(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.enrichWithWorkEmailFromBetterContact")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 250

	records, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsToEnrichWithEmailFromBetterContact(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, record := range records {

		detailsBetterContact, err := s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetByRequestId(ctx, record.FieldStr1)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if detailsBetterContact == nil {
			tracing.TraceErr(span, errors.New("better contact details by request id not found"))

			detailsBetterContact, err = s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetById(ctx, record.FieldStr1)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			if detailsBetterContact == nil {
				tracing.TraceErr(span, errors.New("better contact details by id not found"))
				continue
			}
		}

		if detailsBetterContact.Response == "" {
			continue
		}

		var betterContactResponse entity.BetterContactResponseBody
		if err = json.Unmarshal([]byte(detailsBetterContact.Response), &betterContactResponse); err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if len(betterContactResponse.Data) > 0 && betterContactResponse.Data[0].ContactEmailAddress != "" {

			emailIdIfExists, err := s.commonServices.Neo4jRepositories.EmailReadRepository.GetEmailIdIfExists(ctx, record.Tenant, betterContactResponse.Data[0].ContactEmailAddress)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
			_, err = utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
				contact := model.CONTACT.String()
				return s.commonServices.GrpcClients.EmailClient.UpsertEmail(ctx, &emailpb.UpsertEmailGrpcRequest{
					Tenant:       record.Tenant,
					Id:           emailIdIfExists,
					RawEmail:     betterContactResponse.Data[0].ContactEmailAddress,
					LinkWithType: &contact,
					LinkWithId:   &record.ContactId,
					SourceFields: &commonpb.SourceFields{
						AppSource: constants.AppSourceDataUpkeeper,
					},
				})
			})

			if err != nil {
				tracing.TraceErr(span, err)
				return
			}
		}

		err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, record.Tenant, model.NodeLabelContact, record.ContactId, "techFindWorkEmailWithBetterContactCompletedAt", utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}
}

func (s *contactService) checkBetterContactRequestsWithoutResponse(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.checkBetterContactRequestsWithoutResponse")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	betterContactRequestsWithoutResponse, err := s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetWithoutResponses(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, record := range betterContactRequestsWithoutResponse {

		// Create HTTP client
		client := &http.Client{}

		// Create POST request
		req, err := http.NewRequest("GET", fmt.Sprintf("%s?api_key=%s", s.cfg.BetterContactApi.Url+"/"+record.RequestID, s.cfg.BetterContactApi.ApiKey), nil)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")

		//Perform the request
		resp, err := client.Do(req)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if responseBody == nil || string(responseBody) == "Retry later" {
			return
		}

		// Parse the JSON request body
		var betterContactResponse entity.BetterContactResponseBody
		if err = json.Unmarshal(responseBody, &betterContactResponse); err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if betterContactResponse.Status == "terminated" {
			err = s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.AddResponse(ctx, record.RequestID, string(responseBody))
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}
			// store billable events
			// first check if it was requested externally
			personEnrichmentRequest, err := s.commonServices.PostgresRepositories.CosApiEnrichPersonTempResultRepository.GetByBettercontactRecordId(ctx, betterContactResponse.Id)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to check if bettercontact record was requested from person enrichment"))
			} else if personEnrichmentRequest != nil {
				emailFound, phoneFound := false, false
				for _, item := range betterContactResponse.Data {
					if item.ContactEmailAddress != "" {
						emailFound = true
					}
					if item.ContactPhoneNumber != nil && fmt.Sprintf("%v", item.ContactPhoneNumber) != "" {
						phoneFound = true
					}
				}
				if emailFound {
					_, err = s.commonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, personEnrichmentRequest.Tenant, entity.BillableEventEnrichPersonEmailFound, personEnrichmentRequest.BettercontactRecordId, "generated in upkeeper")
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
					}
				}
				if phoneFound {
					_, err = s.commonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, personEnrichmentRequest.Tenant, entity.BillableEventEnrichPersonPhoneFound, personEnrichmentRequest.BettercontactRecordId, "generated in upkeeper")
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
					}
				}
			}
		}
	}
}

func (s *contactService) EnrichContacts() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	s.enrichContacts(ctx)
}

func (s *contactService) enrichContacts(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.enrichContacts")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 20

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		minutesFromLastContactUpdate := 2
		minutesFromLastContactEnrichAttempt := 1 * 24 * 60 // 1 day
		minutesFromLastFailure := 10 * 24 * 60             // 10 days
		records, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsToEnrich(ctx, minutesFromLastContactUpdate, minutesFromLastContactEnrichAttempt, minutesFromLastFailure, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting socials: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
				return s.commonServices.GrpcClients.ContactClient.EnrichContact(ctx, &contactpb.EnrichContactGrpcRequest{
					Tenant:    record.Tenant,
					ContactId: record.ContactId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error enriching contact {%s}: %s", record.ContactId, err.Error())
			}
			// mark contact with enrich requested
			err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, record.Tenant, model.NodeLabelContact, record.ContactId, string(neo4jentity.ContactPropertyEnrichRequestedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error updating contact' enrich requested: %s", err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}
