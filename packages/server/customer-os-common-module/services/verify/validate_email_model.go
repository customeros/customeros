package verify

import (
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

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
	Status          string                                  `json:"status"`
	Message         string                                  `json:"message,omitempty"`
	InternalMessage string                                  `json:"internalMessage,omitempty"`
	Data            *interfaces.ValidateEmailMailSherpaData `json:"data,omitempty"`
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
	Status  string                                 `json:"status"`
	Message string                                 `json:"message,omitempty"`
	Data    *postgres_entity.TrueInboxResponseBody `json:"data,omitempty"`
}
