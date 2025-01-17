package interfaces

import "context"

type NovuService interface {
	SendNotification(ctx context.Context, notification *NovuNotification) error
}

type NovuNotification struct {
	WorkflowId   string
	TemplateData map[string]string

	To      *NotifiableUser
	Subject string
	Payload map[string]interface{}
}

type NotifiableUser struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	SubscriberID string `json:"subscriberId"` // must be unique uuid for user
}
