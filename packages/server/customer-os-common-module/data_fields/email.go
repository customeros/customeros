package data_fields

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// EmailSignatureContactInfo contains the basic information of a person
type EmailSignatureContactInfo struct {
	Name         string `json:"name"`
	JobTitle     string `json:"jobTitle"`
	Company      string `json:"company"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Mobile       string `json:"mobile"`
	LinkedIn     string `json:"linkedin"`
	GitHub       string `json:"github"`
	CalendarLink string `json:"calendarLink"`
}

// EmailSignatureCompanyInfo contains company information
type EmailSignatureCompanyInfo struct {
	Website   string                `json:"website"`
	LinkedIn  string                `json:"linkedin"`
	Twitter   string                `json:"twitter"`
	Youtube   string                `json:"youtube"`
	Instagram string                `json:"instagram"`
	GitHub    string                `json:"github"`
	Address   EmailSignatureAddress `json:"address"`
}

// EmailSignatureAddress contains address information
type EmailSignatureAddress struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	Region     string `json:"region"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

// EmailSignature represents the complete email signature
type EmailSignature struct {
	ContactInfo EmailSignatureContactInfo `json:"contactInfo"`
	CompanyInfo EmailSignatureCompanyInfo `json:"companyInfo"`
}

// EmailResponse is the top-level response structure
type EmailResponse struct {
	MessageBody  string         `json:"messageBody"`
	HasSignature bool           `json:"hasSignature"`
	Signature    EmailSignature `json:"signature,omitempty"`
}

// EmailSignatureValidator defines validation rules for email signatures
type EmailSignatureValidator struct {
	RequiresMessageBody  bool
	RequireBasicName     bool
	RequireBasicCompany  bool
	RequireContactEmail  bool
	MaxLengthMessageBody int
	MaxLengthName        int
	MaxLengthJobTitle    int
	MaxLengthCompany     int
	MaxLengthDisclaimer  int
	MaxLengthLegalText   int
	EmailRegex           *regexp.Regexp
	PhoneRegex           *regexp.Regexp
	WebsiteRegex         *regexp.Regexp
	LinkedInRegex        *regexp.Regexp
	TwitterRegex         *regexp.Regexp
	YoutubeRegex         *regexp.Regexp
	InstagramRegex       *regexp.Regexp
}

// NewEmailSignatureValidator creates and returns a new validator with default values
func NewEmailSignatureValidator() *EmailSignatureValidator {
	return &EmailSignatureValidator{
		RequiresMessageBody:  true,
		RequireBasicName:     true,
		RequireBasicCompany:  false,
		RequireContactEmail:  false,
		MaxLengthMessageBody: 3000,
		MaxLengthName:        100,
		MaxLengthJobTitle:    200,
		MaxLengthCompany:     100,
		EmailRegex:           regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		PhoneRegex:           regexp.MustCompile(`^(\+)?[0-9\(\)\-\.\s]{6,20}(x\d+)?$`),
		WebsiteRegex:         regexp.MustCompile(`^(http(s)?:\/\/)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}(\/[a-zA-Z0-9-._~:/?#[\]@!$&'()*+,;=]*)?$`),
		LinkedInRegex:        regexp.MustCompile(`^(http(s)?:\/\/)?(www\.)?linkedin\.com\/in\/[a-zA-Z0-9_-]+\/?$`),
		TwitterRegex:         regexp.MustCompile(`^(http(s)?:\/\/)?(www\.)?(twitter\.com|x\.com)\/[a-zA-Z0-9_]+\/?$`),
		YoutubeRegex:         regexp.MustCompile(`^(http(s)?:\/\/)?(www\.)?youtube\.com\/[a-zA-Z0-9_]+\/?$`),
		InstagramRegex:       regexp.MustCompile(`^(http(s)?:\/\/)?(www\.)?instagram\.com\/[a-zA-Z0-9_]+\/?$`),
	}
}

// IsValidJSON checks if the provided JSON string is valid
func (v *EmailSignatureValidator) IsValidJSON(jsonStr string) bool {
	var response EmailResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	return err == nil
}

// ValidateResponse validates the entire email signature response
func (v *EmailSignatureValidator) ValidateResponse(jsonStr string) (*EmailResponse, error) {
	// Parse the JSON
	var response EmailResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	if err := v.validateMessageBody(response.MessageBody); err != nil {
		return nil, err
	}

	if !response.HasSignature {
		return &response, nil
	}

	signature := response.Signature

	// Validate contact info section
	if err := v.validateContactInfo(signature.ContactInfo); err != nil {
		return nil, err
	}

	// Validate company info section
	if err := v.validateCompanyInfo(signature.CompanyInfo); err != nil {
		return nil, err
	}

	return &response, nil
}

func (v *EmailSignatureValidator) validateMessageBody(messageBody string) error {
	if v.RequiresMessageBody && strings.TrimSpace(messageBody) == "" {
		return fmt.Errorf("message body is required")
	}

	if len(messageBody) > v.MaxLengthMessageBody {
		return fmt.Errorf("message body exceeds maximum length: got %d characters, max is %d",
			len(messageBody), v.MaxLengthMessageBody)
	}
	return nil
}

// validateContactInfo validates the contact information section
func (v *EmailSignatureValidator) validateContactInfo(contactInfo EmailSignatureContactInfo) error {
	// Check required fields
	if v.RequireBasicName && strings.TrimSpace(contactInfo.Name) == "" {
		return fmt.Errorf("name is required in contact information section")
	}

	if v.RequireBasicCompany && strings.TrimSpace(contactInfo.Company) == "" {
		return fmt.Errorf("company is required in contact information section")
	}

	// Check field lengths
	if len(contactInfo.Name) > v.MaxLengthName {
		return fmt.Errorf("name exceeds maximum length: got %d characters, max is %d",
			len(contactInfo.Name), v.MaxLengthName)
	}

	if len(contactInfo.JobTitle) > v.MaxLengthJobTitle {
		return fmt.Errorf("job title exceeds maximum length: got %d characters, max is %d",
			len(contactInfo.JobTitle), v.MaxLengthJobTitle)
	}

	if len(contactInfo.Company) > v.MaxLengthCompany {
		return fmt.Errorf("company exceeds maximum length: got %d characters, max is %d",
			len(contactInfo.Company), v.MaxLengthCompany)
	}

	// Validate email if required
	if v.RequireContactEmail && strings.TrimSpace(contactInfo.Email) == "" {
		return fmt.Errorf("email is required in contact information section")
	}

	// Validate email format
	if contactInfo.Email != "" && !v.EmailRegex.MatchString(contactInfo.Email) {
		return fmt.Errorf("invalid email format: %s", contactInfo.Email)
	}

	// Validate phone formats
	if contactInfo.Phone != "" && !v.PhoneRegex.MatchString(contactInfo.Phone) {
		return fmt.Errorf("invalid phone format: %s", contactInfo.Phone)
	}

	if contactInfo.Mobile != "" && !v.PhoneRegex.MatchString(contactInfo.Mobile) {
		return fmt.Errorf("invalid mobile format: %s", contactInfo.Mobile)
	}

	// Validate LinkedIn URL
	if contactInfo.LinkedIn != "" && !v.LinkedInRegex.MatchString(contactInfo.LinkedIn) {
		return fmt.Errorf("invalid LinkedIn URL format: %s", contactInfo.LinkedIn)
	}

	// Validate GitHub (no specific validation for now)

	// Validate Calendar Link (no specific validation for now)

	return nil
}

// validateCompanyInfo validates the company information section
func (v *EmailSignatureValidator) validateCompanyInfo(companyInfo EmailSignatureCompanyInfo) error {
	// Validate website format
	if companyInfo.Website != "" && !v.WebsiteRegex.MatchString(companyInfo.Website) {
		return fmt.Errorf("invalid website format: %s", companyInfo.Website)
	}

	// Validate LinkedIn URL
	if companyInfo.LinkedIn != "" && !v.LinkedInRegex.MatchString(companyInfo.LinkedIn) {
		return fmt.Errorf("invalid company LinkedIn URL format: %s", companyInfo.LinkedIn)
	}

	// Validate Twitter URL
	if companyInfo.Twitter != "" && !v.TwitterRegex.MatchString(companyInfo.Twitter) {
		return fmt.Errorf("invalid Twitter URL format: %s", companyInfo.Twitter)
	}

	// Validate Youtube URL
	if companyInfo.Youtube != "" && !v.YoutubeRegex.MatchString(companyInfo.Youtube) {
		return fmt.Errorf("invalid Youtube URL format: %s", companyInfo.Youtube)
	}

	// Validate Instagram URL
	if companyInfo.Instagram != "" && !v.InstagramRegex.MatchString(companyInfo.Instagram) {
		return fmt.Errorf("invalid Instagram URL format: %s", companyInfo.Instagram)
	}

	// Validate address (minimal validation)
	return v.validateAddress(companyInfo.Address)
}

// validateAddress validates the address section
func (v *EmailSignatureValidator) validateAddress(address EmailSignatureAddress) error {
	// Address validation is minimal since formats vary widely
	// Could add country-specific postal code validation if needed
	return nil
}

// GetExpectedSchema returns a sample schema for documentation
func (v *EmailSignatureValidator) GetExpectedSchema() string {
	return `{
  "messageBody": "This is the clean message body without salutations, greetings, old threads, or messages, formatted in markdown.",
  "hasSignature": true,
  "signature": {
    "contactInfo": {
      "name": "Jane Smith",
      "jobTitle": "Senior Marketing Manager",
      "company": "Acme Corporation",
      "email": "jane.smith@acme.com",
      "phone": "+15551234567",
      "mobile": "+15559876543",
      "linkedin": "https://linkedin.com/in/janesmith",
      "github": "https://github.com/janesmith",
      "calendarLink": "https://calendly.com/janesmith"
    },
    "companyInfo": {
      "website": "https://www.acme.com",
      "linkedin": "https://linkedin.com/company/acme",
      "twitter": "https://twitter.com/acme",
      "youtube": "https://youtube.com/acme",
      "instagram": "https://instagram.com/acme",
      "github": "https://github.com/acme",
      "address": {
        "street": "123 Main Street",
        "city": "San Francisco",
        "region": "CA",
        "postalCode": "94105",
        "country": "USA"
      }
    }
  }
}`
}
