package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

// HandleCommand processes commands and applies the resulting events to the aggregate.
func (a *ContractAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.CreateContractCommand:
		return a.createContract(ctx, c)
	case *command.UpdateContractCommand:
		return a.updateContract(ctx, c)
	case *command.RefreshContractStatusCommand:
		return a.refreshContractStatus(ctx, c)
	case *command.RolloutRenewalOpportunityOnExpirationCommand:
		return a.rolloutRenewalOpportunityOnExpiration(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *ContractAggregate) createContract(ctx context.Context, cmd *command.CreateContractCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.createContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	// Assuming you have a utility function to get the current time if the passed time is nil
	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	// Determine contract status based start and end dates
	status := determineContractStatus(cmd.DataFields.ServiceStartedAt, cmd.DataFields.EndedAt)
	cmd.DataFields.Status = status

	createEvent, err := event.NewContractCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *ContractAggregate) updateContract(ctx context.Context, cmd *command.UpdateContractCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.updateContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	// Validate the dates
	if cmd.DataFields.EndedAt != nil && (cmd.DataFields.SignedAt != nil && cmd.DataFields.EndedAt.Before(*cmd.DataFields.SignedAt) ||
		cmd.DataFields.ServiceStartedAt != nil && cmd.DataFields.EndedAt.Before(*cmd.DataFields.ServiceStartedAt)) {
		return errors.New(constants.FieldValidation + ": endedAt date must be after both signedAt and serviceStartedAt dates")
	}

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	source := utils.StringFirstNonEmpty(cmd.Source.Source, a.Contract.Source.SourceOfTruth)

	// Determine contract status based start and end dates
	status := determineContractStatus(cmd.DataFields.ServiceStartedAt, cmd.DataFields.EndedAt)
	cmd.DataFields.Status = status

	// Determine contract renewal
	if cmd.DataFields.RenewalCycle == model.None {
		cmd.DataFields.RenewalCycle = model.RenewalCycleFromString(a.Contract.RenewalCycle)
	}

	updateEvent, err := event.NewContractUpdateEvent(
		a,
		cmd.DataFields,
		cmd.ExternalSystem,
		source,
		updatedAtNotNil,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *ContractAggregate) refreshContractStatus(ctx context.Context, cmd *command.RefreshContractStatusCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.refreshContractStatus")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	// Determine contract status based start and end dates
	status := determineContractStatus(a.Contract.ServiceStartedAt, a.Contract.EndedAt)

	updateEvent, err := event.NewContractUpdateStatusEvent(a, status.String(), a.Contract.ServiceStartedAt, a.Contract.EndedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractUpdateStatusEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *ContractAggregate) rolloutRenewalOpportunityOnExpiration(ctx context.Context, cmd *command.RolloutRenewalOpportunityOnExpirationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.rolloutRenewalOpportunityOnExpiration")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	// Determine contract status based start and end dates
	updateEvent, err := event.NewRolloutRenewalOpportunityEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewRolloutRenewalOpportunityEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(updateEvent)
}

func determineContractStatus(serviceStartedAt, endedAt *time.Time) model.ContractStatus {
	now := utils.Now()

	// If endedAt is not nil and is in the past, the contract is considered Ended.
	if endedAt != nil && endedAt.Before(now) {
		return model.Ended
	}

	// If serviceStartedAt is nil or in the future, the contract is considered Draft.
	if serviceStartedAt == nil || serviceStartedAt.After(now) {
		return model.Draft
	}

	// Otherwise, the contract is considered Live.
	return model.Live
}
