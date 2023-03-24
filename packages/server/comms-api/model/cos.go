package model

type MailReplyRequest struct {
	Username    string   `json:"username"`
	Content     string   `json:"content"`
	Channel     string   `json:"channel"`
	Source      string   `json:"source"`
	Direction   string   `json:"direction"`
	Destination []string `json:"destination"`
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

type InteractionEventParticipant struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	RawEmail       string `json:"rawEmail,omitempty"`
	FirstName      string `json:"firstName,omitempty"`
	RawPhoneNumber string `json:"rawPhoneNumber,omitempty"`
}

type AnalysisDescriptionInput struct {
	InteractionEventId   *string `json:"interactionEventId,omitempty"`
	InteractionSessionId *string `json:"interactionSessionId,omitempty"`
}

type InteractionEventCreateResponse struct {
	InteractionEventCreate struct {
		Id     string `json:"id"`
		SentBy []struct {
			Typename         string `json:"__typename"`
			EmailParticipant struct {
				Id       string `json:"id"`
				RawEmail string `json:"rawEmail"`
			} `json:"emailParticipant"`
			PhoneNumberParticipant struct {
				ID             string `json:"id"`
				RawPhoneNumber string `json:"rawPhoneNumber"`
			} `json:"phoneNumberParticipant"`
			Type string `json:"type"`
		} `json:"sentBy"`
		SentTo []struct {
			Typename         string `json:"__typename"`
			EmailParticipant struct {
				Id       string `json:"id"`
				RawEmail string `json:"rawEmail"`
			} `json:"emailParticipant"`
			PhoneNumberParticipant struct {
				ID             string `json:"id"`
				RawPhoneNumber string `json:"rawPhoneNumber"`
			} `json:"phoneNumberParticipant"`
			Type string `json:"type"`
		} `json:"sentTo"`
	} `json:"interactionEvent_Create"`
}

type InteractionEventGetResponse struct {
	InteractionEvent struct {
		EventIdentifier string `json:"eventIdentifier"`
		ChannelData     string
		SessionId       string `json:"sessionId"`
		Subject         string `json:"subject"`
	} `json:"interactionEvent_Create"`
}

type EmailChannelData struct {
	Subject   string   `json:"Subject"`
	InReplyTo []string `json:"InReplyTo"`
	Reference []string `json:"Reference"`
}
