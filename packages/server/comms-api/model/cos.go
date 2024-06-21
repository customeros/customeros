package model

import "time"

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

type GetUserByEmailResponse struct {
	UserByEmail struct {
		ID        string  `json:"id"`
		FirstName *string `json:"firstName"`
		LastName  *string `json:"lastName"`
		Name      *string `json:"name"`
	} `json:"user_ByEmail"`
}
