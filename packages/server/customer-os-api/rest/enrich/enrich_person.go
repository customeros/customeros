package restenrich

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	restverify "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/verify"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

const (
	emailTypePersonal       = "personal"
	emailTypeWork           = "work"
	enrichPersonAcceptedUrl = "/enrich/v1/person/results"
)

// EnrichPersonResponse represents the response for the person enrichment API.
// @Description Response structure for the person enrichment API.
// @example 200 {object} EnrichPersonResponse
type EnrichPersonResponse struct {
	Status        string           `json:"status" example:"success"`
	Message       string           `json:"message,omitempty" example:"Enrichment completed"`
	IsComplete    bool             `json:"isComplete" example:"true"`
	PendingFields []string         `json:"pendingFields,omitempty" example:"[\"email\", \"phone number\"]"`
	ResultURL     string           `json:"resultUrl,omitempty" example:"https://api.customeros.ai/enrich/v1/person/results/550e8400-e29b-41d4-a716-446655440000"`
	Data          EnrichPersonData `json:"data"`
}

// EnrichPersonData represents detailed data about a person from enrichment.
// @Description Detailed data about a person from enrichment.
type EnrichPersonData struct {
	Emails       []EnrichPersonEmail       `json:"emails"`
	Jobs         []EnrichPersonJob         `json:"jobs"`
	Location     EnrichPersonLocation      `json:"location"`
	Name         EnrichPersonName          `json:"name"`
	PhoneNumbers []EnrichPersonPhoneNumber `json:"phoneNumbers"`
	ProfilePic   string                    `json:"profilePic" example:"https://example.com/profile.jpg"`
	Social       EnrichPersonSocial        `json:"social"`
}

// EnrichPersonEmail represents the email details of a person.
// @Description Email details of a person.
type EnrichPersonEmail struct {
	Address     string  `json:"address" example:"john.doe@example.com"`
	Deliverable *string `json:"deliverable,omitempty" example:"true"`
	IsRisky     *bool   `json:"isRisky,omitempty" example:"false"`
	Type        *string `json:"type,omitempty" example:"work"`
}

// EnrichPersonJob represents the job details of a person.
// @Description Job details of a person.
type EnrichPersonJob struct {
	Title           string                  `json:"title" example:"Software Engineer"`
	Seniority       string                  `json:"seniority" example:"Senior"`
	Duration        EnrichPersonJobDuration `json:"duration"`
	Company         string                  `json:"company" example:"Tech Corp"`
	CompanyLinkedin string                  `json:"companyLinkedin" example:"https://linkedin.com/company/techcorp"`
	CompanyWebsite  string                  `json:"companyWebsite" example:"https://techcorp.com"`
	IsCurrent       bool                    `json:"isCurrent" example:"true"`
}

// EnrichPersonJobDuration represents the duration of a person's job.
// @Description Job duration of a person.
type EnrichPersonJobDuration struct {
	StartMonth *int `json:"startMonth,omitempty" example:"1"`
	StartYear  *int `json:"startYear,omitempty" example:"2020"`
	EndMonth   *int `json:"endMonth,omitempty" example:"12"`
	EndYear    *int `json:"endYear,omitempty" example:"2023"`
}

// EnrichPersonLocation represents the location details of a person.
// @Description Location details of a person.
type EnrichPersonLocation struct {
	City     string `json:"city" example:"San Francisco"`
	Region   string `json:"region" example:"California"`
	Country  string `json:"country" example:"USA"`
	Timezone string `json:"timezone" example:"PST"`
}

// EnrichPersonName represents the name details of a person.
// @Description Name details of a person.
type EnrichPersonName struct {
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	FullName  string `json:"fullName" example:"John Doe"`
}

// EnrichPersonPhoneNumber represents the phone number details of a person.
// @Description Phone number details of a person.
type EnrichPersonPhoneNumber struct {
	Number string `json:"number" example:"+1234567890"`
	Type   string `json:"type" example:"mobile"`
}

// EnrichPersonSocial represents the social media details of a person.
// @Description Social media details of a person.
type EnrichPersonSocial struct {
	Linkedin EnrichPersonLinkedIn `json:"linkedin"`
	X        EnrichPersonX        `json:"x"`
	Github   EnrichPersonGithub   `json:"github"`
	Discord  EnrichPersonDiscord  `json:"discord"`
}

// EnrichPersonLinkedIn represents the LinkedIn profile details of a person.
// @Description LinkedIn profile details of a person.
type EnrichPersonLinkedIn struct {
	ID            string `json:"id" example:"123456789"`
	PublicID      string `json:"publicId" example:"john-doe"`
	URL           string `json:"url" example:"https://linkedin.com/in/john-doe"`
	FollowerCount int    `json:"followerCount" example:"500"`
}

// EnrichPersonX represents the X (formerly Twitter) profile details of a person.
// @Description X (formerly Twitter) profile details of a person.
type EnrichPersonX struct {
	Handle string `json:"handle" example:"@johndoe"`
	URL    string `json:"url" example:"https://x.com/johndoe"`
}

// EnrichPersonGithub represents the Github profile details of a person.
// @Description Github profile details of a person.
type EnrichPersonGithub struct {
	Username string `json:"username" example:"johndoe"`
	URL      string `json:"url" example:"https://github.com/johndoe"`
}

// EnrichPersonDiscord represents the Discord profile details of a person.
// @Description Discord profile details of a person.
type EnrichPersonDiscord struct {
	Username string `json:"username" example:"johndoe#1234"`
}

// @Summary Enrich Person Information
// @Description Enriches a person's information using LinkedIn URL, email, and other optional details.
// @Tags Enrichment API
// @Param linkedinUrl query string false "LinkedIn profile URL of the person"
// @Param email query string false "Email address of the person"
// @Param firstName query string false "First name of the person"
// @Param lastName query string false "Last name of the person"
// @Param includeMobileNumber query string false "Include mobile phone number in the enrichment result" default(false)
// @Success 200 {object} EnrichPersonResponse "Enrichment results including personal, job, and social data"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Produce json
// @Router /enrich/v1/person [get]
func EnrichPerson(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichPerson", c.Request.Header)
		defer span.Finish()
		commontracing.TagComponentRest(span)
		commontracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing tenant context",
				})
			return
		}

		linkedinUrl := c.Query("linkedinUrl")
		email := c.Query("email")
		firstName := c.Query("firstName")
		lastName := c.Query("lastName")
		enrichPhoneNumberParam := c.Query("includeMobileNumber")
		// convert to boolean
		enrichPhoneNumber := enrichPhoneNumberParam == "true"

		// check linked in or email params are present
		if strings.TrimSpace(linkedinUrl) == "" && strings.TrimSpace(email) == "" {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing required parameters linkedinUrl or email",
				})
			return
		}

		span.LogFields(
			log.String("request.email", email),
			log.String("request.linkedinUrl", linkedinUrl),
			log.Bool("request.includeMobileNumber", enrichPhoneNumber),
			log.String("request.firstName", firstName),
			log.String("request.lastName", lastName))

		// Call enrichPerson API
		enrichPersonApiResponse, err := callApiEnrichPerson(ctx, services, span, email, linkedinUrl, firstName, lastName)
		if err != nil || enrichPersonApiResponse == nil || enrichPersonApiResponse.Status == "error" {
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if enrichPersonApiResponse.Data == nil || enrichPersonApiResponse.PersonFound == false {
			c.JSON(http.StatusOK,
				rest.ErrorResponse{
					Status:  "warning",
					Message: "Person not found",
				})
			return
		}

		// map the response to the customer-os model
		scrapInPersonResponse := postgresentity.ScrapInResponseBody{}
		if enrichPersonApiResponse.Data != nil && enrichPersonApiResponse.Data.PersonProfile != nil {
			scrapInPersonResponse = *enrichPersonApiResponse.Data.PersonProfile
		}
		enrichedPersonData := mapPersonScrapInData(&scrapInPersonResponse)

		// Call findWorkEmail API
		companyName, companyDomain := "", ""
		if enrichPersonApiResponse.Data.PersonProfile.Company != nil {
			companyName = enrichPersonApiResponse.Data.PersonProfile.Company.Name
			companyDomain, _ = services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, enrichPersonApiResponse.Data.PersonProfile.Company.WebsiteUrl)
		}

		findWorkEmailApiResponse, err := services.EnrichmentService.CallApiFindWorkEmail(ctx, enrichPersonApiResponse.Data.PersonProfile.Person.FirstName, enrichPersonApiResponse.Data.PersonProfile.Person.LastName,
			companyName, companyDomain, enrichPersonApiResponse.Data.PersonProfile.Person.LinkedInUrl, enrichPhoneNumber)
		if err != nil || findWorkEmailApiResponse == nil {
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}

		// Compose response
		response := EnrichPersonResponse{
			Status: "success",
		}
		if enrichedPersonData != nil {
			response.Data = *enrichedPersonData
		}

		betterContactResponseBody := findWorkEmailApiResponse.Data
		if betterContactResponseBody == nil {
			dbRecord, err := services.CommonServices.PostgresRepositories.CosApiEnrichPersonTempResultRepository.Create(ctx, postgresentity.CosApiEnrichPersonTempResult{
				Tenant:                tenant,
				ScrapinRecordId:       enrichPersonApiResponse.RecordId,
				BettercontactRecordId: findWorkEmailApiResponse.RecordId,
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to create temp result"))
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Internal error",
					})
				return
			}
			response.IsComplete = false
			response.PendingFields = []string{"email"}
			if enrichPhoneNumber {
				response.PendingFields = append(response.PendingFields, "phone number")
			}
			response.ResultURL = services.Cfg.InternalServices.CustomerOsApiUrl + enrichPersonAcceptedUrl + "/" + dbRecord.ID.String()
		} else {
			response.IsComplete = true
			emailFound, phoneFound := false, false
			for _, item := range betterContactResponseBody.Data {
				if item.ContactEmailAddress != "" {
					emailFound = true
					response.Data.Emails = append(response.Data.Emails, EnrichPersonEmail{
						Address: item.ContactEmailAddress,
					})
				}
				if enrichPhoneNumber {
					if item.ContactPhoneNumber != nil && fmt.Sprintf("%v", item.ContactPhoneNumber) != "" {
						phoneFound = true
						response.Data.PhoneNumbers = append(response.Data.PhoneNumbers, EnrichPersonPhoneNumber{
							Number: fmt.Sprintf("%v", item.ContactPhoneNumber),
							Type:   "mobile",
						})
					}
				}
			}
			if emailFound {
				_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventEnrichPersonEmailFound,
					postgresrepository.BillableEventDetails{
						ExternalID:    betterContactResponseBody.Id,
						ReferenceData: fmt.Sprintf("Email: %s, LinkedIn: %s, FirstName: %s, LastName: %s", email, linkedinUrl, firstName, lastName),
					})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
				}
			}
			if phoneFound {
				_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventEnrichPersonPhoneFound,
					postgresrepository.BillableEventDetails{
						ExternalID:    betterContactResponseBody.Id,
						ReferenceData: fmt.Sprintf("Email: %s, LinkedIn: %s, FirstName: %s, LastName: %s", email, linkedinUrl, firstName, lastName),
					})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
				}
			}
		}

		// call email verify
		for i := range response.Data.Emails {
			emailRecord := &response.Data.Emails[i]
			if emailRecord.Address != "" {
				emailValidationResult, err := restverify.CallApiValidateEmail(ctx, services, emailRecord.Address, true)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
					continue
				}
				if emailValidationResult.Status != "success" {
					tracing.TraceErr(span, errors.New("failed to validate email"))
					continue
				}
				emailRecord.Deliverable = utils.StringPtr(emailValidationResult.Data.EmailData.Deliverable)
				emailRecord.IsRisky = utils.BoolPtr(
					emailValidationResult.Data.DomainData.IsFirewalled ||
						emailValidationResult.Data.EmailData.IsRoleAccount ||
						emailValidationResult.Data.EmailData.IsSystemGenerated ||
						emailValidationResult.Data.EmailData.IsFreeAccount ||
						emailValidationResult.Data.EmailData.IsMailboxFull ||
						!emailValidationResult.Data.DomainData.IsPrimaryDomain)
				if emailValidationResult.Data.EmailData.IsFreeAccount {
					emailRecord.Type = utils.StringPtr(emailTypePersonal)
				} else {
					emailRecord.Type = utils.StringPtr(emailTypeWork)
				}
			}
		}

		responseStatusCode := http.StatusOK
		if !response.IsComplete {
			responseStatusCode = http.StatusAccepted
		}
		c.JSON(responseStatusCode, response)
	}
}

// @Summary Enrich Person Callback
// @Description Retrieves enriched person data from a temporary result based on the given ID.
// @Tags Enrichment API
// @Param id path string true "Temporary result ID"
// @Success 200 {object} EnrichPersonResponse "Enrichment results including personal, job, and social data"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Produce json
// @Router /enrich/v1/person/results/{id} [get]
func EnrichPersonCallback(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichPerson", c.Request.Header)
		defer span.Finish()
		commontracing.TagComponentRest(span)
		commontracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing tenant context",
				})
			return
		}

		tempId := c.Param("id")
		span.LogFields(log.String("request.tempId", tempId))

		getTempRecord, err := services.CommonServices.PostgresRepositories.CosApiEnrichPersonTempResultRepository.GetById(ctx, tempId, tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get temp record"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if getTempRecord == nil {
			c.JSON(http.StatusNotFound,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Record not found",
				})
			return
		}

		// Enrich person data
		scrapInDbRecord, err := services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetById(ctx, getTempRecord.ScrapinRecordId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin record"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if scrapInDbRecord == nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin record"))
			c.JSON(http.StatusNotFound,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Record not found",
				})
			return
		}

		betterContactDbRecord, err := services.CommonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetById(ctx, getTempRecord.BettercontactRecordId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get bettercontact record"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if betterContactDbRecord == nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get bettercontact record"))
			c.JSON(http.StatusNotFound,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Record not found",
				})
			return
		}

		// Compose response
		response := EnrichPersonResponse{
			Status: "success",
		}

		// extract scrapin data
		var scrapInPersonResponse postgresentity.ScrapInResponseBody
		err = json.Unmarshal([]byte(scrapInDbRecord.Data), &scrapInPersonResponse)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin record"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		enrichedPersonData := mapPersonScrapInData(&scrapInPersonResponse)

		if enrichedPersonData != nil {
			response.Data = *enrichedPersonData
		}

		// extract better contact data
		var betterContactResponseBody *postgresentity.BetterContactResponseBody
		if betterContactDbRecord.Response != "" {
			err = json.Unmarshal([]byte(betterContactDbRecord.Response), &betterContactResponseBody)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal bettercontact record"))
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Internal error",
					})
				return
			}
		}

		if betterContactResponseBody == nil {
			response.IsComplete = false
			response.PendingFields = []string{"email"}
			if betterContactDbRecord.EnrichPhoneNumber {
				response.PendingFields = append(response.PendingFields, "phone number")
			}
			response.ResultURL = services.Cfg.InternalServices.CustomerOsApiUrl + enrichPersonAcceptedUrl + "/" + tempId
		} else {
			response.IsComplete = true
			for _, item := range betterContactResponseBody.Data {
				if item.ContactEmailAddress != "" {
					response.Data.Emails = append(response.Data.Emails, EnrichPersonEmail{
						Address: item.ContactEmailAddress,
					})
				}
				if item.ContactPhoneNumber != nil && fmt.Sprintf("%v", item.ContactPhoneNumber) != "" {
					response.Data.PhoneNumbers = append(response.Data.PhoneNumbers, EnrichPersonPhoneNumber{
						Number: fmt.Sprintf("%v", item.ContactPhoneNumber),
						Type:   "mobile",
					})
				}
			}
		}

		// call email verify
		for i := range response.Data.Emails {
			email := &response.Data.Emails[i] // Get a pointer to the email in the slice
			if email.Address != "" {
				emailValidationResult, err := restverify.CallApiValidateEmail(ctx, services, email.Address, true)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
					continue
				}
				if emailValidationResult.Status != "success" {
					tracing.TraceErr(span, errors.New("failed to validate email"))
					continue
				}
				email.Deliverable = utils.StringPtr(emailValidationResult.Data.EmailData.Deliverable)
				email.IsRisky = utils.BoolPtr(
					emailValidationResult.Data.DomainData.IsFirewalled ||
						emailValidationResult.Data.EmailData.IsRoleAccount ||
						emailValidationResult.Data.EmailData.IsSystemGenerated ||
						emailValidationResult.Data.EmailData.IsFreeAccount ||
						emailValidationResult.Data.EmailData.IsMailboxFull ||
						!emailValidationResult.Data.DomainData.IsPrimaryDomain)
				if emailValidationResult.Data.EmailData.IsFreeAccount {
					email.Type = utils.StringPtr(emailTypePersonal)
				} else {
					email.Type = utils.StringPtr(emailTypeWork)
				}
			}
		}

		responseStatusCode := http.StatusOK
		if !response.IsComplete {
			responseStatusCode = http.StatusAccepted
		}
		c.JSON(responseStatusCode, response)
	}
}

func callApiEnrichPerson(ctx context.Context, services *service.Services, span opentracing.Span, email, linkedinUrl, firstName, lastName string) (*enrichmentmodel.EnrichPersonScrapinResponse, error) {
	requestJSON, err := json.Marshal(enrichmentmodel.EnrichPersonRequest{
		Email:       email,
		LinkedinUrl: linkedinUrl,
		FirstName:   firstName,
		LastName:    lastName,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("GET", services.Cfg.InternalServices.EnrichmentApiUrl+"/enrichPerson", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.InternalServices.EnrichmentApiKey)
	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.status.enrichPerson", response.StatusCode))

	var enrichPersonApiResponse enrichmentmodel.EnrichPersonScrapinResponse
	err = json.NewDecoder(response.Body).Decode(&enrichPersonApiResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode enrich person response"))
		return nil, err
	}
	return &enrichPersonApiResponse, nil
}

func mapPersonScrapInData(source *postgresentity.ScrapInResponseBody) *EnrichPersonData {
	if source == nil {
		return nil
	}
	output := EnrichPersonData{}

	// set emails
	if source.Email != "" {
		output.Emails = append(output.Emails, EnrichPersonEmail{
			Address: source.Email,
		})
	}

	// set name
	output.Name = EnrichPersonName{
		FirstName: source.Person.FirstName,
		LastName:  source.Person.LastName,
	}
	if source.Person.FirstName != "" && source.Person.LastName != "" {
		output.Name.FullName = source.Person.FirstName + " " + source.Person.LastName
	}

	// set jobs
	for _, position := range source.Person.Positions.PositionHistory {
		enrichPersonJob := EnrichPersonJob{
			Title:           position.Title,
			Company:         position.CompanyName,
			CompanyLinkedin: position.LinkedInUrl,
			IsCurrent:       position.StartEndDate.End == nil,
			//Seniority:       position.Seniority, // TODO will be implemented later after clarifications
		}
		if position.StartEndDate.Start != nil {
			enrichPersonJob.Duration.StartMonth = &position.StartEndDate.Start.Month
			enrichPersonJob.Duration.StartYear = &position.StartEndDate.Start.Year
		}
		if position.StartEndDate.End != nil {
			enrichPersonJob.Duration.EndMonth = &position.StartEndDate.End.Month
			enrichPersonJob.Duration.EndYear = &position.StartEndDate.End.Year
		}
		if source.Company != nil && source.Company.LinkedInId == position.LinkedInId {
			enrichPersonJob.CompanyWebsite = source.Company.WebsiteUrl
		}
		output.Jobs = append(output.Jobs, enrichPersonJob)
	}

	// set profile picture
	output.ProfilePic = source.Person.PhotoUrl

	// set social
	output.Social = EnrichPersonSocial{
		Linkedin: EnrichPersonLinkedIn{
			ID:            source.Person.LinkedInIdentifier,
			PublicID:      source.Person.PublicIdentifier,
			URL:           source.Person.LinkedInUrl,
			FollowerCount: source.Person.FollowerCount,
		},
		// TODO add X, github, other from Customer OS
	}

	// set location // TODO implement AI lookup to get details
	output.Location = EnrichPersonLocation{
		//City: source.Data.PersonProfile.Person.Location,
		Region: source.Person.Location,
		//Country:  source.Data.PersonProfile.Person.Country,
		//Timezone: source.Data.PersonProfile.Person.Timezone,
	}

	// set phone numbers
	// TODO implement phone numbers
	return &output
}
