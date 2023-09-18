package salesforce

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
			ExternalId:   input.ID,
			CreatedAtStr: input.CreatedDate,
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

	return utils.ToJson(output)
}

func MapContact(inputJson string) (string, error) {
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
			ExternalId:   input.ID,
			CreatedAtStr: input.CreatedDate,
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
