package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

// TODO deprecate and remove all methods
type CustomFieldTemplateRepository interface {
	FindByCustomFieldId(ctx context.Context, fieldId string) (any, error)
}

type customFieldTemplateRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCustomFieldTemplateRepository(driver *neo4j.DriverWithContext, database string) CustomFieldTemplateRepository {
	return &customFieldTemplateRepository{
		driver:   driver,
		database: database,
	}
}

func (r *customFieldTemplateRepository) FindByCustomFieldId(ctx context.Context, customFieldId string) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateRepository.FindByCustomFieldId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	return session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (:CustomField {id:$customFieldId})-[:IS_DEFINED_BY]->(d:CustomFieldTemplate)
					RETURN d`,
			map[string]any{
				"customFieldId": customFieldId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
}
