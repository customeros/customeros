package customerbase

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"

// ContactRecord represents the request structure for creating a contact
// @Description Request to create a contact
type ContactRecord struct {
	// Contact's email address
	Email string `json:"email" example:"john@example.com" csv:"email"`

	// Contact's LinkedIn profile URL
	LinkedInURL string `json:"linkedinUrl" example:"https://linkedin.com/in/john-doe" csv:"linkedin_url"`
}

// ContactResult represents detailed contact information
// @Description Detailed contact information returned by API operations
type ContactResult struct {
	// Inherits standard response fields
	rest.BaseResponse

	// Contact's unique identifier
	ContactId string `json:"contactId,omitempty" example:"contact-123"`

	// Contact's email address
	Email string `json:"email,omitempty" example:"john@example.com"`

	// Contact's LinkedIn profile URL
	LinkedInURL string `json:"linkedinUrl,omitempty" example:"https://linkedin.com/in/john-doe"`
}

// ContactsResponse represents a response containing multiple contacts
// @Description Response structure for operations returning multiple contacts
type ContactsResponse struct {
	// Status of the operation
	Status string `json:"status" example:"success"`

	// Additional information about the operation
	Message string `json:"message,omitempty" example:"Contacts processed successfully"`

	// List of contacts
	Contacts []ContactResult `json:"contacts,omitempty"`
}

// CreateOrganizationRequest represents the request body for creating a new organization
// @Description Request to create an organization
type CreateOrganizationRequest struct {
	// Organization's name
	Name string `json:"name" example:"CustomerOS"`

	// Custom ID provided by the user
	CustomId string `json:"customId" example:"12345"`

	// Organization's website URL
	Website string `json:"website" example:"https://customeros.ai"`

	// Organization's LinkedIn profile URL
	LinkedinUrl string `json:"linkedinUrl" example:"https://linkedin.com/company/openline"`

	// Lead source of the organization
	LeadSource string `json:"leadSource" example:"Web Search"`

	// Relationship status of the organization
	Relationship string `json:"relationship" example:"customer"`

	// Indicates if the organization is an ICP (Ideal Customer Profile) fit
	IcpFit bool `json:"icpFit" example:"true"`
}

// SetPrimaryExternalSystemIdRequest represents the request body for setting a primary external system ID
// @Description Request to set or replace the primary external system ID for an organization
type SetPrimaryExternalSystemIdRequest struct {
	// The ID of the external system to be set as primary
	ExternalId string `json:"externalId" example:"stripe-1234"`
}

// OrganizationResult represents detailed organization information
// @Description Detailed organization information returned by API operations
type OrganizationResult struct {
	// Status indicates the result of the operation
	Status string `json:"status" example:"success"`

	// Message provides additional information
	Message string `json:"message,omitempty" example:"Organization retrieved successfully"`

	// Organization's unique identifier
	ID string `json:"id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Custom ID provided by the user
	CustomId string `json:"customId,omitempty" example:"12345"`

	// CustomerOS unique identifier
	CosId string `json:"cosId,omitempty" example:"C-A12-B45"`

	// Organization name
	Name string `json:"name,omitempty" example:"CustomerOS"`

	// Organization's website URL
	Website string `json:"website,omitempty" example:"https://customeros.ai"`

	// Lead source of the organization
	LeadSource string `json:"leadSource,omitempty" example:"Web Search"`

	// Relationship status with the organization
	Relationship string `json:"relationship,omitempty" example:"customer"`

	// Current stage in the organization lifecycle
	Stage string `json:"stage,omitempty" example:"lead"`

	// ICP fit indicator
	IcpFit bool `json:"icpFit,omitempty" example:"true"`

	// Associated domains
	Domains []string `json:"domains,omitempty" example:"customeros.com,customeros.ai"`

	// External system links
	ExternalLinks []ExternalLink `json:"externalLinks,omitempty"`
}

// OrganizationsResponse represents a response containing multiple organizations
// @Description Response structure for operations returning multiple organizations
type OrganizationsResponse struct {
	// Status of the operation
	Status string `json:"status" example:"success"`

	// Additional information about the operation
	Message string `json:"message,omitempty" example:"Organizations retrieved successfully"`

	// List of organizations
	Organizations []OrganizationResult `json:"organizations,omitempty"`
}

// ExternalSystemResult represents the result of external system operations
// @Description Response structure for external system operations
type ExternalSystemResult struct {
	// Status of the operation
	Status string `json:"status" example:"success"`

	// Additional information
	Message string `json:"message,omitempty" example:"External system ID set successfully"`

	// Associated organization ID
	OrganizationId string `json:"organizationId,omitempty" example:"org-789"`

	// Name of the external system
	ExternalSystem string `json:"externalSystem,omitempty" example:"stripe"`

	// External system identifier
	ExternalId string `json:"externalId,omitempty" example:"stripe-1234"`

	// Indicates if this is the primary link
	Primary bool `json:"primary,omitempty" example:"true"`
}

// ExternalLink represents a link to an external system
// @Description External system link information
type ExternalLink struct {
	// External system name
	Name string `json:"name" example:"stripe"`

	// External system identifier
	Id string `json:"id" example:"cos-12345"`

	// Indicates if this is the primary link
	Primary bool `json:"primary" example:"true"`
}

// SocialLink represents a social media link
// @Description Social media link information
type SocialLink struct {
	// Social media profile URL
	Url string `json:"url" example:"https://linkedin.com/company/openline"`

	// Number of followers
	FollowerCount int `json:"followerCount" example:"1000"`
}
