package models

import (
	"fmt"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type EmailValidation struct {
	ValidationError string `json:"validationError"`
	IsReachable     string `json:"isReachable"`
	AcceptsMail     bool   `json:"acceptsMail"`
	CanConnectSmtp  bool   `json:"canConnectSmtp"`
	HasFullInbox    bool   `json:"hasFullInbox"`
	IsCatchAll      bool   `json:"isCatchAll"`
	IsDeliverable   bool   `json:"isDeliverable"`
	IsDisabled      bool   `json:"isDisabled"`
	Domain          string `json:"domain"`
	IsValidSyntax   bool   `json:"isValidSyntax"`
	Username        string `json:"username"`
}

type Email struct {
	ID              string          `json:"id"`
	RawEmail        string          `json:"rawEmail"`
	Email           string          `json:"email"`
	Source          cmnmod.Source   `json:"source"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	EmailValidation EmailValidation `json:"emailValidation"`
}

func (p *Email) String() string {
	return fmt.Sprintf("Email{ID: %s, RawEmail: %s, Email: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s}", p.ID, p.RawEmail, p.Email, p.Source, p.CreatedAt, p.UpdatedAt)
}

func NewEmail() *Email {
	return &Email{}
}
