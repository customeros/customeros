package flows

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/html"
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

type FathomMeetingSummaryCreatedEvent struct {
	ParticipantEmails *[]string `json:"participantEmails"`
	Content           *string   `json:"content,omitempty"`
}

func (f *FathomZapierPayload) ExternalDomains() []string {
	var domains []string

	for _, domain := range f.Meeting.ExternalDomains {
		domains = append(domains, domain.Domain)
	}
	return domains
}

func (f *FathomZapierPayload) toCleanPayload() error {
	if f.Meeting.ExternalDomainsStr != "" {
		externalDomainsJson := utils.ReplaceSingleQuotesWithDoubleQuotes(f.Meeting.ExternalDomainsStr)
		var externalDomains []ExternalDomain
		err := json.Unmarshal([]byte(externalDomainsJson), &externalDomains)
		if err != nil {
			return err
		}
		f.Meeting.ExternalDomains = externalDomains
	}

	if f.Meeting.InviteesStr != "" {
		inviteesJson := utils.ReplaceSingleQuotesWithDoubleQuotes(f.Meeting.InviteesStr)
		inviteesJson = strings.Replace(inviteesJson, ": True", ": true", -1)
		inviteesJson = strings.Replace(inviteesJson, ": False", ": false", -1)
		var invitees []Invitee
		err := json.Unmarshal([]byte(inviteesJson), &invitees)
		if err != nil {
			return err
		}
		f.Meeting.Invitees = invitees
	}

	return nil
}

func (f *FathomZapierPayload) toMarkdownContent() (string, error) {
	// Convert HTML to clean markdown
	cleanMarkdown, err := f.AISummary.toCleanMarkdown()
	if err != nil {
		return "", fmt.Errorf("error converting HTML to markdown: %w", err)
	}

	// Build the additional sections
	var builder strings.Builder
	builder.WriteString(cleanMarkdown)

	// Add participants section
	builder.WriteString("\n\n### Meeting Participants\n")
	for _, participant := range f.Meeting.Invitees {
		builder.WriteString(fmt.Sprintf("- %s (%s)\n", participant.Name, participant.Email))
	}

	// Add duration
	builder.WriteString("\n### Meeting Duration\n")
	mins, err := strconv.ParseFloat(f.Recording.DurationInMinutes, 64)
	if err == nil {
		builder.WriteString(fmt.Sprintf("%d minutes\n\n", int(mins)))
	}

	// Add recording link
	builder.WriteString(fmt.Sprintf("[View Recording](%s)\n", f.Recording.ShareURL))

	return builder.String(), nil
}

func (a *AISummary) toCleanMarkdown() (string, error) {
	htmlContent := a.HTMLFormatted
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	var process func(*html.Node)

	process = func(n *html.Node) {
		switch n.Type {
		case html.DocumentNode:
			// Start processing from the root node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				process(c)
			}
		case html.TextNode:
			builder.WriteString(n.Data)
		case html.ElementNode:
			switch n.Data {
			case "h1":
				builder.WriteString("\n# ")
			case "h2":
				builder.WriteString("\n## ")
			case "h3":
				builder.WriteString("\n### ")
			case "p":
				builder.WriteString("\n\n")
			case "ul":
				builder.WriteString("\n")
			case "li":
				builder.WriteString("\n- ")
			case "br":
				builder.WriteString("\n")
			}

			// Process child nodes
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				process(c)
			}

			// Close Markdown elements where necessary
			switch n.Data {
			case "p", "h1", "h2", "h3", "ul", "li":
				builder.WriteString("\n")
			}
		}
	}

	process(doc)

	// Clean up the output
	output := builder.String()

	// Remove multiple newlines
	output = regexp.MustCompile(`\n{3,}`).ReplaceAllString(output, "\n\n")
	// Remove leading/trailing whitespace
	output = strings.TrimSpace(output)
	// Ensure consistent newlines between sections
	output = regexp.MustCompile(`\n## `).ReplaceAllString(output, "\n\n## ")
	output = regexp.MustCompile(`\n### `).ReplaceAllString(output, "\n\n### ")

	return output, nil
}
