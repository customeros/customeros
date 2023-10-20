package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CustomFieldRepository interface {
	AddCustomFieldToOrganization(ctx context.Context, tenant, organizationId string, eventData events.OrganizationUpsertCustomField) error
	ExistsById(ctx context.Context, tenant, customFieldId string) (bool, error)
}

type customFieldRepository struct {
	driver *neo4j.DriverWithContext
}

func NewCustomFieldRepository(driver *neo4j.DriverWithContext) CustomFieldRepository {
	return &customFieldRepository{
		driver: driver,
	}
}

func (r *customFieldRepository) AddCustomFieldToOrganization(ctx context.Context, tenant, organizationId string, eventData events.OrganizationUpsertCustomField) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.AddCustomFieldToOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId))

	nodeLabel := entity.NodeLabelForCustomFieldDataType(eventData.CustomFieldDataType)
	propertyName := entity.PropertyNameForCustomFieldDataType(eventData.CustomFieldDataType)

	query := fmt.Sprintf(`
		MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		MERGE (org)-[:HAS_PROPERTY]->(cf:CustomField:CustomField_%s {id:$customFieldId})
		ON CREATE SET
			cf.source=$source,
			cf.sourceOfTruth=$sourceOfTruth,
			cf.appSource=$appSource,
			cf.createdAt=$createdAt,
			cf.updatedAt=$updatedAt,
			cf.name=$name,
			cf.%s=$value,
			cf:%s,
			cf:%s
		WITH t, cf
			OPTIONAL MATCH (t)<-[:CUSTOM_FIELD_TEMPLATE_BELONGS_TO_TENANT]-(cft:CustomFieldTemplate {id:$templateId}) 
				WHERE $templateId <> "" AND NOT $templateId IS NULL
				FOREACH (ignore IN CASE WHEN cft IS NOT NULL THEN [1] ELSE [] END |
    				MERGE (cf)-[:IS_DEFINED_BY]->(cft))
			 `, tenant, propertyName, nodeLabel, nodeLabel+"_"+tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"customFieldId":  eventData.CustomFieldId,
		"source":         eventData.Source,
		"sourceOfTruth":  eventData.SourceOfTruth,
		"appSource":      eventData.AppSource,
		"createdAt":      eventData.CreatedAt,
		"updatedAt":      eventData.UpdatedAt,
		"name":           eventData.CustomFieldName,
		"value":          eventData.CustomFieldValue.RealValue(),
		"templateId":     eventData.TemplateId,
	})
}

func (r *customFieldRepository) ExistsById(ctx context.Context, tenant, customFieldId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.ExistsById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("customFieldId", customFieldId))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (cf:CustomField_%s {id:$customFieldId}) RETURN cf LIMIT 1`, tenant)
	span.LogFields(log.String("query", query))

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"customFieldId": customFieldId,
			}); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil
		}
	})
	if err != nil {
		return false, err
	}
	return result.(bool), err
}

// Common database interaction method
func (r *customFieldRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}
