package model

import "time"

type Reminder struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	UserID         string    `json:"userId"`
	Content        string    `json:"content"`
	DueDate        time.Time `json:"dueDate"`
	Dismissed      bool      `json:"dismissed"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
