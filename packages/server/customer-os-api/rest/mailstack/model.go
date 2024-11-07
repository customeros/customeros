package restmailstack

// RegisterNewDomainRequest defines the request body for registering a new domain
// @Description Request body for domain registration
type RegisterNewDomainRequest struct {
	// Domain is the domain name to be registered
	// Required: true
	Domain string `json:"domain" example:"example.com"`

	// Destination website for permanent redirect
	// Required: true
	Website string `json:"website" example:"www.example.com"`
}

// ConfigureDomainRequest defines the request body for configuring domain
// @Description Request body for domain configuration
type ConfigureDomainRequest struct {
	// Domain is the domain name to be configured
	// Required: true
	Domain string `json:"domain" example:"example.com"`

	// Destination website for permanent redirect
	// Required: true
	Website string `json:"website" example:"www.example.com"`
}

// DomainsResponse defines the response structure for multiple domains in the response
// @Description Response body for all domain details
type DomainsResponse struct {
	// Status indicates the result of the action
	Status string `json:"status,omitempty" example:"success"`

	// Message provides additional information about the action
	Message string `json:"message,omitempty" example:"Domains retrieved successfully"`

	Domains []DomainResponse `json:"domains"`
}

// DomainResponse defines the structure of a domain in the response
// @Description Domain object in the response
type DomainResponse struct {
	// Status indicates the result of the action
	Status string `json:"status,omitempty" example:"success"`

	// Message provides additional information about the action
	Message string `json:"message,omitempty" example:"Domain retrieved successfully"`

	// Domain is the domain name that was registered
	Domain string `json:"domain" example:"example.com"`

	// CreatedDate is the date the domain was registered
	CreatedDate string `json:"createdDate" example:"09/14/2024"`

	// ExpiredDate is the date when the domain registration will expire
	ExpiredDate string `json:"expiredDate" example:"09/14/2025"`

	// Nameservers lists the nameservers associated with the domain
	Nameservers []string `json:"nameservers" example:"['ns1.example.com', 'ns2.example.com']"`
}

// MailboxRequest represents the request body to add and configure a new mailbox
// @Description Request body for adding and configuring a new mailbox
type MailboxRequest struct {
	// Username for the mailbox (e.g., "john.doe")
	// Required: true
	Username string `json:"username" example:"john.doe"`

	// Password for the mailbox (e.g., "SecurePassword123!")
	// Required: false
	Password string `json:"password" example:"SecurePassword123!"`

	// Specifies if email forwarding is enabled
	ForwardingEnabled bool `json:"forwardingEnabled" example:"true"`

	// Email address to forward to (if forwarding is enabled)
	ForwardingTo []string `json:"forwardingTo" example:"['user1@example.com', 'user2@example.com']"`

	// Specifies if webmail access is enabled
	WebmailEnabled bool `json:"webmailEnabled" example:"true"`

	// LinkedUser is the email address of the user to whom new mailbox should be linked. If not provided or not found, mailbox will not be associated with any user
	// Required: false
	LinkedUser string `json:"linkedUser" example:"john.doe@mycompany.com"`
}

// MailboxResponse defines the structure of a mailbox in the response
// @Description Mailbox object in the response
type MailboxResponse struct {
	// Status indicates the result of the action
	Status string `json:"status,omitempty" example:"success"`

	// Message provides additional information about the action
	Message string `json:"message,omitempty" example:"Mailbox setup successful"`

	// Email is the email address for the mailbox
	// Required: true
	Email string `json:"email" example:"user@example.com"`

	// Password is the password for the mailbox
	// Required: false
	Password string `json:"password,omitempty" example:"SecurePassword123!"`

	// ForwardingEnabled indicates if email forwarding is enabled
	ForwardingEnabled bool `json:"forwardingEnabled" example:"true"`

	// ForwardingTo is the email address the mailbox forwards to
	ForwardingTo []string `json:"forwardingTo" example:"['user1@example.com', 'user2@example.com']"`

	// WebmailEnabled indicates if webmail access is enabled
	WebmailEnabled bool `json:"webmailEnabled" example:"true"`
}

// MailboxesResponse defines the response structure for multiple mailboxes in the response
// @Description Response body for all mailbox details
type MailboxesResponse struct {
	// Status indicates the result of the action
	Status string `json:"status,omitempty" example:"success"`

	// Message provides additional information about the action
	Message string `json:"message,omitempty" example:"Mailboxes retrieved successfully"`

	Mailboxes []MailboxResponse `json:"mailboxes"`
}
