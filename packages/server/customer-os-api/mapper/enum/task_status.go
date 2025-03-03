package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

var taskStatusByModel = map[model.TaskStatus]enum.TaskStatus{
	model.TaskStatusDone:       enum.TaskStatusDone,
	model.TaskStatusInProgress: enum.TaskStatusInProgress,
	model.TaskStatusTodo:       enum.TaskStatusTodo,
}

var taskStatusByValue = utils.ReverseMap(taskStatusByModel)

func MapTaskStatusFromModel(input model.TaskStatus) enum.TaskStatus {
	return taskStatusByModel[input]
}

func MapTaskStatusToModel(input enum.TaskStatus) model.TaskStatus {
	return taskStatusByValue[input]
}
