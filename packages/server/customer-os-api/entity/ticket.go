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

func (TicketEntity) Action() {
}

func (TicketEntity) ActionName() string {
	return NodeLabel_Ticket
}

func (TicketEntity) Labels(tenant string) []string {
	return []string{"Ticket", "Action", "Ticket_" + tenant, "Action_" + tenant}
}
