package hubspot

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

const (
	Partner  string = "Partner"
	Customer string = "Customer"
	Reseller string = "Reseller"
	Vendor   string = "Vendor"
)

const (
	Prospect string = "Prospect"
	Live     string = "Live"
)

func MapOrganization(inputJSON string) (string, error) {
	var input struct {
		ID         string `json:"id"`
		Archived   bool   `json:"archived"`
		CreatedAt  string `json:"createdAt"`
		UpdatedAt  string `json:"updatedAt"`
		Properties struct {
			Description       string `json:"description"`
			NumberOfEmployees int    `json:"numberofemployees"`
			Zip               string `json:"zip"`
			City              string `json:"city"`
			Name              string `json:"name"`
			Type              string `json:"type"`
			Phone             string `json:"phone"`
			State             string `json:"state"`
			Domain            string `json:"domain"`
			Address           string `json:"address"`
			Address2          string `json:"address2"`
			Country           string `json:"country"`
			Website           string `json:"website"`
			Industry          string `json:"industry"`
			IsPublic          bool   `json:"is_public"`
			NumNotes          int    `json:"num_notes"`
			Created           string `json:"createdate"`
			About             string `json:"about_us"`
			UpdatedAt         string `json:"hs_lastmodifieddate"`
			HubspotOwnerId    string `json:"hubspot_owner_id"`
		} `json:"properties"`
	}

	// Parse input JSON into temporary structure
	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	// Perform mapping
	output := entity.OrganizationData{
		BaseData: entity.BaseData{
			ExternalId:   input.ID,
			CreatedAtStr: input.CreatedAt,
			UpdatedAtStr: input.UpdatedAt,
		},
		Name:        input.Properties.Name,
		Description: input.Properties.Description,
		Website:     input.Properties.Website,
		Industry:    input.Properties.Industry,
		IsPublic:    input.Properties.IsPublic,
		Employees:   int64(input.Properties.NumberOfEmployees),
		Country:     input.Properties.Country,
		Region:      input.Properties.State,
		Locality:    input.Properties.City,
		Address:     input.Properties.Address,
		Address2:    input.Properties.Address2,
		Zip:         input.Properties.Zip,
		OwnerUser: &entity.ReferencedUser{
			ExternalOwnerId: input.Properties.HubspotOwnerId,
		},
		Domains: []string{input.Properties.Domain},
	}
	if input.Properties.Phone != "" {
		output.PhoneNumbers = []string{input.Properties.Phone}
	}
	switch input.Properties.Type {
	case "PROSPECT":
		output.RelationshipName = Customer
		output.RelationshipStage = Prospect
	case "PARTNER":
		output.RelationshipName = Partner
		output.RelationshipStage = Live
	case "RESELLER":
		output.RelationshipName = Reseller
		output.RelationshipStage = Live
	case "VENDOR":
		output.RelationshipName = Vendor
		output.RelationshipStage = Live
	}

	return utils.ToJson(output)
}

func MapUser(inputJSON string) (string, error) {
	var input struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		UserID    int    `json:"userId"`
		Archived  bool   `json:"archived"`
		LastName  string `json:"lastName"`
		CreatedAt string `json:"createdAt"`
		FirstName string `json:"firstName"`
		UpdatedAt string `json:"updatedAt"`
	}

	// Parse input JSON into temporary structure
	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	// Perform mapping
	output := entity.UserData{
		BaseData: entity.BaseData{
			CreatedAtStr: input.CreatedAt,
			UpdatedAtStr: input.UpdatedAt,
		},
		Name:      fmt.Sprintf("%s %s", input.FirstName, input.LastName),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	if input.UserID != 0 {
		output.ExternalId = fmt.Sprintf("%d", input.UserID)
	}

	// Map the "id" field to "externalOwnerId"
	output.ExternalOwnerId = input.ID

	// Convert output data to JSON
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
}

func MapNote(inputJSON string) (string, error) {
	// Unmarshal into input struct
	var input struct {
		ID         string `json:"id"`
		CreatedAt  string `json:"createdAt"`
		UpdatedAt  string `json:"updatedAt"`
		Properties struct {
			Body      string `json:"hs_note_body"`
			OwnerId   string `json:"hubspot_owner_id"`
			CreatedBy int    `json:"hs_created_by"`
		} `json:"properties"`
		Contacts  []interface{} `json:"contacts"`
		Companies []interface{} `json:"companies"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", err
	}

	// Create output struct
	var output = entity.NoteData{
		BaseData: entity.BaseData{
			ExternalId:   input.ID,
			CreatedAtStr: input.CreatedAt,
			UpdatedAtStr: input.UpdatedAt,
		},
	}

	// Map fields
	output.Content = input.Properties.Body
	output.ContentType = "text/html"
	output.CreatorUserExternalOwnerId = input.Properties.OwnerId
	output.CreatorUserExternalId = strconv.Itoa(input.Properties.CreatedBy)

	// Map contacts
	for _, contact := range input.Contacts {
		id := fmt.Sprint(contact)
		output.NotedContactsExternalIds = append(output.NotedContactsExternalIds, id)
	}

	// Map companies
	for _, company := range input.Companies {
		id := fmt.Sprint(company)
		output.NotedOrganizationsExternalIds = append(output.NotedOrganizationsExternalIds, id)
	}

	// Marshal output to JSON
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJSON), nil
}

func MapMeeting(inputJSON string) (string, error) {
	var input struct {
		ID         string `json:"id"`
		CreatedAt  string `json:"createdAt"`
		UpdatedAt  string `json:"updatedAt"`
		Properties struct {
			Title           string `json:"hs_meeting_title"`
			StartTime       string `json:"hs_meeting_start_time"`
			EndTime         string `json:"hs_meeting_end_time"`
			CreatedByUserId int    `json:"hs_created_by_user_id"`
			Html            string `json:"hs_meeting_body"`
			Text            string `json:"hs_body_preview"`
			Location        string `json:"hs_meeting_location"`
			MeetingUrl      string `json:"hs_meeting_external_url"`
		} `json:"properties"`

		Contacts []interface{} `json:"contacts"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", err
	}

	// Create output
	var output entity.MeetingData

	// Map ID
	output.ExternalId = input.ID
	output.CreatedAtStr = input.CreatedAt
	output.UpdatedAtStr = input.UpdatedAt
	output.StartedAtStr = input.Properties.StartTime
	output.EndedAtStr = input.Properties.EndTime
	output.Name = input.Properties.Title
	output.CreatorUserExternalId = fmt.Sprint(input.Properties.CreatedByUserId)
	output.MeetingUrl = input.Properties.MeetingUrl
	if len(input.Properties.Html) > 0 {
		output.Agenda = input.Properties.Html
		output.ContentType = "text/html"
	} else if len(input.Properties.Text) > 0 {
		output.Agenda = input.Properties.Text
		output.ContentType = "text/plain"
	}
	if len(input.Properties.Location) > 0 {
		if strings.HasPrefix(output.Location, "https://") {
			output.ConferenceUrl = input.Properties.Location
		} else {
			output.Location = input.Properties.Location
		}
	}

	// Map contacts
	for _, contact := range input.Contacts {
		id := fmt.Sprint(contact)
		output.ExternalContactsIds = append(output.ExternalContactsIds, id)
	}

	// Marshal output
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJSON), nil
}

const (
	INBOUND  string = "INBOUND"
	OUTBOUND string = "OUTBOUND"
)

func MapEmailMessage(inputJSON string) (string, error) {
	var input struct {
		ID         string        `json:"id"`
		CreatedAt  string        `json:"createdAt"`
		UpdatedAt  string        `json:"updatedAt"`
		Contacts   []interface{} `json:"contacts"`
		Properties struct {
			Html            string `json:"hs_email_html"`
			Text            string `json:"hs_email_text"`
			Subject         string `json:"hs_email_subject"`
			ThreadId        string `json:"hs_email_thread_id"`
			MessageId       string `json:"hs_email_message_id"`
			Direction       string `json:"hs_email_direction"`
			CreatedByUserId int    `json:"hs_created_by_user_id"`
			FromEmail       string `json:"hs_email_from_email"`
			ToEmail         string `json:"hs_email_to_email"`
			CcEmail         string `json:"hs_email_cc_email"`
			BccEmail        string `json:"hs_email_bcc_email"`
			FromFirstName   string `json:"hs_email_from_firstname"`
			FromLastName    string `json:"hs_email_from_lastname"`
			EmailStatus     string `json:"hs_email_status"`
		} `json:"properties"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", err
	}

	// Create output
	var output entity.EmailMessageData
	if input.Properties.ThreadId == "" || input.Properties.EmailStatus != "SENT" {
		output.Skip = true
		output.SkipReason = "Email is not sent or is not part of a thread"
		outputJSON, err := json.Marshal(output)
		if err != nil {
			return "", err
		}
		return string(outputJSON), nil
	}

	output.ExternalId = input.ID
	output.Html = input.Properties.Html
	output.Text = input.Properties.Text
	output.Subject = input.Properties.Subject
	output.EmailThreadId = input.Properties.ThreadId
	output.EmailMessageId = input.Properties.MessageId
	output.CreatedAtStr = input.CreatedAt
	output.UpdatedAtStr = input.UpdatedAt
	output.FromFirstName = input.Properties.FromFirstName
	output.FromLastName = input.Properties.FromLastName
	output.FromEmail = input.Properties.FromEmail
	output.ToEmail = strings.Split(input.Properties.ToEmail, ";")
	output.CcEmail = strings.Split(input.Properties.CcEmail, ";")
	output.BccEmail = strings.Split(input.Properties.BccEmail, ";")
	if input.Properties.CreatedByUserId != 0 {
		output.ExternalUserId = fmt.Sprintf("%d", input.Properties.CreatedByUserId)
	}

	if input.Properties.Direction == "INCOMING_EMAIL" {
		output.Direction = INBOUND
	} else {
		output.Direction = OUTBOUND
	}

	// Map contacts
	for _, contact := range input.Contacts {
		id := fmt.Sprint(contact)
		output.ExternalContactsIds = append(output.ExternalContactsIds, id)
	}

	// Marshal output
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJSON), nil
}

func MapContact(inputJSON string) (string, error) {
	var input struct {
		ID         string        `json:"id"`
		CreatedAt  string        `json:"createdAt"`
		UpdatedAt  string        `json:"updatedAt"`
		Companies  []interface{} `json:"companies,omitempty"`
		Properties struct {
			FirstName           string `json:"firstname,omitempty"`
			LastName            string `json:"lastname,omitempty"`
			JobTitle            string `json:"jobtitle,omitempty"`
			Email               string `json:"email,omitempty"`
			AdditionalEmails    string `json:"hs_additional_emails,omitempty"`
			PhoneNumber         string `json:"phone,omitempty"`
			AssociatedCompanyId *int   `json:"associatedcompanyid,omitempty"`
			OwnerId             string `json:"hubspot_owner_id,omitempty"`
			Country             string `json:"country,omitempty"`
			State               string `json:"state,omitempty"`
			City                string `json:"city,omitempty"`
			Zipcode             string `json:"zip,omitempty"`
			Address             string `json:"address,omitempty"`
			LifecycleStage      string `json:"lifecyclestage,omitempty"`
			Timezone            string `json:"hs_timezone,omitempty"`
		} `json:"properties"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", err
	}

	// Create output
	var output = entity.ContactData{
		BaseData: entity.BaseData{
			ExternalId:   input.ID,
			CreatedAtStr: input.CreatedAt,
			UpdatedAtStr: input.UpdatedAt,
		},
		FirstName: input.Properties.FirstName,
		LastName:  input.Properties.LastName,
		Email:     input.Properties.Email,
	}
	output.AdditionalEmails = strings.Split(input.Properties.AdditionalEmails, ";")
	output.PhoneNumber = input.Properties.PhoneNumber
	if input.Properties.AssociatedCompanyId != nil {
		output.Organizations = append(output.Organizations, entity.ReferencedOrganization{
			JobTitle:   input.Properties.JobTitle,
			ExternalId: fmt.Sprint(*input.Properties.AssociatedCompanyId),
		})
	}
	output.Country = input.Properties.Country
	output.UserExternalOwnerId = fmt.Sprint(input.Properties.OwnerId)
	output.Region = input.Properties.State
	output.Locality = input.Properties.City
	output.Zip = input.Properties.Zipcode
	output.Address = input.Properties.Address
	output.Timezone = convertToStandardTimezoneFormat(input.Properties.Timezone)

	if len(input.Properties.LifecycleStage) > 0 {
		output.TextCustomFields = append(output.TextCustomFields, entity.TextCustomField{
			Name:         "Hubspot Lifecycle Stage",
			Value:        input.Properties.LifecycleStage,
			CreatedAtStr: input.CreatedAt,
		})
		if isCustomerTag(input.Properties.LifecycleStage) {
			output.Tags = append(output.Tags, "CUSTOMER")
		} else if isProspectTag(input.Properties.LifecycleStage) {
			output.Tags = append(output.Tags, "PROSPECT")
		}
	}

	// Map contacts
	for _, contact := range input.Companies {
		id := fmt.Sprint(contact)
		// Check if organization already exists
		for _, organization := range output.Organizations {
			if organization.ExternalId == id {
				continue
			}
		}
		output.Organizations = append(output.Organizations, entity.ReferencedOrganization{
			ExternalId: id,
		})
	}

	// Marshal output
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJSON), nil
}

func isCustomerTag(hubspotLifecycleStage string) bool {
	customerLifecycleStages := map[string]bool{"customer": true}
	return customerLifecycleStages[hubspotLifecycleStage]
}

func isProspectTag(hubspotLifecycleStage string) bool {
	prospectLifecycleStages := map[string]bool{
		"lead": true, "subscriber": true, "marketingqualifiedlead": true, "salesqualifiedlead": true, "opportunity": true}
	return prospectLifecycleStages[hubspotLifecycleStage]
}

func convertToStandardTimezoneFormat(input string) string {
	if input == "" {
		return ""
	}

	parts := strings.Split(input, "_slash_")
	if len(parts) != 2 {
		return ""
	}

	region, city := parts[0], parts[1]
	region = cases.Title(language.Und).String(strings.ToLower(region))
	city = strings.ReplaceAll(city, "_", " ")
	city = cases.Title(language.Und).String(strings.ToLower(city))
	city = strings.ReplaceAll(city, " ", "_")

	return fmt.Sprintf("%s/%s", region, city)
}
