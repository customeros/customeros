package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type TicketService interface {
	GetContactTickets(context.Context, string) ([]*entity.Ticket, error)
	GetTicketSummaryByStatusForOrganization(ctx context.Context, organizationId string) (map[string]int64, error)
	GetTicketSummaryByStatusForContact(ctx context.Context, contactId string) (map[string]int64, error)
}

type ticketService struct {
	repositories *repository.Repositories
}

func NewTicketService(repositories *repository.Repositories) TicketService {
	return &ticketService{
		repositories: repositories,
	}
}

func (s *ticketService) GetContactTickets(ctx context.Context, contactId string) ([]*entity.Ticket, error) {
	ticketDbNodes, err := s.repositories.TicketRepository.GetForContact(ctx, common.GetTenantFromContext(ctx), contactId)
	if err != nil {
		return nil, err
	}
	ticketEntities := make([]*entity.Ticket, 0)
	for _, dbNodePtr := range ticketDbNodes {
		ticketEntities = append(ticketEntities, s.mapDbNodeToTicket(dbNodePtr))
	}
	return ticketEntities, nil
}

func (s *ticketService) GetTicketSummaryByStatusForOrganization(ctx context.Context, organizationId string) (map[string]int64, error) {
	return s.repositories.TicketRepository.GetTicketCountByStatusForOrganization(ctx, common.GetTenantFromContext(ctx), organizationId)
}

func (s *ticketService) GetTicketSummaryByStatusForContact(ctx context.Context, contactId string) (map[string]int64, error) {
	return s.repositories.TicketRepository.GetTicketCountByStatusForContact(ctx, common.GetTenantFromContext(ctx), contactId)
}

func (s *ticketService) mapDbNodeToTicket(node *dbtype.Node) *entity.Ticket {
	props := utils.GetPropsFromNode(*node)
	ticket := entity.Ticket{
		Id:          utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:   utils.GetTimePropOrNow(props, "createdAt"),
		UpdatedAt:   utils.GetTimePropOrNow(props, "updatedAt"),
		Subject:     utils.GetStringPropOrEmpty(props, "subject"),
		Status:      utils.GetStringPropOrEmpty(props, "status"),
		Priority:    utils.GetStringPropOrEmpty(props, "priority"),
		Description: utils.GetStringPropOrEmpty(props, "description"),
	}
	return &ticket
}
