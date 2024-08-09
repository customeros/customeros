package additional_services

import (
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"net/http"
)

type ScrapInEnrichContactFlow struct {
	Email string `json:"email"`
	Url   string `json:"url"`
}

func (s ScrapInEnrichContactFlow) GetFlow() postgresentity.ScrapInFlow {
	if s.Email != "" {
		return postgresentity.ScrapInFlowPersonSearch
	}
	if s.Url != "" {
		return postgresentity.ScrapInFlowPersonProfile
	}
	return ""
}

func (s ScrapInEnrichContactFlow) GetParam1() string {
	if s.Email != "" {
		return s.Email
	}
	if s.Url != "" {
		return s.Url
	}
	return ""
}

func (s ScrapInEnrichContactFlow) GetTimeLabel() neo4jentity.ContactProperty {
	if s.GetFlow() == postgresentity.ScrapInFlowPersonSearch {
		return neo4jentity.ContactPropertyEnrichedAtScrapInPersonSearch
	} else if s.GetFlow() == postgresentity.ScrapInFlowPersonProfile {
		return neo4jentity.ContactPropertyEnrichedAtScrapInProfile
	}
	return ""
}

func (s ScrapInEnrichContactFlow) GetParamLabel() neo4jentity.ContactProperty {
	if s.GetFlow() == postgresentity.ScrapInFlowPersonSearch {
		return neo4jentity.ContactPropertyEnrichedScrapInPersonSearchParam
	} else if s.GetFlow() == postgresentity.ScrapInFlowPersonProfile {
		return neo4jentity.ContactPropertyEnrichedScrapInProfileParam
	}
	return ""
}

type ScrapInPersonSearchRequest struct {
	FirstName     string `json:"firstName,omitempty"`
	LastName      string `json:"lastName,omitempty"`
	CompanyDomain string `json:"companyDomain,omitempty"`
	Email         string `json:"email,omitempty"`
	LinkedInUrl   string `json:"linkedInUrl,omitempty"`
}

type ScrapInService struct {
	log                logger.Logger
	cfg                *config.Config
	postgresRepository *postgresRepository.Repositories
}

func NewScrapInService(log logger.Logger, cfg *config.Config, postgresRepository *postgresRepository.Repositories) *ScrapInService {
	return &ScrapInService{
		log:                log,
		cfg:                cfg,
		postgresRepository: postgresRepository,
	}
}

// TODO switch to enrich api
// Deprecated
func (h *ScrapInService) ScrapInPersonProfile(ctx context.Context, tenant, linkedInUrl string) (postgresentity.ScrapInPersonResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapInService.ScrapInPersonProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("linkedInUrl", linkedInUrl))

	baseUrl := h.cfg.Services.ScrapInApiUrl
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		tracing.TraceErr(span, err)
		h.log.Errorf("Brandfetch URL not set")
		return postgresentity.ScrapInPersonResponse{}, err
	}
	scrapInApiKey := h.cfg.Services.ScrapInApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		tracing.TraceErr(span, err)
		h.log.Errorf("Scrapin Api key not set")
		return postgresentity.ScrapInPersonResponse{}, err
	}

	url := baseUrl + "/enrichment/profile" + "?apikey=" + scrapInApiKey + "&linkedInUrl=" + linkedInUrl

	body, err := makeScrapInHTTPRequest(url)

	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "makeScrapInHTTPRequest"))
		h.log.Errorf("Error making scrapin HTTP request: %s", err.Error())
		return postgresentity.ScrapInPersonResponse{}, err
	}

	var scrapinResponse postgresentity.ScrapInPersonResponse
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
		span.LogFields(log.String("response.body", string(body)))
		h.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return postgresentity.ScrapInPersonResponse{}, err
	}

	bodyAsString := string(body)
	requestParams := ScrapInPersonSearchRequest{
		LinkedInUrl: linkedInUrl,
	}
	requestJson, err := json.Marshal(requestParams)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "json.Marshal"))
		h.log.Errorf("Error marshalling request params: %s", err.Error())
		return postgresentity.ScrapInPersonResponse{}, err
	}
	_, err = h.postgresRepository.EnrichDetailsScrapInRepository.Create(ctx, postgresentity.EnrichDetailsScrapIn{
		Param1:        linkedInUrl,
		Flow:          postgresentity.ScrapInFlowPersonProfile,
		AllParamsJson: string(requestJson),
		Data:          bodyAsString,
		Success:       scrapinResponse.Success,
		PersonFound:   scrapinResponse.Person != nil,
		CompanyFound:  scrapinResponse.Company != nil,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "EnrichDetailsScrapInRepository.Create"))
		h.log.Errorf("Error saving enriching domain results: %v", err.Error())
		return postgresentity.ScrapInPersonResponse{}, err
	}
	return scrapinResponse, nil
}

func (h *ScrapInService) ScrapInPersonSearch(ctx context.Context, tenant, email, firstName, lastName, domain string) (postgresentity.ScrapInPersonResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapInService.ScrapInPersonSearch")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("email", email), log.String("firstName", firstName), log.String("lastName", lastName), log.String("domain", domain))

	baseUrl := h.cfg.Services.ScrapInApiUrl
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		tracing.TraceErr(span, err)
		h.log.Errorf("Brandfetch URL not set")
		return postgresentity.ScrapInPersonResponse{}, err
	}
	scrapInApiKey := h.cfg.Services.ScrapInApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		tracing.TraceErr(span, err)
		h.log.Errorf("Scrapin Api key not set")
		return postgresentity.ScrapInPersonResponse{}, err
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
		return postgresentity.ScrapInPersonResponse{}, err
	}

	var scrapinResponse postgresentity.ScrapInPersonResponse
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
		span.LogFields(log.String("response.body", string(body)))
		h.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return postgresentity.ScrapInPersonResponse{}, err
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
		return postgresentity.ScrapInPersonResponse{}, err
	}
	_, err = h.postgresRepository.EnrichDetailsScrapInRepository.Create(ctx, postgresentity.EnrichDetailsScrapIn{
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
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "EnrichDetailsScrapInRepository.Create"))
		h.log.Errorf("Error saving enriching domain results: %v", err.Error())
		return postgresentity.ScrapInPersonResponse{}, err
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
