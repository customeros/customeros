package agent_listeners

import (
	"context"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type NewMeetingRecordingListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
	webhookService       interfaces.WebhookService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*NewMeetingRecordingListener)(nil)
	_ interfaces.EventListener        = (*NewMeetingRecordingListener)(nil)
)

func NewNewMeetingRecordingListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *NewMeetingRecordingListener {
	return &NewMeetingRecordingListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.NewMeetingRecording](), // subscribed event
			events.QueueAgents, // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

type NewMeetingRecordingListenerConfig struct {
	MeetingSource           MeetingSourceConfig `json:"meetingSource"`
	MeetingSourceWebhookUrl ConfigSingleValue   `json:"meetingSourceWebhookUrl"`
}

type MeetingSourceConfig struct {
	Value enum.Source `json:"value"`
	Error string      `json:"error"`
}

func (l *NewMeetingRecordingListenerConfig) Validate() bool {
	switch {
	case l.MeetingSource.Value == enum.SourceFathom:
		return true
	case l.MeetingSource.Value == enum.SourceGrain:
		return true
	}
	return false
}

func (l *NewMeetingRecordingListener) Type() enum.AgentListenerEvent {
	return enum.EventNewMeetingRecording
}

func (l *NewMeetingRecordingListener) Name() string {
	return "New meeting recordings"
}

func (l *NewMeetingRecordingListener) DefaultConfig() any {
	return &NewMeetingRecordingListenerConfig{
		MeetingSource: MeetingSourceConfig{
			Value: enum.SourceUnknown,
		},
	}
}

func (l *NewMeetingRecordingListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentMeetingKeeper,
		enum.AgentSupportSpotter,
	}
}

func (l *NewMeetingRecordingListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewMeetingRecordingListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if common.GetUserEmailFromContext(ctx) == "" {
		err := coserrors.ErrUserEmailNotSet
		tracing.TraceErr(span, err)
		return err
	}

	data, err := events.DecodeEventData[dto.NewMeetingRecording](ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if data.Source == enum.SourceUnknown {
		err := errors.New("Meeting source cannot be unknown")
		tracing.TraceErr(span, err)
		return err
	}

	if data.Content == nil {
		err := errors.New("No meeting content")
		tracing.TraceErr(span, err)
		return err
	}

	if data.ParticipantEmails == nil {
		err := errors.New("No meeting participants")
		tracing.TraceErr(span, err)
		return err
	}

	return l.handleExecution(ctx, data)
}

func (l *NewMeetingRecordingListener) handleExecution(ctx context.Context, data dto.NewMeetingRecording) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewMeetingRecordingListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	activeAgents := l.lookupActiveAgentsForUser(ctx)
	if activeAgents == nil || len(activeAgents) == 0 {
		return nil
	}

	var errs error
	for _, agent := range activeAgents {
		// map event data to execution input
		initialParams, err := utils.StructToMap(data)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}

		// run
		_, err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *NewMeetingRecordingListener) lookupActiveAgentsForUser(ctx context.Context) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewMeetingRecordingListener.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByUserAndType(ctx, l.ExecutingAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}

	return agents
}
