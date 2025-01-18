package postgres_entity

import "time"

type MagicLink struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`

	Email string `gorm:"column:email;type:varchar(255);NOT NULL" json:"email"`
	Code  string `gorm:"column:code;type:varchar(255);NOT NULL" json:"code"`
	Url   string `gorm:"column:url;type:varchar(255);NOT NULL" json:"url"`
}

func (MagicLink) TableName() string {
	return "magic_link"
}
