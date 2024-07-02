package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"strconv"
)

func MapWorkflowToModel(entity postgresEntity.Workflow) *model.Workflow {
	workflow := model.Workflow{
		ID:        strconv.Itoa(int(entity.ID)),
		Name:      utils.StringPtrNillable(entity.Name),
		Live:      entity.Live,
		Condition: entity.Condition,
		Type:      mapper.MapWorkflowTypeToModel(entity.WorkflowType),
	}
	return &workflow
}
