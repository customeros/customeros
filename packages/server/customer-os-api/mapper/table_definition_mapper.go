package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"strconv"
)

func MapTableViewDefinitionToModel(entity postgresEntity.TableViewDefinition) *model.TableViewDef {
	var columnsStruct postgresEntity.Columns
	err := json.Unmarshal([]byte(entity.ColumnsJson), &columnsStruct)
	if err != nil {
		fmt.Println(err)
		return &model.TableViewDef{}
	}

	columns := make([]*model.ColumnView, 0, len(columnsStruct.Columns))
	for _, column := range columnsStruct.Columns {
		columns = append(columns, &model.ColumnView{
			ColumnType: model.ColumnViewType(column.ColumnType),
			Width:      column.Width,
			Visible:    column.Visible,
		})
	}
	return &model.TableViewDef{
		ID:        strconv.Itoa(int(entity.ID)),
		Name:      entity.Name,
		TableType: model.TableViewType(entity.TableType),
		TableID:   model.TableIDType(entity.TableId),
		Icon:      entity.Icon,
		Order:     entity.Order,
		Filters:   entity.Filters,
		Sorting:   entity.Sorting,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Columns:   columns,
	}
}

func MapTableViewDefinitionsToModel(entities []postgresEntity.TableViewDefinition) []*model.TableViewDef {
	var tableViewDefs []*model.TableViewDef
	for _, entity := range entities {
		tableViewDefs = append(tableViewDefs, MapTableViewDefinitionToModel(entity))
	}
	return tableViewDefs
}
