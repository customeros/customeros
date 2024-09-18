package restmailstack

// RegisterNewDomainRequest defines the request body for registering a new domain
// @Description Request body for domain registration
type RegisterNewDomainRequest struct {
	// Domain is the domain name to be registered
	// Required: true
	// Example: example.com
	Domain string `json:"domain" example:"example.com"`

	// Destination website for permanent redirect
	// Required: true
	// Example: www.example.com
	Website string `json:"website" example:"www.example.com"`
}

// ConfigureDomainRequest defines the request body for configuring domain
// @Description Request body for domain configuration
type ConfigureDomainRequest struct {
	// Domain is the domain name to be configured
	// Required: true
	// Example: example.com
	Domain string `json:"domain" example:"example.com"`

	// Destination website for permanent redirect
	// Required: true
	// Example: www.example.com
	Website string `json:"website" example:"www.example.com"`
}

// DomainsResponse defines the response structure for multiple domains in the response
// @Description Response body for all domain details
type DomainsResponse struct {
	// Status indicates the result of the action
	// Example: success
	Status string `json:"status,omitempty" example:"success"`

	// Message provides additional information about the action
	// Example: Domain retrieved successfully
	Message string `json:"message,omitempty" example:"Domains retrieved successfully"`

	Domains []DomainResponse `json:"domains"`
}

// DomainResponse defines the structure of a domain in the response
// @Description Domain object in the response
type DomainResponse struct {
	// Status indicates the result of the action
	// Example: success
	Status string `json:"status,omitempty" example:"success"`

	// Message provides additional information about the action
	// Example: Domain registered successfully
	Message string `json:"message,omitempty" example:"Domain retrieved successfully"`

	// Domain is the domain name that was registered
	// Example: example.com
	Domain string `json:"domain" example:"example.com"`

	// CreatedDate is the date the domain was registered
	// Example: 09/14/2024
	CreatedDate string `json:"createdDate" example:"09/14/2024"`

	// ExpiredDate is the date when the domain registration will expire
	// Example: 09/14/2025
	ExpiredDate string `json:"expiredDate" example:"09/14/2025"`

	// Nameservers lists the nameservers associated with the domain
	// Example: [ns1.example.com, ns2.example.com]
	Nameservers []string `json:"nameservers" example:"['ns1.example.com', 'ns2.example.com']"`
}

// MailboxRequest represents the request body to add and configure a new mailbox
// @Description Request body for adding and configuring a new mailbox
type MailboxRequest struct {
	// Username for the mailbox (e.g., "john.doe")
	// Required: true
	// Example: john.doe
	Username string `json:"username" example:"john.doe" binding:"required"`

	// Password for the mailbox (e.g., "SecurePassword123!")
	// Required: true
	// Example: SecurePassword123!
	Password string `json:"password" example:"SecurePassword123!" binding:"required"`

	// Specifies if email forwarding is enabled
	// Example: true
	ForwardingEnabled bool `json:"forwardingEnabled" example:"true"`

	// Email address to forward to (if forwarding is enabled)
	// Example: johndoe.forward@example.com
	ForwardingTo string `json:"forwardingTo" example:"johndoe.forward@example.com"`

	// Specifies if webmail access is enabled
	// Example: true
	WebmailEnabled bool `json:"webmailEnabled" example:"true"`
}

// MailboxResponse defines the structure of a mailbox in the response
// @Description Mailbox object in the response
type MailboxResponse struct {
	// Status indicates the result of the action
	// Example: success
	Status string `json:"status" example:"success"`

	// Message provides additional information about the action
	// Example: Mailbox setup successful
	Message string `json:"message" example:"Mailbox setup successful"`

	// Email is the email address for the mailbox
	// Required: true
	// Example: user@example.com
	Email string `json:"email" example:"user@example.com"`

	// ForwardingEnabled indicates if email forwarding is enabled
	// Example: true
	ForwardingEnabled bool `json:"forwardingEnabled" example:"true"`

	// ForwardingTo is the email address the mailbox forwards to
	// Example: user@forward.com
	ForwardingTo string `json:"forwardingTo" example:"user@forward.com"`

	// WebmailEnabled indicates if webmail access is enabled
	// Example: true
	WebmailEnabled bool `json:"webmailEnabled" example:"true"`
}
