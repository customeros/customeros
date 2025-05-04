package neo4j_repository

import (
	"fmt"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"golang.org/x/net/context"
)

type TagReadRepository interface {
	GetById(ctx context.Context, tenant, tagId string) (*dbtype.Node, error)
	GetAll(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetAllByEntityType(ctx context.Context, tenant string, entityType commonmodel.EntityType) ([]*dbtype.Node, error)
	GetByEntityTypeAndName(ctx context.Context, tenant string, entityType commonmodel.EntityType, name string) (*dbtype.Node, error)
	GetForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetForLogEntries(ctx context.Context, tenant string, logEntryIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeWithRelationAndId, error)
}

type tagReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewTagReadRepository(driver *neo4j.DriverWithContext, database string) TagReadRepository {
	return &tagReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *tagReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *tagReadRepository) GetById(ctx context.Context, tenant, tagId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagReadRepository.GetById")
	defer spans.Finish()

	spans.LogKV("tagId", tagId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$tagId}) return tag`
	params := map[string]any{
		"tenant": tenant,
		"tagId":  tagId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *tagReadRepository) GetAll(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagRepository.GetAll")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag)
			RETURN tag ORDER BY tag.name`
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), err
}

func (r *tagReadRepository) GetForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagReadRepository.GetForIssues")
	defer spans.Finish()

	spans.LogObjectAsJson("issueIds", issueIds)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)-[rel:TAGGED]->(tag:Tag)-[:TAG_BELONGS_TO_TENANT]->(t)
			WHERE i.id IN $issueIds
			RETURN tag, rel, i.id ORDER BY rel.taggedAt, tag.name`
	params := map[string]any{
		"tenant":   tenant,
		"issueIds": issueIds,
	}

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result.([]*utils.DbNodeWithRelationAndId)))
	if len(result.([]*utils.DbNodeWithRelationAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagReadRepository) GetForLogEntries(ctx context.Context, tenant string, logEntryIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagReadRepository.GetForLogEntries")
	defer spans.Finish()

	spans.LogObjectAsJson("logEntryIds", logEntryIds)

	cypher := fmt.Sprintf(`MATCH (l:LogEntry)-[rel:TAGGED]->(tag:Tag)-[:TAG_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			WHERE l.id IN $logEntryIds AND l:LogEntry_%s
			RETURN tag, rel, l.id ORDER BY rel.taggedAt, tag.name`, tenant)
	params := map[string]any{
		"tenant":      tenant,
		"logEntryIds": logEntryIds,
	}

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result.([]*utils.DbNodeWithRelationAndId)))
	if len(result.([]*utils.DbNodeWithRelationAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagReadRepository) GetForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagReadRepository.GetForContacts")
	defer spans.Finish()

	spans.LogObjectAsJson("contactIds", contactIds)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[rel:TAGGED]->(tag:Tag)-[:TAG_BELONGS_TO_TENANT]->(t)
			WHERE c.id IN $contactIds
			RETURN tag, rel, c.id ORDER BY rel.taggedAt, tag.name`
	params := map[string]any{
		"tenant":     tenant,
		"contactIds": contactIds,
	}

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result.([]*utils.DbNodeWithRelationAndId)))
	if len(result.([]*utils.DbNodeWithRelationAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagReadRepository) GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagReadRepository.GetForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("organizationIds", organizationIds)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[rel:TAGGED]->(tag:Tag)-[:TAG_BELONGS_TO_TENANT]->(t)
			WHERE o.id IN $organizationIds
			RETURN tag, rel, o.id ORDER BY rel.taggedAt, tag.name`
	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIds,
	}

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result.([]*utils.DbNodeWithRelationAndId)))
	if len(result.([]*utils.DbNodeWithRelationAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagReadRepository) GetAllByEntityType(ctx context.Context, tenant string, entityType commonmodel.EntityType) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagRepository.GetAll")
	defer spans.Finish()

	spans.LogKV("entityType", entityType.String())

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {entityType:$entityType})
			RETURN tag ORDER BY tag.name`
	params := map[string]any{
		"tenant":     tenant,
		"entityType": entityType.String(),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		resultWithContext, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, resultWithContext, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), err
}

func (r *tagReadRepository) GetByEntityTypeAndName(ctx context.Context, tenant string, entityType commonmodel.EntityType, name string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagReadRepository.GetByEntityTypeAndName")
	defer spans.Finish()

	spans.LogKV("entityType", entityType.String())
	spans.LogKV("name", name)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {entityType:$entityType, name:$name})
			RETURN tag`
	params := map[string]any{
		"tenant":     tenant,
		"entityType": entityType.String(),
		"name":       name,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	if dbRecords == nil || len(dbRecords.([]*dbtype.Node)) == 0 {
		spans.LogKV("result", "not found")
		return nil, nil
	}
	spans.LogKV("result", "found")
	return dbRecords.([]*dbtype.Node)[0], err
}
