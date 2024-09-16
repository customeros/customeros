package restmailstack

// RegisterNewDomainRequest defines the request body for registering a new domain
// @Description Request body for domain registration
type RegisterNewDomainRequest struct {
	// Domain is the domain name to be registered
	// Required: true
	// Example: example.com
	Domain string `json:"domain" example:"example.com"`
}

// ConfigureDomainRequest defines the request body for configuring domain
// @Description Request body for domain configuration
type ConfigureDomainRequest struct {
	// Domain is the domain name to be configured
	// Required: true
	// Example: example.com
	Domain string `json:"domain" example:"example.com"`
}

// MailboxResponse defines the structure of a mailbox in the response
// @Description Mailbox object in the response
type MailboxResponse struct {
	// Email is the email address for the mailbox
	// Example: user@example.com
	Email string `json:"email" example:"user@example.com"`

	// CreatedAt is the date and time the mailbox was created
	// Example: 2021-09-01T12:00:00Z
	CreatedAt string `json:"createdAt" example:"2021-09-01T12:00:00Z"`

	// LastUpdatedAt is the date and time the mailbox was last updated
	// Example: 2021-09-01T12:00:00Z
	LastUpdatedAt string `json:"lastUpdatedAt" example:"2021-09-01T12:00:00Z"`
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

	// DnsRecords provides a list of DNS records associated with the domain
	DnsRecords []DnsRecordResponse `json:"dnsRecords"`
}

// DnsRecordResponse defines the structure of a DNS record in the response
// @Description DNS record object in the response
type DnsRecordResponse struct {
	// Type is the type of the DNS record (e.g., A, MX, TXT)
	// Example: A
	Type string `json:"type" example:"A"`

	// Name is the name of the DNS record
	// Example: example.com
	Name string `json:"name" example:"example.com"`

	// Value is the value of the DNS record
	// Example: 192.0.2.1
	Value string `json:"value" example:"192.0.2.1"`
}
