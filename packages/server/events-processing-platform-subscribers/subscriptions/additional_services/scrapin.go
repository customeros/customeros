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

func (h *ScrapInService) ScrapInPersonProfile(ctx context.Context, tenant, linkedInUrl string) (ScrapInContactResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapInService.ScrapInPersonProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("linkedInUrl", linkedInUrl))

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

	url := baseUrl + "/enrichment/profile" + "?apikey=" + scrapInApiKey + "&linkedInUrl=" + linkedInUrl

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
		LinkedInUrl: linkedInUrl,
	}
	requestJson, err := json.Marshal(requestParams)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "json.Marshal"))
		h.log.Errorf("Error marshalling request params: %s", err.Error())
		return ScrapInContactResponse{}, err
	}
	queryResult := h.postgresRepository.EnrichDetailsScrapInRepository.Add(ctx, postgresentity.EnrichDetailsScrapIn{
		Param1:        linkedInUrl,
		Flow:          postgresentity.ScrapInFlowPersonProfile,
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

func (h *ScrapInService) ScrapInPersonSearch(ctx context.Context, tenant, email, firstName, lastName, domain string) (ScrapInContactResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapInService.ScrapInPersonSearch")
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
	queryResult := h.postgresRepository.EnrichDetailsScrapInRepository.Add(ctx, postgresentity.EnrichDetailsScrapIn{
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
