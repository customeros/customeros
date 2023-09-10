package zendesk_support

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
)

func MapUser(inputJSON string) (string, error) {
	var input struct {
		ID    int64  `json:"id"`
		Name  string `json:"name,omitempty"`
		Role  string `json:"role,omitempty"`
		Email string `json:"email,omitempty"`
		Phone string `json:"phone,omitempty"`

		CreatedAt string `json:"created_at,omitempty"`
		UpdatedAt string `json:"updated_at,omitempty"`

		IanaTimeZone    string `json:"iana_time_zone,omitempty"`
		RestrictedAgent bool   `json:"restricted_agent,omitempty"`
	}

	// Parse input
	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	// Map to output
	output := entity.UserData{
		BaseData: entity.BaseData{
			ExternalId:   fmt.Sprintf("%d", input.ID),
			CreatedAtStr: input.CreatedAt,
			UpdatedAtStr: input.UpdatedAt,
		},
		Name:        input.Name,
		Email:       utils.IfNotNilString(input.Email),
		PhoneNumber: utils.IfNotNilString(input.Phone),
	}

	if input.Role == "end-user" || input.ID == 0 {
		output.Skip = true
		output.SkipReason = "User is not an agent"
	}

	// Return JSON
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
}

func MapOrganization(inputJSON string) (string, error) {
	// Parse into generic map
	var inputMap map[string]interface{}
	if err := json.Unmarshal([]byte(inputJSON), &inputMap); err != nil {
		return "", err
	}

	// Check if "role" key exists
	if _, ok := inputMap["role"]; ok {
		// "role" exists, input is User
		return mapOrganizationFromUser(inputJSON)
	} else {
		// No "role", input is Organization
		return mapOrganizationFromOrg(inputJSON)
	}
}

func mapOrganizationFromOrg(inputJSON string) (string, error) {
	var input struct {
		ID                 int64    `json:"id,omitempty"`
		URL                string   `json:"url,omitempty"`
		Name               string   `json:"name,omitempty"`
		Tags               []string `json:"tags,omitempty"`
		Notes              *string  `json:"notes,omitempty"`
		Details            string   `json:"details,omitempty"`
		GroupID            *int64   `json:"group_id,omitempty"`
		CreatedAt          string   `json:"created_at,omitempty"`
		UpdatedAt          string   `json:"updated_at,omitempty"`
		ExternalID         *string  `json:"external_id,omitempty"`
		DomainNames        []string `json:"domain_names,omitempty"`
		SharedTickets      bool     `json:"shared_tickets,omitempty"`
		SharedComments     bool     `json:"shared_comments,omitempty"`
		OrganizationFields struct{} `json:"organization_fields,omitempty"`
	}

	// Parse input JSON into temporary structure
	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	// Perform mapping
	output := entity.OrganizationData{
		BaseData: entity.BaseData{
			ExternalId:          fmt.Sprintf("%d", input.ID),
			CreatedAtStr:        input.CreatedAt,
			UpdatedAtStr:        input.UpdatedAt,
			ExternalUrl:         input.URL,
			ExternalSourceTable: utils.StringPtr("organizations"),
		},
		Name:    input.Name,
		Domains: input.DomainNames,
	}
	if input.ID == 0 {
		output.Skip = true
		output.SkipReason = "Missing ID"
	}
	if input.Details != "" {
		output.Notes = append(output.Notes, entity.OrganizationNote{
			FieldSource: "details",
			Note:        input.Details,
		})
	}

	// Convert output data to JSON
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
}

func mapOrganizationFromUser(inputJSON string) (string, error) {
	var input struct {
		ID             int64  `json:"id"`
		Name           string `json:"name,omitempty"`
		Role           string `json:"role,omitempty"`
		Email          string `json:"email,omitempty"`
		Phone          string `json:"phone,omitempty"`
		Details        string `json:"details,omitempty"`
		URL            string `json:"url,omitempty"`
		OrganizationId int64  `json:"organization_id,omitempty"`
		Notes          string `json:"notes,omitempty"`
		CreatedAt      string `json:"created_at,omitempty"`
		UpdatedAt      string `json:"updated_at,omitempty"`

		IanaTimeZone    string `json:"iana_time_zone,omitempty"`
		RestrictedAgent bool   `json:"restricted_agent,omitempty"`
	}

	// Parse input JSON into temporary structure
	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	// Perform mapping
	output := entity.OrganizationData{
		BaseData: entity.BaseData{
			ExternalId:          fmt.Sprintf("%d", input.ID),
			CreatedAtStr:        input.CreatedAt,
			UpdatedAtStr:        input.UpdatedAt,
			ExternalUrl:         input.URL,
			ExternalSourceTable: utils.StringPtr("users"),
		},
		Name:        input.Name,
		PhoneNumber: input.Phone,
	}
	if input.Role != "end-user" || input.ID == 0 {
		output.Skip = true
		output.SkipReason = "User is not an agent"
	}
	if input.Details != "" {
		output.Notes = append(output.Notes, entity.OrganizationNote{
			FieldSource: "details",
			Note:        input.Details,
		})
	}
	if len(input.Email) > 0 && !strings.HasSuffix(input.Email, "@without-email.com") {
		output.Email = input.Email
	}
	if input.OrganizationId > 0 {
		output.ParentOrganization = &entity.ParentOrganization{
			Organization: entity.ReferencedOrganization{
				ExternalId: fmt.Sprintf("%d", input.OrganizationId),
			},
			OrganizationRelation: "subsidiary",
			Type:                 "store",
		}
	}
	if input.Notes != "" {
		output.Notes = append(output.Notes, entity.OrganizationNote{
			Note:        input.Notes,
			FieldSource: "notes",
		})
	}

	// Convert output data to JSON
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
}

func MapIssue(inputJSON string) (string, error) {
	var input struct {
		ID              int64    `json:"id"`
		CreatedAt       string   `json:"created_at,omitempty"`
		UpdatedAt       string   `json:"updated_at,omitempty"`
		URL             string   `json:"url,omitempty"`
		Subject         string   `json:"subject,omitempty"`
		Status          string   `json:"status,omitempty"`
		Priority        string   `json:"priority,omitempty"`
		Description     string   `json:"description,omitempty"`
		Tags            []string `json:"tags,omitempty"`
		Type            string   `json:"type,omitempty"`
		RequesterId     int64    `json:"requester_id,omitempty"`
		AssigneeId      int64    `json:"assignee_id,omitempty"`
		CollaboratorIds []int64  `json:"collaborator_ids,omitempty"`
		FollowerIds     []int64  `json:"follower_ids,omitempty"`
	}

	// Parse input
	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	// Map to output
	output := entity.IssueData{
		BaseData: entity.BaseData{
			ExternalId:   fmt.Sprintf("%d", input.ID),
			CreatedAtStr: input.CreatedAt,
			UpdatedAtStr: input.UpdatedAt,
			ExternalUrl:  input.URL,
		},
		Subject:     input.Subject,
		Status:      input.Status,
		Priority:    input.Priority,
		Description: input.Description,
		Tags:        input.Tags,
	}
	if input.ID == 0 {
		output.Skip = true
		output.SkipReason = "Missing issue ID"
	}
	if input.Type != "" {
		output.Tags = append(output.Tags, "type:"+input.Type)
	}
	if input.RequesterId > 0 {
		output.ReporterOrganizationExternalId = fmt.Sprintf("%d", input.RequesterId)
	}
	if input.AssigneeId > 0 {
		output.AssigneeUserExternalId = fmt.Sprintf("%d", input.AssigneeId)
	}
	for _, collaboratorId := range input.CollaboratorIds {
		output.CollaboratorUserExternalIds = append(output.CollaboratorUserExternalIds, fmt.Sprintf("%d", collaboratorId))
	}
	for _, followerId := range input.FollowerIds {
		output.FollowerUserExternalIds = append(output.FollowerUserExternalIds, fmt.Sprintf("%d", followerId))
	}

	// Return JSON
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
}

func MapNote(inputJSON string) (string, error) {
	// Unmarshal into input struct
	var input struct {
		ID        int64  `json:"id"`
		CreatedAt string `json:"created_at,omitempty"`
		Public    bool   `json:"public,omitempty"`
		HtmlBody  string `json:"html_body,omitempty"`
		PlainBody string `json:"plain_body,omitempty"`
		Body      string `json:"body,omitempty"`
		AuthorId  int64  `json:"author_id,omitempty"`
		TicketId  int64  `json:"ticket_id,omitempty"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", err
	}

	// Create output struct
	var output entity.NoteData

	if input.ID == 0 {
		output.Skip = true
		output.SkipReason = "Missing ticket comment ID"
	}
	if input.Public == true {
		output.Skip = true
		output.SkipReason = "Ticket comment is public, it will be synced as interaction event"
	}
	// Map fields
	output.ExternalId = fmt.Sprintf("%d", input.ID)
	output.CreatedAtStr = input.CreatedAt
	output.Content = input.HtmlBody
	output.ContentType = "text/html"
	output.Text = input.Body
	if input.AuthorId > 0 {
		output.CreatorExternalId = fmt.Sprintf("%d", input.AuthorId)
	}
	if input.TicketId > 0 {
		output.MentionedIssueExternalId = fmt.Sprintf("%d", input.TicketId)
	}

	// Marshal output to JSON
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJSON), nil
}

func MapInteractionEvent(inputJSON string) (string, error) {
	// Unmarshal into input struct
	var input struct {
		ID        int64  `json:"id"`
		CreatedAt string `json:"created_at,omitempty"`
		Public    bool   `json:"public,omitempty"`
		HtmlBody  string `json:"html_body,omitempty"`
		PlainBody string `json:"plain_body,omitempty"`
		Body      string `json:"body,omitempty"`
		AuthorId  int64  `json:"author_id,omitempty"`
		TicketId  int64  `json:"ticket_id,omitempty"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", err
	}

	// Create output struct
	var output entity.InteractionEventData

	if input.ID == 0 {
		output.Skip = true
		output.SkipReason = "Missing ticket comment ID"
	}
	if input.Public == false {
		output.Skip = true
		output.SkipReason = "Ticket comment is private, it will be synced as note"
	}
	// Map fields
	output.ExternalId = fmt.Sprintf("%d", input.ID)
	output.CreatedAtStr = input.CreatedAt
	output.Type = "ISSUE"
	if input.HtmlBody != "" {
		output.Content = input.HtmlBody
		output.ContentType = "text/html"
	} else if input.PlainBody != "" {
		output.Content = input.PlainBody
		output.ContentType = "text/plain"
	}
	if input.TicketId > 0 {
		output.PartOfIssue = entity.ReferencedIssue{
			ExternalId: fmt.Sprintf("%d", input.TicketId),
		}
	}
	if input.AuthorId > 0 {
		output.SentBy = entity.InteractionEventParticipant{
			ReferencedParticipant: entity.ReferencedParticipant{
				ExternalId: fmt.Sprintf("%d", input.AuthorId),
			},
		}
	}

	// Marshal output to JSON
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJSON), nil
}
