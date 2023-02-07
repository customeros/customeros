package entity

import (
	"fmt"
	"time"
)

type EmailEntity struct {
	Id            string
	Email         string
	Label         string
	Primary       bool
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (email EmailEntity) ToString() string {
	return fmt.Sprintf("id: %s\nemail: %s\nlabel: %s", email.Id, email.Email, email.Label)
}

type EmailEntities []EmailEntity

func (email EmailEntity) Labels(tenant string) []string {
	return []string{"Email", "Email_" + tenant}
}
