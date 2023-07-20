package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"time"
)

type IssueRepository interface {
	GetMatchedIssue(ctx context.Context, tenant, externalSystem, externalId string) (*dbtype.Node, error)
	MergeIssue(ctx context.Context, tenant string, syncDate time.Time, issue entity.IssueData) error
	LinkIssueWithReporterOrganizationByExternalId(ctx context.Context, tenant, issueId, externalId, externalSystem string) error
	MergeTagForIssue(ctx context.Context, tenant, issueId, tagName, externalSystem string) error
	LinkIssueWithCollaboratorUserByExternalId(ctx context.Context, tenant, issueId, userExternalId, externalSystem string) error
	LinkIssueWithFollowerUserByExternalId(ctx context.Context, tenant, issueId, userExternalId, externalSystem string) error
	LinkIssueWithAssigneeUserByExternalId(ctx context.Context, tenant, issueId, userExternalId, externalSystem string) error
}

type issueRepository struct {
	driver *neo4j.DriverWithContext
}

func NewIssueRepository(driver *neo4j.DriverWithContext) IssueRepository {
	return &issueRepository{
		driver: driver,
	}
}

func (r *issueRepository) GetMatchedIssue(ctx context.Context, tenant, externalSystem, externalId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.GetMatchedIssue")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)-[:IS_LINKED_WITH {externalId:$issueExternalId}]->(e)
				WITH i WHERE i is not null
				return i limit 1`

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":          tenant,
				"externalSystem":  externalSystem,
				"issueExternalId": externalId,
			})
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*dbtype.Node), err
}

func (r *issueRepository) MergeIssue(ctx context.Context, tenant string, syncDate time.Time, issue entity.IssueData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.MergeIssue")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem}) " +
		" MERGE (i:Issue {id:$issueId})-[:ISSUE_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET i.createdAt=$createdAt, " +
		"				i.updatedAt=$updatedAt, " +
		"               i.source=$source, " +
		"				i.sourceOfTruth=$sourceOfTruth, " +
		"				i.appSource=$appSource, " +
		"				i.subject=$subject, " +
		"				i.status=$status, " +
		"				i.priority=$priority, " +
		"				i.description=$description, " +
		"               i:%s," +
		"				i:TimelineEvent," +
		"				i:%s" +
		" ON MATCH SET " +
		"				i.subject = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR i.subject is null OR i.subject = '' THEN $subject ELSE i.subject END, " +
		"				i.description = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR i.description is null OR i.description = '' THEN $description ELSE i.description END, " +
		"				i.status = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR i.status is null OR i.status = '' THEN $status ELSE i.status END, " +
		"				i.priority = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR i.priority is null OR i.priority = '' THEN $priority ELSE i.priority END, " +
		"				i.updatedAt=$now " +
		" WITH i, ext " +
		" MERGE (i)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate, r.externalUrl=$externalUrl " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" WITH i " +
		" FOREACH (x in CASE WHEN i.sourceOfTruth <> $sourceOfTruth THEN [i] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateIssue {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.subject=$subject, alt.status=$status, alt.priority=$priority, alt.description=$description" +
		") RETURN i.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Issue_"+tenant, "TimelineEvent_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"issueId":        issue.Id,
				"externalSystem": issue.ExternalSystem,
				"externalId":     issue.ExternalId,
				"externalUrl":    issue.ExternalUrl,
				"syncDate":       syncDate,
				"createdAt":      utils.TimePtrFirstNonNilNillableAsAny(issue.CreatedAt),
				"updatedAt":      utils.TimePtrFirstNonNilNillableAsAny(issue.UpdatedAt),
				"source":         issue.ExternalSystem,
				"sourceOfTruth":  issue.ExternalSystem,
				"appSource":      constants.AppSourceSyncCustomerOsData,
				"subject":        issue.Subject,
				"description":    issue.Description,
				"status":         issue.Status,
				"priority":       issue.Priority,
				"now":            time.Now().UTC(),
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *issueRepository) MergeTagForIssue(ctx context.Context, tenant, issueId, tagName, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.MergeTagForIssue")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}) " +
		" MERGE (tag:Tag {name:$tagName})-[:TAG_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET tag.id=randomUUID(), " +
		"				tag.createdAt=$now, " +
		"				tag.updatedAt=$now, " +
		"				tag.source=$source," +
		"				tag.appSource=$source," +
		"				tag:%s  " +
		" WITH DISTINCT i, tag " +
		" MERGE (i)-[r:TAGGED]->(tag) " +
		"	ON CREATE SET r.taggedAt=$now " +
		" return r"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Tag_"+tenant),
			map[string]interface{}{
				"tenant":  tenant,
				"issueId": issueId,
				"tagName": tagName,
				"source":  externalSystem,
				"now":     time.Now().UTC(),
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *issueRepository) LinkIssueWithCollaboratorUserByExternalId(ctx context.Context, tenant, issueId, userExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.LinkIssueWithCollaboratorUserByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$userExternalId}]-(u:User)
				MATCH (i:Issue {id:$issueId})-[:IS_LINKED_WITH]->(e)
				MERGE (u)-[:FOLLOWS]->(i)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"issueId":        issueId,
				"userExternalId": userExternalId,
			})
		return nil, err
	})
	return err
}

func (r *issueRepository) LinkIssueWithFollowerUserByExternalId(ctx context.Context, tenant, issueId, userExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.LinkIssueWithFollowerUserByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$userExternalId}]-(u:User)
				MATCH (i:Issue {id:$issueId})-[:IS_LINKED_WITH]->(e)
				MERGE (u)-[:FOLLOWS]->(i)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"issueId":        issueId,
				"userExternalId": userExternalId,
			})
		return nil, err
	})
	return err
}

func (r *issueRepository) LinkIssueWithReporterOrganizationByExternalId(ctx context.Context, tenant, issueId, externalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.LinkIssueWithReporterOrganizationByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$externalId}]-(n)
					WHERE (n:Organization)
				MATCH (i:Issue {id:$issueId})-[:IS_LINKED_WITH]->(e)
				MERGE (n)<-[:REPORTED_BY]-(i)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"issueId":        issueId,
				"externalId":     externalId,
			})
		return nil, err
	})
	return err
}

func (r *issueRepository) LinkIssueWithAssigneeUserByExternalId(ctx context.Context, tenant, issueId, externalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.LinkIssueWithAssigneeUserByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$externalId}]-(n:User)
				MATCH (i:Issue {id:$issueId})-[:IS_LINKED_WITH]->(e)
				MERGE (n)-[:IS_ASSIGNED_TO]->(i)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"issueId":        issueId,
				"externalId":     externalId,
			})
		return nil, err
	})
	return err
}
