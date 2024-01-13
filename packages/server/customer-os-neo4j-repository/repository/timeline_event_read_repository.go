package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

var (
	relationshipsWithOrganization           = []string{"LOGGED", "NOTED", "REPORTED_BY", "SENT_TO", "SENT_BY", "ACTION_ON"}
	relationshipsWithOrganizationProperties = []string{"SENT_TO", "SENT_BY", "PART_OF", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"}
	relationshipsWithContact                = []string{"HAS_ACTION", "PARTICIPATES", "SENT_TO", "SENT_BY", "PART_OF", "REPORTED_BY", "NOTED", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"}
	relationshipsWithContactProperties      = []string{"SENT_TO", "SENT_BY", "PART_OF", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"}
)

type TimelineEventReadRepository interface {
	GetTimelineEvent(ctx context.Context, tenant, id string) (*dbtype.Node, error)
	CalculateAndGetLastTouchPoint(ctx context.Context, tenant, organizationId string) (*time.Time, string, error)
}

type timelineEventReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewTimelineEventReadRepository(driver *neo4j.DriverWithContext, database string) TimelineEventReadRepository {
	return &timelineEventReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *timelineEventReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *timelineEventReadRepository) GetTimelineEvent(ctx context.Context, tenant, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventReadRepository.GetTimelineEventsWithIds")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("id", id))

	cypher := fmt.Sprintf(`MATCH (a:TimelineEvent{id: $id}) WHERE a:TimelineEvent_%s RETURN a`, tenant)
	params := map[string]any{
		"id": id,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *timelineEventReadRepository) CalculateAndGetLastTouchPoint(ctx context.Context, tenant, organizationId string) (*time.Time, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventReadRepository.CalculateAndGetLastTouchPoint")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("organizationId", organizationId))

	params := map[string]any{
		"tenant":                                  tenant,
		"organizationId":                          organizationId,
		"nodeLabels":                              []string{neo4jutil.NodeLabelInteractionSession, neo4jutil.NodeLabelIssue, neo4jutil.NodeLabelInteractionEvent, neo4jutil.NodeLabelMeeting, neo4jutil.NodeLabelLogEntry},
		"excludeInteractionEventContentType":      []string{"x-openline-transcript-element"},
		"relationshipsWithOrganization":           relationshipsWithOrganization,
		"relationshipsWithContact":                relationshipsWithContact,
		"relationshipsWithContactProperties":      relationshipsWithContactProperties,
		"relationshipsWithOrganizationProperties": relationshipsWithOrganizationProperties,
		"now": utils.Now(),
	}

	cypher := `MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
		CALL { ` +
		// get all timeline events for the organization contacts
		` WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact), 
		p = (c)-[*1..2]-(a:TimelineEvent) 
		WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithContact)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 
		AND (NOT "InteractionEvent" in labels(a) or "InteractionEvent" in labels(a) AND NOT a.contentType IN $excludeInteractionEventContentType)
		AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
		RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		UNION ` +
		// get all timeline events directly for the organization
		` WITH o MATCH (o), 
		p = (o)-[*1]-(a:TimelineEvent) 
		WHERE all(r IN relationships(p) WHERE type(r) in $relationshipsWithOrganization)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 
		AND (NOT "InteractionEvent" in labels(a) or "InteractionEvent" in labels(a) AND NOT a.contentType IN $excludeInteractionEventContentType)
		AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
		RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		UNION ` +
		// get all timeline events for the organization contacts' emails and phone numbers
		` WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact)-[:HAS]->(e), 
		p = (e)-[*1..2]-(a:TimelineEvent) 
		WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) 
		AND all(r IN relationships(p) WHERE type(r) in $relationshipsWithContactProperties)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
		RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		UNION ` +
		// get all timeline events for the organization emails, phone numbers and job roles
		` WITH o MATCH (o)-[:HAS|ROLE_IN]-(e), 
		p = (e)-[*1..2]-(a:TimelineEvent) 
		WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e) OR 'JobRole' in labels(e)) 
		AND all(r IN relationships(p) WHERE type(r) in $relationshipsWithOrganizationProperties)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
	 	RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		} 
		RETURN coalesce(timelineEvent.startedAt, timelineEvent.createdAt), timelineEvent.id ORDER BY coalesce(timelineEvent.startedAt, timelineEvent.createdAt) DESC LIMIT 1`

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, "", err
	}

	if len(records.([]*neo4j.Record)) > 0 {
		// Try to assert the value to time.Time
		if t, ok := records.([]*neo4j.Record)[0].Values[0].(time.Time); ok {
			// If assertion is successful, proceed
			return utils.TimePtr(t), records.([]*neo4j.Record)[0].Values[1].(string), nil
		} else {
			err = errors.New(fmt.Sprintf("Value %v associated to timeline event id %s is not of type time.Time", records.([]*neo4j.Record)[0].Values[0], records.([]*neo4j.Record)[0].Values[1].(string)))
			tracing.TraceErr(span, err)
			return nil, "", nil
		}
	}
	return nil, "", nil
}
