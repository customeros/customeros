package entity

import (
	"time"
)

type TicketComment struct {
	Id                          int64     `gorm:"column:id"`
	AirbyteAbId                 string    `gorm:"column:_airbyte_ab_id"`
	AirbyteTicketCommentsHashid string    `gorm:"column:_airbyte_ticket_comments_hashid"`
	CreateDate                  time.Time `gorm:"column:created_at"`
	TicketId                    int64     `gorm:"column:ticket_id"`
	AuthorId                    int64     `gorm:"column:author_id"`
	HtmlBody                    string    `gorm:"column:html_body"`
	PlainBody                   string    `gorm:"column:plain_body"`
	Body                        string    `gorm:"column:body"`
	Public                      bool      `gorm:"column:public"`
}

type TicketComments []TicketComment

func (TicketComment) TableName() string {
	return "ticket_comments"
}
