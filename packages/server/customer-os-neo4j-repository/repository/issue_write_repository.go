package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type IssueWriteRepository interface {
	Create(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId string, data data_fields.IssueFields) error
	Update(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId string, data data_fields.IssueFields) error
	AddUserAssignee(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId, userId string) error
	RemoveUserAssignee(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId, userId string) error
	AddUserFollower(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId, userId string) error
	RemoveUserFollower(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId, userId string) error

	ReportedByOrganizationWithGroupId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, groupId string) error
	RemoveReportedByOrganizationWithGroupId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, groupId string) error

	LinkUnthreadIssuesToOrganizationByGroupId(ctx context.Context) error
}

type issueWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewIssueWriteRepository(driver *neo4j.DriverWithContext, database string) IssueWriteRepository {
	return &issueWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *issueWriteRepository) Create(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId string, data data_fields.IssueFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.Create")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, issueId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}) 
							ON CREATE SET 
								i:Issue_%s,
								i:TimelineEvent,
								i:TimelineEvent_%s,
								i.createdAt=$createdAt,
								i.updatedAt=datetime(),
								i.source=$source,
								i.appSource=$appSource,
								i.groupId=$groupId,
								i.subject=$subject,
								i.description=$description,	
								i.status=$status,	
								i.priority=$priority
							WITH i, t
							OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$reportedByOrganizationId}) 
							WHERE $reportedByOrganizationId <> ""
							FOREACH (ignore IN CASE WHEN o IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:REPORTED_BY]->(o))
							WITH i, t
							OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$submittedByOrganizationId}) 
							WHERE $submittedByOrganizationId <> ""
							FOREACH (ignore IN CASE WHEN o IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:SUBMITTED_BY]->(o))
							WITH i, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$submittedByUserId}) 
							WHERE $submittedByUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:SUBMITTED_BY]->(u))
							`, tenant, tenant)
	params := map[string]any{
		"tenant":                    tenant,
		"issueId":                   issueId,
		"createdAt":                 utils.IfNotNilTimeWithDefault(data.CreatedAt, utils.Now()),
		"source":                    utils.IfNotNilString(data.Source),
		"appSource":                 utils.IfNotNilString(data.AppSource),
		"groupId":                   utils.IfNotNilString(data.GroupId),
		"subject":                   utils.IfNotNilString(data.Subject),
		"description":               utils.IfNotNilString(data.Description),
		"status":                    utils.IfNotNilString(data.Status),
		"priority":                  utils.IfNotNilString(data.Priority),
		"reportedByOrganizationId":  utils.IfNotNilString(data.ReportedByOrganizationId),
		"submittedByOrganizationId": utils.IfNotNilString(data.SubmittedByOrganizationId),
		"submittedByUserId":         utils.IfNotNilString(data.SubmittedByUserId),
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

func (r *issueWriteRepository) Update(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId string, data data_fields.IssueFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.Create")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, issueId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId})
		 	SET	i.updatedAt = datetime() `
	params := map[string]any{
		"tenant":  tenant,
		"issueId": issueId,
	}

	if data.GroupId != nil {
		cypher += `, i.groupId = $groupId`
		params["groupId"] = *data.GroupId
	}
	if data.Subject != nil {
		cypher += `, i.subject = $subject`
		params["subject"] = *data.Subject
	}
	if data.Description != nil {
		cypher += `, i.description = $description`
		params["description"] = *data.Description
	}
	if data.Status != nil {
		cypher += `, i.status = $status`
		params["status"] = *data.Status
	}
	if data.Priority != nil {
		cypher += `, i.priority = $priority`
		params["priority"] = *data.Priority
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *issueWriteRepository) AddUserAssignee(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.AddUserAssignee")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, issueId)
	span.LogKV("userId", userId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}),
				(t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
		 	MERGE (i)-[:ASSIGNED_TO]->(u)
				ON CREATE SET i.updatedAt = datetime()`
	params := map[string]any{
		"tenant":  tenant,
		"issueId": issueId,
		"userId":  userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *issueWriteRepository) AddUserFollower(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.AddUserFollower")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, issueId)
	span.LogFields(log.String("userId", userId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}),
				(t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
		 	MERGE (i)-[:FOLLOWED_BY]->(u)
				ON CREATE SET i.updatedAt = datetime()`
	params := map[string]any{
		"tenant":  tenant,
		"issueId": issueId,
		"userId":  userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *issueWriteRepository) RemoveUserAssignee(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.RemoveUserAssignee")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, issueId)
	span.LogFields(log.String("userId", userId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}),
				(t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(i)-[r:ASSIGNED_TO]->(u)
				SET i.updatedAt = datetime()
		 		DELETE r`
	params := map[string]any{
		"tenant":  tenant,
		"issueId": issueId,
		"userId":  userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *issueWriteRepository) RemoveUserFollower(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.RemoveUserFollower")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, issueId)
	span.LogFields(log.String("userId", userId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}),
				(t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(i)-[r:FOLLOWED_BY]->(u)
				SET i.updatedAt = datetime()
		 		DELETE r`
	params := map[string]any{
		"tenant":  tenant,
		"issueId": issueId,
		"userId":  userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *issueWriteRepository) ReportedByOrganizationWithGroupId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, groupId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.ReportedByOrganizationWithGroupId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("organizationId", organizationId))
	span.LogFields(log.String("groupId", groupId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId}) 
			   OPTIONAL MATCH (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {groupId:$groupId}) 
			   MERGE (i)-[:REPORTED_BY]->(o)`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"groupId":        groupId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *issueWriteRepository) RemoveReportedByOrganizationWithGroupId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, groupId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.RemoveReportedByOrganizationWithGroupId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("organizationId", organizationId))
	span.LogFields(log.String("groupId", groupId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId})<-[r:REPORTED_BY]-(i:Issue{groupId:$groupId}) 
			   DELETE r`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"groupId":        groupId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *issueWriteRepository) LinkUnthreadIssuesToOrganizationByGroupId(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.LinkUnthreadIssuesToOrganizationByGroupId")
	defer span.Finish()

	cypher := `match (t:Tenant)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem{id:"unthread"})<-[:IS_LINKED_WITH]-(i:Issue)
			   with t, i
			   match (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {slackChannelId: i.groupId})
			   where not (i)-[:REPORTED_BY]->(o)
			   MERGE (i)-[:REPORTED_BY]->(o)
				SET o.updatedAt = datetime()`
	params := map[string]any{}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
