package aggregate

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *IssueAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	switch c := cmd.(type) {
	case *command.UpsertIssueCommand:
		if c.IsCreateCommand {
			return a.createIssue(ctx, c)
		} else {
			return a.updateIssue(ctx, c)
		}
	case *command.AddUserAssigneeCommand:
		return a.addUserAssignee(ctx, c)
	case *command.RemoveUserAssigneeCommand:
		return a.removeUserAssignee(ctx, c)
	case *command.AddUserFollowerCommand:
		return a.addUserFollower(ctx, c)
	case *command.RemoveUserFollowerCommand:
		return a.removeUserFollower(ctx, c)
	default:
		return errors.New("invalid command type")
	}
}

func (a *IssueAggregate) createIssue(ctx context.Context, cmd *command.UpsertIssueCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "IssueAggregate.createIssue")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	createEvent, err := event.NewIssueCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewIssueCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.Metadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *IssueAggregate) updateIssue(ctx context.Context, cmd *command.UpsertIssueCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "IssueAggregate.updateIssue")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	source := cmd.Source.Source
	if source == "" {
		source = a.Issue.Source.Source
	}

	updateEvent, err := event.NewIssueUpdateEvent(a, cmd.DataFields, cmd.Source.Source, cmd.ExternalSystem, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewIssueUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.Metadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *IssueAggregate) addUserAssignee(ctx context.Context, cmd *command.AddUserAssigneeCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "IssueAggregate.addUserAssignee")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	atNotNil := utils.IfNotNilTimeWithDefault(cmd.At, utils.Now())

	addUserAssigneeEvent, err := event.NewIssueAddUserAssigneeEvent(a, cmd.AssigneeId, atNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewIssueAddUserAssigneeEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&addUserAssigneeEvent, span, aggregate.Metadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(addUserAssigneeEvent)
}

func (a *IssueAggregate) removeUserAssignee(ctx context.Context, cmd *command.RemoveUserAssigneeCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "IssueAggregate.removeUserAssignee")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	atNotNil := utils.IfNotNilTimeWithDefault(cmd.At, utils.Now())

	removeUserAssigneeEvent, err := event.NewIssueRemoveUserAssigneeEvent(a, cmd.AssigneeId, atNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewIssueRemoveUserAssigneeEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&removeUserAssigneeEvent, span, aggregate.Metadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(removeUserAssigneeEvent)
}

func (a *IssueAggregate) addUserFollower(ctx context.Context, cmd *command.AddUserFollowerCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "IssueAggregate.addUserFollower")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	atNotNil := utils.IfNotNilTimeWithDefault(cmd.At, utils.Now())

	addUserFollowerEvent, err := event.NewIssueAddUserFollowerEvent(a, cmd.FollowerId, atNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewIssueAddUserFollowerEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&addUserFollowerEvent, span, aggregate.Metadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(addUserFollowerEvent)
}

func (a *IssueAggregate) removeUserFollower(ctx context.Context, cmd *command.RemoveUserFollowerCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "IssueAggregate.removeUserFollower")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	atNotNil := utils.IfNotNilTimeWithDefault(cmd.At, utils.Now())

	removeUserFollowerEvent, err := event.NewIssueRemoveUserFollowerEvent(a, cmd.FollowerId, atNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewIssueRemoveUserFollowerEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&removeUserFollowerEvent, span, aggregate.Metadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(removeUserFollowerEvent)
}
