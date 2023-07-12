package hubspot

import (
	"encoding/json"
	"fmt"
)

type OutputData struct {
	Name            string   `json:"name"`
	FirstName       string   `json:"firstName"`
	LastName        string   `json:"lastName"`
	Email           string   `json:"email"`
	PhoneNumber     string   `json:"phoneNumber"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
	ExternalID      string   `json:"externalId,omitempty"`
	ExternalOwnerID string   `json:"externalOwnerId,omitempty"`
	Description     string   `json:"description"`
	Domains         []string `json:"domains"`
	Notes           []struct {
		FieldSource string `json:"fieldSource"`
		Note        string `json:"note"`
	} `json:"notes"`
	Website             string `json:"website"`
	Industry            string `json:"industry"`
	IsPublic            bool   `json:"isPublic"`
	Employees           int    `json:"employees"`
	ExternalUrl         string `json:"externalUrl"`
	ExternalSourceTable string `json:"externalSourceTable"`
	UserExternalOwnerId string `json:"userExternalOwnerId"`
	ExternalSystem      string `json:"externalSystem"`
	ExternalSyncId      string `json:"externalSyncId"`
	LocationName        string `json:"locationName"`
	Country             string `json:"country"`
	Region              string `json:"region"`
	Locality            string `json:"locality"`
	Address             string `json:"address"`
	Address2            string `json:"address2"`
	Zip                 string `json:"zip"`
	RelationshipName    string `json:"relationshipName"`
	RelationshipStage   string `json:"relationshipStage"`
	ParentOrganization  struct {
		ExternalId           string `json:"externalId"`
		OrganizationRelation string `json:"organizationRelation"`
		Type                 string `json:"type"`
	} `json:"parentOrganization"`
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
	outputData := OutputData{
		ExternalID:          temp.ID,
		CreatedAt:           temp.CreatedAt,
		UpdatedAt:           temp.UpdatedAt,
		Name:                temp.Properties.Name,
		Description:         temp.Properties.Description,
		Website:             temp.Properties.Website,
		Industry:            temp.Properties.Industry,
		IsPublic:            temp.Properties.IsPublic,
		Employees:           temp.Properties.NumberOfEmployees,
		PhoneNumber:         temp.Properties.Phone,
		Country:             temp.Properties.Country,
		Region:              temp.Properties.State,
		Locality:            temp.Properties.City,
		Address:             temp.Properties.Address,
		Address2:            temp.Properties.Address2,
		Zip:                 temp.Properties.Zip,
		UserExternalOwnerId: temp.Properties.HubspotOwnerId,
		Domains:             []string{temp.Properties.Domain},
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
	var temp struct {
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
	err := json.Unmarshal([]byte(inputJSON), &temp)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	// Perform mapping
	outputData := OutputData{
		Name:      fmt.Sprintf("%s %s", temp.FirstName, temp.LastName),
		FirstName: temp.FirstName,
		LastName:  temp.LastName,
		Email:     temp.Email,
		CreatedAt: temp.CreatedAt,
		UpdatedAt: temp.UpdatedAt,
	}

	if temp.UserID != 0 {
		outputData.ExternalID = fmt.Sprintf("%d", temp.UserID)
	}

	// Map the "id" field to "externalOwnerId"
	outputData.ExternalOwnerID = temp.ID

	// Convert output data to JSON
	outputJSON, err := json.Marshal(outputData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
}
