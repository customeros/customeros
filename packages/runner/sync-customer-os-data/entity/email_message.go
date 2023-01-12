package entity

import "time"

type EmailMessageData struct {
	Id                  string
	Html                string
	CreatedAt           time.Time
	ContactsExternalIds []string
	ExternalId          string
	ExternalSystem      string
}
