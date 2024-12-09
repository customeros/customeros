package entity

import "time"

type ReminderEntity struct {
	Id             string    `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	UserId         string    `json:"userId"`
	OrganizationId string    `json:"organizationId"`
	Content        string    `json:"content"`
	DueDate        time.Time `json:"dueDate"`
	Dismissed      bool      `json:"dismissed"`
	Sent           bool      `json:"sent"`
}
