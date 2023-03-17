package entity

import "time"

type TicketEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Subject     string
	Status      string
	Priority    string
	Description string
}

func (TicketEntity) IsTimelineEvent() {
}

func (TicketEntity) TimelineEventLabel() string {
	return NodeLabel_Ticket
}

func (TicketEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Ticket,
		NodeLabel_Ticket + "_" + tenant,
		NodeLabel_Action,
		NodeLabel_Action + "_" + tenant,
	}
}
