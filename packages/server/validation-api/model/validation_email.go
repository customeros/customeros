package model

import postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

type EmailDeliverableStatus string

const (
	EmailDeliverableStatusDeliverable   EmailDeliverableStatus = "true"
	EmailDeliverableStatusUndeliverable EmailDeliverableStatus = "false"
	EmailDeliverableStatusUnknown       EmailDeliverableStatus = "unknown"
)

type ValidateEmailRequest struct {
	Email string `json:"email"`
}

type ValidateEmailRequestWithOptions struct {
	Email   string                      `json:"email"`
	Options ValidateEmailRequestOptions `json:"options"`
}

type ValidateEmailRequestOptions struct {
	VerifyCatchAll      bool `json:"verifyCatchAll"`
	ExtendedWaitingTime bool `json:"extendedWaitingTime"`
}

type ValidateEmailResponse struct {
	Status          string                       `json:"status"`
	Message         string                       `json:"message,omitempty"`
	InternalMessage string                       `json:"internalMessage,omitempty"`
	Data            *ValidateEmailMailSherpaData `json:"data,omitempty"`
}

type ValidateEmailMailSherpaData struct {
	Email  string `json:"email"`
	Syntax struct {
		IsValid    bool   `json:"isValid"`
		User       string `json:"user"`
		Domain     string `json:"domain"`
		CleanEmail string `json:"cleanEmail"`
	} `json:"syntax"`
	DomainData struct {
		IsFirewalled          bool   `json:"isFirewalled"`
		Provider              string `json:"provider"`
		SecureGatewayProvider string `json:"secureGatewayProvider"`
		IsCatchAll            bool   `json:"isCatchAll"`
		CanConnectSMTP        bool   `json:"canConnectSMTP"`
		HasMXRecord           bool   `json:"hasMXRecord"`
		HasSPFRecord          bool   `json:"hasSPFRecord"`
		TLSRequired           bool   `json:"tlsRequired"`
		ResponseCode          string `json:"responseCode"`
		ErrorCode             string `json:"errorCode"`
		Description           string `json:"description"`
		IsPrimaryDomain       bool   `json:"isPrimaryDomain"`
		PrimaryDomain         string `json:"primaryDomain"`
	} `json:"domainData"`
	EmailData struct {
		SkippedValidation bool   `json:"skippedValidation"` // if true, email validation was skipped
		Deliverable       string `json:"deliverable"`
		IsMailboxFull     bool   `json:"isMailboxFull"`
		IsRoleAccount     bool   `json:"isRoleAccount"`
		IsSystemGenerated bool   `json:"isSystemGenerated"`
		IsFreeAccount     bool   `json:"isFreeAccount"`
		SmtpSuccess       bool   `json:"smtpSuccess"`
		ResponseCode      string `json:"responseCode"`
		ErrorCode         string `json:"errorCode"`
		Description       string `json:"description"`
		RetryValidation   bool   `json:"retryValidation"` // if true, email validation should be retried
		TLSRequired       bool   `json:"tlsRequired"`
		AlternateEmail    string `json:"alternateEmail"`
	} `json:"emailData"`
}

type ValidateEmailWithScrubbyResponse struct {
	Status          string `json:"status"`
	Message         string `json:"message,omitempty"`
	InternalMessage string `json:"internalMessage,omitempty"`
	EmailIsValid    bool   `json:"emailIsValid"`
	EmailIsInvalid  bool   `json:"emailIsInvalid"`
	EmailIsUnknown  bool   `json:"emailIsUnknown"`
	EmailIsPending  bool   `json:"emailIsPending"`
}

type ValidateEmailWithTrueinboxResponse struct {
	Status  string                                `json:"status"`
	Message string                                `json:"message,omitempty"`
	Data    *postgresentity.TrueInboxResponseBody `json:"data,omitempty"`
}
