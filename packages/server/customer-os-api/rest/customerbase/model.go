package customerbase

// CreateOrganizationRequest represents the request body for creating a new organization
// @Description Request to create an organization
type CreateOrganizationRequest struct {
	// Organization's name.
	Name string `json:"name" example:"CustomerOS"`

	// Custom ID provided by the user.
	CustomId string `json:"customId" example:"12345"`

	// Organization's website URL.
	Website string `json:"website" example:"https://customeros.ai"`

	// Organization's LinkedIn profile URL.
	LinkedinUrl string `json:"linkedinUrl" example:"https://linkedin.com/company/openline"`

	// Lead source of the organization.
	LeadSource string `json:"leadSource" example:"Web Search"`

	// Relationship status of the organization.
	Relationship string `json:"relationship" example:"customer"`

	// Indicates if the organization is an ICP (Ideal Customer Profile) fit.
	IcpFit bool `json:"icpFit" example:"true"`
}

// CreateOrganizationResponse represents the response returned after creating an organization
// @Description The response structure after an organization is successfully created.
// @example 201 {object} CreateOrganizationResponse
type CreateOrganizationResponse struct {
	// Status indicates the status of the creation process (e.g., "success" or "partial_success").
	Status string `json:"status" example:"success"`

	// Message provides additional information regarding the creation process.
	Message string `json:"message,omitempty" example:"Organization created successfully"`

	// ID is the unique identifier of the created organization.
	ID string `json:"id" example:"1234567890"`

	// PartialSuccess indicates whether the creation process encountered partial success (e.g., when some fields failed to process).
	PartialSuccess bool `json:"partialSuccess,omitempty" example:"false"`
}

// Organization represents the structure of an organization
// @Description Organization details
type OrganizationResponse struct {
	// Status indicates the result of the action.
	Status string `json:"status,omitempty" example:"success"`

	// Message provides additional information about the action.
	Message string `json:"message,omitempty" example:"Organization retrieved successfully"`

	// ID is the unique identifier for the organization, uuid format.
	ID string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Custom ID provided by the user.
	CustomId string `json:"customId" example:"12345"`

	// COS ID is the unique identifier for the organization in the Customer OS system.
	CosId string `json:"cosId" example:"C-A12-B45"`

	// Organization name.
	Name string `json:"name" example:"CustomerOS"`

	// Organization's website URL.
	Website string `json:"website" example:"https://customeros.ai"`

	// Domains associated with the organization.
	Domains []string `json:"domains" example:"customeros.com, customeros.ai"`

	// Lead source of the organization.
	LeadSource string `json:"leadSource" example:"Web Search"`

	// Relationship status of the organization.
	Relationship string `json:"relationship" example:"customer"`

	// Stage of the organization.
	Stage string `json:"stage" example:"lead"`

	// Indicates if the organization is an ICP (Ideal Customer Profile) fit.
	IcpFit bool `json:"icpFit" example:"true"`

	// External links associated with the organization.
	ExternalLinks []ExternalLink `json:"externalLinks"`
}

type ExternalLink struct {
	// External app identified
	Name string `json:"name" example:"stripe"`

	// External app id
	Id string `json:"id" example:"cos-12345"`
}

type SocialLink struct {
	// Social media url
	Url string `json:"url" example:"https://linkedin.com/company/openline"`

	// Follower count on the social media platform
	FollowerCount int `json:"followerCount" example:"1000"`
}

// SetPrimaryExternalSystemIdRequest represents the request body for setting a primary external system ID
// @Description Request to set or replace the primary external system ID for an organization
type SetPrimaryExternalSystemIdRequest struct {
	// The ID of the external system to be set as primary.
	ExternalId string `json:"externalId" example:"stripe-1234"`
}

// SetPrimaryExternalSystemIdResponse represents the response body when setting a primary external system ID
// @Description Response after setting the primary external system ID for an organization
type SetPrimaryExternalSystemIdResponse struct {
	// Status indicates the status of the creation process (e.g., "success" or "partial_success").
	Status string `json:"status" example:"success"`

	// Message provides additional information regarding the API.
	Message string `json:"message,omitempty" example:"Primary ID set successfully"`

	// The organization ID for which the primary ID was set.
	OrganizationId string `json:"organizationId" example:"org-789"`

	// The name of the external system.
	ExternalSystem string `json:"externalSystem" example:"stripe"`

	// The primary external ID that was set.
	PrimaryExternalId string `json:"primaryExternalId" example:"stripe-1234"`
}
