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
	output.PhoneNumbers = []string{input.Phone}
	output.PhoneNumbers = append(output.PhoneNumbers, input.MobilePhone)

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
		output.PhoneNumbers = []string{input.Phone}
	}

	return utils.ToJson(output)
}

//func MapContact(inputJson string) (string, error) {
//	var input struct {
//		ID               string `json:"id,omitempty"`
//		CreatedAt        int64  `json:"created_at,omitempty"`
//		UpdatedAt        int64  `json:"updated_at,omitempty"`
//		Email            string `json:"email,omitempty"`
//		Phone            string `json:"phone,omitempty"`
//		Name             string `json:"name,omitempty"`
//		Role             string `json:"role,omitempty"` // not used yet in sync
//		Type             string `json:"type,omitempty"` // not used yet in sync
//		Avatar           string `json:"avatar,omitempty"`
//		CustomAttributes struct {
//			JobTitle string `json:"job_title,omitempty"`
//		} `json:"custom_attributes,omitempty"`
//		Location struct {
//			City          string `json:"city,omitempty"`
//			Type          string `json:"type,omitempty"`
//			Region        string `json:"region,omitempty"`
//			Country       string `json:"country,omitempty"`
//			CountryCodeA3 string `json:"country_code,omitempty"`
//			ContinentCode string `json:"continent_code,omitempty"`
//		} `json:"location,omitempty"`
//	}
//
//	err := json.Unmarshal([]byte(inputJson), &input)
//	if err != nil {
//		return "", fmt.Errorf("failed to parse input JSON: %v", err)
//	}
//	if input.Email == "" {
//		output := entity.BaseData{
//			Skip:       true,
//			SkipReason: "Missing email",
//		}
//		return utils.ToJson(output)
//	}
//	if input.ID == "" {
//		output := entity.BaseData{
//			Skip:       true,
//			SkipReason: "Missing id",
//		}
//		return utils.ToJson(output)
//	}
//	emailDomain := extractDomain(input.Email)
//	if emailDomain == "" {
//		output := entity.BaseData{
//			Skip:       true,
//			SkipReason: "Missing email domain",
//		}
//		return utils.ToJson(output)
//	}
//
//	output := entity.ContactData{
//		BaseData: entity.BaseData{
//			ExternalId: input.ID,
//		},
//		PhoneNumber:          input.Phone,
//		Name:                 input.Name,
//		Email:                input.Email,
//		ProfilePhotoUrl:      input.Avatar,
//		OrganizationRequired: true,
//	}
//	output.Organizations = append(output.Organizations, entity.ReferencedOrganization{
//		Domain:   emailDomain,
//		JobTitle: input.CustomAttributes.JobTitle,
//	})
//	if input.CreatedAt != 0 {
//		output.CreatedAtStr = tsStrToRFC3339(input.CreatedAt)
//	}
//	if input.UpdatedAt != 0 {
//		output.UpdatedAtStr = tsStrToRFC3339(input.UpdatedAt)
//	}
//
//	return utils.ToJson(output)
//}
//
//func extractDomain(email string) string {
//	parts := strings.Split(email, "@")
//	if len(parts) != 2 {
//		return ""
//	}
//	return parts[1]
//}
//
//func tsStrToRFC3339(timestamp int64) string {
//	t := time.Unix(timestamp, 0).UTC()
//	layout := "2006-01-02T15:04:05Z"
//	return t.Format(layout)
//}
//
//func MapInteractionEvent(inputJson string) (string, error) {
//	var data map[string]interface{}
//	err := json.Unmarshal([]byte(inputJson), &data)
//	if err != nil {
//		return "", fmt.Errorf("failed to parse input JSON: %v", err)
//	}
//	if _, exists := data["conversation_id"]; exists {
//		return mapInteractionEventFromConversationPart(inputJson)
//	} else {
//		return mapInteractionEventFromConversation(inputJson)
//	}
//}
//
//func mapInteractionEventFromConversation(inputJson string) (string, error) {
//	var input struct {
//		ID        string `json:"id,omitempty"`
//		State     string `json:"state,omitempty"`
//		Title     string `json:"title,omitempty"`
//		CreatedAt int64  `json:"created_at,omitempty"`
//		UpdatedAt int64  `json:"updated_at,omitempty"`
//		Source    struct {
//			ID     string `json:"id,omitempty"`
//			Body   string `json:"body,omitempty"`
//			Type   string `json:"type,omitempty"`
//			Author struct {
//				ID    string `json:"id,omitempty"`
//				Name  string `json:"name,omitempty"`
//				Type  string `json:"type,omitempty"`
//				Email string `json:"email,omitempty"`
//			} `json:"author,omitempty"`
//		} `json:"source,omitempty"`
//		Contacts struct {
//			Contacts []struct {
//				ID   string `json:"id,omitempty"`
//				Type string `json:"type,omitempty"`
//			} `json:"contacts,omitempty"`
//		} `json:"contacts,omitempty"`
//	}
//
//	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
//		return "", err
//	}
//
//	if input.ID == "" {
//		output := entity.BaseData{
//			Skip:       true,
//			SkipReason: "Missing id",
//		}
//		return utils.ToJson(output)
//	}
//
//	output := entity.InteractionEventData{
//		BaseData: entity.BaseData{
//			ExternalId:          input.ID,
//			ExternalSourceTable: utils.StringPtr("conversations"),
//			CreatedAtStr:        tsStrToRFC3339(input.CreatedAt),
//			UpdatedAtStr:        tsStrToRFC3339(input.UpdatedAt),
//		},
//		Channel:         "CHAT",
//		ContentType:     "text/html",
//		Content:         input.Source.Body,
//		Hide:            false,
//		ContactRequired: true,
//	}
//	if input.Source.Type == "email" {
//		output.Type = "EMAIL"
//	} else {
//		output.Type = "MESSAGE"
//	}
//
//	output.SessionDetails.Name = input.Title
//	output.SessionDetails.Channel = "CHAT"
//	output.SessionDetails.Type = "THREAD"
//	output.SessionDetails.CreatedAtStr = tsStrToRFC3339(input.CreatedAt)
//	output.SessionDetails.ExternalId = "session/" + input.ID
//	if input.State == "closed" {
//		output.SessionDetails.Status = "INACTIVE"
//	} else {
//		output.SessionDetails.Status = "ACTIVE"
//	}
//
//	if input.Source.Author.Type == "admin" || input.Source.Author.Type == "team" {
//		output.SentBy = entity.InteractionEventParticipant{
//			ReferencedUser: entity.ReferencedUser{
//				ExternalId: input.Source.Author.ID,
//			},
//		}
//	} else {
//		output.SentBy = entity.InteractionEventParticipant{
//			ReferencedContact: entity.ReferencedContact{
//				ExternalId: input.Source.Author.ID,
//			},
//		}
//	}
//
//	for _, contact := range input.Contacts.Contacts {
//		output.SentTo = append(output.SentTo, entity.InteractionEventParticipant{
//			ReferencedContact: entity.ReferencedContact{
//				ExternalId: contact.ID,
//			},
//		})
//	}
//
//	return utils.ToJson(output)
//}
//
//func mapInteractionEventFromConversationPart(inputJson string) (string, error) {
//	var input struct {
//		ID     string `json:"id,omitempty"`
//		Body   string `json:"body,omitempty"`
//		Type   string `json:"type,omitempty"`
//		Author struct {
//			ID    string `json:"id,omitempty"`
//			Name  string `json:"name,omitempty"`
//			Type  string `json:"type,omitempty"`
//			Email string `json:"email,omitempty"`
//		} `json:"author,omitempty"`
//		CreatedAt      int64  `json:"created_at,omitempty"`
//		UpdatedAt      int64  `json:"updated_at,omitempty"`
//		PartType       string `json:"part_type,omitempty"`
//		ConversationId string `json:"conversation_id,omitempty"`
//	}
//
//	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
//		return "", err
//	}
//
//	if input.ID == "" {
//		output := entity.BaseData{
//			Skip:       true,
//			SkipReason: "Missing id",
//		}
//		return utils.ToJson(output)
//	}
//	if input.ConversationId == "" {
//		output := entity.BaseData{
//			Skip:       true,
//			SkipReason: "Missing conversation id (interaction session)",
//		}
//		return utils.ToJson(output)
//	}
//	if input.Body == "" {
//		output := entity.BaseData{
//			Skip:       true,
//			SkipReason: "Missing conversation part body",
//		}
//		return utils.ToJson(output)
//	}
//
//	output := entity.InteractionEventData{
//		BaseData: entity.BaseData{
//			ExternalId:          input.ID,
//			ExternalSourceTable: utils.StringPtr("conversation_parts"),
//			CreatedAtStr:        tsStrToRFC3339(input.CreatedAt),
//			UpdatedAtStr:        tsStrToRFC3339(input.UpdatedAt),
//		},
//		Channel:         "CHAT",
//		Type:            "MESSAGE",
//		ContentType:     "text/html",
//		Content:         input.Body,
//		Hide:            true,
//		ContactRequired: false,
//		SessionRequired: true,
//	}
//
//	output.PartOfSession = entity.ReferencedInteractionSession{
//		"session/" + input.ConversationId,
//	}
//
//	if input.Author.Type == "admin" || input.Author.Type == "team" {
//		output.SentBy = entity.InteractionEventParticipant{
//			ReferencedUser: entity.ReferencedUser{
//				ExternalId: input.Author.ID,
//			},
//		}
//	} else {
//		output.SentBy = entity.InteractionEventParticipant{
//			ReferencedContact: entity.ReferencedContact{
//				ExternalId: input.Author.ID,
//			},
//		}
//	}
//
//	return utils.ToJson(output)
//}
