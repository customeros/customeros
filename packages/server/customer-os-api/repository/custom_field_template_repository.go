package repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CustomFieldTemplateRepository.FindByCustomFieldId")
	defer spans.Finish()

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
