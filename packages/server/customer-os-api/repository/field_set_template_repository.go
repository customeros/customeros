package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type FieldSetTemplateRepository interface {
	createFieldSetTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, entityTemplateId string, entity *entity.FieldSetTemplateEntity) error
	FindAllByEntityTemplateId(ctx context.Context, entityTemplateId string) (any, error)
	FindByFieldSetId(ctx context.Context, fieldSetId string) (any, error)
}

type fieldSetTemplateRepository struct {
	driver       *neo4j.DriverWithContext
	repositories *Repositories
}

func NewFieldSetTemplateRepository(driver *neo4j.DriverWithContext, repositories *Repositories) FieldSetTemplateRepository {
	return &fieldSetTemplateRepository{
		driver:       driver,
		repositories: repositories,
	}
}

func (r *fieldSetTemplateRepository) createFieldSetTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, entityTemplateId string, entity *entity.FieldSetTemplateEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FieldSetTemplateRepository.createFieldSetTemplateInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (e:EntityTemplate {id:$entityTemplateId}) " +
		" MERGE (e)-[:CONTAINS]->(f:FieldSetTemplate {id:randomUUID(), name:$name}) " +
		" ON CREATE SET f:%s, " +
		"				f.order=$order, " +
		"				f.createdAt=$now, " +
		"				f.updatedAt=datetime() " +
		" RETURN f"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "FieldSetTemplate_"+tenant),
		map[string]any{
			"entityTemplateId": entityTemplateId,
			"name":             entity.Name,
			"order":            entity.Order,
			"now":              utils.Now(),
		})

	record, err := queryResult.Single(ctx)
	if err != nil {
		return err
	}
	fieldSetTemplateId := utils.GetPropsFromNode(record.Values[0].(dbtype.Node))["id"].(string)
	for _, v := range entity.CustomFields {
		err := r.repositories.CustomFieldTemplateRepository.createCustomFieldTemplateForFieldSetInTx(ctx, tx, tenant, fieldSetTemplateId, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *fieldSetTemplateRepository) FindAllByEntityTemplateId(ctx context.Context, entityTemplateId string) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FieldSetTemplateRepository.FindAllByEntityTemplateId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	return session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (:EntityTemplate {id:$entityTemplateId})-[:CONTAINS]->(f:FieldSetTemplate) RETURN f ORDER BY f.order`,
			map[string]any{
				"entityTemplateId": entityTemplateId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
}

func (r *fieldSetTemplateRepository) FindByFieldSetId(ctx context.Context, fieldSetId string) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FieldSetTemplateRepository.FindByFieldSetId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	return session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (:FieldSet {id:$fieldSetId})-[:IS_DEFINED_BY]->(d:FieldSetTemplate)
					RETURN d`,
			map[string]any{
				"fieldSetId": fieldSetId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
}
