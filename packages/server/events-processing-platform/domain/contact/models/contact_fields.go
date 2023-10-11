package models

import "time"

type ContactDataFields struct {
	FirstName       string
	LastName        string
	Name            string
	Prefix          string
	Description     string
	Timezone        string
	ProfilePhotoUrl string
}

type JobRoleFields struct {
	JobTitle    string     `json:"jobTitle"`
	Description string     `json:"description"`
	Primary     bool       `json:"primary"`
	StartedAt   *time.Time `json:"startedAt"`
	EndedAt     *time.Time `json:"endedAt"`
}
