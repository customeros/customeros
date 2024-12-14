// @openapi 3.0.0
package customerbase

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
)

// Contact Types

// ContactRecord represents the request structure for creating a contact
// @Description Request to create a contact
type ContactRecord struct {
	// Contact's unique identifier
	// example: contact-123
	ContactId string `json:"contactId,omitempty"`

	// Contact's email address
	// example: john@example.com
	Email string `json:"email" csv:"email"`

	// Contact's LinkedIn profile URL
	// example: https://linkedin.com/in/john-doe
	LinkedInURL string `json:"linkedinUrl" csv:"linkedin_url"`
}

// SingleContactResponse represents a response containing a single contact
// @Description Response structure for single contact operations
type SingleContactResponse struct {
	// Inherits standard response fields
	enum.BaseResponse
	// The contact information
	Contact ContactRecord `json:"contact,omitempty"`
}

// ContactsResponse represents a response containing multiple contacts
// @Description Response structure for multiple contact operations
type ContactsResponse struct {
	// Inherits standard response fields
	enum.BaseResponse
	// List of contacts
	Contacts []ContactRecord `json:"contacts,omitempty"`
}

// Organization Types

// CreateOrganizationRequest represents the request body for creating a new organization
// @Description Request to create an organization
type CreateOrganizationRequest struct {
	// Organization's name
	// required: true
	// example: CustomerOS
	Name string `json:"name"`

	// Custom ID provided by the user
	// example: 12345
	CustomId string `json:"customId"`

	// Organization's website URL
	// example: https://customeros.ai
	Website string `json:"website"`

	// Organization's LinkedIn profile URL
	// example: https://linkedin.com/company/openline
	LinkedinUrl string `json:"linkedinUrl"`

	// Lead source of the organization
	// example: Web Search
	LeadSource string `json:"leadSource"`

	// Relationship status of the organization
	// example: customer
	Relationship string `json:"relationship"`

	// Indicates if the organization is an ICP (Ideal Customer Profile) fit
	// example: true
	IcpFit bool `json:"icpFit"`
}

// OrganizationRecord represents detailed organization information
// @Description Detailed organization information returned by API operations
type OrganizationRecord struct {
	// Organization's unique identifier
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id,omitempty"`

	// Custom ID provided by the user
	// example: 12345
	CustomId string `json:"customId,omitempty"`

	// CustomerOS unique identifier
	// example: C-A12-B45
	CosId string `json:"cosId,omitempty"`

	// Organization name
	// example: CustomerOS
	Name string `json:"name,omitempty"`

	// Organization's website URL
	// example: https://customeros.ai
	Website string `json:"website,omitempty"`

	// Organization's LinkedIn profile URL
	// example: https://linkedin.com/company/openline
	LinkedinUrl string `json:"linkedinUrl"`

	// Lead source of the organization
	// example: Web Search
	LeadSource string `json:"leadSource,omitempty"`

	// Relationship status with the organization
	// example: customer
	Relationship string `json:"relationship,omitempty"`

	// Current stage in the organization lifecycle
	// example: lead
	Stage string `json:"stage,omitempty"`

	// ICP fit indicator
	// example: true
	IcpFit bool `json:"icpFit,omitempty"`

	// Associated domains
	// example: ["customeros.com","customeros.ai"]
	Domains []string `json:"domains,omitempty"`

	// External system links
	ExternalLinks []ExternalLink `json:"externalLinks,omitempty"`
}

// OrganizationResponse represents a response containing a single organization
// @Description Response structure for single organization operations
type OrganizationResponse struct {
	// Inherits standard response fields
	enum.BaseResponse
	// The organization information
	Organization OrganizationRecord `json:"organization,omitempty"`
}

// OrganizationsResponse represents a response containing multiple organizations
// @Description Response structure for multiple organization operations
type OrganizationsResponse struct {
	// Inherits standard response fields
	enum.BaseResponse
	// List of organizations
	Organizations []OrganizationRecord `json:"organizations,omitempty"`
}

// External System Types

// SetPrimaryExternalSystemIdRequest represents the request to set a primary external system ID
// @Description Request to set or replace the primary external system ID
type SetPrimaryExternalSystemIdRequest struct {
	// The ID of the external system to be set as primary
	// required: true
	// example: stripe-1234
	ExternalId string `json:"externalId"`
}

// ExternalSystemRecord represents external system information
// @Description External system information and its relationship to an organization
type ExternalSystemRecord struct {
	// Associated organization ID
	// example: org-789
	OrganizationId string `json:"organizationId,omitempty"`

	// Name of the external system
	// example: stripe
	ExternalSystem string `json:"externalSystem,omitempty"`

	// External system identifier
	// example: stripe-1234
	ExternalId string `json:"externalId,omitempty"`

	// Indicates if this is the primary link
	// example: true
	Primary bool `json:"primary,omitempty"`
}

// ExternalSystemResponse represents a response containing external system information
// @Description Response structure for external system operations
type ExternalSystemResponse struct {
	// Inherits standard response fields
	enum.BaseResponse
	// The external system information
	Organization ExternalSystemRecord `json:"organization,omitempty"`
}

// Link Types

// ExternalLink represents a link to an external system
// @Description External system link information
type ExternalLink struct {
	// External system name
	// example: stripe
	Name string `json:"name"`

	// External system identifier
	// example: cos-12345
	Id string `json:"id"`

	// Indicates if this is the primary link
	// example: true
	Primary bool `json:"primary"`
}

// SocialLink represents a social media link
// @Description Social media link information
type SocialLink struct {
	// Social media profile URL
	// example: https://linkedin.com/company/openline
	Url string `json:"url"`

	// Number of followers
	// example: 1000
	FollowerCount int `json:"followerCount"`
}
