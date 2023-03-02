package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToTicket(entity *entity.Ticket) *model.Ticket {
	return &model.Ticket{
		ID:          entity.Id,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Subject:     &entity.Subject,
		Status:      &entity.Status,
		Priority:    &entity.Priority,
		Description: &entity.Description,
	}
}

func MapEntitiesToTickets(entities []*entity.Ticket) []*model.Ticket {
	var tags []*model.Ticket
	for _, ticketEntity := range entities {
		tags = append(tags, MapEntityToTicket(ticketEntity))
	}
	return tags
}
