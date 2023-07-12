package hubspot

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Output struct {
	Name              string   `json:"name,omitempty"`
	FirstName         string   `json:"firstName,omitempty"`
	LastName          string   `json:"lastName,omitempty"`
	Email             string   `json:"email,omitempty"`
	PhoneNumber       string   `json:"phoneNumber,omitempty"`
	CreatedAt         string   `json:"createdAt,omitempty"`
	UpdatedAt         string   `json:"updatedAt,omitempty"`
	ExternalId        string   `json:"externalId,omitempty"`
	ExternalOwnerId   string   `json:"externalOwnerId,omitempty"`
	ExternalUserId    string   `json:"externalUserId,omitempty"`
	ExternalCreatorId string   `json:"externalCreatorId,omitempty"`
	Description       string   `json:"description,omitempty"`
	Domains           []string `json:"domains,omitempty"`
	Notes             []struct {
		FieldSource string `json:"fieldSource,omitempty"`
		Note        string `json:"note,omitempty"`
	} `json:"notes,omitempty"`
	Website             string `json:"website,omitempty"`
	Industry            string `json:"industry,omitempty"`
	IsPublic            bool   `json:"isPublic,omitempty"`
	Employees           int    `json:"employees,omitempty"`
	ExternalUrl         string `json:"externalUrl,omitempty"`
	ExternalSourceTable string `json:"externalSourceTable,omitempty"`
	ExternalSystem      string `json:"externalSystem,omitempty"`
	ExternalSyncId      string `json:"externalSyncId,omitempty"`
	LocationName        string `json:"locationName,omitempty"`
	Country             string `json:"country,omitempty"`
	Region              string `json:"region,omitempty"`
	Locality            string `json:"locality,omitempty"`
	Address             string `json:"address,omitempty"`
	Address2            string `json:"address2,omitempty"`
	Zip                 string `json:"zip,omitempty"`
	RelationshipName    string `json:"relationshipName,omitempty"`
	RelationshipStage   string `json:"relationshipStage,omitempty"`
	ParentOrganization  struct {
		ExternalId           string `json:"externalId,omitempty"`
		OrganizationRelation string `json:"organizationRelation,omitempty"`
		Type                 string `json:"type,omitempty"`
	} `json:"parentOrganization,omitempty"`
	Html                     string   `json:"html,omitempty"`
	Text                     string   `json:"text,omitempty"`
	ContactsExternalIds      []string `json:"contactsExternalIds,omitempty"`
	OrganizationsExternalIds []string `json:"organizationsExternalIds,omitempty"`
	MentionedTags            []string `json:"mentionedTags,omitempty"`
}

func MapOrganization(inputJSON string) (string, error) {
	var temp struct {
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
	err := json.Unmarshal([]byte(inputJSON), &temp)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	// Perform mapping
	outputData := Output{
		ExternalId:      temp.ID,
		CreatedAt:       temp.CreatedAt,
		UpdatedAt:       temp.UpdatedAt,
		Name:            temp.Properties.Name,
		Description:     temp.Properties.Description,
		Website:         temp.Properties.Website,
		Industry:        temp.Properties.Industry,
		IsPublic:        temp.Properties.IsPublic,
		Employees:       temp.Properties.NumberOfEmployees,
		PhoneNumber:     temp.Properties.Phone,
		Country:         temp.Properties.Country,
		Region:          temp.Properties.State,
		Locality:        temp.Properties.City,
		Address:         temp.Properties.Address,
		Address2:        temp.Properties.Address2,
		Zip:             temp.Properties.Zip,
		ExternalOwnerId: temp.Properties.HubspotOwnerId,
		Domains:         []string{temp.Properties.Domain},
	}
	switch temp.Properties.Type {
	case "PROSPECT":
		outputData.RelationshipName = "Customer"
		outputData.RelationshipStage = "Prospect"
	case "PARTNER":
		outputData.RelationshipName = "Partner"
		outputData.RelationshipStage = "Live"
	case "RESELLER":
		outputData.RelationshipName = "Reseller"
		outputData.RelationshipStage = "Live"
	case "VENDOR":
		outputData.RelationshipName = "Vendor"
		outputData.RelationshipStage = "Live"
	}

	// Convert output data to JSON
	outputJSON, err := json.Marshal(outputData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
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
	outputData := Output{
		Name:      fmt.Sprintf("%s %s", input.FirstName, input.LastName),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}

	if input.UserID != 0 {
		outputData.ExternalId = fmt.Sprintf("%d", input.UserID)
	}

	// Map the "id" field to "externalOwnerId"
	outputData.ExternalOwnerId = input.ID

	// Convert output data to JSON
	outputJSON, err := json.Marshal(outputData)
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
	var output Output

	// Map fields
	output.ExternalId = input.ID
	output.CreatedAt = input.CreatedAt
	output.UpdatedAt = input.UpdatedAt
	output.Html = input.Properties.Body
	output.ExternalOwnerId = input.Properties.OwnerId
	output.ExternalUserId = strconv.Itoa(input.Properties.CreatedBy)

	// Map contacts
	for _, contact := range input.Contacts {
		id := fmt.Sprint(contact)
		output.ContactsExternalIds = append(output.ContactsExternalIds, id)
	}

	// Map companies
	for _, company := range input.Companies {
		id := fmt.Sprint(company)
		output.OrganizationsExternalIds = append(output.OrganizationsExternalIds, id)
	}

	// Marshal output to JSON
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJSON), nil
}
