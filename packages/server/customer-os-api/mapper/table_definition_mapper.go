package mapper

import (
	"encoding/json"
	"strconv"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func MapTableViewDefinitionToModel(entity postgresEntity.TableViewDefinition, span opentracing.Span) *model.TableViewDef {
	var columnsStruct postgresEntity.Columns
	err := json.Unmarshal([]byte(entity.ColumnsJson), &columnsStruct)
	if err != nil {
		span.LogFields(log.String("columnsJson", entity.ColumnsJson))
		tracing.TraceErr(span, err)
	}

	columns := make([]*model.ColumnView, 0, len(columnsStruct.Columns))
	for _, column := range columnsStruct.Columns {
		columns = append(columns, &model.ColumnView{
			ColumnID:   column.ColumnId,
			ColumnType: model.ColumnViewType(column.ColumnType),
			Width:      column.Width,
			Visible:    column.Visible,
			Name:       column.Name,
			Filter:     column.Filter,
		})
	}
	return &model.TableViewDef{
		ID:             strconv.Itoa(int(entity.ID)),
		Name:           entity.Name,
		TableType:      model.TableViewType(entity.TableType),
		TableID:        model.TableIDType(entity.TableId),
		Icon:           entity.Icon,
		Order:          entity.Order,
		Filters:        entity.Filters,
		DefaultFilters: entity.DefaultFilters,
		Sorting:        entity.Sorting,
		IsPreset:       entity.IsPreset,
		IsShared:       entity.IsShared,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		Columns:        columns,
	}
}

func MapTableViewDefinitionsToModel(entities []postgresEntity.TableViewDefinition, span opentracing.Span) []*model.TableViewDef {
	var tableViewDefs []*model.TableViewDef
	for _, entity := range entities {
		tableViewDefs = append(tableViewDefs, MapTableViewDefinitionToModel(entity, span))
	}
	return tableViewDefs
}
