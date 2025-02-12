package agent_capability

import (
	"context"
	"fmt"
	"strings"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type AddMeetingNotesToCompanyCapability struct {
	organizationService  interfaces.OrganizationService
	markdownEventService interfaces.MarkdownEventService
}

func NewAddMeetingNotesToCompanyCapability(orgService interfaces.OrganizationService, mdEventService interfaces.MarkdownEventService) *AddMeetingNotesToCompanyCapability {
	return &AddMeetingNotesToCompanyCapability{
		organizationService:  orgService,
		markdownEventService: mdEventService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[AddMeetingNotesToCompanyInput, AddMeetingNotesToCompanyOutput, postgres_entity.NoConfig] = (*AddMeetingNotesToCompanyCapability)(nil)
)

func (c *AddMeetingNotesToCompanyCapability) Type() enum.AgentCapability {
	return enum.CapabilityAddMeetingNotesToCompany
}

func (c *AddMeetingNotesToCompanyCapability) Name() string {
	return "Add meeting notes to company timeline"
}

func (c *AddMeetingNotesToCompanyCapability) NewInput() AddMeetingNotesToCompanyInput {
	return AddMeetingNotesToCompanyInput{}
}

func (c *AddMeetingNotesToCompanyCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *AddMeetingNotesToCompanyCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *AddMeetingNotesToCompanyCapability) ValidateInput(input AddMeetingNotesToCompanyInput) error {
	switch {
	case input.MeetingTimestamp.IsZero():
		return errors.New("Meeting timestamp cannot be empty")
	case input.MeetingRecordingUrl == "":
		return errors.New("Meeting recording url cannot be empty")
	case len(input.MeetingParticipantEmails) == 0:
		return errors.New("Meeting participant emails cannot be empty")
	case len(input.MeetingParticipantEmailsTenant) == 0:
		return errors.New("Meeting participant emails tenant cannot be empty")
	case input.MeetingSummary == "":
		return errors.New("Meeging summary cannot be empty")
	}
	return nil
}

func (c *AddMeetingNotesToCompanyCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

type AddMeetingNotesToCompanyInput struct {
	MeetingSource                  string    `json:"meetingSource"`
	MeetingTitle                   string    `json:"MeetingTitle"`
	MeetingTimestamp               time.Time `json:"meetingTimestamp"`
	MeetingRecordingUrl            string    `json:"meetingRecordingUrl"`
	MeetingParticipantEmails       []string  `json:"meetingParticipantEmails"`
	MeetingParticipantEmailsTenant []string  `json:"meetingParticipantEmailsTenant"`
	MeetingSummary                 string    `json:"meetingSummary"`
	ActionItems                    []string  `json:"actionItems"`
}

type AddMeetingNotesToCompanyOutput struct {
	OrganizationIDs []string `json:"organizationIds"`
}

func (c *AddMeetingNotesToCompanyCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[AddMeetingNotesToCompanyInput, postgres_entity.NoConfig]) (bool, AddMeetingNotesToCompanyOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AddMeetingNotesToCompanyCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := AddMeetingNotesToCompanyOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return false, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return false, result, err
	}

	timelineEvent := c.createMarkdownTimelineEvent(ctx, executionContainer.InputData)

	created := make([]string, 0)
	var err error
	for _, email := range executionContainer.InputData.MeetingParticipantEmails {
		if !utils.IsStringInSlice(email, executionContainer.InputData.MeetingParticipantEmailsTenant) {
			created, err = c.processMeetingParticipant(ctx, email, timelineEvent, executionContainer.InputData, created)
			if err != nil {
				tracing.TraceErr(span, err)
			}
		}
	}

	tracing.LogObjectAsJson(span, "result", result)
	return true, result, nil
}

func (c *AddMeetingNotesToCompanyCapability) processMeetingParticipant(ctx context.Context, email, timelineEvent string, input AddMeetingNotesToCompanyInput, created []string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AddMeetingNotesToCompanyCapability.processMeetingParticipant")
	defer span.Finish()
	tracing.TagComponentService(span)

	// get org by email
	var (
		ok    bool
		orgID string
		err   error
	)
	ok, orgID, err = c.organizationService.CheckOrganizationExistsWithEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return created, err
	}
	if !ok {
		validation := mailvalidate.ValidateEmailSyntax(email)
		orgID, err = c.organizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
			Domains: []string{validation.Domain},
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return created, err
		}
	}

	if utils.IsStringInSlice(orgID, created) {
		return created, nil
	}

	// create timeline event
	source := enum.DecodeSource(input.MeetingSource)
	_, err = c.markdownEventService.Save(ctx, nil, nil, data_fields.MarkdownEventFields{
		Source:         &source,
		CreatedAt:      &input.MeetingTimestamp,
		OrganizationId: &orgID,
		Content:        &timelineEvent,
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return append(created, orgID), err
}

func (c *AddMeetingNotesToCompanyCapability) createMarkdownTimelineEvent(ctx context.Context, input AddMeetingNotesToCompanyInput) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AddMeetingNotesToCompanyCapability.createMarkdownTimelineEvent")
	defer span.Finish()
	tracing.TagComponentService(span)

	var b strings.Builder

	// Header
	b.WriteString(fmt.Sprintf("# %s\n", input.MeetingTitle))

	// Participants
	b.WriteString("## Who was there\n")
	for _, p := range input.MeetingParticipantEmails {
		b.WriteString(fmt.Sprintf("* %s\n", p))
	}
	b.WriteString("\n")

	// Summary
	b.WriteString("## In a nutshell\n")
	b.WriteString(input.MeetingSummary)
	b.WriteString("\n\n")

	// Action Items by Person
	if len(input.ActionItems) > 0 {
		b.WriteString("## Who's doing what\n")
		for _, p := range input.ActionItems {
			b.WriteString(fmt.Sprintf("* %s\n", p))
			b.WriteString("\n")
		}
	}

	// Recording
	if input.MeetingRecordingUrl != "" {
		b.WriteString(fmt.Sprintf("[Watch the recording](%s)", input.MeetingRecordingUrl))
	}

	return b.String()
}
