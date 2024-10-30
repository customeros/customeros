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

type IssueCreateFields struct {
	CreatedAt                 time.Time          `json:"createdAt"`
	SourceFields              model.SourceFields `json:"sourceFields"`
	GroupId                   string             `json:"groupId"`
	Subject                   string             `json:"subject"`
	Description               string             `json:"description"`
	Status                    string             `json:"status"`
	Priority                  string             `json:"priority"`
	ReportedByOrganizationId  string             `json:"reportedByOrganizationId"`
	SubmittedByOrganizationId string             `json:"submittedByOrganizationId"`
	SubmittedByUserId         string             `json:"submittedByUserId"`
}

type IssueUpdateFields struct {
	GroupId     string `json:"groupId"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Source      string `json:"source"`
}

type IssueWriteRepository interface {
	Create(ctx context.Context, tenant, issueId string, data IssueCreateFields) error
	Update(ctx context.Context, tenant, issueId string, data IssueUpdateFields) error
	AddUserAssignee(ctx context.Context, tenant, issueId, userId string) error
	RemoveUserAssignee(ctx context.Context, tenant, issueId, userId string) error
	AddUserFollower(ctx context.Context, tenant, issueId, userId string) error
	RemoveUserFollower(ctx context.Context, tenant, issueId, userId string) error

	ReportedByOrganizationWithGroupId(ctx context.Context, tenant, organizationId, groupId string) error
	RemoveReportedByOrganizationWithGroupId(ctx context.Context, tenant, organizationId, groupId string) error

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

func (r *issueWriteRepository) Create(ctx context.Context, tenant, issueId string, data IssueCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.Create")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, issueId)
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
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.groupId=$groupId,
								i.subject=$subject,
								i.description=$description,	
								i.status=$status,	
								i.priority=$priority
							ON MATCH SET 	
								i.groupId = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.groupId is null OR i.groupId = '' THEN $groupId ELSE i.groupId END,
								i.subject = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.subject is null OR i.subject = '' THEN $subject ELSE i.subject END,
								i.description = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.description is null OR i.description = '' THEN $description ELSE i.description END,
								i.status = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.status is null OR i.status = '' THEN $status ELSE i.status END,
								i.priority = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.priority is null OR i.priority = '' THEN $priority ELSE i.priority END,
								i.updatedAt = datetime(),
								i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END
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
		"createdAt":                 data.CreatedAt,
		"source":                    data.SourceFields.Source,
		"sourceOfTruth":             data.SourceFields.SourceOfTruth,
		"appSource":                 data.SourceFields.AppSource,
		"groupId":                   data.GroupId,
		"subject":                   data.Subject,
		"description":               data.Description,
		"status":                    data.Status,
		"priority":                  data.Priority,
		"reportedByOrganizationId":  data.ReportedByOrganizationId,
		"submittedByOrganizationId": data.SubmittedByOrganizationId,
		"submittedByUserId":         data.SubmittedByUserId,
		"overwrite":                 data.SourceFields.Source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *issueWriteRepository) Update(ctx context.Context, tenant, issueId string, data IssueUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.Update")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, issueId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId})
		 	SET	
				i.groupId = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.groupId is null OR i.groupId = '' THEN $groupId ELSE i.groupId END,
				i.subject = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.subject is null OR i.subject = '' THEN $subject ELSE i.subject END,
				i.description = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.description is null OR i.description = '' THEN $description ELSE i.description END,
				i.status = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.status is null OR i.status = '' THEN $status ELSE i.status END,
				i.priority = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.priority is null OR i.priority = '' THEN $priority ELSE i.priority END,
				i.updatedAt = datetime(),
				i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END`
	params := map[string]any{
		"tenant":        tenant,
		"issueId":       issueId,
		"groupId":       data.GroupId,
		"subject":       data.Subject,
		"description":   data.Description,
		"status":        data.Status,
		"priority":      data.Priority,
		"sourceOfTruth": data.Source,
		"overwrite":     data.Source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *issueWriteRepository) AddUserAssignee(ctx context.Context, tenant, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.AddUserAssignee")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, issueId)
	span.LogFields(log.String("userId", userId))

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

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *issueWriteRepository) AddUserFollower(ctx context.Context, tenant, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.AddUserFollower")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, issueId)
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

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *issueWriteRepository) RemoveUserAssignee(ctx context.Context, tenant, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.RemoveUserAssignee")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, issueId)
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

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *issueWriteRepository) RemoveUserFollower(ctx context.Context, tenant, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueWriteRepository.RemoveUserFollower")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, issueId)
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

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *issueWriteRepository) ReportedByOrganizationWithGroupId(ctx context.Context, tenant, organizationId, groupId string) error {
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

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *issueWriteRepository) RemoveReportedByOrganizationWithGroupId(ctx context.Context, tenant, organizationId, groupId string) error {
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

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
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
