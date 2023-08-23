package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type TimelineEventRepository interface {
	CalculateAndGetLastTouchpointInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, organizationId string) (*time.Time, string, error)
}

type timelineEventRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTimelineEventRepository(driver *neo4j.DriverWithContext) TimelineEventRepository {
	return &timelineEventRepository{
		driver: driver,
	}
}

func (r *timelineEventRepository) CalculateAndGetLastTouchpointInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, organizationId string) (*time.Time, string, error) {
	params := map[string]any{
		"tenant":                             tenant,
		"organizationId":                     organizationId,
		"nodeLabels":                         []string{entity.NodeLabel_InteractionSession, entity.NodeLabel_Issue, entity.NodeLabel_Conversation, entity.NodeLabel_InteractionEvent, entity.NodeLabel_Meeting},
		"excludeInteractionEventContentType": []string{"x-openline-transcript-element"},
		"contactRelationTypes":               []string{"HAS_ACTION", "PARTICIPATES", "SENT_TO", "SENT_BY", "PART_OF", "REPORTED_BY", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"},
		"organizationRelationTypes":          []string{"REPORTED_BY", "SENT_TO", "SENT_BY"},
		"emailAndPhoneRelationTypes":         []string{"SENT_TO", "SENT_BY", "PART_OF", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"},
		"now":                                utils.Now(),
	}

	query := `MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
		CALL { ` +
		// get all timeline events for the organization contatcs
		` WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact), 
		p = (c)-[*1..2]-(a:TimelineEvent) 
		WHERE all(r IN relationships(p) WHERE type(r) in $contactRelationTypes)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 
		AND (NOT "InteractionEvent" in labels(a) or "InteractionEvent" in labels(a) AND NOT a.contentType IN $excludeInteractionEventContentType)
		AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
		RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		UNION ` +
		// get all timeline events directly for the organization
		` WITH o MATCH (o), 
		p = (o)-[*1]-(a:TimelineEvent) 
		WHERE all(r IN relationships(p) WHERE type(r) in $organizationRelationTypes)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 
		AND (NOT "InteractionEvent" in labels(a) or "InteractionEvent" in labels(a) AND NOT a.contentType IN $excludeInteractionEventContentType)
		AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
		RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		UNION ` +
		// get all timeline events for the organization contacts' emails and phone numbers
		` WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact)-[:HAS]->(e), 
		p = (e)-[*1..2]-(a:TimelineEvent) 
		WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) 
		AND all(r IN relationships(p) WHERE type(r) in $emailAndPhoneRelationTypes)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
		RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		UNION ` +
		// get all timeline events for the organization emails, phone numbers and job roles
		` WITH o MATCH (o)-[:HAS|ROLE_IN]-(e), 
		p = (e)-[*1..2]-(a:TimelineEvent) 
		WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e) OR 'JobRole' in labels(e)) 
		AND all(r IN relationships(p) WHERE type(r) in $emailAndPhoneRelationTypes)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
	 	RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		} 
		RETURN coalesce(timelineEvent.startedAt, timelineEvent.createdAt), timelineEvent.id ORDER BY coalesce(timelineEvent.startedAt, timelineEvent.createdAt) DESC LIMIT 1`

	queryResult, err := tx.Run(ctx, query, params)
	if err != nil {
		return nil, "", err
	}

	records, err := queryResult.Collect(ctx)
	if err != nil {
		return nil, "", err
	}

	if len(records) > 0 {
		return utils.TimePtr(records[0].Values[0].(time.Time)), records[0].Values[1].(string), nil
	}

	return nil, "", nil
}
