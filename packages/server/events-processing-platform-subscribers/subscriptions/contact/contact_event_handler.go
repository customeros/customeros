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
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type ScrapInContactResponse struct {
	Success bool `json:"success"`
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
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting contact with id %s: %s", contactId, err.Error())
		return nil
	}
	contactEntity := neo4jmapper.MapDbNodeToContactEntity(contactDbNode)

	if contactEntity.EnrichDetails.EnrichedAt != nil {
		h.log.Infof("Contact %s already enriched", contactId)
		return nil
	}

	// get email from contact
	email, err := h.getContactEmail(ctx, tenant, contactId)
	if err != nil {
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return err
	}

	// if enrich details found and not older than 1 year, use the data, otherwise scrape it
	daysAgo365 := utils.Now().Add(-time.Hour * 24 * 365)
	if record != nil && record.CreatedAt.After(daysAgo365) && record.Data != "" {
		// enrich contact with the data found
		span.LogFields(log.Bool("result.email already enriched", true))

		var scrapinContactResponse ScrapInContactResponse
		err := json.Unmarshal([]byte(record.Data), &scrapinContactResponse)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
			return err
		}
		return h.enrichContactWithScrapInEnrichDetails(ctx, tenant, contactEntity, scrapinContactResponse)
	}

	scrapinContactResponse, err := h.scrapInPersonSearch(ctx, tenant, email, contactEntity.FirstName, contactEntity.LastName)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return h.enrichContactWithScrapInEnrichDetails(ctx, tenant, contactEntity, scrapinContactResponse)
}

func (h *ContactEventHandler) scrapInPersonSearch(ctx context.Context, tenant, email, firstName, lastName string) (ScrapInContactResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.scrapInPersonSearch")
	defer span.Finish()

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

	body, err := makeScrapInHTTPRequest(url)

	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error making scrapin HTTP request: %s", err.Error())
		return ScrapInContactResponse{}, err
	}

	var scrapinResponse ScrapInContactResponse
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return ScrapInContactResponse{}, err
	}

	if scrapinResponse.Success == false {
		err = errors.New("ScrapIn person search failed, returned success false")
		tracing.TraceErr(span, err)
		h.log.Errorf("ScrapIn person search failed, returned success false")
		return ScrapInContactResponse{}, err
	}

	bodyAsString := string(body)
	requestParams := ScrapInPersonSearchRequest{
		FirstName:     firstName,
		LastName:      lastName,
		CompanyDomain: "",
		Email:         email,
	}
	requestJson, err := json.Marshal(requestParams)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error marshalling request params: %s", err.Error())
		return ScrapInContactResponse{}, err
	}
	queryResult := h.repositories.PostgresRepositories.EnrichDetailsScrapInRepository.Add(ctx, postgresentity.EnrichDetailsScrapIn{
		Param1:        email,
		Flow:          postgresentity.ScrapInFlowPersonSearch,
		AllParamsJson: string(requestJson),
		Data:          bodyAsString,
	})
	if queryResult.Error != nil {
		tracing.TraceErr(span, queryResult.Error)
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
		tracing.TraceErr(span, err)
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

func (h *ContactEventHandler) enrichContactWithScrapInEnrichDetails(ctx context.Context, tenant string, entity *neo4jentity.ContactEntity, scrapinContactResponse ScrapInContactResponse) error {
	// TODO enrich contact with data
	return nil
}
