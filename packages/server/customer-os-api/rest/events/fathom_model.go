// @openapi 3.0.0
package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
	"time"
)

// FathomAISummaryZapier represents the webhook payload from Fathom via Zapier
type RawFathomAISummaryZapier struct {
	// AI Summary fields
	AISummaryHTMLFormatted             string `json:"ai_summary_html_formatted"`
	AISummaryPlaintextFormatted        string `json:"ai_summary_plaintext_formatted"`
	AISummarySectionMarkdownFormatted  string `json:"ai_summary_section_markdown_formatted"`
	AISummarySectionPlaintextFormatted string `json:"ai_summary_section_plaintext_formatted"`
	AISummarySectionsDepth             string `json:"ai_summary_sections_depth"`
	AISummarySectionsHTMLFormatted     string `json:"ai_summary_sections_html_formatted"`
	AISummarySectionsSlackFormatted    string `json:"ai_summary_sections_slack_formatted"`
	AISummarySlackFormatted            string `json:"ai_summary_slack_formatted"`
	AISummaryTemplateName              string `json:"ai_summary_template_name"`
	AISummaryTitles                    string `json:"ai_summary_titles"`

	// User fields
	FathomUserEmail string `json:"fathom_user_email"`
	FathomUserName  string `json:"fathom_user_name"`

	// Meeting metadata
	ID                        string    `json:"id"`
	HasExternal               string    `json:"has_external"`
	InviteesExternal          string    `json:"invitees_external"`
	MeetingExternalDomains    string    `json:"meeting_external_domains"`
	MeetingInviteeEmails      string    `json:"meeting_invitee_emails"`
	MeetingJoinURL            string    `json:"meeting_join_url"`
	MeetingScheduledEndTime   time.Time `json:"meeting_scheduled_end_time"`
	MeetingScheduledStartTime time.Time `json:"meeting_scheduled_start_time"`
	MeetingTitle              string    `json:"meeting_title"`

	// Recording details
	RecordingDuration string `json:"recording_duration"`
	RecordingShareURL string `json:"recording_share_url"`
	RecordingURL      string `json:"recording_url"`
}

// Helper method to check if meeting has external participants
func (f *RawFathomAISummaryZapier) HasExternalParticipants() bool {
	return strings.ToLower(f.HasExternal) == "true"
}

// Helper method to get meeting duration in minutes
func (f *RawFathomAISummaryZapier) GetDurationMinutes() float64 {
	if f.RecordingDuration == "" {
		return 0
	}
	return utils.IfNotNilFloat64(utils.ParseStringToFloat(f.RecordingDuration))
}

// Helper method to get external domains
func (f *RawFathomAISummaryZapier) GetExternalDomains() []string {
	if f.MeetingExternalDomains == "" {
		return nil
	}
	return strings.Split(f.MeetingExternalDomains, ",")
}

// Helper method to get invitee emails
func (f *RawFathomAISummaryZapier) GetInviteeEmails() []string {
	if f.MeetingInviteeEmails == "" {
		return nil
	}
	return strings.Split(f.MeetingInviteeEmails, ",")
}
