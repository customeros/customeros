package shopify

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapOrganization(inputJson string) (string, error) {
	var input struct {
		Email        string `json:"email"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		ID           int64  `json:"id"`
		CreatedAtStr string `json:"created_at"`
		Addresses    []struct {
			Company string `json:"company"`
		}
	}

	err := json.Unmarshal([]byte(inputJson), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	isCustomer := false
	var name string

	if input.Addresses != nil && len(input.Addresses) > 0 && input.Addresses[0].Company != "" {
		name = input.Addresses[0].Company
		isCustomer = true
	}

	output := entity.OrganizationData{
		BaseData: entity.BaseData{
			ExternalId:           fmt.Sprintf("%d", input.ID),
			ExternalSourceEntity: "customers",
			ExternalSystem:       "shopify",
		},
		Name:        name,
		IsCustomer:  isCustomer,
		Whitelisted: true,
	}
	output.CreatedAtStr = input.CreatedAtStr

	return utils.ToJson(output)
}

func MapOrder(inputJson string) (string, error) {
	var input struct {
		Email        string `json:"email,omitempty"`
		FirstName    string `json:"first_name,omitempty"`
		LastName     string `json:"last_name,omitempty"`
		ID           int64  `json:"id,omitempty"`
		CreatedAtStr string `json:"created_at,omitempty"`
	}

	err := json.Unmarshal([]byte(inputJson), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %v", err)
	}

	var name string
	if input.FirstName != "" && input.LastName != "" {
		name = input.FirstName + " " + input.LastName
	} else {
		name = input.Email
	}

	output := entity.OrganizationData{
		BaseData: entity.BaseData{
			ExternalId:           fmt.Sprintf("%d", input.ID),
			ExternalSourceEntity: "customers",
			ExternalSystem:       "shopify",
		},
		Name:        name,
		IsCustomer:  true,
		Whitelisted: true,
	}
	output.CreatedAtStr = input.CreatedAtStr

	return utils.ToJson(output)
}
