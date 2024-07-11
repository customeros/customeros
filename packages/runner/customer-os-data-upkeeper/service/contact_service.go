package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/forPelevin/gomoji"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	cosClient "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/client"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ContactService interface {
	UpkeepContacts()
	AskForWorkEmailOnBetterContact()
	EnrichWithWorkEmailFromBetterContact()
	CheckBetterContactRequestsWithoutResponse()
	EnrichContactsByEmail()
	SyncWeConnectContacts()
	LinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn()
}

type contactService struct {
	cfg                 *config.Config
	log                 logger.Logger
	commonServices      *commonService.Services
	customerOSApiClient cosClient.CustomerOSApiClient
}

func NewContactService(cfg *config.Config, log logger.Logger, commonServices *commonService.Services, customerOSApiClient cosClient.CustomerOSApiClient) ContactService {
	return &contactService{
		cfg:                 cfg,
		log:                 log,
		commonServices:      commonServices,
		customerOSApiClient: customerOSApiClient,
	}
}

func (s *contactService) UpkeepContacts() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	s.removeEmptySocials(ctx)
	s.removeDuplicatedSocials(ctx)
	s.hideContactsWithGroupEmail(ctx)
}

func (s *contactService) removeEmptySocials(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.removeEmptySocials")
	defer span.Finish()

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
		records, err := s.commonServices.Neo4jRepositories.SocialReadRepository.GetEmptySocialsForEntityType(ctx, neo4jutil.NodeLabelContact, minutesSinceLastUpdate, limit)
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
		records, err := s.commonServices.Neo4jRepositories.SocialReadRepository.GetDuplicatedSocialsForEntityType(ctx, neo4jutil.NodeLabelContact, minutesSinceLastUpdate, limit)
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

func (s *contactService) SyncWeConnectContacts() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.syncWeConnectContacts(ctx)
}

func (s *contactService) LinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn(ctx)
}

type WeConnectContactResponse struct {
	Linkedin           string `json:"linkedin"`
	LinkedinProfileUrl string `json:"linkedin_profile_url"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Company            string `json:"company"`
	Title              string `json:"title"`
	Industry           string `json:"industry"`
	Email              string `json:"email"`
	Education          string `json:"education"`
	Location           string `json:"location"`
	Connections        string `json:"connections"`
	Experience         []struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	} `json:"experience"`
	Campaigns            []string      `json:"campaigns"`
	CustomFields         []interface{} `json:"custom_fields"`
	ConnectedAt          string        `json:"connected_at"`
	TimestampConnectedAt int           `json:"timestamp_connected_at"`
	Id                   string        `json:"id"`
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

func (s *contactService) syncWeConnectContacts(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.syncWeConnectContacts")
	defer span.Finish()

	weConnectIntegrations, err := s.commonServices.PostgresRepositories.PersonalIntegrationRepository.FindByIntegration("weconnect")
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	span.LogFields(log.Int("integrationsCount", len(weConnectIntegrations)))

	for _, integration := range weConnectIntegrations {

		tenant := integration.TenantName

		page := 0

		total := 0
		skippedEmptyEmail := 0
		skippedExisting := 0
		created := 0
		addedSocial := 0

		for {

			select {
			case <-ctx.Done():
				s.log.Infof("Context cancelled, stopping")
				return
			default:
				// continue as normal
			}

			// Create new request
			req, err := http.NewRequest("GET", "https://api-us-1.we-connect.io/api/v1/connections?api_key="+integration.Secret+"&page="+fmt.Sprintf("%d", page), nil)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			req.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}
			defer res.Body.Close()

			responseBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			contactList := make([]WeConnectContactResponse, 0)

			// Convert response to map
			err = json.Unmarshal(responseBody, &contactList)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			if len(contactList) == 0 {
				break
			} else {
				page++
				total += len(contactList)
			}

			for _, contact := range contactList {

				linkedinProfileUrl := contact.LinkedinProfileUrl
				if linkedinProfileUrl != "" && linkedinProfileUrl[len(linkedinProfileUrl)-1] != '/' {
					linkedinProfileUrl = linkedinProfileUrl + "/"
				}

				linkedinProfileUrl = utils.NormalizeString(linkedinProfileUrl)

				contactsWithLinkedin, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsWithSocialUrl(ctx, tenant, linkedinProfileUrl)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				if len(contactsWithLinkedin) == 0 {
					created++
					contactInput := cosModel.ContactInput{
						FirstName: &contact.FirstName,
						LastName:  &contact.LastName,
						SocialURL: &linkedinProfileUrl,
						ExternalReference: &cosModel.ExternalSystemReferenceInput{
							Type:       "WECONNECT",
							ExternalID: contact.Id,
						},
					}

					_, err := s.customerOSApiClient.CreateContact(tenant, "", contactInput)
					if err != nil {
						tracing.TraceErr(span, err)
						return
					}
				}
			}

		}

		span.LogFields(log.Int("total", total), log.Int("created", created), log.Int("addedSocial", addedSocial), log.Int("skippedEmptyEmail", skippedEmptyEmail), log.Int("skippedExisting", skippedExisting))
	}

}

func (s *contactService) linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn")
	defer span.Finish()

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

		scrapIn, err := s.commonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, orpanContact.LinkedInUrl, entity.ScrapInFlowPersonProfile)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if scrapIn != nil && scrapIn.Success && scrapIn.CompanyFound {

			var scrapinContactResponse entity.ScrapInContactResponse
			err := json.Unmarshal([]byte(scrapIn.Data), &scrapinContactResponse)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
				return
			}

			if scrapinContactResponse.Company.WebsiteUrl == "" {
				continue
			}

			domain := utils.ExtractDomain(scrapinContactResponse.Company.WebsiteUrl)

			organizationByDomainNode, err := s.commonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationWithDomain(ctx, tenant, domain)
			if err != nil {
				//TODO uncomment when data is fixed in DB
				//tracing.TraceErr(span, errors.Wrap(err, "OrganizationReadRepository.GetOrganizationWithDomain"))
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

	if s.cfg.BetterContactApi.ApiKey == "" {
		err := errors.New("BetterContact API key is not set")
		tracing.TraceErr(span, err)
		s.log.Errorf("BetterContact API key is not set: %s", err.Error())
		return
	}
	if s.cfg.BetterContactApi.Url == "" {
		err := errors.New("BetterContact API URL is not set")
		tracing.TraceErr(span, err)
		s.log.Errorf("BetterContact API URL is not set: %s", err.Error())
		return
	}

	// Better contact is limited to 60 requests per minute
	// https://bettercontact.notion.site/Documentation-API-e8e1b352a0d647ee9ff898609bf1a168
	limit := 1

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
			requestId, err := s.requestBetterContactToFindEmail(ctx, record)
			if err != nil {
				tracing.TraceErr(span, err)
			} else {
				// mark contact with enrich requested
				err = s.commonServices.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, record.Tenant, record.ContactId, "techFindWorkEmailWithBetterContactRequestId", requestId)
				if err != nil {
					tracing.TraceErr(span, err)
				}
				err = s.commonServices.Neo4jRepositories.ContactWriteRepository.UpdateTimeProperty(ctx, record.Tenant, record.ContactId, "techFindWorkEmailWithBetterContactRequestedAt", utils.NowPtr())
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

// check if there is already a request for the same contact ( by linkedin url or by name and company )
// if data exists, mark as completed
// if data doesn't exist, create a new request
func (s *contactService) requestBetterContactToFindEmail(ctx context.Context, details neo4jrepository.ContactsEnrichWorkEmail) (string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.requestBetterContactToFindEmail")
	defer span.Finish()

	// replace special characters
	details.ContactFirstName = utils.NormalizeString(details.ContactFirstName)
	details.ContactLastName = utils.NormalizeString(details.ContactLastName)
	details.OrganizationName = utils.NormalizeString(details.OrganizationName)
	details.OrganizationDomain = utils.NormalizeString(details.OrganizationDomain)

	// strip special characters
	details.ContactFirstName = gomoji.RemoveEmojis(details.ContactFirstName)
	details.ContactLastName = gomoji.RemoveEmojis(details.ContactLastName)
	details.OrganizationName = gomoji.RemoveEmojis(details.OrganizationName)
	details.OrganizationDomain = gomoji.RemoveEmojis(details.OrganizationDomain)

	details.ContactFirstName = strings.TrimSpace(details.ContactFirstName)
	details.ContactLastName = strings.TrimSpace(details.ContactLastName)
	details.OrganizationName = strings.TrimSpace(details.OrganizationName)
	details.OrganizationDomain = strings.TrimSpace(details.OrganizationDomain)

	var existingBetterContactData *entity.EnrichDetailsBetterContact

	if details.LinkedInUrl != "" {
		betterContactByLinkedInUrl, err := s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetByLinkedInUrl(ctx, details.LinkedInUrl)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", fmt.Errorf("failed to get better contact details: %v", err)
		}

		if betterContactByLinkedInUrl != nil {
			existingBetterContactData = betterContactByLinkedInUrl
		}
	} else {
		detailsBetterContactList, err := s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetBy(ctx, details.ContactFirstName, details.ContactLastName, details.OrganizationName, details.OrganizationDomain)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", fmt.Errorf("failed to get better contact details: %v", err)
		}

		if detailsBetterContactList != nil && len(detailsBetterContactList) > 0 {
			existingBetterContactData = detailsBetterContactList[0]
		}
	}

	if existingBetterContactData != nil {
		return existingBetterContactData.ID.String(), nil
	}

	requestBodyDtls := BetterContactRequestBody{}

	requestBodyDtls.Data = []BetterContactData{
		{
			FirstName:     details.ContactFirstName,
			LastName:      details.ContactLastName,
			LinkedInUrl:   details.LinkedInUrl,
			Company:       details.OrganizationName,
			CompanyDomain: details.OrganizationDomain,
		},
	}

	requestBodyDtls.Webhook = s.cfg.BetterContactApi.CallbackUrl

	// Marshal request body to JSON
	requestBody, err := json.Marshal(requestBodyDtls)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?api_key=%s", s.cfg.BetterContactApi.Url, s.cfg.BetterContactApi.ApiKey), bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	//Decode response body
	var responseBody BetterContactResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to decode response body: %v", err)
	}

	err = s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.RegisterRequest(ctx, entity.EnrichDetailsBetterContact{
		RequestID:          responseBody.ID,
		ContactFirstName:   details.ContactFirstName,
		ContactLastName:    details.ContactLastName,
		ContactLinkedInUrl: details.LinkedInUrl,
		CompanyName:        details.OrganizationName,
		CompanyDomain:      details.OrganizationDomain,
		Request:            string(requestBody),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return responseBody.ID, nil
}

func (s *contactService) enrichWithWorkEmailFromBetterContact(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.enrichWithWorkEmailFromBetterContact")
	defer span.Finish()

	records, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsToEnrichWithEmailFromBetterContact(ctx, 1)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, record := range records {

		detailsBetterContact, err := s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetByRequestId(ctx, record.RequestId)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if detailsBetterContact == nil {
			tracing.TraceErr(span, errors.New("better contact details not found"))
			continue
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
				contact := enum.CONTACT.String()
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

		err = s.commonServices.Neo4jRepositories.ContactWriteRepository.UpdateTimeProperty(ctx, record.Tenant, record.ContactId, "techFindWorkEmailWithBetterContactCompletedAt", utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}
}

func (s *contactService) checkBetterContactRequestsWithoutResponse(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.checkBetterContactRequestsWithoutResponse")
	defer span.Finish()

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

		requestBody, err := io.ReadAll(resp.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		// Parse the JSON request body
		var betterContactResponse entity.BetterContactResponseBody
		if err = json.Unmarshal(requestBody, &betterContactResponse); err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if betterContactResponse.Status == "terminated" {
			err = s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.AddResponse(ctx, record.RequestID, string(requestBody))
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}
		}
	}
}

func (s *contactService) EnrichContactsByEmail() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	now := utils.Now()

	s.enrichContactsByEmail(ctx, now)
}

func (s *contactService) enrichContactsByEmail(ctx context.Context, now time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "enrichContactsByEmail")
	defer span.Finish()

	limit := 100

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
		minutesFromLastFailure := 7 * 24 * 60              // 7 days
		records, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsToEnrichByEmail(ctx, minutesFromLastContactUpdate, minutesFromLastContactEnrichAttempt, minutesFromLastFailure, limit)
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
			err = s.commonServices.Neo4jRepositories.ContactWriteRepository.UpdateTimeProperty(ctx, record.Tenant, record.ContactId, neo4jentity.ContactPropertyEnrichRequestedAt, utils.NowPtr())
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
