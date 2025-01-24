package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	enummapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

func MapAgentToModel(entity *postgresEntity.Agents) *model.Agent {
	if entity == nil {
		return nil
	}
	//var columnsStruct postgresEntity.Columns
	//err := json.Unmarshal([]byte(entity.ColumnsJson), &columnsStruct)
	//if err != nil {
	//	span.LogFields(log.String("columnsJson", entity.ColumnsJson))
	//	tracing.TraceErr(span, err)
	//}
	//
	//columns := make([]*model.ColumnView, 0, len(columnsStruct.Columns))
	//for _, column := range columnsStruct.Columns {
	//	columns = append(columns, &model.ColumnView{
	//		ColumnID:   column.ColumnId,
	//		ColumnType: postgresEntity.ColumnViewType(column.ColumnType),
	//		Width:      column.Width,
	//		Visible:    column.Visible,
	//		Name:       column.Name,
	//		Filter:     column.Filter,
	//	})
	//}
	return &model.Agent{
		ID:        entity.ID,
		Name:      entity.Name,
		Icon:      entity.Icon,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: utils.IfNotNilTimeWithDefault(entity.UpdatedAt, entity.CreatedAt),
		Type:      enummapper.MapAgentTypeToModel(entity.Type),
		Color:     entity.Color,
		Goal:      entity.Goal,
		IsActive:  entity.IsActive,
		Visible:   entity.VisibleInUI,
		//Columns:   columns, TODO Capabilities
	}
}

func MapAgentsToModel(entities []postgresEntity.Agents) []*model.Agent {
	var agents []*model.Agent
	for _, entity := range entities {
		agents = append(agents, MapAgentToModel(entity))
	}
	return agents
}
