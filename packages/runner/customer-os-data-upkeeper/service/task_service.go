package service

import (
	"context"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
)

type TaskService interface {
	DeleteArchivedTasks()
}

type taskService struct {
	log            logger.Logger
	commonServices *commonservice.CommonServices
}

func NewTaskService(log logger.Logger, commonServices *commonservice.CommonServices) TaskService {
	return &taskService{
		log:            log,
		commonServices: commonServices,
	}
}

func (s *taskService) DeleteArchivedTasks() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	spans, ctx := telemetry.StartCronSpan(ctx, "TaskService.DeleteArchivedTasks")
	defer spans.Finish()

	delitionDelayDays := 30

	// Get ids of tasks to be deleted
	taskIds, err := s.commonServices.Neo4jRepositories.TaskReadRepository.GetHiddenTasks(ctx, delitionDelayDays)
	if err != nil {
		s.log.Error("Error getting archived tasks", "error", err)
		return
	}

	spans.LogKV("taskIds.count", len(taskIds))
	spans.LogObjectAsJson("taskIds", taskIds)
	s.log.Infof("Deleting %d archived tasks: %v", len(taskIds), taskIds)

	// Delete tasks
	err = s.commonServices.Neo4jRepositories.TaskWriteRepository.PermanentDeleteHiddenTasks(ctx, nil, taskIds)
	if err != nil {
		spans.TraceError(err)
		s.log.Error("Error deleting tasks", "error", err)
		return
	}

	s.log.Infof("Deleted %d archived tasks: %v", len(taskIds), taskIds)
}
