package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToTicket(entity *entity.TicketEntity) *model.Ticket {
	return &model.Ticket{
		ID:          entity.Id,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Subject:     utils.StringPtr(entity.Subject),
		Status:      entity.Status,
		Priority:    utils.StringPtr(entity.Priority),
		Description: utils.StringPtr(entity.Description),
	}
}

func MapEntitiesToTickets(entities []*entity.TicketEntity) []*model.Ticket {
	var tags []*model.Ticket
	for _, ticketEntity := range entities {
		tags = append(tags, MapEntityToTicket(ticketEntity))
	}
	return tags
}
