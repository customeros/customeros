package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TagRepository interface {
	Merge(ctx context.Context, tenant string, tag entity.TagEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tenant string, tag entity.TagEntity) (*dbtype.Node, error)
	UnlinkAndDelete(ctx context.Context, tenant string, tagId string) error
	GetAll(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error)
	GetForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetForLogEntries(ctx context.Context, tenant string, logEntryIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetById(ctx context.Context, tagId string) (*dbtype.Node, error)
	GetByName(ctx context.Context, tagName string) (*dbtype.Node, error)
}

type tagRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTagRepository(driver *neo4j.DriverWithContext) TagRepository {
	return &tagRepository{
		driver: driver,
	}
}

func (r *tagRepository) Merge(ctx context.Context, tenant string, tag entity.TagEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:$name}) 
		 ON CREATE SET 
		  tag.id=randomUUID(),
		  tag.createdAt=$now,
		  tag.updatedAt=$now,
		  tag.source=$source,
		  tag.sourceOfTruth=$sourceOfTruth,
		  tag.appSource=$appSource,
		  tag:Tag_%s
		 RETURN tag`, tenant)

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":        tenant,
				"name":          tag.Name,
				"source":        tag.Source,
				"sourceOfTruth": tag.SourceOfTruth,
				"appSource":     tag.AppSource,
				"now":           utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tagRepository) Update(ctx context.Context, tenant string, tag entity.TagEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.Update")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
			SET tag.name=$name, tag.updatedAt=$now
			RETURN tag`,
			map[string]any{
				"tenant": tenant,
				"id":     tag.Id,
				"name":   tag.Name,
				"now":    utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tagRepository) UnlinkAndDelete(ctx context.Context, tenant string, tagId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.UnlinkAndDelete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
			DETACH DELETE tag`,
			map[string]any{
				"tenant": tenant,
				"id":     tagId,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *tagRepository) GetAll(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.GetAll")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag)
			RETURN tag ORDER BY tag.name`,
			map[string]any{
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}

func (r *tagRepository) GetForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.GetForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})-[rel:TAGGED]->(tag:Tag)
			RETURN tag ORDER BY rel.taggedAt, tag.name`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}

func (r *tagRepository) GetForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.GetForContacts")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[rel:TAGGED]->(tag:Tag)-[:TAG_BELONGS_TO_TENANT]->(t)
			WHERE c.id IN $contactIds
			RETURN tag, rel, c.id ORDER BY rel.taggedAt, tag.name`,
			map[string]any{
				"tenant":     tenant,
				"contactIds": contactIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagRepository) GetForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.GetForIssues")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)-[rel:TAGGED]->(tag:Tag)-[:TAG_BELONGS_TO_TENANT]->(t)
			WHERE i.id IN $issueIds
			RETURN tag, rel, i.id ORDER BY rel.taggedAt, tag.name`,
			map[string]any{
				"tenant":   tenant,
				"issueIds": issueIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagRepository) GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.GetForOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[rel:TAGGED]->(tag:Tag)-[:TAG_BELONGS_TO_TENANT]->(t)
			WHERE o.id IN $organizationIds
			RETURN tag, rel, o.id ORDER BY rel.taggedAt, tag.name`,
			map[string]any{
				"tenant":          tenant,
				"organizationIds": organizationIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagRepository) GetForLogEntries(ctx context.Context, tenant string, logEntryIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.GetForLogEntries")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (l:LogEntry_%s)-[rel:TAGGED]->(tag:Tag)-[:TAG_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			WHERE l.id IN $logEntryIds
			RETURN tag, rel, l.id ORDER BY rel.taggedAt, tag.name`, tenant)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":      tenant,
				"logEntryIds": logEntryIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *tagRepository) GetById(ctx context.Context, tagId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("tagId", tagId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id}) return tag`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": common.GetTenantFromContext(ctx),
				"id":     tagId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tagRepository) GetByName(ctx context.Context, tagName string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryRepository.GetByName")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("tagName", tagName))

	query := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:$name}) return tag limit 1`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": common.GetTenantFromContext(ctx),
				"name":   tagName,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
