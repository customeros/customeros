package salesforce

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
)

func MapUser(inputJson string) (string, error) {
	var input struct {
		ID          string `json:"Id,omitempty"`
		FirstName   string `json:"FirstName,omitempty"`
		LastName    string `json:"LastName,omitempty"`
		Name        string `json:"Name,omitempty"`
		CreatedDate string `json:"CreatedDate,omitempty"`
		Email       string `json:"Email,omitempty"`
		Phone       string `json:"Phone,omitempty"`
		MobilePhone string `json:"MobilePhone,omitempty"`
	}

	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
		return "", err
	}

	if input.ID == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing id",
		}
		return utils.ToJson(output)
	}

	output := entity.UserData{
		BaseData: entity.BaseData{
			ExternalId:   input.ID,
			CreatedAtStr: input.CreatedDate,
		},
		Name:      input.Name,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}
	output.PhoneNumbers = append(output.PhoneNumbers, entity.PhoneNumber{
		Number:  input.Phone,
		Primary: true,
	})
	output.PhoneNumbers = append(output.PhoneNumbers, entity.PhoneNumber{
		Number:  input.MobilePhone,
		Primary: false,
		Label:   "MOBILE",
	})

	return utils.ToJson(output)
}

func MapOrganization(inputJson string) (string, error) {
	var input struct {
		Attributes struct {
			Type string `json:"type,omitempty"`
		} `json:"attributes,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJson), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}
	if input.Attributes.Type == "Account" {
		return mapOrganizationFromAccount(inputJson)
	} else if input.Attributes.Type == "Lead" {
		return mapOrganizationFromLead(inputJson)
	} else {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Attributes type not equal Account or Lead",
		}
		return utils.ToJson(output)
	}
}

func mapOrganizationFromAccount(inputJson string) (string, error) {
	var input struct {
		ID          string `json:"Id,omitempty"`
		Name        string `json:"Name,omitempty"`
		Phone       string `json:"Phone,omitempty"`
		CreatedDate string `json:"CreatedDate,omitempty"`
		Description string `json:"Description,omitempty"`
		Industry    string `json:"Industry,omitempty"`
		Website     string `json:"Website,omitempty"`
		OwnerId     string `json:"OwnerId,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJson), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}
	if input.ID == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing id",
		}
		return utils.ToJson(output)
	}
	if input.Name == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing name",
		}
		return utils.ToJson(output)
	}

	output := entity.OrganizationData{
		BaseData: entity.BaseData{
			ExternalId:          input.ID,
			CreatedAtStr:        input.CreatedDate,
			ExternalSourceTable: utils.StringPtr("account"),
		},
		Name:        input.Name,
		Description: input.Description,
		Website:     input.Website,
		Industry:    input.Industry,
	}
	if input.Phone != "" {
		output.PhoneNumbers = []entity.PhoneNumber{
			{
				Number:  input.Phone,
				Primary: true,
			},
		}
	}
	if input.OwnerId != "" {
		output.OwnerUser = &entity.ReferencedUser{
			ExternalId: input.OwnerId,
		}
	}

	return utils.ToJson(output)
}

func mapOrganizationFromLead(inputJson string) (string, error) {
	var input struct {
		ID          string `json:"Id,omitempty"`
		Email       string `json:"Email,omitempty"`
		Industry    string `json:"Industry,omitempty"`
		Company     string `json:"Company,omitempty"`
		CreatedDate string `json:"CreatedDate,omitempty"`
		IsConverted bool   `json:"IsConverted,omitempty"`
		OwnerId     string `json:"OwnerId,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJson), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}
	if input.IsConverted {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Lead already converted into account",
		}
		return utils.ToJson(output)
	}
	if input.ID == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing id",
		}
		return utils.ToJson(output)
	}
	domain := extractDomain(input.Email)
	if domain == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing email domain",
		}
		return utils.ToJson(output)
	}

	output := entity.OrganizationData{
		BaseData: entity.BaseData{
			ExternalId:          input.ID,
			CreatedAtStr:        input.CreatedDate,
			ExternalSourceTable: utils.StringPtr("lead"),
		},
		CreateByDomain: true,
		Domains:        []string{domain},
		Industry:       input.Industry,
	}
	if input.OwnerId != "" {
		output.OwnerUser = &entity.ReferencedUser{
			ExternalId: input.OwnerId,
		}
	}

	return utils.ToJson(output)
}

func MapContact(inputJson string) (string, error) {
	var input struct {
		Attributes struct {
			Type string `json:"type,omitempty"`
		} `json:"attributes,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJson), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}
	if input.Attributes.Type == "Contact" {
		return mapContactFromContact(inputJson)
	} else if input.Attributes.Type == "Lead" {
		return mapContactFromLead(inputJson)
	} else {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Attributes type not equal Account or Lead",
		}
		return utils.ToJson(output)
	}
}

func mapContactFromContact(inputJson string) (string, error) {
	var input struct {
		ID             string `json:"Id,omitempty"`
		AccountId      string `json:"AccountId,omitempty"`
		Name           string `json:"Name,omitempty"`
		FirstName      string `json:"FirstName,omitempty"`
		LastName       string `json:"LastName,omitempty"`
		Email          string `json:"Email,omitempty"`
		Phone          string `json:"Phone,omitempty"`
		HomePhone      string `json:"HomePhone,omitempty"`
		MobilePhone    string `json:"MobilePhone,omitempty"`
		OtherPhone     string `json:"OtherPhone,omitempty"`
		AssistantPhone string `json:"AssistantPhone,omitempty"`
		Fax            string `json:"Fax,omitempty"`
		Title          string `json:"Title,omitempty"`
		CreatedDate    string `json:"CreatedDate,omitempty"`
		MailingAddress struct {
			City       string `json:"city,omitempty"`
			Country    string `json:"country,omitempty"`
			State      string `json:"state,omitempty"`
			Street     string `json:"street,omitempty"`
			PostalCode string `json:"postalCode,omitempty"`
		} `json:"MailingAddress,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJson), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}
	if input.ID == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing id",
		}
		return utils.ToJson(output)
	}

	output := entity.ContactData{
		BaseData: entity.BaseData{
			ExternalId:          input.ID,
			CreatedAtStr:        input.CreatedDate,
			ExternalSourceTable: utils.StringPtr("contact"),
		},
		Name:       input.Name,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Email:      input.Email,
		Country:    input.MailingAddress.Country,
		Region:     input.MailingAddress.State,
		Locality:   input.MailingAddress.City,
		Street:     input.MailingAddress.Street,
		PostalCode: input.MailingAddress.PostalCode,
	}
	if input.Phone != "" {
		output.PhoneNumbers = []entity.PhoneNumber{
			{
				Number:  input.Phone,
				Primary: true,
			},
		}
	}
	output.PhoneNumbers = append(output.PhoneNumbers, entity.PhoneNumber{
		Number:  input.MobilePhone,
		Primary: false,
		Label:   "MOBILE",
	}, entity.PhoneNumber{
		Number:  input.HomePhone,
		Primary: false,
		Label:   "HOME",
	}, entity.PhoneNumber{
		Number:  input.OtherPhone,
		Primary: false,
		Label:   "OTHER",
	}, entity.PhoneNumber{
		Number:  input.AssistantPhone,
		Primary: false,
		Label:   "ASSISTANT",
	}, entity.PhoneNumber{
		Number:  input.Fax,
		Primary: false,
		Label:   "FAX",
	})
	if input.AccountId != "" {
		output.Organizations = []entity.ReferencedOrganization{
			{
				ExternalId: input.AccountId,
				JobTitle:   input.Title,
			},
		}
	}

	return utils.ToJson(output)
}

func mapContactFromLead(inputJson string) (string, error) {
	var input struct {
		ID          string `json:"Id,omitempty"`
		IsConverted bool   `json:"IsConverted,omitempty"`
		Name        string `json:"Name,omitempty"`
		FirstName   string `json:"FirstName,omitempty"`
		LastName    string `json:"LastName,omitempty"`
		Email       string `json:"Email,omitempty"`
		Phone       string `json:"Phone,omitempty"`
		MobilePhone string `json:"MobilePhone,omitempty"`
		Title       string `json:"Title,omitempty"`
		CreatedDate string `json:"CreatedDate,omitempty"`
		City        string `json:"City,omitempty"`
		Country     string `json:"Country,omitempty"`
		Street      string `json:"Street,omitempty"`
		State       string `json:"State,omitempty"`
		PostalCode  string `json:"PostalCode,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJson), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	if input.ID == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing id",
		}
		return utils.ToJson(output)
	}
	if input.IsConverted {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Lead already converted",
		}
		return utils.ToJson(output)
	}
	emailDomain := extractDomain(input.Email)
	if emailDomain == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing email domain",
		}
		return utils.ToJson(output)
	}

	output := entity.ContactData{
		BaseData: entity.BaseData{
			ExternalId:          input.ID,
			CreatedAtStr:        input.CreatedDate,
			ExternalSourceTable: utils.StringPtr("lead"),
		},
		Name:                 input.Name,
		FirstName:            input.FirstName,
		LastName:             input.LastName,
		Email:                input.Email,
		Country:              input.Country,
		Region:               input.State,
		Locality:             input.City,
		Street:               input.Street,
		PostalCode:           input.PostalCode,
		OrganizationRequired: true,
	}
	output.PhoneNumbers = []entity.PhoneNumber{
		{
			Number:  input.Phone,
			Primary: true,
		},
		{
			Number:  input.MobilePhone,
			Primary: false,
			Label:   "MOBILE",
		},
	}
	output.Organizations = append(output.Organizations, entity.ReferencedOrganization{
		Domain:   emailDomain,
		JobTitle: input.Title,
	})

	return utils.ToJson(output)
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func MapLogEntry(inputJson string) (string, error) {
	var input struct {
		ID          string `json:"Id,omitempty"`
		Body        string `json:"Body,omitempty"`
		CreatedDate string `json:"CreatedDate,omitempty"`
		Type        string `json:"Type,omitempty"`
		Status      string `json:"Status,omitempty"`
		IsRichText  bool   `json:"IsRichText,omitempty"`
		CreatedById string `json:"CreatedById,omitempty"`
		ParentId    string `json:"ParentId,omitempty"`
	}

	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
		return "", err
	}

	if input.ID == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing id",
		}
		return utils.ToJson(output)
	}
	if input.Body == "" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Missing body",
		}
		return utils.ToJson(output)
	}
	if input.Status != "Published" {
		output := entity.BaseData{
			Skip:       true,
			SkipReason: "Not published",
		}
		return utils.ToJson(output)
	}

	output := entity.LogEntryData{
		BaseData: entity.BaseData{
			ExternalId:          input.ID,
			CreatedAtStr:        input.CreatedDate,
			ExternalSourceTable: utils.StringPtr("feeditem"),
		},
		Content:      input.Body,
		StartedAtStr: input.CreatedDate,
		AuthorUser: entity.ReferencedUser{
			ExternalId: input.CreatedById,
		},
		LoggedOrganization: entity.ReferencedOrganization{
			ExternalId: input.ParentId,
		},
		LoggedEntityRequired: true,
	}
	if input.IsRichText {
		output.ContentType = "text/html"
	} else {
		output.ContentType = "text/plain"
	}

	return utils.ToJson(output)
}
