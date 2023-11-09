package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	contactevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type JobRoleRepository interface {
	LinkWithUser(ctx context.Context, tenant, userId, jobRoleId string, updatedAt time.Time) error
	CreateJobRole(ctx context.Context, tenant, jobRoleId string, event events.JobRoleCreateEvent) error
	LinkContactWithOrganization(ctx context.Context, tenant, contactId string, data contactevent.ContactLinkWithOrganizationEvent) error
}

type jobRoleRepository struct {
	driver *neo4j.DriverWithContext
}

func NewJobRoleRepository(driver *neo4j.DriverWithContext) JobRoleRepository {
	return &jobRoleRepository{
		driver: driver,
	}
}

func (r *jobRoleRepository) LinkWithUser(ctx context.Context, tenant, userId, jobRoleId string, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.LinkWithUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("userId", userId), log.String("jobRoleId", jobRoleId))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (u:User_%s {id: $userId})
              MERGE (jr:JobRole:JobRole_%s {id: $jobRoleId})
              ON CREATE SET jr.syncedWithEventStore = true
              MERGE (u)-[r:WORKS_AS]->(jr)
			  SET u.updatedAt = $updatedAt`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if _, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]any{
				"userId":    userId,
				"jobRoleId": jobRoleId,
				"updatedAt": updatedAt,
			}); err != nil {
			return nil, err
		} else {
			return nil, nil
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *jobRoleRepository) CreateJobRole(ctx context.Context, tenant, jobRoleId string, event events.JobRoleCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.CreateJobRole")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("jobRoleId", jobRoleId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `  MATCH (t:Tenant {name:$tenant})
				MERGE (jr:JobRole:JobRole_%s {id:$id}) 
				SET 	jr.jobTitle = $jobTitle,
						jr.description = $description,
						jr.createdAt = $createdAt,
						jr.updatedAt = $updatedAt,
						jr.startedAt = $startedAt,
 						jr.endedAt = $endedAt,
						jr.sourceOfTruth = $sourceOfTruth,
						jr.source = $source,
						jr.appSource = $appSource,
						jr.syncedWithEventStore = true
`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"id":            jobRoleId,
				"jobTitle":      event.JobTitle,
				"description":   event.Description,
				"tenant":        event.Tenant,
				"startedAt":     utils.TimePtrFirstNonNilNillableAsAny(event.StartedAt),
				"endedAt":       utils.TimePtrFirstNonNilNillableAsAny(event.EndedAt),
				"sourceOfTruth": event.SourceOfTruth,
				"source":        event.Source,
				"appSource":     event.AppSource,
				"createdAt":     event.CreatedAt,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *jobRoleRepository) LinkContactWithOrganization(ctx context.Context, tenant, contactId string, eventData contactevent.ContactLinkWithOrganizationEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.LinkContactWithOrganizationByInternalId")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("contactId", contactId))

	query := fmt.Sprintf(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), 
		  								(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}) 
		 MERGE (c)-[:WORKS_AS]->(jr:JobRole)-[:ROLE_IN]->(org) 
		 ON CREATE SET 	jr.id=randomUUID(), 
						jr.source=$source, 
						jr.sourceOfTruth=$sourceOfTruth, 
						jr.appSource=$appSource, 
						jr.jobTitle=$jobTitle, 
						jr.description=$description,
						jr.startedAt=$startedAt,	
						jr.endedAt=$endedAt,
						jr.primary=$primary,
						jr.createdAt=$createdAt, 
						jr.updatedAt=$updatedAt, 
						jr:JobRole_%s,
						c.updatedAt = $updatedAt
		 ON MATCH SET 	jr.jobTitle = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true  OR jr.jobTitle is null OR jr.jobTitle = '' THEN $jobTitle ELSE jr.jobTitle END,
						jr.description = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true  OR jr.description is null OR jr.description = '' THEN $description ELSE jr.description END,
						jr.primary = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true THEN $primary ELSE jr.primary END,
						jr.startedAt = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true THEN $startedAt ELSE jr.startedAt END,
						jr.endedAt = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true THEN $endedAt ELSE jr.endedAt END,
						jr.sourceOfTruth = case WHEN $overwrite=true THEN $source ELSE jr.sourceOfTruth END,
						jr.updatedAt = $updatedAt,
						c.updatedAt = $updatedAt`, tenant)
	params := map[string]interface{}{
		"tenant":         tenant,
		"contactId":      contactId,
		"organizationId": eventData.OrganizationId,
		"source":         helper.GetSource(eventData.SourceFields.Source),
		"sourceOfTruth":  helper.GetSourceOfTruth(eventData.SourceFields.SourceOfTruth),
		"appSource":      helper.GetAppSource(eventData.SourceFields.AppSource),
		"jobTitle":       eventData.JobTitle,
		"description":    eventData.Description,
		"createdAt":      eventData.CreatedAt,
		"updatedAt":      eventData.UpdatedAt,
		"startedAt":      utils.TimePtrFirstNonNilNillableAsAny(eventData.StartedAt),
		"endedAt":        utils.TimePtrFirstNonNilNillableAsAny(eventData.EndedAt),
		"primary":        eventData.Primary,
		"overwrite":      helper.GetSourceOfTruth(eventData.SourceFields.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}
