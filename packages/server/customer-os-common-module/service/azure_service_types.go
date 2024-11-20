package service

import (
	"time"
)

type MicrosoftRawEmailsResponse struct {
	OdataNextLink string                      `json:"@odata.nextLink"`
	Value         []MicrosoftRawEmailResponse `json:"value"`
}

type MicrosoftRawEmailResponse struct {
	Id                      string    `json:"id"`
	SentDateTime            time.Time `json:"sentDateTime"`
	InternetMessageId       string    `json:"internetMessageId"`
	Subject                 string    `json:"subject"`
	ConversationId          string    `json:"conversationId"`
	ConversationIndex       string    `json:"conversationIndex"`
	InferenceClassification string    `json:"inferenceClassification"`
	Body                    struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	} `json:"body"`
	Sender struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"sender"`
	From struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"from"`
	ToRecipients []struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"toRecipients"`
	CcRecipients []struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"ccRecipients"`
	BccRecipients []struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"bccRecipients"`
	ReplyTo []interface{} `json:"replyTo"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type MailRequest struct {
	Subject string `json:"subject"`
	Body    struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	} `json:"body"`
	From struct {
		EmailAddress struct {
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"from"`
	ToRecipients  []Recipient `json:"toRecipients"`
	CcRecipients  []Recipient `json:"ccRecipients,omitempty"`
	BccRecipients []Recipient `json:"bccRecipients,omitempty"`
}

type Recipient struct {
	EmailAddress struct {
		Address string `json:"address"`
	} `json:"emailAddress"`
}

type ReplyRequest struct {
	Message MailRequest `json:"message"`
	Comment string      `json:"comment"`
}

type DraftResponse struct {
	Id string `json:"id"`
}

type SendDraftRequest struct {
	SaveToSentItems bool `json:"saveToSentItems"`
}
