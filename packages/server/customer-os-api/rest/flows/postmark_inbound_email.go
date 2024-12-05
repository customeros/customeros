package flows

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
)

const EXTERNAL_SYSTEM = "mailstack"

func PostmarkInboundEmail(c *gin.Context, s *service.Services) {
	_, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.PostmarkInboundEmail", c.Request.Header)
	defer span.Finish()
	tracing.TagComponentRest(span)

	// Validate Postmark User-Agent
	if c.Request.UserAgent() == "" || !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
		tracing.TraceErr(span, fmt.Errorf("Invalid user agent %s", c.Request.UserAgent()))
		rest.SendError(c, span, http.StatusForbidden, rest.ErrForbidden)
		return
	}

	// Parse email data
	emailData, err := parseInboundEmail(c)
	if err != nil {
		tracing.LogObjectAsJson(span, "body", c.Request.Body)
		tracing.TraceErr(span, err)
		rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
		return
	}

	// Process email asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				err := fmt.Errorf("panic recovered in email processing: %v\n%s", r, stack)
				tracing.TraceErr(span, err)
			}
		}()

		if err := processInboundEmail(c, s, &emailData); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to process inbound email from Postmark"))
		}
	}()
}

func parseInboundEmail(c *gin.Context) (PostmarkInboundEmailData, error) {
	var emailData PostmarkInboundEmailData
	err := c.BindJSON(&emailData)
	if err != nil {
		return emailData, err
	}

	return emailData, nil
}

func processInboundEmail(c *gin.Context, s *service.Services, emailData *PostmarkInboundEmailData) error {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.processInboundEmail")
	defer span.Finish()
	tracing.TagComponentRest(span)

	tenant, err := getTenantFromEmail(c, s, emailData)
	if err != nil {
		return err
	}

	// update context with tenant, pass this where tenant is needed
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceCustomerOsApiRest,
	})

	participants := emailData.AllParticipantEmails()
	if len(participants) == 0 {
		err := errors.New("no email participants")
		tracing.TraceErr(span, err)
		return err
	}

	username, err := getUsername(ctx, s, participants)
	if err != nil || username == "" {
		span.LogFields(tracingLog.Bool("mailbox.found", false))
		tracing.TraceErr(span, err)
		return err
	}

	messageId := emailData.GetHeaderValue("Message-Id")
	if messageId == "" {
		err := errors.New("email messageId is empty")
		tracing.TraceErr(span, err)
		return err
	}

	// probably don't need this if we drop db and go events
	emailExistsInDb, err := s.CommonServices.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(
		ctx, EXTERNAL_SYSTEM, tenant, username, messageId,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("Unable to determine if email exists in db for messageId %s: %v", messageId, err)
	}

	if emailExistsInDb {
		return nil
	}

	emailMessage := emailData.ToEmailMessageData()

	emailAnalysis := s.MailService.ProcessEmailCheck(ctx, tenant, &emailMessage)
	if emailAnalysis.IsBulkMail {
		return nil
	}

	return publishEmailEvents(ctx, &emailMessage)
}

// Filter spam

// Determine what type of email it is
// Create event
// email.bounced
// email.autoreply
// email.reply
// email.new_thread
// contact.create

// Check to see if email is a reply to a flow.  If so, mark as complete.
// This should be handled in the email processor common service, not here.
// Same with goal achieved.
// Move slack notifications to processing service

func getTenantFromEmail(c *gin.Context, s *service.Services, emailData *PostmarkInboundEmailData) (string, error) {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.getTenant")
	defer span.Finish()
	tracing.TagComponentRest(span)

	nameFromBcc := emailData.TenantFromBcc()

	n, err := s.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantByNameIgnoreCase(ctx, nameFromBcc)
	if err != nil {
		span.LogFields(tracingLog.Bool("tenant.found", false))
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("Cannot identify tenant %s: %v", nameFromBcc, err)
	}

	if n == nil {
		span.LogFields(tracingLog.Bool("tenant.found", false))
		return "", fmt.Errorf("No valid tenant %s: %v", nameFromBcc, err)
	}

	tenant := mapper.MapDbNodeToTenantEntity(n)
	span.LogFields(tracingLog.Bool("tenant.found", true))
	span.LogFields(tracingLog.String("tenant.name", tenant.Name))

	return tenant.Name, nil
}

func getUsername(ctx context.Context, s *service.Services, EmailParticipants []string) (string, error) {
	span, _ := tracing.StartTracerSpan(ctx, "Flows.getUsername")
	defer span.Finish()
	tracing.TagComponentRest(span)

	tenant := common.GetTenantFromContext(ctx)

	for _, p := range EmailParticipants {
		userByEmail, err := s.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(
			ctx, tenant, p)
		if err != nil {
			return "", fmt.Errorf("Error getting username from email %s: %v", p, err)
		}
		if userByEmail != nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("Unable to find user amongst email participants")
}
