package data_fields

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// EmailSignatureBasics contains the basic information of a person
type EmailSignatureBasics struct {
	Name     string `json:"name"`
	JobTitle string `json:"jobTitle"`
	Company  string `json:"company"`
}

// EmailSignatureContact contains contact information
type EmailSignatureContact struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Mobile  string `json:"mobile"`
	Website string `json:"website"`
}

// EmailSignatureAddress contains address information
type EmailSignatureAddress struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	Region     string `json:"region"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

// EmailSignatureSocial contains social media information
type EmailSignatureSocial struct {
	LinkedIn  string `json:"linkedin"`
	Twitter   string `json:"twitter"`
	Youtube   string `json:"youtube"`
	Instagram string `json:"instagram"`
	GitHub    string `json:"github"`
}

// EmailSignatureAdditional contains additional information
type EmailSignatureAdditional struct {
	Disclaimer   string `json:"disclaimer"`
	LegalText    string `json:"legalText"`
	CalendarLink string `json:"calendarLink"`
}

// EmailSignature represents the complete email signature
type EmailSignature struct {
	Basics     EmailSignatureBasics     `json:"basics"`
	Contact    EmailSignatureContact    `json:"contact"`
	Address    EmailSignatureAddress    `json:"address"`
	Social     EmailSignatureSocial     `json:"social"`
	Additional EmailSignatureAdditional `json:"additional"`
}

// EmailSignatureResponse is the top-level response structure
type EmailResponse struct {
	MessageBody string         `json:"messageBody"`
	Signature   bool           `json:"signature"`
	Details     EmailSignature `json:"details"`
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
		RequireContactEmail:  true,
		MaxLengthMessageBody: 3000,
		MaxLengthName:        100,
		MaxLengthJobTitle:    200,
		MaxLengthCompany:     100,
		MaxLengthDisclaimer:  1000,
		MaxLengthLegalText:   1000,
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

	if !response.Signature {
		return &response, nil
	}

	signature := response.Details

	// Validate basics section
	if err := v.validateBasics(signature.Basics); err != nil {
		return nil, err
	}

	// Validate contact section
	if err := v.validateContact(signature.Contact); err != nil {
		return nil, err
	}

	// Validate address section
	if err := v.validateAddress(signature.Address); err != nil {
		return nil, err
	}

	// Validate social section
	if err := v.validateSocial(signature.Social); err != nil {
		return nil, err
	}

	// Validate additional section
	if err := v.validateAdditional(signature.Additional); err != nil {
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

// validateBasics validates the basics section
func (v *EmailSignatureValidator) validateBasics(basics EmailSignatureBasics) error {
	// Check required fields
	if v.RequireBasicName && strings.TrimSpace(basics.Name) == "" {
		return fmt.Errorf("name is required in basics section")
	}

	if v.RequireBasicCompany && strings.TrimSpace(basics.Company) == "" {
		return fmt.Errorf("company is required in basics section")
	}

	// Check field lengths
	if len(basics.Name) > v.MaxLengthName {
		return fmt.Errorf("name exceeds maximum length: got %d characters, max is %d",
			len(basics.Name), v.MaxLengthName)
	}

	if len(basics.JobTitle) > v.MaxLengthJobTitle {
		return fmt.Errorf("job title exceeds maximum length: got %d characters, max is %d",
			len(basics.JobTitle), v.MaxLengthJobTitle)
	}

	if len(basics.Company) > v.MaxLengthCompany {
		return fmt.Errorf("company exceeds maximum length: got %d characters, max is %d",
			len(basics.Company), v.MaxLengthCompany)
	}

	return nil
}

// validateContact validates the contact section
func (v *EmailSignatureValidator) validateContact(contact EmailSignatureContact) error {
	// Check required fields
	if v.RequireContactEmail && strings.TrimSpace(contact.Email) == "" {
		return fmt.Errorf("email is required in contact section")
	}

	// Validate email format
	if contact.Email != "" && !v.EmailRegex.MatchString(contact.Email) {
		return fmt.Errorf("invalid email format: %s", contact.Email)
	}

	// Validate phone formats
	if contact.Phone != "" && !v.PhoneRegex.MatchString(contact.Phone) {
		return fmt.Errorf("invalid phone format: %s", contact.Phone)
	}

	if contact.Mobile != "" && !v.PhoneRegex.MatchString(contact.Mobile) {
		return fmt.Errorf("invalid mobile format: %s", contact.Mobile)
	}

	// Validate website format
	if contact.Website != "" && !v.WebsiteRegex.MatchString(contact.Website) {
		return fmt.Errorf("invalid website format: %s", contact.Website)
	}

	return nil
}

// validateAddress validates the address section
func (v *EmailSignatureValidator) validateAddress(address EmailSignatureAddress) error {
	// Address validation is minimal since formats vary widely
	// Could add country-specific postal code validation if needed
	return nil
}

// validateSocial validates the social section
func (v *EmailSignatureValidator) validateSocial(social EmailSignatureSocial) error {
	// Validate LinkedIn URL
	if social.LinkedIn != "" && !v.LinkedInRegex.MatchString(social.LinkedIn) {
		return fmt.Errorf("invalid LinkedIn URL format: %s", social.LinkedIn)
	}

	// Validate Twitter URL
	if social.Twitter != "" && !v.TwitterRegex.MatchString(social.Twitter) {
		return fmt.Errorf("invalid Twitter URL format: %s", social.Twitter)
	}

	// Validate Youtube URL
	if social.Youtube != "" && !v.YoutubeRegex.MatchString(social.Twitter) {
		return fmt.Errorf("invalid Youtube URL format: %s", social.Twitter)
	}

	// Validate Instagram URL
	if social.Instagram != "" && !v.InstagramRegex.MatchString(social.Twitter) {
		return fmt.Errorf("invalid Instagram URL format: %s", social.Twitter)
	}

	return nil
}

// validateAdditional validates the additional section
func (v *EmailSignatureValidator) validateAdditional(additional EmailSignatureAdditional) error {
	// Check field lengths
	if len(additional.Disclaimer) > v.MaxLengthDisclaimer {
		return fmt.Errorf("disclaimer exceeds maximum length: got %d characters, max is %d",
			len(additional.Disclaimer), v.MaxLengthDisclaimer)
	}

	if len(additional.LegalText) > v.MaxLengthLegalText {
		return fmt.Errorf("legal text exceeds maximum length: got %d characters, max is %d",
			len(additional.LegalText), v.MaxLengthLegalText)
	}

	return nil
}

// GetExpectedSchema returns a sample schema for documentation
func (v *EmailSignatureValidator) GetExpectedSchema() string {
	return `{
  "messageBody": "This is the clean message body without salutations, greetings, old threads, or messages, formatted in markdown.",
  "signature": true,
  "details": {
    "basics": {
      "name": "Jane Smith",
      "jobTitle": "Senior Marketing Manager",
      "company": "Acme Corporation",
    },
    "contact": {
      "email": "jane.smith@acme.com",
      "phone": "+15551234567",
      "mobile": "+15559876543",
      "website": "https://www.acme.com"
    },
    "address": {
      "street": "123 Main Street",
      "city": "San Francisco",
      "region": "CA",
      "postalCode": "94105",
      "country": "USA"
    },
    "social": {
      "linkedin": "https://linkedin.com/in/janesmith",
      "twitter": "https://twitter.com/janesmith",
      "youtube": "https://youtube.com/janesmith",
      "instagram": "https://instagram.com/janesmith",
      "github": "https://github.com/janesmith",
    },
    "additional": {
      "disclaimer": "This email and any files transmitted with it are confidential and intended solely for the use of the individual or entity to whom they are addressed.",
      "legalText": "Acme Corporation is a registered trademark. Registration No. 12345.",
      "calendarLink": "https://calendly.com/janesmith",
    }
  }
}`
}
