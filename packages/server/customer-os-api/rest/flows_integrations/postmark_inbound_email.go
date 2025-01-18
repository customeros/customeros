package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
)

const EXTERNAL_SYSTEM = "mailstack"

func PostmarkInboundEmail(c *gin.Context, s *cosapi_services.Services) {
	_, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.PostmarkInboundEmail", c.Request.Header)
	defer span.Finish()
	tracing.TagComponentRest(span)

	// Validate Postmark User-Agent
	if c.Request.UserAgent() == "" || !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
		tracing.TraceErr(span, fmt.Errorf("Invalid user agent %s", c.Request.UserAgent()))
		handlers.SendError(c, span, http.StatusForbidden, enum.ErrForbidden)
		return
	}

	// Parse email data
	emailData, err := parseInboundEmail(c)
	if err != nil {
		tracing.LogObjectAsJson(span, "body", c.Request.Body)
		tracing.TraceErr(span, err)
		handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest)
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

func processInboundEmail(c *gin.Context, s *cosapi_services.Services, emailData *PostmarkInboundEmailData) error {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.processInboundEmail")
	defer span.Finish()
	tracing.TagComponentRest(span)

	tenant, err := getTenant(c, s, emailData)
	if err != nil {
		return err
	}

	// update context with tenant, pass this where tenant is needed
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceCustomerOsApiRest,
	})

	participants := emailData.AllParticipantEmails()
	username, err := getUsername(ctx, s, participants)
	if err != nil || username == "" {
		span.LogFields(tracingLog.Bool("mailbox.found", false))
		tracing.TraceErr(span, err)
		return err
	}

	messageId := emailData.GetHeaderValue("Message-Id")
	emailExistsInDb, err := s.Repositories.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(
		ctx, EXTERNAL_SYSTEM, tenant, username, messageId,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("Unable to determine if email exists in db for messageId %s: %v", messageId, err)
	}

	if emailExistsInDb {
		return nil
	}

	dbEmailEntity := emailData.ToRawDbObject()
	jsonEmailEntity, err := json.Marshal(dbEmailEntity)
	if err != nil {
		return fmt.Errorf("Unable to produce JSON email object for db: %v", err)
	}

	dbErr := s.Repositories.PostgresRepositories.RawEmailRepository.Store(
		ctx, EXTERNAL_SYSTEM, tenant, username, dbEmailEntity.ProviderMessageId, messageId, string(jsonEmailEntity), dbEmailEntity.Sent, entity.REAL_TIME,
	)
	if dbErr != nil {
		span.LogFields(tracingLog.Object("raw_email", jsonEmailEntity))
		tracing.TraceErr(span, err)
		return fmt.Errorf("Error writing email to db: %v", err)
	}

	// Check to see if email is a reply to a flow.  If so, mark as complete.
	// This should be handled in the email processor common service, not here.
	// Same with goal achieved.
	// Move slack notifications to processing service
	// Handle attachments

	return nil
}

func getTenant(c *gin.Context, s *cosapi_services.Services, emailData *PostmarkInboundEmailData) (string, error) {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.getTenant")
	defer span.Finish()
	tracing.TagComponentRest(span)

	nameFromBcc := emailData.TenantFromBcc()

	n, err := s.Repositories.Neo4jRepositories.TenantReadRepository.GetTenantByNameIgnoreCase(ctx, nameFromBcc)
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

func getUsername(ctx context.Context, s *cosapi_services.Services, EmailParticipants []string) (string, error) {
	span, _ := tracing.StartTracerSpan(ctx, "Flows.getUsername")
	defer span.Finish()
	tracing.TagComponentRest(span)

	tenant := common.GetTenantFromContext(ctx)

	for _, p := range EmailParticipants {
		userByEmail, err := s.Repositories.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(
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
