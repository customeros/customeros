package entity

import "time"

type EmailMessageData struct {
	Html      string
	Subject   string
	CreatedAt time.Time

	ContactsExternalIds []string
	UserExternalId      string
	EmailMessageId      string
	EmailThreadId       string
	ExternalId          string
	ExternalSystem      string

	FromEmail string
	ToEmail   []string
	CcEmail   []string
	BccEmail  []string

	Direction Direction

	FromFirstName string
	FromLastName  string
}
