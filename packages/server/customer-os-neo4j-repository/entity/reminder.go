package entity

import "time"

type ReminderEntity struct {
	Id            string     `json:"id"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Source        DataSource `json:"source"`
	SourceOfTruth DataSource `json:"sourceOfTruth"`
	AppSource     string     `json:"appSource"`
	Content       string     `json:"content"`
	DueDate       time.Time  `json:"dueDate"`
	Dismissed     bool       `json:"dismissed"`
}
