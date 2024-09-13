package restmailstack

// RegisterNewDomainRequest defines the request body for registering a new domain
// @Description Request body for domain registration
type RegisterNewDomainRequest struct {
	// Domain is the domain name to be registered
	// Required: true
	// Example: example.com
	Domain string `json:"domain" example:"example.com"`
}

// RegisterNewDomainResponse defines the response structure for a successful domain registration
// @Description Response body for a successful domain registration
// @example 201 {object} RegisterNewDomainResponse
type RegisterNewDomainResponse struct {
	// Status indicates the result of the domain registration
	// Example: success
	Status string `json:"status" example:"success"`

	// Message provides additional information about the registration
	// Example: Domain registered successfully
	Message string `json:"message" example:"Domain registered successfully"`

	// Domain is the domain name that was registered
	// Example: example.com
	Domain string `json:"domain" example:"example.com"`

	// CreatedAt is the date and time the domain was registered
	// Example: 2021-09-01T12:00:00Z
	CreatedAt string `json:"createdAt" example:"2021-09-01T12:00:00Z"`

	// ExpiresAt is the date and time when the domain registration will expire
	// Example: 2022-09-01T12:00:00Z
	ExpiresAt string `json:"expiresAt" example:"2022-09-01T12:00:00Z"`

	// LastUpdatedAt is the date and time the domain registration was last updated
	// Example: 2021-09-01T12:00:00Z
	LastUpdatedAt string `json:"lastUpdatedAt" example:"2021-09-01T12:00:00Z"`

	// AutoRenew indicates whether the domain will be automatically renewed upon expiration
	// Example: true
	AutoRenew bool `json:"autoRenew" example:"true"`

	// Nameservers lists the nameservers associated with the domain
	// Example: [ns1.example.com, ns2.example.com]
	Nameservers []string `json:"nameservers" example:"[\"ns1.example.com\", \"ns2.example.com\"]"`

	// DnsRecords provides a list of DNS records associated with the domain
	DnsRecords []DnsRecordResponse `json:"dnsRecords"`
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
