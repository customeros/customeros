package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

var workflowTypeByModel = map[model.WorkflowType]postgresentity.WorkflowType{
	model.WorkflowTypeIdealContactPersona:  postgresentity.WorkflowTypeIdealContactPersona,
	model.WorkflowTypeIdealCustomerProfile: postgresentity.WorkflowTypeIdealCustomerProfile,
}

var workflowTypeByValue = utils.ReverseMap(workflowTypeByModel)

func MapWorkflowTypeToModel(input postgresentity.WorkflowType) model.WorkflowType {
	if v, exists := workflowTypeByValue[input]; exists {
		return v
	} else {
		return ""
	}
}

func MapWorkflowTypeFromModel(input model.WorkflowType) postgresentity.WorkflowType {
	if v, exists := workflowTypeByModel[input]; exists {
		return v
	} else {
		return ""
	}
}
