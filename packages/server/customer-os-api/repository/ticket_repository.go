package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type TicketRepository interface {
	GetForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error)
}

type ticketRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTicketRepository(driver *neo4j.DriverWithContext) TicketRepository {
	return &ticketRepository{
		driver: driver,
	}
}

func (r *ticketRepository) GetForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:TICKET_BELONGS_TO_TENANT]-(tt:Ticket)<-[:REQUESTED]-(c:Contact {id:$contactId})
			RETURN tt ORDER BY tt.createdAt DESC`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsNodePtrs(ctx, queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}
