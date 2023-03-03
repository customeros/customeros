package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type TicketRepository interface {
	GetForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error)
	GetTicketCountByStatusForContact(ctx context.Context, tenant, contactId string) (map[string]int64, error)
	GetTicketCountByStatusForOrganization(ctx context.Context, tenant, organizationId string) (map[string]int64, error)
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

func (r *ticketRepository) GetTicketCountByStatusForContact(ctx context.Context, tenant, contactId string) (map[string]int64, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(:Contact {id:$contactId})--(tt:Ticket)
			WITH DISTINCT tt
			RETURN tt.status AS status, COUNT(tt) AS count`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	output := make(map[string]int64)
	for _, v := range result.([]*neo4j.Record) {
		status := ""
		if v.Values[0] != nil {
			status = v.Values[0].(string)
		}
		output[status] = v.Values[1].(int64)
	}
	return output, err
}

func (r *ticketRepository) GetTicketCountByStatusForOrganization(ctx context.Context, tenant, organizationId string) (map[string]int64, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$organizationId})--(:Contact)--(tt:Ticket)
			WITH DISTINCT tt
			RETURN tt.status AS status, COUNT(tt) AS count`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	output := make(map[string]int64)
	for _, v := range result.([]*neo4j.Record) {
		status := ""
		if v.Values[0] != nil {
			status = v.Values[0].(string)
		}
		output[status] = v.Values[1].(int64)
	}
	return output, err
}
