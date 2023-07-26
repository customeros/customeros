package entity

import "time"

type EmailMessageData struct {
	Html      string
	Text      string
	Subject   string
	CreatedAt time.Time

	ContactsExternalIds []string
	UserExternalId      string
	EmailThreadId       string
	ExternalId          string
	ExternalSystem      string

	FromEmail string
	ToEmail   []string
	CcEmail   []string
	BccEmail  []string

	Direction Direction

	Channel     string
	ChannelData *string

	FromFirstName string
	FromLastName  string
}
