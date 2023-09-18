package pipedrive

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
)

func MapUser(inputJson string) (string, error) {
	var input struct {
		ID        int64  `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Email     string `json:"email,omitempty"`
		Phone     string `json:"phone,omitempty"`
		CreatedAt string `json:"created,omitempty"`
		Modified  string `json:"modified,omitempty"`
	}

	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
		return "", err
	}

	output := entity.UserData{
		BaseData: entity.BaseData{
			ExternalId:   fmt.Sprintf("%d", input.ID),
			CreatedAtStr: input.CreatedAt,
			UpdatedAtStr: input.Modified,
		},
		Name:  input.Name,
		Email: input.Email,
		PhoneNumbers: []entity.PhoneNumber{
			{
				Number: input.Phone,
			},
		},
	}
	if input.ID == 0 {
		output.Skip = true
		output.SkipReason = "Missing external id"
	}

	outputJson, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJson), nil
}

func MapOrganization(inputJSON string) (string, error) {
	var input struct {
		ID                int64  `json:"id,omitempty"`
		Name              string `json:"name,omitempty"`
		Address           string `json:"address,omitempty"`
		AddTime           string `json:"add_time,omitempty"`
		UpdateTime        string `json:"update_time,omitempty"`
		OwnerID           int64  `json:"owner_id,omitempty"`
		PeopleCount       int    `json:"people_count,omitempty"`
		AddressCountry    string `json:"address_country,omitempty"`
		CountryCode       string `json:"country_code,omitempty"`
		AddressLocality   string `json:"address_locality,omitempty"`
		AddressPostalCode string `json:"address_postal_code,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	output := entity.OrganizationData{
		BaseData: entity.BaseData{
			ExternalId:   fmt.Sprintf("%d", input.ID),
			CreatedAtStr: input.AddTime,
			UpdatedAtStr: input.UpdateTime,
		},
		Name:    input.Name,
		Address: input.Address,
		OwnerUser: &entity.ReferencedUser{
			ExternalId: fmt.Sprintf("%d", input.OwnerID),
		},
		Employees: int64(input.PeopleCount),
		Country:   utils.StringFirstNonEmpty(input.AddressCountry, input.CountryCode),
		Locality:  input.AddressLocality,
		Zip:       input.AddressPostalCode,
	}
	if input.ID == 0 {
		output.Skip = true
		output.SkipReason = "Missing external id"
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
}

func MapContact(inputJSON string) (string, error) {
	var input struct {
		ID         int64  `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
		FirstName  string `json:"first_name,omitempty"`
		LastName   string `json:"last_name,omitempty"`
		Active     bool   `json:"active_flag,omitempty"`
		AddTime    string `json:"add_time,omitempty"`
		UpdateTime string `json:"update_time,omitempty"`
		OrgId      int64  `json:"org_id,omitempty"`
		OwnerId    int64  `json:"owner_id,omitempty"`
		Emails     []struct {
			Value   string `json:"value,omitempty"`
			Primary bool   `json:"primary,omitempty"`
		} `json:"email,omitempty"`
		Phones []struct {
			Value   string `json:"value,omitempty"`
			Primary bool   `json:"primary,omitempty"`
		} `json:"phone,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	output := entity.ContactData{
		BaseData: entity.BaseData{
			ExternalId:   fmt.Sprintf("%d", input.ID),
			CreatedAtStr: input.AddTime,
			UpdatedAtStr: input.UpdateTime,
		},
		Name:      input.Name,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	if input.ID == 0 {
		output.Skip = true
		output.SkipReason = "Missing external id"
	}
	if input.OrgId != 0 {
		output.Organizations = append(output.Organizations, entity.ReferencedOrganization{
			ExternalId: fmt.Sprintf("%d", input.OrgId),
		})
	}
	if input.OwnerId != 0 {
		output.UserExternalUserId = fmt.Sprintf("%d", input.OwnerId)
	}

	var primaryEmailFound = false
	for _, email := range input.Emails {
		if email.Value != "" {
			if email.Primary && !primaryEmailFound {
				output.Email = email.Value
				primaryEmailFound = true
			} else {
				output.AdditionalEmails = append(output.AdditionalEmails, email.Value)
			}
		}
	}
	for _, phone := range input.Phones {
		if phone.Value != "" {
			output.PhoneNumbers = append(output.PhoneNumbers, entity.PhoneNumber{
				Number:  phone.Value,
				Primary: phone.Primary,
			})
		}
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
}

func MapNote(inputJSON string) (string, error) {
	var input struct {
		ID         int64  `json:"id,omitempty"`
		Content    string `json:"content,omitempty"`
		UserId     int64  `json:"user_id,omitempty"`
		AddTime    string `json:"add_time,omitempty"`
		UpdateTime string `json:"update_time,omitempty"`
		PersonId   int64  `json:"person_id,omitempty"`
		OrgId      int64  `json:"org_id,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	output := entity.NoteData{
		BaseData: entity.BaseData{
			ExternalId:   fmt.Sprintf("%d", input.ID),
			CreatedAtStr: input.AddTime,
			UpdatedAtStr: input.UpdateTime,
		},
	}
	if input.ID == 0 {
		output.Skip = true
		output.SkipReason = "Missing external id"
	}
	if input.UserId != 0 {
		output.CreatorUserExternalId = fmt.Sprintf("%d", input.UserId)
	}
	if input.PersonId != 0 {
		output.NotedContactsExternalIds = append(output.NotedContactsExternalIds, fmt.Sprintf("%d", input.PersonId))
	}
	if input.OrgId != 0 {
		output.NotedOrganizationsExternalIds = append(output.NotedOrganizationsExternalIds, fmt.Sprintf("%d", input.OrgId))
	}
	if strings.Contains(input.Content, "<") {
		output.Content = input.Content
		output.ContentType = "text/html"
	} else {
		output.Content = input.Content
		output.ContentType = "text/plain"
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output JSON: %v", err)
	}

	return string(outputJSON), nil
}
