package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type JobRoleWriteRepository interface {
	CreateJobRoleForContact(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, jobRoleId, contactId string, data data_fields.JobRoleFields) error
	LinkContactWithOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, jobRoleId, contactId, organizationId string, data data_fields.JobRoleFields) error
	LinkWithUser(ctx context.Context, tenant, userId, jobRoleId string) error
	DeleteJobRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string) error
	SetOtherJobRolesForContactNonPrimaryInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId, skipRoleId string) error
	SetJobRolePrimaryInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, roleId string) error
	UpdateJobRoleDetails(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, roleId string, data data_fields.JobRoleFields) error
	LinkWithOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, roleId, organizationId string) error
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

func (r *jobRoleWriteRepository) LinkContactWithOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, jobRoleId, contactId, organizationId string, data data_fields.JobRoleFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleWriteRepository.LinkContactWithOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("contactId", contactId), log.String("organizationId", organizationId))
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), 
		  								(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}) 
		 MERGE (c)-[:WORKS_AS]->(jr:JobRole)-[:ROLE_IN]->(org) 
		 ON CREATE SET 	jr.id=$jobRoleId, 
						jr.source=$source, 
						jr.appSource=$appSource, 
						jr.jobTitle=$jobTitle, 
						jr.description=$description,
						jr.company=$company,
						jr.startedAt=$startedAt,	
						jr.endedAt=$endedAt,
						jr.primary=$primary,
						jr.createdAt=datetime(), 
						jr.updatedAt=datetime(), 
						jr:JobRole_%s,
						c.updatedAt = datetime(),
						org.updatedAt = datetime()`, tenant)
	params := map[string]interface{}{
		"tenant":         tenant,
		"jobRoleId":      jobRoleId,
		"contactId":      contactId,
		"organizationId": organizationId,
		"source":         utils.IfNotNilString(data.Source),
		"appSource":      utils.IfNotNilString(data.AppSource),
		"jobTitle":       utils.IfNotNilString(data.JobTitle),
		"description":    utils.IfNotNilString(data.Description),
		"company":        utils.IfNotNilString(data.Company),
		"startedAt":      utils.TimePtrAsAny(data.StartedAt),
		"endedAt":        utils.TimePtrAsAny(data.EndedAt),
		"primary":        utils.IfNotNilBool(data.Primary),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *jobRoleWriteRepository) CreateJobRoleForContact(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, jobRoleId, contactId string, data data_fields.JobRoleFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleWriteRepository.CreateJobRoleForContact")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("contactId", contactId))
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
		 MERGE (c)-[:WORKS_AS]->(jr:JobRole {id:$jobRoleId}) 
		 ON CREATE SET 	jr.source=$source, 
						jr.appSource=$appSource, 
						jr.jobTitle=$jobTitle, 
						jr.description=$description,
						jr.company=$company,
						jr.startedAt=$startedAt,	
						jr.endedAt=$endedAt,
						jr.primary=$primary,
						jr.createdAt=datetime(), 
						jr.updatedAt=datetime(), 
						jr:JobRole_%s,
						c.updatedAt = datetime()`, tenant)
	params := map[string]interface{}{
		"tenant":      tenant,
		"jobRoleId":   jobRoleId,
		"contactId":   contactId,
		"source":      utils.IfNotNilString(data.Source),
		"appSource":   utils.IfNotNilString(data.AppSource),
		"jobTitle":    utils.IfNotNilString(data.JobTitle),
		"description": utils.IfNotNilString(data.Description),
		"company":     utils.IfNotNilString(data.Company),
		"startedAt":   utils.TimePtrAsAny(data.StartedAt),
		"endedAt":     utils.TimePtrAsAny(data.EndedAt),
		"primary":     utils.IfNotNilBool(data.Primary),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *jobRoleWriteRepository) UpdateJobRoleDetails(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, roleId string, data data_fields.JobRoleFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.UpdateJobRoleDetails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (jr:JobRole_%s {id:$roleId})
			SET jr.updatedAt=datetime()`, tenant)
	params := map[string]interface{}{
		"roleId": roleId,
	}

	if data.JobTitle != nil {
		cypher += ", jr.jobTitle=$jobTitle"
		params["jobTitle"] = *data.JobTitle
	}

	if data.Description != nil {
		cypher += ", jr.description=$description"
		params["description"] = *data.Description
	}

	if data.Company != nil {
		cypher += ", jr.company=$company"
		params["company"] = *data.Company
	}

	if data.Primary != nil {
		cypher += ", jr.primary=$primary"
		params["primary"] = *data.Primary
	}

	if data.StartedAt != nil {
		cypher += ", jr.startedAt=$startedAt"
		params["startedAt"] = *data.StartedAt
	}

	if data.EndedAt != nil {
		cypher += ", jr.endedAt=$endedAt"
		params["endedAt"] = *data.EndedAt
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
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

func (r *jobRoleWriteRepository) SetOtherJobRolesForContactNonPrimaryInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId, skipRoleId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.SetOtherJobRolesForContactNonPrimaryInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				 (c)-[:WORKS_AS]->(r:JobRole)
			WHERE r.id <> $skipRoleId
            SET r.primary=false,
				r.updatedAt=datetime({timezone: 'UTC'})`
	params := map[string]interface{}{
		"tenant":     tenant,
		"contactId":  contactId,
		"skipRoleId": skipRoleId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *jobRoleWriteRepository) LinkWithOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, roleId string, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(r:JobRole {id:$roleId})<-[:WORKS_AS]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (r)-[rel:ROLE_IN]->(org2:Organization)
				WHERE org2.id <> org.id
			DELETE rel
			WITH r, org
			MERGE (r)-[:ROLE_IN]->(org)
			SET r.updatedAt=datetime()`
	params := map[string]interface{}{
		"tenant":         tenant,
		"roleId":         roleId,
		"organizationId": organizationId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *jobRoleWriteRepository) SetJobRolePrimaryInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, roleId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.SetJobRolePrimaryInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (j:JobRole_%s {id:$roleId})
			SET j.primary=true,
				j.updatedAt=datetime()`, tenant)
	params := map[string]interface{}{
		"roleId": roleId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}
