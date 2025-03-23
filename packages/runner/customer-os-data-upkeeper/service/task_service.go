package service

import (
	"context"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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

	span, ctx := tracing.StartTracerSpan(ctx, "TaskService.DeleteArchivedTasks")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	delitionDelayDays := 30

	// Get ids of tasks to be deleted
	taskIds, err := s.commonServices.Neo4jRepositories.TaskReadRepository.GetHiddenTasks(ctx, delitionDelayDays)
	if err != nil {
		s.log.Error("Error getting archived tasks", "error", err)
		return
	}

	span.LogKV("taskIds.count", len(taskIds))
	tracing.LogObjectAsJson(span, "taskIds", taskIds)
	s.log.Infof("Deleting %d archived tasks: %v", len(taskIds), taskIds)

	// Delete tasks
	err = s.commonServices.Neo4jRepositories.TaskWriteRepository.PermanentDeleteHiddenTasks(ctx, nil, taskIds)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error("Error deleting tasks", "error", err)
		return
	}

	s.log.Infof("Deleted %d archived tasks: %v", len(taskIds), taskIds)
}
