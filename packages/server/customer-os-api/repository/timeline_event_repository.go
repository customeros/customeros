package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type TimelineEventRepository interface {
	GetTimelineEventsForContact(ctx context.Context, tenant, contactId string, startingDate time.Time, size int, labels []string) ([]*dbtype.Node, error)
	GetTimelineEventsForOrganization(ctx context.Context, tenant, organizationId string, startingDate time.Time, size int, labels []string) ([]*dbtype.Node, error)
	GetTimelineEventsTotalCountForContact(ctx context.Context, tenant string, id string, labels []string) (int64, error)
	GetTimelineEventsTotalCountForOrganization(ctx context.Context, tenant string, id string, labels []string) (int64, error)
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
	query := fmt.Sprintf("MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
		" CALL {"+
		// get all timeline events for the contact
		" WITH c MATCH (c), "+
		" p = (c)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['HAS_ACTION','PARTICIPATES','SENT_TO','SENT_BY','PART_OF','REQUESTED','NOTED'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the contact's emails and phone numbers
		" WITH c MATCH (c)-[:HAS]->(e),"+
		" p = (e)-[*1]-(a:TimelineEvent) "+
		" WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) "+
		" AND all(r IN relationships(p) WHERE type(r) in ['SENT_TO','SENT_BY'])"+
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

func (r *timelineEventRepository) GetTimelineEventsTotalCountForContact(ctx context.Context, tenant string, contactId string, labels []string) (int64, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	query := fmt.Sprintf("MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
		" CALL {"+
		// get all timeline events for the contact
		" WITH c MATCH (c), "+
		" p = (c)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['HAS_ACTION','PARTICIPATES','SENT_TO','SENT_BY','PART_OF','REQUESTED','NOTED']) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the contact's emails and phone numbers
		" WITH c MATCH (c)-[:HAS]->(e),"+
		" p = (e)-[*1]-(a:TimelineEvent) "+
		" WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) "+
		" AND all(r IN relationships(p) WHERE type(r) in ['SENT_TO','SENT_BY'])"+
		" %s "+
		" return a as timelineEvent "+
		" } "+
		" RETURN count(distinct timelineEvent)",
		filterByTypeCypherFragment, filterByTypeCypherFragment)

	record, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Single(ctx)
	})
	if err != nil {
		return int64(0), err
	}
	return record.(*db.Record).Values[0].(int64), nil
}

func (r *timelineEventRepository) GetTimelineEventsForOrganization(ctx context.Context, tenant, organizationId string, startingDate time.Time, size int, labels []string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"startingDate":   startingDate,
		"size":           size,
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	query := fmt.Sprintf("MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
		" CALL { "+
		// get all timeline events for the organization contatcs
		" WITH o MATCH (o)--(c:Contact), "+
		" p = (c)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['HAS_ACTION','PARTICIPATES','SENT_TO','SENT_BY','PART_OF','REQUESTED','NOTED'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events directly for the organization
		" WITH o MATCH (o), "+
		" p = (o)-[*1]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['NOTED'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the organization contacts' emails and phone numbers
		" WITH o MATCH (o)--(c:Contact)-[:HAS]->(e), "+
		" p = (e)-[*1]-(a:TimelineEvent) "+
		" WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) "+
		" AND all(r IN relationships(p) WHERE type(r) in ['SENT_TO','SENT_BY'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the organization emails and phone numbers
		" WITH o MATCH (o)-[:HAS]->(e), "+
		" p = (e)-[*1]-(a:TimelineEvent) "+
		" WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) "+
		" AND all(r IN relationships(p) WHERE type(r) in ['SENT_TO','SENT_BY'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" %s "+
		" return a as timelineEvent "+
		" } "+
		" RETURN distinct timelineEvent ORDER BY coalesce(timelineEvent.startedAt, timelineEvent.createdAt) DESC LIMIT $size",
		filterByTypeCypherFragment, filterByTypeCypherFragment, filterByTypeCypherFragment, filterByTypeCypherFragment)

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

func (r *timelineEventRepository) GetTimelineEventsTotalCountForOrganization(ctx context.Context, tenant string, organizationId string, labels []string) (int64, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	query := fmt.Sprintf("MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
		" CALL { "+
		// get all timeline events for the organization' contatcs
		" WITH o MATCH (o)--(c:Contact), "+
		" p = (c)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['HAS_ACTION','PARTICIPATES','SENT_TO','SENT_BY','PART_OF','REQUESTED','NOTED'])"+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events directly for the organization
		" WITH o MATCH (o), "+
		" p = (o)-[*1]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['NOTED'])"+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the organization contacts' emails and phone numbers
		" WITH o MATCH (o)--(c:Contact)-[:HAS]->(e), "+
		" p = (e)-[*1]-(a:TimelineEvent) "+
		" WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) "+
		" AND all(r IN relationships(p) WHERE type(r) in ['SENT_TO','SENT_BY'])"+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the organization emails and phone numbers
		" WITH o MATCH (o)-[:HAS]->(e), "+
		" p = (e)-[*1]-(a:TimelineEvent) "+
		" WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) "+
		" AND all(r IN relationships(p) WHERE type(r) in ['SENT_TO','SENT_BY'])"+
		" %s "+
		" return a as timelineEvent "+
		" } "+
		" RETURN count(distinct timelineEvent)",
		filterByTypeCypherFragment, filterByTypeCypherFragment, filterByTypeCypherFragment, filterByTypeCypherFragment)

	record, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Single(ctx)
	})
	if err != nil {
		return int64(0), err
	}
	return record.(*db.Record).Values[0].(int64), nil
}
