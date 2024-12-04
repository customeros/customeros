package flows

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

// RecordingRes the top-level response structure
type GrainRecordingData struct {
	RecordingData GrainRecording `json:"data"`
	Type          string         `json:"type"`
	UserID        string         `json:"user_id"`
}

// RecordingData contains the details of the recording
type GrainRecording struct {
	EndDatetime         time.Time          `json:"end_datetime"`
	ICalUID             string             `json:"ical_uid"`
	ID                  string             `json:"id"`
	IntelligenceNotesMD string             `json:"intelligence_notes_md"`
	OwnersStr           string             `json:"owners"`
	Owners              []string           `json:"ownersArr"`
	ParticipantsStr     string             `json:"participants"`
	Participants        []GrainParticipant `json:"participantsArray"`
	PublicThumbnailURL  string             `json:"public_thumbnail_url"`
	PublicURL           string             `json:"public_url"`
	StartDatetime       time.Time          `json:"start_datetime"`
	TagsStr             string             `json:"tags"`
	Tags                []string           `json:"tagsArray"`
	Title               string             `json:"title"`
	URL                 string             `json:"url"`
}

// Participant represents an individual participant in the recording
type GrainParticipant struct {
	Name              string  `json:"name"`
	Scope             string  `json:"scope"`
	Email             *string `json:"email,omitempty"`
	ConfirmedAttendee bool    `json:"confirmed_attendee"`
}

func (g *GrainRecording) participantEmails() []string {
	emails := make([]string, len(g.Participants))

	for v, participant := range g.Participants {
		emails[v] = *participant.Email
	}
	return emails
}

func (g *GrainRecordingData) cleanPayload() error {
	if g.RecordingData.OwnersStr != "" {
		ownersJson := utils.ReplaceSingleQuotesWithDoubleQuotes(g.RecordingData.OwnersStr)
		var owners []string
		err := json.Unmarshal([]byte(ownersJson), &owners)
		if err != nil {
			return err
		}
		g.RecordingData.Owners = owners
	}

	if g.RecordingData.ParticipantsStr != "" {
		participantsJson := utils.ReplaceSingleQuotesWithDoubleQuotes(g.RecordingData.ParticipantsStr)

		participantsJson = strings.Replace(participantsJson, " True", " true", -1)
		participantsJson = strings.Replace(participantsJson, " False", " false", -1)
		participantsJson = strings.Replace(participantsJson, " None", " \"\"", -1)

		var participants []GrainParticipant
		err := json.Unmarshal([]byte(participantsJson), &participants)
		if err != nil {
			return err
		}
		g.RecordingData.Participants = participants
	}

	if g.RecordingData.TagsStr != "" {
		tagsJson := utils.ReplaceSingleQuotesWithDoubleQuotes(g.RecordingData.TagsStr)
		var tags []string
		err := json.Unmarshal([]byte(tagsJson), &tags)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GrainRecordingData) MeetingNoteContent() string {
	if g == nil || g.RecordingData.IntelligenceNotesMD == "" {
		return ""
	}

	// Remove hyperlinks using regex
	linkPattern := regexp.MustCompile(`\[\([^)]+\)\]\([^)]+\)`)
	cleanNotes := linkPattern.ReplaceAllStringFunc(g.RecordingData.IntelligenceNotesMD, func(match string) string {
		// Extract just the text between the brackets
		textStart := strings.Index(match, "(") + 1
		textEnd := strings.Index(match, ")")
		return match[textStart:textEnd]
	})

	// Calculate meeting duration
	duration := g.RecordingData.EndDatetime.Sub(g.RecordingData.StartDatetime)
	durationMinutes := int(duration.Minutes())

	// Build participants section
	var attendees, declines []string
	for _, p := range g.RecordingData.Participants {
		email := "No email"
		if p.Email != nil {
			email = *p.Email
		}
		var participant string
		if email == "No email" {
			participant = p.Name
		} else {
			participant = fmt.Sprintf("- %s (%s)", p.Name, email)
		}
		if p.ConfirmedAttendee {
			attendees = append(attendees, participant)
		} else {
			declines = append(declines, participant)
		}
	}

	// Build the final markdown
	var sb strings.Builder

	// Add meeting details
	sb.WriteString(fmt.Sprintf("# %s\n\n", g.RecordingData.Title))
	sb.WriteString(fmt.Sprintf("**Meeting Date:** %s\n", g.RecordingData.StartDatetime.Format("January 2, 2006 3:04 PM MST")))
	sb.WriteString(fmt.Sprintf("**Duration:** %d minutes\n\n", durationMinutes))

	// Add participants
	sb.WriteString("**Attended**\n")
	sb.WriteString(strings.Join(attendees, "\n"))
	sb.WriteString("\n\n")
	if len(declines) > 0 {
		sb.WriteString("**Declined/No Show**\n")
		sb.WriteString(strings.Join(declines, "\n"))
		sb.WriteString("\n\n")
	}

	// Add meeting notes
	sb.WriteString("## Meeting Notes\n")
	sb.WriteString(cleanNotes)
	sb.WriteString("\n\n")

	// Add recording link
	sb.WriteString(fmt.Sprintf("[View Recording](%s)\n", g.RecordingData.PublicURL))

	return sb.String()
}

func (r *GrainRecording) getDomains() []string {
	var domains []string

	for _, participant := range r.Participants {
		emailData := mailvalidate.ValidateEmailSyntax(*participant.Email)
		if !emailData.IsValid {
			continue
		}
		domains = append(domains, emailData.Domain)
	}
	return domains
}
