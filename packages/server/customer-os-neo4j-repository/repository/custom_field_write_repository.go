package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type CustomFieldCreateFields struct {
	CreatedAt           time.Time              `json:"createdAt"`
	ExistsInEventStore  bool                   `json:"existsInEventStore"`
	TemplateId          *string                `json:"templateId,omitempty"`
	CustomFieldId       string                 `json:"customFieldId"`
	CustomFieldName     string                 `json:"customFieldName"`
	CustomFieldDataType string                 `json:"customFieldDataType"`
	CustomFieldValue    model.CustomFieldValue `json:"customFieldValue"`
	SourceFields        model.SourceFields     `json:"sourceFields,omitempty"`
}

type CustomFieldWriteRepository interface {
	AddCustomFieldToOrganization(ctx context.Context, tenant, organizationId string, data CustomFieldCreateFields) error
}

type customFieldWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCustomFieldWriteRepository(driver *neo4j.DriverWithContext, database string) CustomFieldWriteRepository {
	return &customFieldWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *customFieldWriteRepository) AddCustomFieldToOrganization(ctx context.Context, tenant, organizationId string, data CustomFieldCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldWriteRepository.AddCustomFieldToOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("organizationId", organizationId))
	tracing.LogObjectAsJson(span, "data", data)

	nodeLabel := enum.NodeLabelForCustomFieldDataType(data.CustomFieldDataType)
	propertyName := enum.PropertyNameForCustomFieldDataType(data.CustomFieldDataType)

	cypher := fmt.Sprintf(`
		MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		MERGE (org)-[:HAS_PROPERTY]->(cf:CustomField:CustomField_%s {id:$customFieldId})
		ON CREATE SET
		    org.updatedAt = datetime(),
			cf.source=$source,
			cf.sourceOfTruth=$sourceOfTruth,
			cf.appSource=$appSource,
			cf.createdAt=$createdAt,
			cf.updatedAt=datetime(),
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
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"customFieldId":  data.CustomFieldId,
		"source":         data.SourceFields.Source,
		"sourceOfTruth":  data.SourceFields.SourceOfTruth,
		"appSource":      data.SourceFields.AppSource,
		"createdAt":      data.CreatedAt,
		"name":           data.CustomFieldName,
		"value":          data.CustomFieldValue.RealValue(),
		"templateId":     data.TemplateId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
