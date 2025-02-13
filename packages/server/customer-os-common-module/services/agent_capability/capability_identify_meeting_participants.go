package agent_capability

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type IdentifyMeetingParticipantsCapability struct {
	workspaceService interfaces.WorkspaceService
}

type IdentifyMeetingParticipantsInput struct {
	MeetingParticipantEmails []string `json:"meetingParticipantEmails"`
}

type IdentifyMeetingParticipantsOutput struct {
	MeetingParticipantEmailsTenant []string `json:"meetingParticipantEmailsTenant"`
	ContactEmails                  []string `json:"contactEmails"`
}

func NewIdentifyMeetingParticipantsCapability(workspaceService interfaces.WorkspaceService) *IdentifyMeetingParticipantsCapability {
	return &IdentifyMeetingParticipantsCapability{
		workspaceService: workspaceService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[IdentifyMeetingParticipantsInput, IdentifyMeetingParticipantsOutput, postgres_entity.NoConfig] = (*IdentifyMeetingParticipantsCapability)(nil)
)

func (c *IdentifyMeetingParticipantsCapability) Type() enum.AgentCapability {
	return enum.CapabilityIdentifyMeetingParticipants
}

func (c *IdentifyMeetingParticipantsCapability) Name() string {
	return "Identify meeting participants"
}

func (c *IdentifyMeetingParticipantsCapability) NewInput() IdentifyMeetingParticipantsInput {
	return IdentifyMeetingParticipantsInput{}
}

func (c *IdentifyMeetingParticipantsCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *IdentifyMeetingParticipantsCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *IdentifyMeetingParticipantsCapability) ValidateConfig(config postgres_entity.NoConfig) error {
	return nil
}

func (c *IdentifyMeetingParticipantsCapability) ValidateInput(input IdentifyMeetingParticipantsInput) error {
	if len(input.MeetingParticipantEmails) == 0 {
		return errors.New("MeetingParticipantEmails cannot be empty")
	}
	return nil
}

func (c *IdentifyMeetingParticipantsCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[IdentifyMeetingParticipantsInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, IdentifyMeetingParticipantsOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyMeetingParticipantsCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := IdentifyMeetingParticipantsOutput{}

	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}

	// check each email to determine if they belong to tenant or external participant
	workspaceDomains, err := c.workspaceService.GetWorkspaceDomainsForTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	for _, email := range executionContainer.InputData.MeetingParticipantEmails {
		if utils.IsStringInSlice(email, workspaceDomains) {
			result.MeetingParticipantEmailsTenant = append(result.MeetingParticipantEmailsTenant, email)
		} else {
			result.ContactEmails = append(result.ContactEmails, email)
		}
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
