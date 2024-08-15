package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type JobRoleCreateFields struct {
	Description  string       `json:"description"`
	JobTitle     string       `json:"jobTitle"`
	StartedAt    *time.Time   `json:"startedAt"`
	EndedAt      *time.Time   `json:"endedAt"`
	SourceFields model.Source `json:"sourceFields"`
	CreatedAt    time.Time    `json:"createdAt"`
	Primary      bool         `json:"primary"`
}

type JobRoleWriteRepository interface {
	CreateJobRole(ctx context.Context, tenant, jobRoleId string, data JobRoleCreateFields) error
	LinkWithUser(ctx context.Context, tenant, userId, jobRoleId string) error
	LinkContactWithOrganization(ctx context.Context, tenant, contactId, organizationId string, data JobRoleCreateFields) error
}

type jobRoleWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewJobRoleWriteRepository(driver *neo4j.DriverWithContext, database string) JobRoleWriteRepository {
	return &jobRoleWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *jobRoleWriteRepository) CreateJobRole(ctx context.Context, tenant, jobRoleId string, data JobRoleCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleWriteRepository.CreateJobRole")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, jobRoleId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
				MERGE (jr:JobRole:JobRole_%s {id:$id}) 
				SET 	jr.jobTitle = $jobTitle,
						jr.description = $description,
						jr.createdAt = $createdAt,
						jr.updatedAt = datetime(),
						jr.startedAt = $startedAt,
 						jr.endedAt = $endedAt,
						jr.sourceOfTruth = $sourceOfTruth,
						jr.source = $source,
						jr.appSource = $appSource,
						jr.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"id":            jobRoleId,
		"jobTitle":      data.JobTitle,
		"description":   data.Description,
		"tenant":        tenant,
		"startedAt":     utils.TimePtrAsAny(data.StartedAt),
		"endedAt":       utils.TimePtrAsAny(data.EndedAt),
		"sourceOfTruth": data.SourceFields.SourceOfTruth,
		"source":        data.SourceFields.Source,
		"appSource":     data.SourceFields.AppSource,
		"createdAt":     data.CreatedAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *jobRoleWriteRepository) LinkWithUser(ctx context.Context, tenant, userId, jobRoleId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleWriteRepository.LinkWithUser")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("userId", userId), log.String("jobRoleId", jobRoleId))

	cypher := fmt.Sprintf(`MATCH (u:User_%s {id: $userId})
              MERGE (jr:JobRole:JobRole_%s {id: $jobRoleId})
              ON CREATE SET jr.syncedWithEventStore = true
              MERGE (u)-[r:WORKS_AS]->(jr)
			  SET u.updatedAt = datetime()`, tenant, tenant)
	params := map[string]any{
		"userId":    userId,
		"jobRoleId": jobRoleId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *jobRoleWriteRepository) LinkContactWithOrganization(ctx context.Context, tenant, contactId, organizationId string, data JobRoleCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.LinkContactWithOrganizationByInternalId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("contactId", contactId), log.String("organizationId", organizationId))
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), 
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
						jr.updatedAt=datetime(), 
						jr:JobRole_%s,
						c.updatedAt = datetime()
		 ON MATCH SET 	jr.jobTitle = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true  OR jr.jobTitle is null OR jr.jobTitle = '' THEN $jobTitle ELSE jr.jobTitle END,
						jr.description = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true  OR jr.description is null OR jr.description = '' THEN $description ELSE jr.description END,
						jr.primary = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true THEN $primary ELSE jr.primary END,
						jr.startedAt = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true THEN $startedAt ELSE jr.startedAt END,
						jr.endedAt = CASE WHEN jr.sourceOfTruth=$source OR $overwrite=true THEN $endedAt ELSE jr.endedAt END,
						jr.sourceOfTruth = case WHEN $overwrite=true THEN $source ELSE jr.sourceOfTruth END,
						jr.updatedAt = datetime(),
						c.updatedAt = datetime()`, tenant)
	params := map[string]interface{}{
		"tenant":         tenant,
		"contactId":      contactId,
		"organizationId": organizationId,
		"source":         data.SourceFields.Source,
		"sourceOfTruth":  data.SourceFields.SourceOfTruth,
		"appSource":      data.SourceFields.AppSource,
		"jobTitle":       data.JobTitle,
		"description":    data.Description,
		"createdAt":      data.CreatedAt,
		"startedAt":      utils.TimePtrAsAny(data.StartedAt),
		"endedAt":        utils.TimePtrAsAny(data.EndedAt),
		"primary":        data.Primary,
		"overwrite":      data.SourceFields.Source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
