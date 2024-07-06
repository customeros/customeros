package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	cosClient "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/client"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type ContactService interface {
	UpkeepContacts()
	FindEmails()
	EnrichContacts()
	SyncWeConnectContacts()
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

func (s *contactService) FindEmails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	s.findEmailsWithBetterContact(ctx)
}

func (s *contactService) SyncWeConnectContacts() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.syncWeConnectContacts(ctx)
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
	Data           []BetterContactData `json:"data"`
	Webhook        string              `json:"webhook"`
	VerifyCatchAll bool                `json:"verify_catch_all"`
}

type BetterContactData struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Company       string `json:"company"`
	CompanyDomain string `json:"company_domain"`
	LinkedinUrl   string `json:"linkedin_url"`
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
			req, err := http.NewRequest("GET", "https://api-us-1.we-connect.io/api/v1/contacts_with_email?api_key="+integration.Secret+"&page="+fmt.Sprintf("%d", page), nil)
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

				contactsWithLinkedin, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsWithSocialUrl(ctx, tenant, contact.LinkedinProfileUrl)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				if len(contactsWithLinkedin) == 0 {
					created++
					contactInput := cosModel.ContactInput{
						FirstName: &contact.FirstName,
						LastName:  &contact.LastName,
						SocialURL: &contact.LinkedinProfileUrl,
						ExternalReference: &cosModel.ExternalSystemReferenceInput{
							Type:       "WECONNECT",
							ExternalID: contact.Id,
						},
					}

					if contact.Email != "" {
						contactInput.Email = &cosModel.EmailInput{
							Email: contact.Email,
						}
					}

					_, err := s.customerOSApiClient.CreateContact(tenant, integration.Email, contactInput)
					if err != nil {
						tracing.TraceErr(span, err)
						return
					}
				} else if len(contactsWithLinkedin) == 1 {
					contactEntity := neo4jmapper.MapDbNodeToContactEntity(contactsWithLinkedin[0])
					socialEntities, err := s.commonServices.SocialService.GetAllForEntities(ctx, tenant, neo4jenum.CONTACT, []string{contactEntity.Id})
					if err != nil {
						tracing.TraceErr(span, err)
						return
					}

					hasLinkedIn := false

					for _, social := range *socialEntities {
						if social.Url == contact.LinkedinProfileUrl {
							hasLinkedIn = true
							break
						}
					}

					if !hasLinkedIn {
						_, err := s.customerOSApiClient.AddSocialToContact(tenant, contactEntity.Id, cosModel.SocialInput{
							URL: contact.LinkedinProfileUrl,
						})
						if err != nil {
							tracing.TraceErr(span, err)
							return
						}
						addedSocial++
					} else {
						skippedExisting++
					}
				}
			}

		}

		span.LogFields(log.Int("total", total), log.Int("created", created), log.Int("addedSocial", addedSocial), log.Int("skippedEmptyEmail", skippedEmptyEmail), log.Int("skippedExisting", skippedExisting))
	}

}

func (s *contactService) findEmailsWithBetterContact(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.findEmailsWithBetterContact")
	defer span.Finish()

	//TODO alexb unblock this when implementation is ready
	return

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
		records, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsToFindEmail(ctx, minutesFromLastContactUpdate, limit)
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
			err = s.requestBetterContactToFindEmail(ctx, record)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error requesting better contact to find email: %s", err.Error())
			} else {
				// mark contact with enrich requested
				err = s.commonServices.Neo4jRepositories.ContactWriteRepository.UpdateTimeProperty(ctx, record.Tenant, record.ContractId, "techFindEmailRequestedAt", utils.NowPtr())
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error updating contact' find email requested: %s", err.Error())
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

func (s *contactService) requestBetterContactToFindEmail(ctx context.Context, details neo4jrepository.TenantAndContactDetails) error {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.requestBetterContactToFindEmail")
	defer span.Finish()

	// TODO alexb implement it to get contact name, company and linked in
	requestBodyDtls := BetterContactRequestBody{}

	// Marshal request body to JSON
	requestBody, err := json.Marshal(requestBodyDtls)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?api_key=%s", s.cfg.BetterContactApi.Url, s.cfg.BetterContactApi.ApiKey), bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	// Decode response body
	var responseBody BetterContactResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to decode response body: %v", err)
	}

	result := s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.RegisterRequest(ctx, entity.EnrichDetailsBetterContact{
		Tenant:    details.Tenant,
		ContactID: details.ContractId,
		RequestID: responseBody.ID,
		Request:   string(requestBody),
	})
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return fmt.Errorf("failed to register better contact request: %s", result.Error.Error())
	}

	return nil
}

func (s *contactService) EnrichContacts() {
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
					ContactId: record.ContractId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error enriching contact {%s}: %s", record.ContractId, err.Error())
			}
			// mark contact with enrich requested
			err = s.commonServices.Neo4jRepositories.ContactWriteRepository.UpdateTimeProperty(ctx, record.Tenant, record.ContractId, neo4jentity.ContactPropertyEnrichRequestedAt, utils.NowPtr())
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
