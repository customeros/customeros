package reminder

import "time"

type Reminder struct {
	OrganizationID string     `json:"organizationId"`
	UserID         string     `json:"userId"`
	Content        string     `json:"content"`
	DueDate        time.Time  `json:"dueDate"`
	Dismissed      bool       `json:"dismissed"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	SentAt         *time.Time `json:"sentAt"`
}
