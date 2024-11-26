package events

import (
	"time"
)

// FathomWebhookPayload represents the root webhook payload
type FathomZapierPayload struct {
	AISummary  AISummary  `json:"ai_summary"`
	FathomUser FathomUser `json:"fathom_user"`
	ID         string     `json:"id"`
	Meeting    Meeting    `json:"meeting"`
	Recording  Recording  `json:"recording"`
}

// AISummary represents the AI summary section
type AISummary struct {
	HTMLFormatted      string `json:"html_formatted"`
	MarkdownFormatted  string `json:"markdown_formatted"`
	PlaintextFormatted string `json:"plaintext_formatted"`
	Sections           string `json:"sections"`
	SlackFormatted     string `json:"slack_formatted"`
	TemplateName       string `json:"template_name"`
}

// FathomUser represents the user information
type FathomUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Meeting represents the meeting information
type Meeting struct {
	ExternalDomainsStr       string           `json:"external_domains"`
	ExternalDomains          []ExternalDomain `json:"external_domains_array"`
	HasExternalInvitees      string           `json:"has_external_invitees"`
	InviteesStr              string           `json:"invitees"`
	Invitees                 []Invitee        `json:"invitees_array"`
	JoinURL                  string           `json:"join_url"`
	ScheduledDurationMinutes string           `json:"scheduled_duration_in_minutes"`
	ScheduledEndTime         time.Time        `json:"scheduled_end_time"`
	ScheduledStartTime       time.Time        `json:"scheduled_start_time"`
	Title                    string           `json:"title"`
}

type ExternalDomain struct {
	Domain string `json:"domain_name"`
}

// Invitee represents meeting participant information
type Invitee struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	IsExternal bool   `json:"is_external"`
}

// Recording represents the recording information
type Recording struct {
	DurationInMinutes string `json:"duration_in_minutes"`
	ShareURL          string `json:"share_url"`
	URL               string `json:"url"`
}

func (f *FathomZapierPayload) ExternalDomains() []string {
	var domains []string

	for _, domain := range f.Meeting.ExternalDomains {
		domains = append(domains, domain.Domain)
	}
	return domains
}
