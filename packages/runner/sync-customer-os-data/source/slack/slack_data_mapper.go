package slack

import (
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const unknownUserName = "Unknown User"

func MapUser(inputJson string) (string, error) {
	var input struct {
		ID      string `json:"id,omitempty"`
		Profile struct {
			Email     string `json:"email,omitempty"`
			Phone     string `json:"phone,omitempty"`
			FirstName string `json:"first_name,omitempty"`
			LastName  string `json:"last_name,omitempty"`
			Name      string `json:"real_name_normalized,omitempty"`
			Image192  string `json:"image_192,omitempty"`
		} `json:"profile"`
	}

	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
		return "", err
	}

	output := model.Output{
		ExternalId:  input.ID,
		Email:       input.Profile.Email,
		PhoneNumber: input.Profile.Phone,
		FirstName:   input.Profile.FirstName,
		LastName:    input.Profile.LastName,
		Name:        input.Profile.Name,
	}
	if !strings.HasPrefix(input.Profile.Image192, "https://secure.gravatar.com") {
		output.ProfilePhotoUrl = input.Profile.Image192
	}

	outputJson, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJson), nil
}

func MapContact(inputJson string) (string, error) {
	var input struct {
		ID                     string `json:"id,omitempty"`
		Timezone               string `json:"tz,omitempty"`
		OpenlineOrganizationId string `json:"openline_organization_id,omitempty"`
		Profile                struct {
			Email     string `json:"email,omitempty"`
			Phone     string `json:"phone,omitempty"`
			FirstName string `json:"first_name,omitempty"`
			LastName  string `json:"last_name,omitempty"`
			Name      string `json:"real_name_normalized,omitempty"`
			Image192  string `json:"image_192,omitempty"`
		} `json:"profile"`
	}

	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
		return "", err
	}

	output := model.Output{
		ExternalId:             input.ID,
		Email:                  input.Profile.Email,
		PhoneNumber:            input.Profile.Phone,
		FirstName:              input.Profile.FirstName,
		LastName:               input.Profile.LastName,
		Name:                   input.Profile.Name,
		Timezone:               input.Timezone,
		OpenlineOrganizationId: input.OpenlineOrganizationId,
	}
	if !strings.HasPrefix(input.Profile.Image192, "https://secure.gravatar.com") {
		output.ProfilePhotoUrl = input.Profile.Image192
	}

	outputJson, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJson), nil
}

type OutputContent struct {
	Text   string `json:"text,omitempty"`
	Blocks []any  `json:"blocks,omitempty"`
}

func MapInteractionEvent(inputJson string) (string, error) {
	var input struct {
		Ts                     string            `json:"ts,omitempty"`
		ChannelId              string            `json:"channel_id,omitempty"`
		ChannelName            string            `json:"channel_name,omitempty"`
		Type                   string            `json:"type,omitempty"`
		SenderUser             string            `json:"user,omitempty"`
		Text                   string            `json:"text,omitempty"`
		UserIds                []string          `json:"channel_user_ids,omitempty"`
		UserNamesById          map[string]string `json:"channel_user_names,omitempty"`
		ThreadTs               string            `json:"thread_ts,omitempty"`
		OpenlineOrganizationId string            `json:"openline_organization_id,omitempty"`
		Blocks                 []any             `json:"blocks,omitempty"`
	}

	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
		return "", err
	}

	output := model.Output{
		ExternalId:  input.ChannelId + "/" + input.Ts,
		CreatedAt:   tsStrToRFC3339Nanos(input.Ts),
		ContentType: "application/json",
		Type:        "MESSAGE",
		Channel:     "SLACK",
	}
	outputContent := OutputContent{
		Text:   replaceUserMentionsInText(input.Text, input.UserNamesById),
		Blocks: addUserNameInBlocks(input.Blocks, input.UserNamesById),
	}
	outputContentJson, err := json.Marshal(outputContent)
	if err != nil {
		return "", err
	}
	output.Content = string(outputContentJson)

	output.SentBy = struct {
		OpenlineId                string `json:"openlineId,omitempty"`
		ExternalId                string `json:"externalId,omitempty"`
		ParticipantType           string `json:"participantType,omitempty"`
		RelationType              string `json:"relationType,omitempty"`
		ReplaceContactWithJobRole bool   `json:"replaceContactWithJobRole,omitempty"`
		OrganizationId            string `json:"organizationId,omitempty"`
	}{
		ExternalId:                input.SenderUser,
		ReplaceContactWithJobRole: true,
		OrganizationId:            input.OpenlineOrganizationId,
	}
	output.PartOfSession.Channel = "SLACK"
	output.PartOfSession.Type = "THREAD"
	output.PartOfSession.Status = "ACTIVE"
	output.PartOfSession.Name = input.ChannelName
	if input.ThreadTs != "" {
		if input.ThreadTs != input.Ts {
			output.Hide = true
		}
		output.PartOfSession.ExternalId = "session/" + input.ChannelId + "/" + input.ThreadTs
		output.PartOfSession.CreatedAt = tsStrToRFC3339Nanos(input.ThreadTs)
		output.PartOfSession.Identifier = input.ChannelId + "/" + input.ThreadTs
	} else {
		output.PartOfSession.ExternalId = "session/" + input.ChannelId + "/" + input.Ts
		output.PartOfSession.CreatedAt = tsStrToRFC3339Nanos(input.Ts)
		output.PartOfSession.Identifier = input.ChannelId + "/" + input.Ts
	}

	for _, user := range input.UserIds {
		if user != input.SenderUser {
			output.SentTo = append(output.SentTo,
				struct {
					OpenlineId                string `json:"openlineId,omitempty"`
					ExternalId                string `json:"externalId,omitempty"`
					ParticipantType           string `json:"participantType,omitempty"`
					RelationType              string `json:"relationType,omitempty"`
					ReplaceContactWithJobRole bool   `json:"replaceContactWithJobRole,omitempty"`
					OrganizationId            string `json:"organizationId,omitempty"`
				}{
					ExternalId:                user,
					ReplaceContactWithJobRole: true,
					OrganizationId:            input.OpenlineOrganizationId,
				})
		}
	}

	output.SentTo = append(output.SentTo,
		struct {
			OpenlineId                string `json:"openlineId,omitempty"`
			ExternalId                string `json:"externalId,omitempty"`
			ParticipantType           string `json:"participantType,omitempty"`
			RelationType              string `json:"relationType,omitempty"`
			ReplaceContactWithJobRole bool   `json:"replaceContactWithJobRole,omitempty"`
			OrganizationId            string `json:"organizationId,omitempty"`
		}{
			OpenlineId:      input.OpenlineOrganizationId,
			ParticipantType: "ORGANIZATION",
		})

	if input.Type != "message" {
		output.Skip = true
		output.SkipReason = "Not a message type. Type: " + input.Type
	}

	outputJson, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputJson), nil
}

func tsStrToRFC3339Nanos(ts string) string {
	parts := strings.Split(ts, ".")
	secs, _ := strconv.ParseInt(parts[0], 10, 64)
	millis, _ := strconv.ParseInt(parts[1], 10, 64)
	t := time.Unix(secs, millis*1000).UTC()
	layout := "2006-01-02T15:04:05.000000Z"
	return t.Format(layout)
}

func replaceUserMentionsInText(text string, userNames map[string]string) string {
	re := regexp.MustCompile("<@(U[A-Z0-9]+)>")
	replaced := re.ReplaceAllStringFunc(text, func(mention string) string {
		id := mention[2 : len(mention)-1]
		name, ok := userNames[id]
		if !ok || name == "" {
			return unknownUserName
		}
		return name
	})
	return replaced
}

func addUserNameInBlocks(blocks []any, userNamesById map[string]string) []any {
	for _, block := range blocks {
		blockMap, ok := block.(map[string]any)
		if !ok {
			continue
		}

		if elements, exists := blockMap["elements"]; exists {
			elementsSlice, ok := elements.([]any)
			if !ok {
				continue
			}

			for _, element := range elementsSlice {
				elementMap, ok := element.(map[string]any)
				if !ok {
					continue
				}

				if innerElements, exists := elementMap["elements"]; exists {
					innerElementsSlice, ok := innerElements.([]any)
					if !ok {
						continue
					}
					for _, innerElement := range innerElementsSlice {
						innerElementMap, ok := innerElement.(map[string]any)
						if !ok {
							continue
						}
						if innerElementMap["type"] == "user" {
							userID := innerElementMap["user_id"].(string)
							if userName, exists := userNamesById[userID]; exists {
								innerElementMap["user_name"] = userName
							} else {
								innerElementMap["user_name"] = unknownUserName
							}
						}
					}
				}
			}
		}
	}
	return blocks
}
