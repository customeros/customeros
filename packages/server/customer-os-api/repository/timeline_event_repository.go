package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

var (
	relationshipsWithOrganization           = []string{"LOGGED", "REPORTED_BY", "SENT_TO", "SENT_BY", "ACTION_ON"}
	relationshipsWithOrganizationProperties = []string{"SENT_TO", "SENT_BY", "PART_OF", "DESCRIBES", "ATTENDED_BY", "CREATED_BY", "ACTION_ON"}
	relationshipsWithContact                = []string{"HAS_ACTION", "PARTICIPATES", "SENT_TO", "SENT_BY", "PART_OF", "REPORTED_BY", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"}
	relationshipsWithContactProperties      = []string{"SENT_TO", "SENT_BY", "PART_OF", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"}
)

type TimelineEventRepository interface {
	GetTimelineEventsForContact(ctx context.Context, tenant, contactId string, startingDate time.Time, size int, labels []string) ([]*dbtype.Node, error)
	GetTimelineEventsForOrganization(ctx context.Context, tenant, organizationId string, startingDate time.Time, size int, labels []string) ([]*dbtype.Node, error)
	GetTimelineEventsTotalCountForContact(ctx context.Context, tenant string, id string, labels []string) (int64, error)
	GetTimelineEventsTotalCountForOrganization(ctx context.Context, tenant string, id string, labels []string) (int64, error)
	GetTimelineEventsWithIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventRepository.GetTimelineEventsForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("startingDate", startingDate.String()), log.Int("size", size))

	params := map[string]any{
		"tenant":       tenant,
		"contactId":    contactId,
		"startingDate": startingDate,
		"size":         size,
		"skipDeleted":  "deleted",
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	cypher := fmt.Sprintf("MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
		" CALL {"+
		// get all timeline events for the contact
		" WITH c MATCH (c), "+
		" p = (c)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['WORKS_AS','HAS_ACTION','PARTICIPATES','SENT_TO','SENT_BY','PART_OF','REPORTED_BY', 'DESCRIBES', 'ATTENDED_BY', 'CREATED_BY'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" AND (a.status IS NULL OR a.status <> $skipDeleted) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the contact's emails and phone numbers
		" WITH c MATCH (c)-[:HAS]->(e),"+
		" p = (e)-[*1..2]-(a:TimelineEvent) "+
		" WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) "+
		" AND all(r IN relationships(p) WHERE type(r) in ['SENT_TO','SENT_BY', 'PART_OF', 'DESCRIBES', 'ATTENDED_BY', 'CREATED_BY'])"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" %s "+
		" return a as timelineEvent "+
		" } "+
		" RETURN distinct timelineEvent ORDER BY coalesce(timelineEvent.startedAt, timelineEvent.createdAt) DESC LIMIT $size",
		filterByTypeCypherFragment, filterByTypeCypherFragment)

	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventRepository.GetTimelineEventsTotalCountForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	params := map[string]any{
		"tenant":      tenant,
		"contactId":   contactId,
		"skipDeleted": "deleted",
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	cypher := fmt.Sprintf("MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
		" CALL {"+
		// get all timeline events for the contact
		" WITH c MATCH (c), "+
		" p = (c)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in ['WORKS_AS','HAS_ACTION','PARTICIPATES','SENT_TO','SENT_BY','PART_OF','REPORTED_BY', 'DESCRIBES', 'ATTENDED_BY', 'CREATED_BY']) "+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" AND (a.status IS NULL OR a.status <> $skipDeleted) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the contact's emails and phone numbers
		" WITH c MATCH (c)-[:HAS]->(e),"+
		" p = (e)-[*1..2]-(a:TimelineEvent) "+
		" WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) "+
		" AND all(r IN relationships(p) WHERE type(r) in ['SENT_TO','SENT_BY', 'PART_OF', 'DESCRIBES', 'ATTENDED_BY', 'CREATED_BY'])"+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" %s "+
		" return a as timelineEvent "+
		" } "+
		" RETURN count(distinct timelineEvent)",
		filterByTypeCypherFragment, filterByTypeCypherFragment)

	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	record, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventRepository.GetTimelineEventsForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.String("startingDate", startingDate.String()), log.Int("size", size))

	params := map[string]any{
		"tenant":                        tenant,
		"organizationId":                organizationId,
		"startingDate":                  startingDate,
		"size":                          size,
		"relationshipsWithOrganization": relationshipsWithOrganization,
		"relationshipsWithOrganizationProperties": relationshipsWithOrganizationProperties,
		"relationshipsWithContact":                relationshipsWithContact,
		"relationshipsWithContactProperties":      relationshipsWithContactProperties,
		"skipStatus":                              "deleted",
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	cypher := fmt.Sprintf("MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
		" CALL { "+
		// get all timeline events for the organization contacts
		" WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact), "+
		" p = (c)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithContact)"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" AND (a.status IS NULL OR a.status <> $skipStatus) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events directly for the organization
		" WITH o MATCH (o), "+
		" p = (o)-[*1]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithOrganization)"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" AND (a.status IS NULL OR a.status <> $skipStatus) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the organization contacts' emails and phone numbers
		" WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact)-[:HAS]->(e:Email|PhoneNumber), "+
		" p = (e)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithContactProperties)"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the organization emails, phone numbers, job roles or contracts
		" WITH o MATCH (o)-[:HAS|ROLE_IN|HAS_CONTRACT]-(e:Email|PhoneNumber|JobRole|Contract), "+
		" p = (e)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithOrganizationProperties)"+
		" AND coalesce(a.startedAt, a.createdAt) < datetime($startingDate) "+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" %s "+
		" return a as timelineEvent "+
		" } "+
		" RETURN distinct timelineEvent ORDER BY coalesce(timelineEvent.startedAt, timelineEvent.createdAt) DESC LIMIT $size",
		filterByTypeCypherFragment, filterByTypeCypherFragment, filterByTypeCypherFragment, filterByTypeCypherFragment)

	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventRepository.GetTimelineEventsTotalCountForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	params := map[string]any{
		"tenant":                        tenant,
		"organizationId":                organizationId,
		"relationshipsWithOrganization": relationshipsWithOrganization,
		"relationshipsWithOrganizationProperties": relationshipsWithOrganizationProperties,
		"relationshipsWithContact":                relationshipsWithContact,
		"relationshipsWithContactProperties":      relationshipsWithContactProperties,
		"skipStatus":                              "deleted",
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	cypher := fmt.Sprintf("MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
		" CALL { "+
		// get all timeline events for the organization' contacts
		" WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact), "+
		" p = (c)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithContact)"+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" AND (a.status IS NULL OR a.status <> $skipStatus) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events directly for the organization
		" WITH o MATCH (o), "+
		" p = (o)-[*1]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithOrganization)"+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" AND (a.status IS NULL OR a.status <> $skipStatus) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the organization contacts' emails and phone numbers
		" WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact)-[:HAS]->(e:Email|PhoneNumber), "+
		" p = (e)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithContactProperties)"+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" %s "+
		" return a as timelineEvent "+
		" UNION "+
		// get all timeline events for the organization emails, phone numbers job roles and contract
		" WITH o MATCH (o)-[:HAS|ROLE_IN|HAS_CONTRACT]-(e:Email|PhoneNumber|JobRole|Contract), "+
		" p = (e)-[*1..2]-(a:TimelineEvent) "+
		" WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithOrganizationProperties)"+
		" AND (a.hide IS NULL OR a.hide = false) "+
		" %s "+
		" return a as timelineEvent "+
		" } "+
		" RETURN count(distinct timelineEvent)",
		filterByTypeCypherFragment, filterByTypeCypherFragment, filterByTypeCypherFragment, filterByTypeCypherFragment)

	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	record, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
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

func (r *timelineEventRepository) GetTimelineEventsWithIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventRepository.GetTimelineEventsWithIds")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (a:TimelineEvent) WHERE a.id in $ids AND a:TimelineEvent_%s RETURN a`, tenant)
	params := map[string]any{
		"ids": ids,
	}
	span.LogFields(log.String("query", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	return records.([]*dbtype.Node), err
}
