package data_fields

import "time"

type EmailValidationFields struct {
	EmailAddress      string    `json:"emailAddress"`
	Domain            string    `json:"domain"`
	IsCatchAll        bool      `json:"isCatchAll"`
	Deliverable       string    `json:"deliverable"`
	IsValidSyntax     bool      `json:"isValidSyntax"`
	Username          string    `json:"username"`
	ValidatedAt       time.Time `json:"validatedAt"`
	IsRoleAccount     bool      `json:"isRoleAccount"`
	IsSystemGenerated bool      `json:"isSystemGenerated"`
	IsRisky           bool      `json:"isRisky"`
	IsFirewalled      bool      `json:"isFirewalled"`
	Provider          string    `json:"provider"`
	Firewall          string    `json:"firewall"`
	IsMailboxFull     bool      `json:"isMailboxFull"`
	IsFreeAccount     bool      `json:"isFreeAccount"`
	SmtpSuccess       bool      `json:"smtpSuccess"`
	ResponseCode      string    `json:"responseCode"`
	ErrorCode         string    `json:"errorCode"`
	Description       string    `json:"description"`
	IsPrimaryDomain   bool      `json:"isPrimaryDomain"`
	PrimaryDomain     string    `json:"primaryDomain"`
	AlternateEmail    string    `json:"alternateEmail"`
	RetryValidation   bool      `json:"retryValidation"`
}
