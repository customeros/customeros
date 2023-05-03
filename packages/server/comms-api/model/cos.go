package model

import "time"

type MailReplyRequest struct {
	Username    string   `json:"username"`
	Content     string   `json:"content"`
	Channel     string   `json:"channel"`
	Source      string   `json:"source"`
	Direction   string   `json:"direction"`
	Destination []string `json:"destination"`
	Subject     *string  `json:"subject"`
	ReplyTo     *string  `json:"replyTo,omitempty"`
}

type MailFwdRequest struct {
	Sender     string `json:"sender"`
	RawMessage string `json:"rawMessage"`
	Subject    string `json:"subject"`
	ApiKey     string `json:"api-key"`
	Tenant     string `json:"X-Openline-TENANT"`
}

type InteractionEventParticipantInput struct {
	Email           *string `json:"email,omitempty"`
	PhoneNumber     *string `json:"phoneNumber,omitempty"`
	ContactID       *string `json:"contactID,omitempty"`
	UserID          *string `json:"userID,omitempty"`
	ParticipantType *string `json:"type,omitempty"`
}

type InteractionSessionParticipantInput struct {
	Email           *string `json:"email,omitempty"`
	PhoneNumber     *string `json:"phoneNumber,omitempty"`
	ContactID       *string `json:"contactID,omitempty"`
	UserID          *string `json:"userID,omitempty"`
	ParticipantType *string `json:"type,omitempty"`
}

type AnalysisDescriptionInput struct {
	InteractionEventId   *string `json:"interactionEventId,omitempty"`
	InteractionSessionId *string `json:"interactionSessionId,omitempty"`
	MeetingId            *string `json:"meetingId,omitempty"`
}

type InteractionEventCreate struct {
	Channel            string    `json:"channel"`
	Content            string    `json:"content"`
	ContentType        string    `json:"contentType"`
	CreatedAt          time.Time `json:"createdAt"`
	Id                 string    `json:"id"`
	InteractionSession struct {
		Name string `json:"name"`
	} `json:"interactionSession"`
	SentBy []struct {
		Typename         string `json:"__typename"`
		EmailParticipant struct {
			Contacts []interface{} `json:"contacts"`
			Id       string        `json:"id"`
			Email    string        `json:"email"`
		} `json:"emailParticipant"`
		Type interface{} `json:"type"`
	} `json:"sentBy"`
	SentTo []struct {
		Typename         string `json:"__typename"`
		EmailParticipant struct {
			Contacts []struct {
				Id string `json:"id"`
			} `json:"contacts"`
			Id    string `json:"id"`
			Email string `json:"email"`
		} `json:"emailParticipant"`
		Type string `json:"type"`
	} `json:"sentTo"`
}

type InteractionEventCreateResponse struct {
	InteractionEventCreate `json:"interactionEvent_Create"`
}

type TenantResponse struct {
	Tenant string `json:"tenant"`
}

type InteractionEventGetResponse struct {
	InteractionEvent struct {
		EventIdentifier    string `json:"eventIdentifier"`
		ChannelData        string
		InteractionSession struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"InteractionSession"`
	} `json:"interactionEvent"`
}

type EmailChannelData struct {
	Subject   string   `json:"Subject"`
	InReplyTo []string `json:"InReplyTo"`
	Reference []string `json:"Reference"`
}
