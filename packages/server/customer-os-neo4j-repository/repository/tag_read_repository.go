package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type TagReadRepository interface {
	GetById(ctx context.Context, tenant, tagId string) (*dbtype.Node, error)
	GetAll(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetByNameOptional(ctx context.Context, tenant, name string) (*dbtype.Node, error)
	GetForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagReadRepository.GetById")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("tagId", tagId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$tagId}) return tag`
	params := map[string]any{
		"tenant": tenant,
		"tagId":  tagId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.GetAll")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag)
			RETURN tag ORDER BY tag.name`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), err
}

func (r *tagReadRepository) GetByNameOptional(ctx context.Context, tenant, name string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagReadRepository.GetByNameOptional")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:$name}) return tag limit 1`
	params := map[string]any{
		"name":   name,
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *tagReadRepository) GetForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagReadRepository.GetForIssues")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("issueIds", issueIds))

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
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeWithRelationAndId))))
	if len(result.([]*utils.DbNodeWithRelationAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagReadRepository) GetForLogEntries(ctx context.Context, tenant string, logEntryIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagReadRepository.GetForLogEntries")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("logEntryIds", logEntryIds))

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
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeWithRelationAndId))))
	if len(result.([]*utils.DbNodeWithRelationAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagReadRepository) GetForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagReadRepository.GetForContacts")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("contactIds", contactIds))

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
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeWithRelationAndId))))
	if len(result.([]*utils.DbNodeWithRelationAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagReadRepository) GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagReadRepository.GetForOrganizations")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("organizationIds", organizationIds))

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
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeWithRelationAndId))))
	if len(result.([]*utils.DbNodeWithRelationAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagReadRepository) GetForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.GetAll")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("contactId", contactId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})-[rel:TAGGED]->(tag:Tag)
			RETURN tag ORDER BY rel.taggedAt, tag.name`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), err
}
