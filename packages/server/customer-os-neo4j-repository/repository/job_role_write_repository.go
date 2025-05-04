package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "JobRoleWriteRepository.LinkWithUser")
	defer spans.Finish()

	spans.LogKV("userId", userId)
	spans.LogKV("jobRoleId", jobRoleId)

	cypher := fmt.Sprintf(`MATCH (u:User_%s {id: $userId})
              MERGE (jr:JobRole:JobRole_%s {id: $jobRoleId})
              MERGE (u)-[r:WORKS_AS]->(jr)
			  SET u.updatedAt = datetime()`, tenant, tenant)
	params := map[string]any{
		"userId":    userId,
		"jobRoleId": jobRoleId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *jobRoleWriteRepository) LinkContactWithOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, jobRoleId, contactId, organizationId string, data data_fields.JobRoleFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "JobRoleWriteRepository.LinkContactWithOrganization")
	defer spans.Finish()

	spans.LogKV("contactId", contactId)
	spans.LogKV("organizationId", organizationId)
	spans.LogObjectAsJson("data", data)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *jobRoleWriteRepository) CreateJobRoleForContact(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, jobRoleId, contactId string, data data_fields.JobRoleFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "JobRoleWriteRepository.CreateJobRoleForContact")
	defer spans.Finish()

	spans.LogKV("contactId", contactId)
	spans.LogObjectAsJson("data", data)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *jobRoleWriteRepository) UpdateJobRoleDetails(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, roleId string, data data_fields.JobRoleFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "JobRoleRepository.UpdateJobRoleDetails")
	defer spans.Finish()

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

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *jobRoleWriteRepository) DeleteJobRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "JobRoleRepository.DeleteJobRoleInTx")
	defer spans.Finish()

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "JobRoleRepository.SetOtherJobRolesForContactNonPrimaryInTx")
	defer spans.Finish()

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

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *jobRoleWriteRepository) LinkWithOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, roleId string, organizationId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "JobRoleRepository.LinkWithOrganization")
	defer spans.Finish()

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

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *jobRoleWriteRepository) SetJobRolePrimaryInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, roleId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "JobRoleRepository.SetJobRolePrimaryInTx")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (j:JobRole_%s {id:$roleId})
			SET j.primary=true,
				j.updatedAt=datetime()`, tenant)
	params := map[string]interface{}{
		"roleId": roleId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}
