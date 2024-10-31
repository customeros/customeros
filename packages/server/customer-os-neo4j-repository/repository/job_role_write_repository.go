package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type JobRoleFields struct {
	Description  string             `json:"description"`
	JobTitle     string             `json:"jobTitle"`
	StartedAt    *time.Time         `json:"startedAt"`
	EndedAt      *time.Time         `json:"endedAt"`
	SourceFields model.SourceFields `json:"sourceFields"`
	Primary      bool               `json:"primary"`
}

type JobRoleWriteRepository interface {
	CreateJobRole(ctx context.Context, tenant, jobRoleId string, data JobRoleFields) error
	CreateJobRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, input entity.JobRoleEntity) (*dbtype.Node, error)
	LinkWithUser(ctx context.Context, tenant, userId, jobRoleId string) error
	LinkContactWithOrganization(ctx context.Context, tenant, contactId, organizationId string, data JobRoleFields) error

	DeleteJobRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string) error
	SetOtherJobRolesForContactNonPrimaryInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, skipRoleId string) error
	UpdateJobRoleDetails(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string, input entity.JobRoleEntity) (*dbtype.Node, error)
	LinkWithOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, roleId, organizationId string) error
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

func (r *jobRoleWriteRepository) CreateJobRole(ctx context.Context, tenant, jobRoleId string, data JobRoleFields) error {
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
						jr.createdAt = datetime(),
						jr.updatedAt = datetime(),
						jr.startedAt = $startedAt,
 						jr.endedAt = $endedAt,
						jr.source = $source,
						jr.appSource = $appSource`, tenant)
	params := map[string]any{
		"id":          jobRoleId,
		"jobTitle":    data.JobTitle,
		"description": data.Description,
		"tenant":      tenant,
		"startedAt":   utils.TimePtrAsAny(data.StartedAt),
		"endedAt":     utils.TimePtrAsAny(data.EndedAt),
		"source":      data.SourceFields.Source,
		"appSource":   data.SourceFields.AppSource,
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

func (r *jobRoleWriteRepository) LinkContactWithOrganization(ctx context.Context, tenant, contactId, organizationId string, data JobRoleFields) error {
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
						jr.appSource=$appSource, 
						jr.jobTitle=$jobTitle, 
						jr.description=$description,
						jr.startedAt=$startedAt,	
						jr.endedAt=$endedAt,
						jr.primary=$primary,
						jr.createdAt=datetime(), 
						jr.updatedAt=datetime(), 
						jr:JobRole_%s,
						c.updatedAt = datetime(),
						org.updatedAt = datetime()
		 ON MATCH SET 	jr.jobTitle = CASE WHEN $overwrite=true OR jr.jobTitle is null OR jr.jobTitle = '' THEN $jobTitle ELSE jr.jobTitle END,
						jr.description = CASE WHEN $overwrite=true OR jr.description is null OR jr.description = '' THEN $description ELSE jr.description END,
						jr.primary = CASE WHEN OR $overwrite=true THEN $primary ELSE jr.primary END,
						jr.startedAt = CASE WHEN OR $overwrite=true THEN $startedAt ELSE jr.startedAt END,
						jr.endedAt = CASE WHEN OR $overwrite=true THEN $endedAt ELSE jr.endedAt END,
						jr.updatedAt = datetime(),
						c.updatedAt = datetime(),
						org.updatedAt = datetime()`, tenant)
	params := map[string]interface{}{
		"tenant":         tenant,
		"contactId":      contactId,
		"organizationId": organizationId,
		"source":         data.SourceFields.Source,
		"appSource":      data.SourceFields.AppSource,
		"jobTitle":       data.JobTitle,
		"description":    data.Description,
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

func (r *jobRoleWriteRepository) CreateJobRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, contactId string, input entity.JobRoleEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.CreateJobRole")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[:WORKS_AS]->(r:JobRole {id:randomUUID()}) " +
		" ON CREATE SET r.jobTitle=$jobTitle, " +
		"				r.primary=$primary, " +
		"				r.description=$description, " +
		"				r.company=$company, " +
		"				r.source=$source, " +
		"				r.appSource=$appSource, " +
		"				r.createdAt=$now, " +
		"				r.updatedAt=datetime(), " +
		"				r.startedAt=$startedAt, " +
		"				r.endedAt=$endedAt, " +
		"				r:%s " +
		" RETURN r"

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "JobRole_"+tenant),
		map[string]interface{}{
			"tenant":      tenant,
			"contactId":   contactId,
			"jobTitle":    input.JobTitle,
			"description": input.Description,
			"company":     input.Company,
			"primary":     input.Primary,
			"source":      input.Source,
			"appSource":   input.AppSource,
			"startedAt":   utils.TimePtrAsAny(input.StartedAt, utils.TimePtr(utils.Now())),
			"endedAt":     utils.TimePtrAsAny(input.EndedAt),
			"now":         utils.Now(),
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *jobRoleWriteRepository) UpdateJobRoleDetails(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string, input entity.JobRoleEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.UpdateJobRoleDetails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[:WORKS_AS]->(r:JobRole {id:$roleId})
			SET r.jobTitle=$jobTitle, 
				r.primary=$primary,
				r.description=$description,
				r.company=$company,
				r.startedAt=$startedAt,
				r.endedAt=$endedAt,
				r.updatedAt=datetime()
			RETURN r`,
		map[string]interface{}{
			"tenant":      tenant,
			"contactId":   contactId,
			"roleId":      roleId,
			"jobTitle":    input.JobTitle,
			"description": input.Description,
			"company":     input.Company,
			"primary":     input.Primary,
			"now":         utils.Now(),
			"startedAt":   utils.TimePtrAsAny(input.StartedAt),
			"endedAt":     utils.TimePtrAsAny(input.EndedAt),
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *jobRoleWriteRepository) DeleteJobRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.DeleteJobRoleInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)-[:WORKS_AS]->(r:JobRole {id:$roleId})
			DETACH DELETE r`,
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
			"roleId":    roleId,
		})
	return err
}

func (r *jobRoleWriteRepository) SetOtherJobRolesForContactNonPrimaryInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, skipRoleId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.SetOtherJobRolesForContactNonPrimaryInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				 (c)-[:WORKS_AS]->(r:JobRole)
			WHERE r.id <> $skipRoleId
            SET r.primary=false,
				r.updatedAt=datetime({timezone: 'UTC'})`,
		map[string]interface{}{
			"tenant":     tenant,
			"contactId":  contactId,
			"skipRoleId": skipRoleId,
		})
	return err
}

func (r *jobRoleWriteRepository) LinkWithOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, roleId string, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(r:JobRole {id:$roleId})<-[:WORKS_AS]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (r)-[rel:ROLE_IN]->(org2:Organization)
				WHERE org2.id <> org.id
			DELETE rel
			WITH r, org
			MERGE (r)-[:ROLE_IN]->(org)
			SET r.updatedAt=datetime()`,
		map[string]interface{}{
			"tenant":         tenant,
			"roleId":         roleId,
			"organizationId": organizationId,
		})
	return err
}
