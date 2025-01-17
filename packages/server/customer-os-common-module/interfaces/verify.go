package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	extract "github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
)

type VerifyService interface {
	Threats(ctx context.Context, ipAddress string) (*IpThreats, error)
	IdentifyCompanyDomain(ctx context.Context, ipAddress string) (*string, error)
	ValidatePhoneNumber(ctx context.Context, countryCodeA2 string, phoneNumber string) (*string, *string, error)
	ValidateEmail(ctx context.Context, email string) (*ValidateEmailMailSherpaData, error)
	ValidateEmailWithMailSherpa(ctx context.Context, email string) (*ValidateEmailMailSherpaData, error)
	ValidateEmailScrubby(ctx context.Context, email string) (string, error)
	ValidateEmailWithTrueinbox(ctx context.Context, email string) (*entity.TrueInboxResponseBody, error)
	ValidateEmailWithEnrow(ctx context.Context, email string, extendedWaitingTimeForResponse bool) (string, error)

	LookupIp(ctx context.Context, ip string) (*entity.IPDataResponseBody, error)

	ValidateUsAddress(address string) (*extract.Lookup, error)
	ValidateInternationalAddress(address, country string) (*international_street.Lookup, error)
}

type IpThreats struct {
	IsThreat      bool
	IsAnonymous   bool
	IsBogon       bool
	IsDatacenter  bool
	IsICloudRelay bool
	IsKnownAbuser bool
	IsProxy       bool
	IsTor         bool
	IsVpn         bool
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
