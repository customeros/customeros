package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CustomFieldTemplateSaveFields struct {
	Name        string           `json:"name"`
	EntityType  model.EntityType `json:"entityType"`
	Type        string           `json:"type"`
	ValidValues []string         `json:"validValues"`
	Order       *int64           `json:"order,omitempty"`
	Required    *bool            `json:"required,omitempty"`
	Length      *int64           `json:"length,omitempty"`
	Min         *int64           `json:"min,omitempty"`
	Max         *int64           `json:"max,omitempty"`

	UpdateName        bool `json:"updateName"`
	UpdateType        bool `json:"updateType"`
	UpdateValidValues bool `json:"updateValidValues"`
	UpdateOrder       bool `json:"updateOrder"`
	UpdateRequired    bool `json:"updateRequired"`
	UpdateLength      bool `json:"updateLength"`
	UpdateMin         bool `json:"updateMin"`
	UpdateMax         bool `json:"updateMax"`
}

type CustomFieldTemplateWriteRepository interface {
	Save(ctx context.Context, tenant, customFieldTemplateId string, data CustomFieldTemplateSaveFields) error
	Delete(ctx context.Context, tenant, customFieldTemplateId string) error
}

type customFieldTemplateWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCustomFieldTemplateWriteRepository(driver *neo4j.DriverWithContext, database string) CustomFieldTemplateWriteRepository {
	return &customFieldTemplateWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *customFieldTemplateWriteRepository) Save(ctx context.Context, tenant, customFieldTemplateId string, data CustomFieldTemplateSaveFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateWriteRepository.Save")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`
		MATCH (t:Tenant {name:$tenant})
		MERGE (t)<-[:CUSTOM_FIELD_TEMPLATE_BELONGS_TO_TENANT]-(cft:CustomFieldTemplate: {id:$customFieldTemplateId})
		ON CREATE SET
			cft:CustomFieldTemplate_%s,
			cft.createdAt=datetime(),
			cft.entityType=$entityType
		WITH cft
			cft.updatedAt=datetime()
		 `, tenant)
	params := map[string]any{
		"tenant":                tenant,
		"customFieldTemplateId": customFieldTemplateId,
		"entityType":            data.EntityType.String(),
	}

	if data.UpdateName {
		cypher += ", cft.name=$name"
		params["name"] = data.Name
	}
	if data.UpdateType {
		cypher += ", cft.type=$type"
		params["type"] = data.Type
	}
	if data.UpdateValidValues {
		cypher += ", cft.validValues=$validValues"
		params["validValues"] = data.ValidValues
	}
	if data.UpdateOrder {
		cypher += ", cft.order=$order"
		params["order"] = data.Order
	}
	if data.UpdateRequired {
		cypher += ", cft.required=$required"
		params["required"] = data.Required
	}
	if data.UpdateLength {
		cypher += ", cft.length=$length"
		params["length"] = data.Length
	}
	if data.UpdateMin {
		cypher += ", cft.min=$min"
		params["min"] = data.Min
	}
	if data.UpdateMax {
		cypher += ", cft.max=$max"
		params["max"] = data.Max
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *customFieldTemplateWriteRepository) Delete(ctx context.Context, tenant, customFieldTemplateId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateWriteRepository.Delete")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, customFieldTemplateId)

	cypher := fmt.Sprintf(`
		MATCH (t:Tenant {name:$tenant})<-[rel:CUSTOM_FIELD_TEMPLATE_BELONGS_TO_TENANT]-(cft:CustomFieldTemplate {id:$customFieldTemplateId})
		DELETE rel, cft
	`)
	params := map[string]any{
		"tenant":                tenant,
		"customFieldTemplateId": customFieldTemplateId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
