package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type ActionRepository interface {
	GetContactActions(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, from, to time.Time, labels []string) ([]*dbtype.Node, error)
}

type actionRepository struct {
	driver *neo4j.DriverWithContext
}

func NewActionRepository(driver *neo4j.DriverWithContext) ActionRepository {
	return &actionRepository{
		driver: driver,
	}
}

func (r *actionRepository) GetContactActions(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, from, to time.Time, labels []string) ([]*dbtype.Node, error) {
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
		"from":      from,
		"to":        to,
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	query := fmt.Sprintf("MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), "+
		" p = (c)-[*1..2]-(a:Action) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['HAS_ACTION','PARTICIPATES','SENT_BY','SENT_TO','PART_OF','REQUESTED','NOTED'])"+
		" AND coalesce(a.startedAt, a.createdAt) >= datetime($from) AND coalesce(a.startedAt, a.createdAt) <= datetime($to) "+
		" %s "+
		" RETURN distinct a ORDER BY coalesce(a.startedAt, a.createdAt) DESC", filterByTypeCypherFragment)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	var actionDbNodes []*dbtype.Node
	for _, v := range records.([]*neo4j.Record) {
		if v.Values[0] != nil {
			actionDbNodes = append(actionDbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return actionDbNodes, err
}
