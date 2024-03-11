package shopify

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
	"time"
)

func MapOrganization(inputJson string) (string, error) {
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

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func tsStrToRFC3339(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	layout := "2013-06-27T08:48:27-04:00"
	return t.Format(layout)
}
