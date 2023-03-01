package entity

type TicketField struct {
	Id    int64  `gorm:"column:id"`
	Type  string `gorm:"column:type"`
	Title string `gorm:"column:title"`
}

type TicketFields []TicketField

func (TicketField) TableName() string {
	return "ticket_fields"
}
