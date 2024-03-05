package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/command/base"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type CreateJobRoleCommandHander interface {
	Handle(ctx context.Context, command *model.CreateJobRoleCommand) error
}

type createJobRoleCommandHandler struct {
	base.BaseCommandHandler
}

func NewCreateJobRoleCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *createJobRoleCommandHandler {
	handler := createJobRoleCommandHandler{}
	handler.BaseCommandHandler = *base.NewBaseCommandHandler(log, cfg, es)
	return &handler
}

func (c *createJobRoleCommandHandler) Handle(ctx context.Context, command *model.CreateJobRoleCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createJobRoleCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	jobRoleAggregate := aggregate.NewJobRoleAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := jobRoleAggregate.CreateJobRole(ctx, command)
	if err != nil {
		return errors.Wrap(err, "createJobRoleCommandHandler.Handle")
	}
	return c.BaseCommandHandler.Es.Save(ctx, jobRoleAggregate)

}
