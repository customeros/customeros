package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
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

type TimelineEventRepository interface {
	GetTimelineEvent(ctx context.Context, tenant string, id string) (*dbtype.Node, error)
	CalculateAndGetLastTouchpoint(ctx context.Context, tenant string, organizationId string) (*time.Time, string, error)
}

type timelineEventRepository struct {
	driver *neo4j.DriverWithContext
	log    logger.Logger
}

func NewTimelineEventRepository(driver *neo4j.DriverWithContext, log logger.Logger) TimelineEventRepository {
	return &timelineEventRepository{
		driver: driver,
		log:    log,
	}
}

func (r *timelineEventRepository) GetTimelineEvent(ctx context.Context, tenant string, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventRepository.GetTimelineEventsWithIds")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)

	cypher := fmt.Sprintf(`MATCH (a:TimelineEvent{id: $id}) WHERE a:TimelineEvent_%s RETURN a`, tenant)
	params := map[string]any{
		"id": id,
	}
	span.LogFields(log.String("query", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	return records.(*dbtype.Node), err
}

func (r *timelineEventRepository) CalculateAndGetLastTouchpoint(ctx context.Context, tenant string, organizationId string) (*time.Time, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventRepository.CalculateAndGetLastTouchpoint")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId))

	params := map[string]any{
		"tenant":                                  tenant,
		"organizationId":                          organizationId,
		"nodeLabels":                              []string{entity.NodeLabel_InteractionSession, entity.NodeLabel_Issue, entity.NodeLabel_InteractionEvent, entity.NodeLabel_Meeting, entity.NodeLabel_LogEntry},
		"excludeInteractionEventContentType":      []string{"x-openline-transcript-element"},
		"relationshipsWithOrganization":           relationshipsWithOrganization,
		"relationshipsWithContact":                relationshipsWithContact,
		"relationshipsWithContactProperties":      relationshipsWithContactProperties,
		"relationshipsWithOrganizationProperties": relationshipsWithOrganizationProperties,
		"now": utils.Now(),
	}

	query := `MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
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

	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
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
			r.log.Warnf("Failed to get last touchpoint. Skip setting it to organization: %v", err.Error())
			return nil, "", nil
		}
	}
	return nil, "", nil
}
