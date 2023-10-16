package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type OrganizationRepository interface {
	GetMatchedOrganizationId(ctx context.Context, tenant, externalSystem, externalId string, domains []string) (string, error)
	CalculateAndGetLastTouchpoint(ctx context.Context, tenant string, organizationId string) (*time.Time, string, error)
	UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error
	GetOrganizationIdsForContact(ctx context.Context, tenant, contactId string) ([]string, error)
	GetOrganizationIdsForContactByExternalId(ctx context.Context, tenant, contactExternalId, externalSystem string) ([]string, error)
	GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error)
	GetAllDomainLinksCrossTenantsNotSynced(ctx context.Context, size int) ([]*neo4j.Record, error)
	GetOrganizationIdById(ctx context.Context, tenant, id string) (string, error)
	GetOrganizationIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	GetOrganizationIdByDomain(ctx context.Context, tenant, domain string) (string, error)
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
	log    logger.Logger
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext, log logger.Logger) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
		log:    log,
	}
}

func (r *organizationRepository) GetMatchedOrganizationId(ctx context.Context, tenant, externalSystem, externalId string, domains []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetMatchedOrganizationId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o1:Organization)-[:IS_LINKED_WITH {externalId:$organizationExternalId}]->(e)
				OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o2:Organization)-[:HAS_DOMAIN]->(d:Domain)
					WHERE d.domain in $domains
				with coalesce(o1, o2) as organization
				where organization is not null
				return organization.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":                 tenant,
				"externalSystem":         externalSystem,
				"organizationExternalId": externalId,
				"domains":                domains,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	orgIDs := dbRecords.([]*db.Record)
	if len(orgIDs) == 1 {
		return orgIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *organizationRepository) CalculateAndGetLastTouchpoint(ctx context.Context, tenant string, organizationId string) (*time.Time, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.CalculateAndGetLastTouchpoint")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	params := map[string]any{
		"tenant":                             tenant,
		"now":                                utils.Now(),
		"organizationId":                     organizationId,
		"nodeLabels":                         []string{"InteractionSession", "Issue", "Conversation", "InteractionEvent", "Meeting"},
		"excludeInteractionEventContentType": []string{"x-openline-transcript-element"},
		"contactRelationTypes":               []string{"HAS_ACTION", "PARTICIPATES", "SENT_TO", "SENT_BY", "PART_OF", "REPORTED_BY", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"},
		"organizationRelationTypes":          []string{"REPORTED_BY", "SENT_TO", "SENT_BY"},
		"emailAndPhoneRelationTypes":         []string{"SENT_TO", "SENT_BY", "PART_OF", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"},
	}

	query := `MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
		CALL { ` +
		// get all timeline events for the organization contacts
		` WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact), 
		p = (c)-[*1..2]-(a:TimelineEvent) 
		WHERE all(r IN relationships(p) WHERE type(r) in $contactRelationTypes)
		AND (NOT "InteractionEvent" in labels(a) or "InteractionEvent" in labels(a) AND NOT a.contentType IN $excludeInteractionEventContentType)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 
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
		RETURN coalesce(timelineEvent.startedAt, timelineEvent.updatedAt, timelineEvent.createdAt), timelineEvent.id ORDER BY coalesce(timelineEvent.startedAt, timelineEvent.updatedAt, timelineEvent.createdAt) DESC LIMIT 1`

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

func (r *organizationRepository) UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateLastTouchpoint")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.lastTouchpointAt=$touchpointAt, org.lastTouchpointId=$touchpointId`

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
				"touchpointAt":   touchpointAt,
				"touchpointId":   touchpointId,
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) GetOrganizationIdsForContact(ctx context.Context, tenant, contactId string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdsForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)--(:JobRole)--(:Contact {id:$contactId})
		RETURN org.id`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return []string{}, err
	}
	orgIDs := make([]string, 0)
	for _, v := range dbRecords.([]*db.Record) {
		orgIDs = append(orgIDs, v.Values[0].(string))
	}
	return orgIDs, nil
}

func (r *organizationRepository) GetOrganizationIdsForContactByExternalId(ctx context.Context, tenant, contactExternalId, externalSystem string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdsForContactByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$contactExternalId}]-(c:Contact)--(:JobRole)--(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
		RETURN org.id`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":            tenant,
				"externalSystem":    externalSystem,
				"contactExternalId": contactExternalId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return []string{}, err
	}
	orgIDs := make([]string, 0)
	for _, v := range dbRecords.([]*db.Record) {
		orgIDs = append(orgIDs, v.Values[0].(string))
	}
	return orgIDs, nil
}

func (r *organizationRepository) GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllCrossTenantsNotSynced")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant)
 			WHERE (org.syncedWithEventStore is null or org.syncedWithEventStore=false)
			RETURN org, t.name limit $size`,
			map[string]any{
				"size": size,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationRepository) GetAllDomainLinksCrossTenantsNotSynced(ctx context.Context, size int) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllCrossTenantsNotSynced")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[rel:HAS_DOMAIN]->(d:Domain)
 			WHERE (rel.syncedWithEventStore is null or rel.syncedWithEventStore=false) AND org.syncedWithEventStore=true AND d.domain <> "" 
			RETURN org.id, t.name, d.domain limit $size`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"size": size,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*neo4j.Record), err
}

func (r *organizationRepository) GetOrganizationIdById(ctx context.Context, tenant, id string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				return org.id order by org.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":         tenant,
			"organizationId": id,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *organizationRepository) GetOrganizationIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				return org.id order by org.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":           tenant,
			"externalId":       externalId,
			"externalSystemId": externalSystemId,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *organizationRepository) GetOrganizationIdByDomain(ctx context.Context, tenant, domain string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdByDomain")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_DOMAIN]->(d:Domain {domain:$domain})
				return org.id order by org.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant": tenant,
			"domain": domain,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}
