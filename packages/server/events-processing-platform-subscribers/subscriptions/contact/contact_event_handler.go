package contact

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
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
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	socialpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/social"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type ScrapInContactResponse struct {
	Success       bool                   `json:"success"`
	Email         string                 `json:"email"`
	EmailType     string                 `json:"emailType"`
	CreditsLeft   int                    `json:"credits_left"`
	RateLimitLeft int                    `json:"rate_limit_left"`
	Person        *ScrapinPersonDetails  `json:"person,omitempty"`
	Company       *ScrapinCompanyDetails `json:"company,omitempty"`
}

type ScrapinPersonDetails struct {
	PublicIdentifier   string `json:"publicIdentifier"`
	LinkedInIdentifier string `json:"linkedInIdentifier"`
	LinkedInUrl        string `json:"linkedInUrl"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Headline           string `json:"headline"`
	Location           string `json:"location"`
	Summary            string `json:"summary"`
	PhotoUrl           string `json:"photoUrl"`
	CreationDate       struct {
		Month int `json:"month"`
		Year  int `json:"year"`
	} `json:"creationDate"`
	FollowerCount int `json:"followerCount"`
	Positions     struct {
		PositionsCount  int `json:"positionsCount"`
		PositionHistory []struct {
			Title        string `json:"title"`
			CompanyName  string `json:"companyName"`
			Description  string `json:"description"`
			StartEndDate struct {
				Start struct {
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"start"`
				End struct {
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"end"`
			} `json:"startEndDate"`
			CompanyLogo string `json:"companyLogo"`
			LinkedInUrl string `json:"linkedInUrl"`
			LinkedInId  string `json:"linkedInId"`
		} `json:"positionHistory"`
	} `json:"positions"`
	Schools struct {
		EducationsCount  int `json:"educationsCount"`
		EducationHistory []struct {
			DegreeName   string      `json:"degreeName"`
			FieldOfStudy string      `json:"fieldOfStudy"`
			Description  interface{} `json:"description"` // Can be null, so use interface{}
			LinkedInUrl  string      `json:"linkedInUrl"`
			SchoolLogo   string      `json:"schoolLogo"`
			SchoolName   string      `json:"schoolName"`
			StartEndDate struct {
				Start struct {
					Month *int `json:"month"` // Can be null, so use pointer
					Year  *int `json:"year"`  // Can be null, so use pointer
				} `json:"start"`
				End struct {
					Month *int `json:"month"` // Can be null, so use pointer
					Year  *int `json:"year"`  // Can be null, so use pointer
				} `json:"end"`
			} `json:"startEndDate"`
		} `json:"educationHistory"`
	} `json:"schools"`
	Skills    []interface{} `json:"skills"`    // Can be empty, so use interface{}
	Languages []interface{} `json:"languages"` // Can be empty, so use interface{}
}

type ScrapinCompanyDetails struct {
	LinkedInId         string `json:"linkedInId"`
	Name               string `json:"name"`
	UniversalName      string `json:"universalName"`
	LinkedInUrl        string `json:"linkedInUrl"`
	EmployeeCount      int    `json:"employeeCount"`
	EmployeeCountRange struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"employeeCountRange"`
	WebsiteUrl    string      `json:"websiteUrl"`
	Tagline       interface{} `json:"tagline"` // Can be null, so use interface{}
	Description   string      `json:"description"`
	Industry      string      `json:"industry"`
	Phone         interface{} `json:"phone"` // Can be null, so use interface{}
	Specialities  []string    `json:"specialities"`
	FollowerCount int         `json:"followerCount"`
	Headquarter   struct {
		City           string      `json:"city"`
		Country        string      `json:"country"`
		PostalCode     string      `json:"postalCode"`
		GeographicArea string      `json:"geographicArea"`
		Street1        string      `json:"street1"`
		Street2        interface{} `json:"street2"` // Can be null, so use interface{}
	} `json:"headquarter"`
	Logo string `json:"logo"`
}

type ScrapInPersonSearchRequest struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	CompanyDomain string `json:"companyDomain"`
	Email         string `json:"email"`
}

type ScrapInPersonProfileRequest struct {
}

type ContactEventHandler struct {
	repositories *repository.Repositories
	log          logger.Logger
	cfg          *config.Config
	caches       caches.Cache
	grpcClients  *grpc_client.Clients
}

func NewContactEventHandler(repositories *repository.Repositories, log logger.Logger, cfg *config.Config, caches caches.Cache, grpcClients *grpc_client.Clients) *ContactEventHandler {
	return &ContactEventHandler{
		repositories: repositories,
		log:          log,
		cfg:          cfg,
		caches:       caches,
		grpcClients:  grpcClients,
	}
}

func (h *ContactEventHandler) OnEnrichContactRequested(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnEnrichContactRequested")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.ContactRequestEnrich
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	return h.enrichContact(ctx, eventData.Tenant, contactId)
}

func (h *ContactEventHandler) enrichContact(ctx context.Context, tenant, contactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.enrichContact")
	defer span.Finish()

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

	// get email from contact
	email, err := h.getContactEmail(ctx, tenant, contactId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "getContactEmail"))
		h.log.Errorf("Error getting contact email: %s", err.Error())
		return err
	}
	if email != "" {
		return h.enrichContactByEmail(ctx, tenant, email, contactEntity)
	}

	return nil
}

func (h *ContactEventHandler) enrichContactByEmail(ctx context.Context, tenant, email string, contactEntity *neo4jentity.ContactEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.enrichContactByEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("email", email))
	span.LogFields(log.String("contactId", contactEntity.Id))

	// get enrich details for email
	record, err := h.repositories.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, email, postgresentity.ScrapInFlowPersonSearch)
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

		var scrapinContactResponse ScrapInContactResponse
		err := json.Unmarshal([]byte(record.Data), &scrapinContactResponse)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
			h.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
			return err
		}
		if scrapinContactResponse.Success || record.CreatedAt.After(daysAgo30) {
			return h.enrichContactWithScrapInEnrichDetails(ctx, tenant, email, contactEntity, scrapinContactResponse)
		}
	}

	domains, err := h.repositories.Neo4jRepositories.ContactReadRepository.GetLinkedOrgDomains(ctx, tenant, contactEntity.Id)
	domain := ""
	emailDomain := utils.ExtractDomainFromEmail(email)
	if utils.Contains(domains, emailDomain) {
		domain = emailDomain
	} else if len(domains) > 0 {
		domain = domains[0]
	}
	if domain == "" {
		domain = emailDomain
	}
	firstName, lastName := contactEntity.DeriveFirstAndLastNames()
	scrapinContactResponse, err := h.scrapInPersonSearch(ctx, tenant, email, firstName, lastName, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "scrapInPersonSearch"))
		return err
	}
	return h.enrichContactWithScrapInEnrichDetails(ctx, tenant, email, contactEntity, scrapinContactResponse)
}

func (h *ContactEventHandler) scrapInPersonSearch(ctx context.Context, tenant, email, firstName, lastName, domain string) (ScrapInContactResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.scrapInPersonSearch")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("email", email), log.String("firstName", firstName), log.String("lastName", lastName), log.String("domain", domain))

	baseUrl := h.cfg.Services.ScrapInApiUrl
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		tracing.TraceErr(span, err)
		h.log.Errorf("Brandfetch URL not set")
		return ScrapInContactResponse{}, err
	}
	scrapInApiKey := h.cfg.Services.ScrapInApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		tracing.TraceErr(span, err)
		h.log.Errorf("Scrapin Api key not set")
		return ScrapInContactResponse{}, err
	}

	url := baseUrl + "/enrichment" + "?apikey=" + scrapInApiKey + "&email=" + email
	if firstName != "" {
		url += "&firstName=" + firstName
	}
	if lastName != "" {
		url += "&lastName=" + lastName
	}
	if domain != "" {
		url += "&companyDomain=" + domain
	}

	body, err := makeScrapInHTTPRequest(url)

	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "makeScrapInHTTPRequest"))
		h.log.Errorf("Error making scrapin HTTP request: %s", err.Error())
		return ScrapInContactResponse{}, err
	}

	var scrapinResponse ScrapInContactResponse
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
		span.LogFields(log.String("response.body", string(body)))
		h.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return ScrapInContactResponse{}, err
	}

	bodyAsString := string(body)
	requestParams := ScrapInPersonSearchRequest{
		FirstName:     firstName,
		LastName:      lastName,
		CompanyDomain: domain,
		Email:         email,
	}
	requestJson, err := json.Marshal(requestParams)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "json.Marshal"))
		h.log.Errorf("Error marshalling request params: %s", err.Error())
		return ScrapInContactResponse{}, err
	}
	queryResult := h.repositories.PostgresRepositories.EnrichDetailsScrapInRepository.Add(ctx, postgresentity.EnrichDetailsScrapIn{
		Param1:        email,
		Param2:        firstName,
		Param3:        lastName,
		Param4:        domain,
		Flow:          postgresentity.ScrapInFlowPersonSearch,
		AllParamsJson: string(requestJson),
		Data:          bodyAsString,
		Success:       scrapinResponse.Success,
		PersonFound:   scrapinResponse.Person != nil,
		CompanyFound:  scrapinResponse.Company != nil,
	})
	if queryResult.Error != nil {
		tracing.TraceErr(span, errors.Wrap(queryResult.Error, "EnrichDetailsScrapInRepository.Add"))
		h.log.Errorf("Error saving enriching domain results: %v", queryResult.Error.Error())
		return ScrapInContactResponse{}, queryResult.Error
	}
	return scrapinResponse, nil
}

func makeScrapInHTTPRequest(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	return body, err
}

func (h *ContactEventHandler) getContactEmail(ctx context.Context, tenant string, contact string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.getContactEmail")
	defer span.Finish()

	records, err := h.repositories.Neo4jRepositories.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, neo4jenum.CONTACT, []string{contact})
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

func (h *ContactEventHandler) enrichContactWithScrapInEnrichDetails(ctx context.Context, tenant, email string, contact *neo4jentity.ContactEntity, scrapinContactResponse ScrapInContactResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.enrichContactWithScrapInEnrichDetails")
	defer span.Finish()

	if !scrapinContactResponse.Success {
		// mark contact as failed to enrich
		err := h.repositories.Neo4jRepositories.ContactWriteRepository.UpdateTimeProperty(ctx, tenant, contact.Id, neo4jentity.ContactPropertyEnrichedFailedAtScrapInPersonSearch, utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateTimeProperty"))
			h.log.Errorf("Error updating enriched at scrap in person search property: %s", err.Error())
		}

		err = h.repositories.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contact.Id, neo4jentity.ContactPropertyEnrichedScrapInPersonSearchParam, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateAnyProperty"))
			h.log.Errorf("Error updating enriched scrap in person search param property: %s", err.Error())
		}

		return nil
	}

	// if person is not found, return
	if scrapinContactResponse.Person == nil {
		span.LogFields(log.String("result", "person not found"))
		h.log.Infof("Person not found for email %s", email)
		return nil
	}

	// update contact
	tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	upsertContactGrpcRequest := contactpb.UpsertContactGrpcRequest{
		Id:     contact.Id,
		Tenant: tenant,
		SourceFields: &commonpb.SourceFields{
			Source:    constants.SourceOpenline,
			AppSource: "scrapin",
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

	// add social profiles
	if scrapinContactResponse.Person.LinkedInUrl != "" {
		_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
			return h.grpcClients.ContactClient.AddSocial(ctx, &contactpb.ContactAddSocialGrpcRequest{
				ContactId:      contact.Id,
				Tenant:         tenant,
				Url:            scrapinContactResponse.Person.LinkedInUrl,
				Alias:          scrapinContactResponse.Person.PublicIdentifier,
				FollowersCount: int64(scrapinContactResponse.Person.FollowerCount),
				SourceFields: &commonpb.SourceFields{
					Source:    constants.SourceOpenline,
					AppSource: "scrapin",
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ContactClient.AddSocial"))
			h.log.Errorf("Error adding social profile: %s", err.Error())
		}
	}

	// mark contact as enriched
	nowPtr := utils.NowPtr()
	err := h.repositories.Neo4jRepositories.ContactWriteRepository.UpdateTimeProperty(ctx, tenant, contact.Id, neo4jentity.ContactPropertyEnrichedAt, nowPtr)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateTimeProperty"))
		h.log.Errorf("Error updating enriched at property: %s", err.Error())
	}

	err = h.repositories.Neo4jRepositories.ContactWriteRepository.UpdateTimeProperty(ctx, tenant, contact.Id, neo4jentity.ContactPropertyEnrichedAtScrapInPersonSearch, nowPtr)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error updating enriched at scrap in person search property: %s", err.Error())
	}

	err = h.repositories.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contact.Id, neo4jentity.ContactPropertyEnrichedScrapInPersonSearchParam, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ContactWriteRepository.UpdateAnyProperty"))
		h.log.Errorf("Error updating enriched scrap in person search param property: %s", err.Error())
	}

	return nil
}
