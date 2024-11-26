// @openapi 3.0.0
package restmailstack

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"

// RegisterNewDomainRequest represents the domain registration request
// @Description Request payload for registering a new domain for mail services
type RegisterNewDomainRequest struct {
	// Domain name to register
	// required: true
	// pattern: ^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$
	// example: example.com
	Domain string `json:"domain"`

	// Website URL for domain configuration
	// required: true
	// format: uri
	// example: https://www.example.com
	Website string `json:"website"`
}

// ConfigureDomainRequest represents the domain configuration request
// @Description Request payload for configuring domain DNS and mail services
type ConfigureDomainRequest struct {
	// Domain name to configure
	// required: true
	// pattern: ^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$
	// example: example.com
	Domain string `json:"domain"`

	// Website URL for domain configuration
	// required: true
	// format: uri
	// example: https://www.example.com
	Website string `json:"website"`
}

// DomainResponse represents a single domain response
// @Description Response containing domain details and status
type DomainResponse struct {
	// Inherits standard response fields
	rest.BaseResponse
	// Domain information
	// required: true
	Domain DomainRecord `json:"domain"`
}

// DomainsResponse represents multiple domains response
// @Description Response containing list of domains and status
type DomainsResponse struct {
	// Inherits standard response fields
	rest.BaseResponse
	// List of domains
	// required: true
	Domains []DomainRecord `json:"domains"`
}

// DomainRecord represents detailed domain information
// @Description Comprehensive domain record information
type DomainRecord struct {
	// Registered domain name
	// required: true
	// example: example.com
	Domain string `json:"domain"`

	// Domain registration date
	// required: true
	// format: date
	// example: 2024-09-14
	CreatedDate string `json:"createdDate"`

	// Domain expiration date
	// required: true
	// format: date
	// example: 2025-09-14
	ExpiredDate string `json:"expiredDate"`

	// List of assigned nameservers
	// required: true
	// minItems: 2
	// example: ["ns1.example.com","ns2.example.com"]
	Nameservers []string `json:"nameservers"`
}

// MailboxRequest represents mailbox creation request
// @Description Request payload for creating and configuring a new mailbox
type MailboxRequest struct {
	// Username for the mailbox
	// required: true
	// pattern: ^[a-zA-Z0-9._%+-]+$
	// minLength: 3
	// maxLength: 64
	// example: john.doe
	Username string `json:"username"`

	// Password for mailbox access
	// required: false
	// minLength: 8
	// maxLength: 64
	// example: SecurePassword123!
	Password string `json:"password"`

	// List of email addresses to forward to
	// required: false
	// maxItems: 10
	// example: ["user1@example.com","user2@example.com"]
	ForwardingTo []string `json:"forwardingTo"`

	// Enable webmail access
	// required: false
	// default: false
	WebmailEnabled bool `json:"webmailEnabled"`

	// Associated user's email address
	// required: false
	// format: email
	// example: john.doe@mycompany.com
	LinkedUser string `json:"linkedUser"`
}

// MailboxResponse represents single mailbox response
// @Description Response containing mailbox details and status
type MailboxResponse struct {
	// Inherits standard response fields
	rest.BaseResponse
	// Mailbox information
	// required: false
	Mailbox MailboxRecord `json:"mailbox,omitempty"`
}

// MailboxesResponse represents multiple mailboxes response
// @Description Response containing list of mailboxes and status
type MailboxesResponse struct {
	// Inherits standard response fields
	rest.BaseResponse
	// List of mailboxes
	// required: false
	Mailboxes []MailboxRecord `json:"mailboxes,omitempty"`
}

// MailboxRecord represents detailed mailbox information
// @Description Comprehensive mailbox configuration and status
type MailboxRecord struct {
	// Email address for the mailbox
	// required: true
	// format: email
	// example: user@example.com
	Email string `json:"email"`

	// Mailbox password (only included in specific responses)
	// required: false
	// minLength: 8
	// maxLength: 64
	Password string `json:"password,omitempty"`

	// Email forwarding status
	// required: true
	// default: false
	ForwardingEnabled bool `json:"forwardingEnabled"`

	// List of forwarding email addresses
	// required: false
	// maxItems: 10
	ForwardingTo []string `json:"forwardingTo"`

	// Webmail access status
	// required: true
	// default: false
	WebmailEnabled bool `json:"webmailEnabled"`
}
