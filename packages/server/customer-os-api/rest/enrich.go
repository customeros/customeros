package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
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
	emailTypeProfessional   = "professional"
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
	Type        *string `json:"type,omitempty" example:"professional"`
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
	StartMonth int  `json:"startMonth" example:"1"`
	StartYear  int  `json:"startYear" example:"2020"`
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

// EnrichOrganizationResponse represents the response for the organization enrichment API.
// @Description Response structure for the organization enrichment API.
// @example 200 {object} EnrichOrganizationResponse
type EnrichOrganizationResponse struct {
	// Status of the response.
	// Example: success
	Status string `json:"status" example:"success"`
	// Message for the response.
	// Example: Enrichment completed
	Message string                 `json:"message,omitempty" example:"Enrichment completed"`
	Data    EnrichOrganizationData `json:"data"`
}

// EnrichOrganizationData represents the detailed data about the organization from enrichment.
// @Description Detailed data about an organization from enrichment.
type EnrichOrganizationData struct {
	// Name of the organization.
	// Example: Acme Corporation
	Name string `json:"name" example:"Acme Corporation"`

	// Domain of the organization's website.
	// Example: acme.com
	Domain string `json:"domain" example:"acme.com"`

	// Short description of the organization.
	// Example: A global leader in innovative solutions.
	ShortDescription string `json:"description" example:"A global leader in innovative solutions"`

	// Long description of the organization.
	// Example: Acme Corporation provides cutting-edge technology solutions across the globe.
	LongDescription string `json:"longDescription" example:"Acme Corporation provides cutting-edge technology solutions across the globe."`

	// Website URL of the organization.
	// Example: https://acme.com
	Website string `json:"website" example:"https://acme.com"`

	// Number of employees in the organization.
	// Example: 5000
	Employees int `json:"employees" example:"5000"`

	// Year the organization was founded.
	// Example: 1995
	FoundedYear int `json:"foundedYear" example:"1995"`

	// Indicates whether the organization is publicly traded.
	// Example: true
	Public bool `json:"public,omitempty" example:"true"`

	// List of logo URLs for the organization.
	// Example: ["https://acme.com/logo.png"]
	Logos []string `json:"logos" example:"https://acme.com/logo.png"`

	// List of icon URLs for the organization.
	// Example: ["https://acme.com/icon.png"]
	Icons []string `json:"icons" example:"https://acme.com/icon.png"`

	// Industry in which the organization operates.
	Industry EnrichOrganizationIndustry `json:"industry"`

	// List of social media URLs for the organization.
	// Example: ["https://linkedin.com/company/acme"]
	Socials []string `json:"socials" example:"https://linkedin.com/company/acme"`

	// Location information about the organization.
	Location EnrichOrganizationLocation `json:"location"`
}

type EnrichOrganizationIndustry struct {
	// Industry in which the organization operates.
	// Example: Technology
	Industry string `json:"industry" example:"Technology"`
}

// EnrichOrganizationLocation represents the location details of an organization.
// @Description Location details of an organization.
type EnrichOrganizationLocation struct {
	// Indicates if the location is the headquarters.
	// Example: true
	IsHeadquarter bool `json:"isHeadquarter" example:"true"`

	// Country of the organization.
	// Example: United States
	Country string `json:"country" example:"United States"`

	// ISO Alpha-2 code of the country.
	// Example: US
	CountryCodeA2 string `json:"countryCodeA2" example:"US"`

	// City or locality of the organization.
	// Example: San Francisco
	City string `json:"city" example:"San Francisco"`

	// Region or state of the organization.
	// Example: California
	Region string `json:"region" example:"California"`

	// Postal code of the organization's location.
	// Example: 94105
	PostalCode string `json:"postalCode" example:"94105"`

	// Address line 1 of the organization.
	// Example: 123 Main St
	AddressLine1 string `json:"addressLine1" example:"123 Main St"`

	// Address line 2 of the organization (optional).
	// Example: Suite 100
	AddressLine2 string `json:"addressLine2" example:"Suite 100"`
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

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				ErrorResponse{
					Status:  "error",
					Message: "Missing tenant context",
				})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))

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
				ErrorResponse{
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
				ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if enrichPersonApiResponse.Data == nil || enrichPersonApiResponse.PersonFound == false {
			c.JSON(http.StatusOK,
				ErrorResponse{
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
		findWorkEmailApiResponse, err := callApiFindWorkEmail(ctx, services, span, *enrichPersonApiResponse, enrichPhoneNumber)
		if err != nil || findWorkEmailApiResponse == nil {
			c.JSON(http.StatusInternalServerError,
				ErrorResponse{
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
					ErrorResponse{
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
			response.ResultURL = services.Cfg.Services.CustomerOsApiUrl + enrichPersonAcceptedUrl + "/" + dbRecord.ID.String()
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
				_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventEnrichPersonEmailFound, betterContactResponseBody.Id,
					fmt.Sprintf("Email: %s, LinkedIn: %s, FirstName: %s, LastName: %s", email, linkedinUrl, firstName, lastName))
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
				}
			}
			if phoneFound {
				_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventEnrichPersonPhoneFound, betterContactResponseBody.Id,
					fmt.Sprintf("Email: %s, LinkedIn: %s, FirstName: %s, LastName: %s", email, linkedinUrl, firstName, lastName))
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
				}
			}
		}

		// call email verify
		for i := range response.Data.Emails {
			emailRecord := &response.Data.Emails[i]
			if emailRecord.Address != "" {
				emailValidationResult, err := callApiValidateEmail(ctx, services, emailRecord.Address, true)
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
						emailValidationResult.Data.EmailData.IsFreeAccount ||
						emailValidationResult.Data.EmailData.IsMailboxFull ||
						!emailValidationResult.Data.DomainData.IsPrimaryDomain)
				if emailValidationResult.Data.EmailData.IsFreeAccount {
					emailRecord.Type = utils.StringPtr(emailTypePersonal)
				} else {
					emailRecord.Type = utils.StringPtr(emailTypeProfessional)
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

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				ErrorResponse{
					Status:  "error",
					Message: "Missing tenant context",
				})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))

		tempId := c.Param("id")
		span.LogFields(log.String("request.tempId", tempId))

		getTempRecord, err := services.CommonServices.PostgresRepositories.CosApiEnrichPersonTempResultRepository.GetById(ctx, tempId, tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get temp record"))
			c.JSON(http.StatusInternalServerError,
				ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if getTempRecord == nil {
			c.JSON(http.StatusNotFound,
				ErrorResponse{
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
				ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if scrapInDbRecord == nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin record"))
			c.JSON(http.StatusNotFound,
				ErrorResponse{
					Status:  "error",
					Message: "Record not found",
				})
			return
		}

		betterContactDbRecord, err := services.CommonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetById(ctx, getTempRecord.BettercontactRecordId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get bettercontact record"))
			c.JSON(http.StatusInternalServerError,
				ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if betterContactDbRecord == nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get bettercontact record"))
			c.JSON(http.StatusNotFound,
				ErrorResponse{
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
				ErrorResponse{
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
					ErrorResponse{
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
			response.ResultURL = services.Cfg.Services.CustomerOsApiUrl + enrichPersonAcceptedUrl + "/" + tempId
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
				emailValidationResult, err := callApiValidateEmail(ctx, services, email.Address, true)
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
						emailValidationResult.Data.EmailData.IsFreeAccount ||
						emailValidationResult.Data.EmailData.IsMailboxFull ||
						!emailValidationResult.Data.DomainData.IsPrimaryDomain)
				if emailValidationResult.Data.EmailData.IsFreeAccount {
					email.Type = utils.StringPtr(emailTypePersonal)
				} else {
					email.Type = utils.StringPtr(emailTypeProfessional)
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

// @Summary Enrich Organization Information
// @Description Enriches an organization's information using the domain or other details.
// @Tags Enrichment API
// @Param linkedinUrl query string false "LinkedIn profile URL of the organization"
// @Param domain query string false "Domain of the organization"
// @Success 200 {object} EnrichOrganizationResponse "Enrichment results including organizational data"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Produce json
// @Router /enrich/v1/organization [get]
func EnrichOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichOrganization", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				ErrorResponse{
					Status:  "error",
					Message: "Missing tenant context",
				})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))

		linkedinUrl := c.Query("linkedinUrl")
		domain := c.Query("domain")

		// check linked in or email params are present
		if strings.TrimSpace(linkedinUrl) == "" && strings.TrimSpace(domain) == "" {
			c.JSON(http.StatusBadRequest,
				ErrorResponse{
					Status:  "error",
					Message: "Missing required parameters linkedinUrl or domain",
				})
			return
		}
		span.LogFields(
			log.String("request.domain", domain),
			log.String("request.linkedinUrl", linkedinUrl))

		// Call enrichPerson API
		enrichOrganizationApiResponse, err := callApiEnrichOrganization(ctx, services, span, linkedinUrl, domain)
		if err != nil || enrichOrganizationApiResponse.Status == "error" {
			c.JSON(http.StatusInternalServerError,
				ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if enrichOrganizationApiResponse.Success == false {
			c.JSON(http.StatusOK,
				ErrorResponse{
					Status:  "warning",
					Message: "Organization not found",
				})
			return
		}
		// Compose the response
		response := EnrichOrganizationResponse{
			Status:  "success",
			Message: "Enrichment completed",
			Data: EnrichOrganizationData{
				Name:             enrichOrganizationApiResponse.Data.Name,
				Domain:           enrichOrganizationApiResponse.Data.Domain,
				ShortDescription: enrichOrganizationApiResponse.Data.ShortDescription,
				LongDescription:  enrichOrganizationApiResponse.Data.LongDescription,
				Website:          enrichOrganizationApiResponse.Data.Website,
				Employees:        int(enrichOrganizationApiResponse.Data.Employees),
				FoundedYear:      int(enrichOrganizationApiResponse.Data.FoundedYear),
				Public:           utils.BoolDefaultIfNil(enrichOrganizationApiResponse.Data.Public, false),
				Logos:            enrichOrganizationApiResponse.Data.Logos,
				Icons:            enrichOrganizationApiResponse.Data.Icons,
				Industry: EnrichOrganizationIndustry{
					Industry: enrichOrganizationApiResponse.Data.Industry,
				},
				Socials: enrichOrganizationApiResponse.Data.Socials,
				Location: EnrichOrganizationLocation{
					IsHeadquarter: utils.BoolDefaultIfNil(enrichOrganizationApiResponse.Data.Location.IsHeadquarter, false),
					Country:       enrichOrganizationApiResponse.Data.Location.Country,
					CountryCodeA2: enrichOrganizationApiResponse.Data.Location.CountryCodeA2,
					City:          enrichOrganizationApiResponse.Data.Location.Locality,
					Region:        enrichOrganizationApiResponse.Data.Location.Region,
					PostalCode:    enrichOrganizationApiResponse.Data.Location.PostalCode,
					AddressLine1:  enrichOrganizationApiResponse.Data.Location.AddressLine1,
					AddressLine2:  enrichOrganizationApiResponse.Data.Location.AddressLine2,
				},
			},
		}

		if enrichOrganizationApiResponse.Success == true {
			_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventEnrichOrganizationSuccess, "",
				fmt.Sprintf("LinkedIn URL: %s, Domain: %s", linkedinUrl, domain))
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to register billable event"))
			}
		}

		c.JSON(http.StatusOK, response)
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
	req, err := http.NewRequest("GET", services.Cfg.Services.EnrichmentApiUrl+"/enrichPerson", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.Services.EnrichmentApiKey)
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

func callApiFindWorkEmail(ctx context.Context, services *service.Services, span opentracing.Span, enrichPersonApiResponse enrichmentmodel.EnrichPersonScrapinResponse, enrichPhoneNumber bool) (*enrichmentmodel.FindWorkEmailResponse, error) {
	companyName, companyDomain := "", ""
	if enrichPersonApiResponse.Data.PersonProfile.Company != nil {
		companyName = enrichPersonApiResponse.Data.PersonProfile.Company.Name
		companyDomain = services.CommonServices.DomainService.ExtractDomainFromOrganizationWebsite(ctx, enrichPersonApiResponse.Data.PersonProfile.Company.WebsiteUrl)
	}
	requestJSON, err := json.Marshal(enrichmentmodel.FindWorkEmailRequest{
		LinkedinUrl:       enrichPersonApiResponse.Data.PersonProfile.Person.LinkedInUrl,
		FirstName:         enrichPersonApiResponse.Data.PersonProfile.Person.FirstName,
		LastName:          enrichPersonApiResponse.Data.PersonProfile.Person.LastName,
		CompanyName:       companyName,
		CompanyDomain:     companyDomain,
		EnrichPhoneNumber: enrichPhoneNumber,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("GET", services.Cfg.Services.EnrichmentApiUrl+"/findWorkEmail", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.Services.EnrichmentApiKey)
	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

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

func callApiEnrichOrganization(ctx context.Context, services *service.Services, span opentracing.Span, linkedinUrl, domain string) (*enrichmentmodel.EnrichOrganizationResponse, error) {
	requestJSON, err := json.Marshal(enrichmentmodel.EnrichOrganizationRequest{
		Domain:      domain,
		LinkedinUrl: linkedinUrl,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("GET", services.Cfg.Services.EnrichmentApiUrl+"/enrichOrganization", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.Services.EnrichmentApiKey)
	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.status.enrichOrganization", response.StatusCode))

	var enrichOrganizationApiResponse enrichmentmodel.EnrichOrganizationResponse
	err = json.NewDecoder(response.Body).Decode(&enrichOrganizationApiResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode enrich organization response"))
		return nil, err
	}
	return &enrichOrganizationApiResponse, nil
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
			Duration: EnrichPersonJobDuration{
				StartMonth: position.StartEndDate.Start.Month,
				StartYear:  position.StartEndDate.Start.Year,
			},
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
