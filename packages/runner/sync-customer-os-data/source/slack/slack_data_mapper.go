package slack

import (
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common/model"
	"strconv"
	"strings"
	"time"
)

func MapUser(inputJson string) (string, error) {
	var input struct {
		ID      string `json:"id,omitempty"`
		Profile struct {
			Email     string `json:"email,omitempty"`
			Phone     string `json:"phone,omitempty"`
			FirstName string `json:"first_name,omitempty"`
			LastName  string `json:"last_name,omitempty"`
			Name      string `json:"real_name,omitempty"`
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

func MapInteractionEvent(inputJson string) (string, error) {
	var input struct {
		Ts         string   `json:"ts,omitempty"`
		ChannelId  string   `json:"channel_id,omitempty"`
		Type       string   `json:"type,omitempty"`
		SenderUser string   `json:"user,omitempty"`
		Text       string   `json:"text,omitempty"`
		UserIds    []string `json:"channel_user_ids,omitempty"`
		ThreadTs   string   `json:"thread_ts,omitempty"`
	}

	if err := json.Unmarshal([]byte(inputJson), &input); err != nil {
		return "", err
	}

	output := model.Output{
		ExternalId:  input.ChannelId + "/" + input.Ts,
		CreatedAt:   TsStrToRFC3339Nanos(input.Ts),
		Content:     input.Text,
		ContentType: "text/plain",
		Type:        "MESSAGE",
		Channel:     "SLACK",
	}
	output.SentBy = struct {
		ExternalId      string `json:"externalId,omitempty"`
		ParticipantType string `json:"participantType,omitempty"`
		RelationType    string `json:"relationType,omitempty"`
	}{
		ExternalId:      input.SenderUser,
		ParticipantType: "",
		RelationType:    "",
	}
	output.PartOfSession.Channel = "SLACK"
	output.PartOfSession.Type = "THREAD"
	output.PartOfSession.Status = "ACTIVE"
	if input.ThreadTs != "" {
		output.Hide = true
		output.PartOfSession.ExternalId = "session/" + input.ChannelId + "/" + input.ThreadTs
		output.PartOfSession.CreatedAt = TsStrToRFC3339Nanos(input.ThreadTs)
		output.PartOfSession.Identifier = input.ChannelId + "/" + input.ThreadTs
	} else {
		output.PartOfSession.ExternalId = "session/" + input.ChannelId + "/" + input.Ts
		output.PartOfSession.CreatedAt = TsStrToRFC3339Nanos(input.Ts)
		output.PartOfSession.Identifier = input.ChannelId + "/" + input.Ts
	}

	for _, user := range input.UserIds {
		if user != input.SenderUser {
			output.SentTo = append(output.SentTo,
				struct {
					ExternalId      string `json:"externalId,omitempty"`
					ParticipantType string `json:"participantType,omitempty"`
					RelationType    string `json:"relationType,omitempty"`
				}{
					ExternalId:      user,
					ParticipantType: "",
					RelationType:    "",
				})
		}
	}

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

func TsStrToRFC3339Nanos(ts string) string {
	parts := strings.Split(ts, ".")
	secs, _ := strconv.ParseInt(parts[0], 10, 64)
	millis, _ := strconv.ParseInt(parts[1], 10, 64)
	t := time.Unix(secs, millis*1000).UTC()
	layout := "2006-01-02T15:04:05.000000Z"
	return t.Format(layout)
}
