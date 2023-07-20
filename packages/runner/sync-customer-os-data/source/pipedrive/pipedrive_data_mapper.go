package pipedrive

import (
	"encoding/json"
	"fmt"
)

func MapUser(input string) (string, error) {
	var inputUser struct {
		ID        int64  `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Email     string `json:"email,omitempty"`
		Phone     string `json:"phone,omitempty"`
		CreatedAt string `json:"created,omitempty"`
		Modified  string `json:"modified,omitempty"`
	}

	if err := json.Unmarshal([]byte(input), &inputUser); err != nil {
		return "", err
	}

	outputUser := map[string]interface{}{
		"externalId":  fmt.Sprintf("%d", inputUser.ID),
		"name":        inputUser.Name,
		"email":       inputUser.Email,
		"phoneNumber": inputUser.Phone,
		"createdAt":   inputUser.CreatedAt,
		"updatedAt":   inputUser.Modified,
	}

	outputJson, err := json.Marshal(outputUser)
	if err != nil {
		return "", err
	}

	return string(outputJson), nil
}
