package entity

import "time"

type BaseEntity struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null" json:"updatedAt,omitempty"`
}
