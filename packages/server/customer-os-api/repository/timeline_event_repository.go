package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type TimelineEventRepository interface {
	GetTimelineEventsForContact(ctx context.Context, tenant, contactId string, startingDate time.Time, size int, labels []string) ([]*dbtype.Node, error)
	GetTimelineEventsForOrganization(ctx context.Context, tenant, organizationId string, startingDate time.Time, size int, labels []string) ([]*dbtype.Node, error)
}

type timelineEventRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTimelineEventRepository(driver *neo4j.DriverWithContext) TimelineEventRepository {
	return &timelineEventRepository{
		driver: driver,
	}
}

func (r *timelineEventRepository) GetTimelineEventsForContact(ctx context.Context, tenant, contactId string, startingDate time.Time, size int, labels []string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	params := map[string]any{
		"tenant":       tenant,
		"contactId":    contactId,
		"startingDate": startingDate,
		"size":         size,
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	query := fmt.Sprintf("MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), "+
		" p = (c)-[*1..2]-(a:Action) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['HAS_ACTION','PARTICIPATES','SENT','SENT_TO','PART_OF','REQUESTED','NOTED'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" %s "+
		" RETURN distinct a ORDER BY coalesce(a.startedAt, a.createdAt) DESC LIMIT $size", filterByTypeCypherFragment)

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

func (r *timelineEventRepository) GetTimelineEventsForOrganization(ctx context.Context, tenant, contactId string, startingDate time.Time, size int, labels []string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	params := map[string]any{
		"tenant":       tenant,
		"contactId":    contactId,
		"startingDate": startingDate,
		"size":         size,
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	query := fmt.Sprintf("CALL { "+
		"MATCH (o:Organization {id:$contactId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), "+
		" (o)--(c:Contact), "+
		" p = (c)-[*1..2]-(a:Action) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['HAS_ACTION','PARTICIPATES','SENT','SENT_TO','PART_OF','REQUESTED','NOTED'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		"MATCH (o:Organization {id:$contactId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), "+
		" p = (o)-[*1]-(a:Action) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['NOTED'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" %s "+
		" return a as timelineEvent "+
		" } "+
		" RETURN distinct timelineEvent ORDER BY coalesce(timelineEvent.startedAt, timelineEvent.createdAt) DESC LIMIT $size",
		filterByTypeCypherFragment, filterByTypeCypherFragment)

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
