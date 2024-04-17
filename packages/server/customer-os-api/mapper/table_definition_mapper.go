package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"strconv"
	"strings"
)

func MapTableViewDefinitionToModel(entity postgresEntity.TableViewDefinition) *model.TableViewDef {
	columnNames := strings.Split(entity.Columns, ",")
	columns := make([]model.ColumnViewType, 0, len(columnNames))
	for _, column := range columnNames {
		columns = append(columns, model.ColumnViewType(column))
	}
	return &model.TableViewDef{
		ID:        strconv.Itoa(int(entity.ID)),
		Name:      entity.Name,
		TableType: model.TableViewType(entity.TableType),
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
