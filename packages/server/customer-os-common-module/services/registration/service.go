package registration

import (
	"context"
	"fmt"
	"strings"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	common_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	mailbov_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/mailbox"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type registrationService struct {
	events   *events.EventsService
	postgres *repository.Repositories
	neo4j    *neoRepo.Repositories
	contact  interfaces.ContactService
	email    interfaces.EmailService
	flow     interfaces.FlowService
	mailbox  interfaces.MailboxService
	org      interfaces.OrganizationService
	postmark interfaces.PostmarkService
	user     interfaces.UserService
}

func NewRegistrationService(events *events.EventsService, postgres *repository.Repositories, neo4j *neoRepo.Repositories, contact interfaces.ContactService, email interfaces.EmailService, flow interfaces.FlowService, mailbox interfaces.MailboxService, org interfaces.OrganizationService, postmark interfaces.PostmarkService, user interfaces.UserService) interfaces.RegistrationService {
	return &registrationService{
		events:   events,
		postgres: postgres,
		neo4j:    neo4j,
		contact:  contact,
		email:    email,
		flow:     flow,
		mailbox:  mailbox,
		org:      org,
		postmark: postmark,
		user:     user,
	}
}

func (s *registrationService) SetContactService(contact interfaces.ContactService) {
	s.contact = contact
}

func (s *registrationService) SetEmailService(email interfaces.EmailService) {
	s.email = email
}

func (s *registrationService) SetFlowService(flow interfaces.FlowService) {
	s.flow = flow
}

func (s *registrationService) SetMailboxService(mailbox interfaces.MailboxService) {
	s.mailbox = mailbox
}

func (s *registrationService) SetOrganizationService(org interfaces.OrganizationService) {
	s.org = org
}

func (s *registrationService) IsInitialized() bool {
	return utils.IsInitialized(s)
}

func (s *registrationService) PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error {
	span, ctx := s.initializeTracing(ctx, "RegistrationService.PrepareDefaultTenantSetup", map[string]interface{}{
		"loggedInUserEmail": loggedInUserEmail,
	})
	defer span.Finish()

	if err := common.ValidateTenant(ctx); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	testUser, err := s.ConfigureTestMailbox(ctx)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error configuring test mailbox during tenant onboarding"))
	}

	if err = s.ConfigureDefaultFlowData(ctx, testUser); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error configuring test flow data during tenant onboarding"))
	}

	if err = s.CreatePostmarkServer(ctx); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating postmark server during tenant onboarding"))
	}

	return nil
}

func (s *registrationService) ConfigureDefaultFlowData(ctx context.Context, testUser *interfaces.TestUserSetup) error {
	span, ctx := s.initializeTracing(ctx, "RegistrationService.ConfigureDefaultFlowData", nil)
	defer span.Finish()

	if err := common.ValidateTenant(ctx); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowList, err := s.flow.FlowGetList(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	defaultFlowFound := false
	for _, flow := range *flowList {
		if flow.DefaultName == "Cold Outbound Example" {
			defaultFlowFound = true
			break
		}
	}

	if defaultFlowFound {
		return nil
	}

	tenant := common.GetTenantFromContext(ctx)

	organizationId, err := s.org.Save(ctx, nil, nil, data_fields.OrganizationFields{
		Name:         utils.StringPtr("Example Inc."),
		IndustryCode: utils.StringPtr("513210"),
		Employees:    utils.Int64Ptr(int64(100)),
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error saving organization during tenant onboarding"))
		return err
	}

	contactId, err := s.contact.Save(ctx, nil, nil, data_fields.ContactFields{
		FirstName: utils.StringPtr("Justin"),
		LastName:  utils.StringPtr("Example"),
	}, false)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error saving contact during tenant onboarding"))
		return err
	}

	_, err = s.email.Merge(ctx, nil, tenant, interfaces.EmailFields{
		Email: fmt.Sprintf("%s@%s", tenant, mailbov_srv.TEST_MAILBOX_DOMAIN),
	},
		&common_srv.LinkWith{
			Id:   contactId,
			Type: model.CONTACT,
		})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error saving email during tenant onboarding"))
		return err
	}

	err = s.contact.LinkContactWithOrganization(ctx, nil, contactId, organizationId, "Chief Testing Officer", "", constants.AppSourceUserAdminApi, true, utils.TimePtr(utils.Now()), nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error linking contact with organization during tenant onboarding"))
		return err
	}

	flow, err := s.flow.FlowMerge(ctx, nil, &neo4jentity.FlowEntity{
		Name:        "Cold Outbound Example",
		DefaultName: "Cold Outbound Example",
		Nodes: `
[
   {
      "$H":497,
      "data":{
         "action":"FLOW_START",
         "entity":"CONTACT",
         "triggerType":"RecordAddedManually"
      },
      "height":83,
      "id":"tn-1",
      "internalId":"cdff90d2-e357-4458-80e0-9dbfdab64b08",
      "measured":{
         "height":83,
         "width":300
      },
      "position":{
         "x":12,
         "y":12
      },
      "properties":{
         "org.eclipse.elk.portConstraints":"FIXED_ORDER"
      },
      "selected":false,
      "sourcePosition":"bottom",
      "targetPosition":"top",
      "type":"trigger",
      "width":300,
      "x":12,
      "y":12
   },
   {
      "$H":499,
      "data":{
         "action":"FLOW_END"
      },
      "height":56,
      "id":"tn-2",
      "internalId":"ec113fca-0a3d-4af7-9ae4-765491a827fd",
      "measured":{
         "height":56,
         "width":156
      },
      "position":{
         "x":84,
         "y":1131
      },
      "properties":{
         "org.eclipse.elk.portConstraints":"FIXED_ORDER"
      },
      "selected":false,
      "sourcePosition":"bottom",
      "targetPosition":"top",
      "type":"control",
      "width":156,
      "x":84,
      "y":1131
   },
   {
      "$H":345,
      "data":{
         "action":"WAIT",
         "fe_waitDurationUnit":"minutes",
         "isEditing":false,
         "nextStepId":"EMAIL_NEW-5bab768c-722c-4bf6-a315-8419f48013f0",
         "waitDuration":10
      },
      "height":56,
      "id":"WAIT-3bcac2b1-dcd8-4472-829d-1f1a74142c1e",
      "measured":{
         "height":56,
         "width":156
      },
      "position":{
         "x":84,
         "y":195
      },
      "selected":false,
      "type":"wait",
      "width":156,
      "x":84,
      "y":195
   },
   {
      "$H":347,
      "data":{
         "action":"EMAIL_NEW",
         "bodyTemplate":"<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\"><!--$-->\n<html dir=\"ltr\" lang=\"en\">\n\n  <body style=\"margin:auto;font-family:ui-sans-serif, system-ui, sans-serif\"><pre spellcheck=\"false\" data-language=\"plaintext\" data-highlight-language=\"plaintext\"><span style=\"white-space:pre-wrap\">Hello </span><span data-lexical-variable=\"true\">{{contact_first_name}}</span><span style=\"white-space:pre-wrap\">.</span><br/><br/><span style=\"white-space:pre-wrap\">Watch out, if you have outbound emails planned for </span><span data-lexical-variable=\"true\">{{organization_name}}</span><span style=\"white-space:pre-wrap\"> after Jan 1 2025: Two new EU regulations could damage your domain reputation.  </span><br/><span style=\"white-space:pre-wrap\">Rule 3A-12 comes down hard on tracking pixels, so we can forget about open rates (which have been broken for a while anyway).  </span><br/><span style=\"white-space:pre-wrap\">Rule 5C restricts how we can collect opt-in, and requires strict audit trails.  </span><br/><br/><span style=\"white-space:pre-wrap\">I&#x27;ve talked to over 12 CROs and CMOs in the last week, and there&#x27;s a lot of confusion as to what this means. With their feedback, I&#x27;ve put together a one-page breakdown of what we can, should, and must do by January. It&#x27;s short, but important.  </span><br/><br/><span style=\"white-space:pre-wrap\">Can I send you a copy?  </span><br/><br/><span style=\"white-space:pre-wrap\">Cheers,  </span><br/><span data-lexical-variable=\"true\">{{sender_first_name}}</span><code spellcheck=\"false\" style=\"white-space:pre-wrap\"><span class=\"editor-textCode\">  </span></code></pre>\n  </body>\n\n</html><!--/$-->",
         "fe_waitDurationUnit":"minutes",
         "isEditing":false,
         "subject":"Europe's new rules for email outbound",
         "waitBefore":10,
         "waitDuration":30,
         "waitStepId":"WAIT-3bcac2b1-dcd8-4472-829d-1f1a74142c1e"
      },
      "height":56,
      "id":"EMAIL_NEW-5bab768c-722c-4bf6-a315-8419f48013f0",
      "internalId":"279178d0-4d5a-4377-b6c6-5efadb6aa64f",
      "measured":{
         "height":56,
         "width":300
      },
      "position":{
         "x":12,
         "y":351
      },
      "selected":false,
      "type":"action",
      "width":300,
      "x":12,
      "y":351
   },
   {
      "$H":379,
      "data":{
         "action":"WAIT",
         "fe_waitDurationUnit":"days",
         "nextStepId":"EMAIL_REPLY-663dfdf6-7a9c-48cc-8cad-8f8a70ec95dd",
         "waitDuration":1440
      },
      "height":56,
      "id":"WAIT-24213077-fc7f-4726-bc6f-a73cc6c8beef",
      "measured":{
         "height":56,
         "width":156
      },
      "position":{
         "x":84,
         "y":507
      },
      "selected":false,
      "type":"wait",
      "width":156,
      "x":84,
      "y":507
   },
   {
      "$H":381,
      "data":{
         "action":"EMAIL_REPLY",
         "bodyTemplate":"<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\"><!--$-->\n<html dir=\"ltr\" lang=\"en\">\n\n  <body style=\"margin:auto;font-family:ui-sans-serif, system-ui, sans-serif\"><pre spellcheck=\"false\" data-language=\"plaintext\" data-highlight-language=\"plaintext\"><span style=\"white-space:pre-wrap\">Hi </span><span data-lexical-variable=\"true\">{{contact_first_name}}</span><span style=\"white-space:pre-wrap\">.</span><br/><br/><span style=\"white-space:pre-wrap\">A quick update on the new EU regulation coming online on Jan 1 2025: Rule 5C&#x27;s opt-in clause was amended, with far-reaching impacts.</span><br/><span style=\"white-space:pre-wrap\">There&#x27;s been a lot of discussion about this online, and a few of our clients shared a clever and empowering approach to compliance with this regulation. </span><br/><br/><span style=\"white-space:pre-wrap\">Can I share it with you?</span><br/><br/><span style=\"white-space:pre-wrap\">Cheers,  </span><br/><span data-lexical-variable=\"true\">{{sender_first_name}}</span></pre>\n  </body>\n\n</html><!--/$-->",
         "isEditing":false,
         "replyTo":"EMAIL_NEW-5bab768c-722c-4bf6-a315-8419f48013f0",
         "subject":"RE: Europe's new rules for email outbound",
         "waitBefore":1440,
         "waitStepId":"WAIT-24213077-fc7f-4726-bc6f-a73cc6c8beef"
      },
      "height":56,
      "id":"EMAIL_REPLY-663dfdf6-7a9c-48cc-8cad-8f8a70ec95dd",
      "internalId":"00af4bbe-af89-4d54-97ee-87b4d60ce8ec",
      "measured":{
         "height":56,
         "width":300
      },
      "position":{
         "x":12,
         "y":663
      },
      "selected":false,
      "type":"action",
      "width":300,
      "x":12,
      "y":663
   },
   {
      "$H":429,
      "data":{
         "action":"WAIT",
         "fe_waitDurationUnit":"days",
         "isEditing":false,
         "nextStepId":"EMAIL_REPLY-17f282b8-bd4d-47d2-9c6e-00add6b5114a",
         "waitDuration":2880
      },
      "height":56,
      "id":"WAIT-b4fd333b-6575-47d8-93b8-4c478b19b3fb",
      "measured":{
         "height":56,
         "width":156
      },
      "position":{
         "x":84,
         "y":819
      },
      "selected":false,
      "type":"wait",
      "width":156,
      "x":84,
      "y":819
   },
   {
      "$H":431,
      "data":{
         "action":"EMAIL_REPLY",
         "bodyTemplate":"<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\"><!--$-->\n<html dir=\"ltr\" lang=\"en\">\n\n  <body style=\"margin:auto;font-family:ui-sans-serif, system-ui, sans-serif\"><pre spellcheck=\"false\" data-language=\"plaintext\" data-highlight-language=\"plaintext\"><span style=\"white-space:pre-wrap\">Hi </span><span data-lexical-variable=\"true\">{{contact_first_name}}</span><span style=\"white-space:pre-wrap\">.  </span><br/><br/><span style=\"white-space:pre-wrap\">There&#x27;s been a lot of confusion about the new EU rules coming online Jan 1, so we&#x27;ve put together a diagnostic tool to simplify coverage and compliance, based on work done with over 100 clients. It&#x27;s in alpha testing right now. </span><br/><span style=\"white-space:pre-wrap\">I believe you do outbound emailing in the EU market, so you&#x27;d qualify to test it. Would you like to be part of the alpha group?</span><br/><br/><span style=\"white-space:pre-wrap\">Cheers,  </span><br/><span data-lexical-variable=\"true\">{{sender_first_name}}</span></pre>\n  </body>\n\n</html><!--/$-->",
         "isEditing":false,
         "replyTo":"EMAIL_NEW-5bab768c-722c-4bf6-a315-8419f48013f0",
         "subject":"RE: Europe's new rules for email outbound",
         "waitBefore":2880,
         "waitStepId":"WAIT-b4fd333b-6575-47d8-93b8-4c478b19b3fb"
      },
      "height":56,
      "id":"EMAIL_REPLY-17f282b8-bd4d-47d2-9c6e-00add6b5114a",
      "internalId":"14a310f5-1fbb-42ca-b8f9-473b3ae1ad79",
      "measured":{
         "height":56,
         "width":300
      },
      "position":{
         "x":12,
         "y":987
      },
      "selected":false,
      "type":"action",
      "width":300,
      "x":12,
      "y":987
   }
]
`,
		Edges: "[{\"id\":\"etn-1-WAIT-3bcac2b1-dcd8-4472-829d-1f1a74142c1e\",\"source\":\"tn-1\",\"target\":\"WAIT-3bcac2b1-dcd8-4472-829d-1f1a74142c1e\",\"type\":\"baseEdge\",\"markerEnd\":{\"type\":\"arrow\",\"width\":20,\"height\":20}},{\"id\":\"eWAIT-3bcac2b1-dcd8-4472-829d-1f1a74142c1e-EMAIL_NEW-5bab768c-722c-4bf6-a315-8419f48013f0\",\"source\":\"WAIT-3bcac2b1-dcd8-4472-829d-1f1a74142c1e\",\"target\":\"EMAIL_NEW-5bab768c-722c-4bf6-a315-8419f48013f0\",\"type\":\"baseEdge\",\"markerEnd\":{\"type\":\"arrow\",\"width\":20,\"height\":20},\"data\":{\"isHovered\":false}},{\"id\":\"eEMAIL_NEW-5bab768c-722c-4bf6-a315-8419f48013f0-WAIT-24213077-fc7f-4726-bc6f-a73cc6c8beef\",\"source\":\"EMAIL_NEW-5bab768c-722c-4bf6-a315-8419f48013f0\",\"target\":\"WAIT-24213077-fc7f-4726-bc6f-a73cc6c8beef\",\"type\":\"baseEdge\",\"markerEnd\":{\"type\":\"arrow\",\"width\":20,\"height\":20}},{\"id\":\"eWAIT-24213077-fc7f-4726-bc6f-a73cc6c8beef-EMAIL_REPLY-663dfdf6-7a9c-48cc-8cad-8f8a70ec95dd\",\"source\":\"WAIT-24213077-fc7f-4726-bc6f-a73cc6c8beef\",\"target\":\"EMAIL_REPLY-663dfdf6-7a9c-48cc-8cad-8f8a70ec95dd\",\"type\":\"baseEdge\",\"markerEnd\":{\"type\":\"arrow\",\"width\":20,\"height\":20}},{\"id\":\"eEMAIL_REPLY-663dfdf6-7a9c-48cc-8cad-8f8a70ec95dd-WAIT-b4fd333b-6575-47d8-93b8-4c478b19b3fb\",\"source\":\"EMAIL_REPLY-663dfdf6-7a9c-48cc-8cad-8f8a70ec95dd\",\"target\":\"WAIT-b4fd333b-6575-47d8-93b8-4c478b19b3fb\",\"type\":\"baseEdge\",\"markerEnd\":{\"type\":\"arrow\",\"width\":20,\"height\":20}},{\"id\":\"eWAIT-b4fd333b-6575-47d8-93b8-4c478b19b3fb-EMAIL_REPLY-17f282b8-bd4d-47d2-9c6e-00add6b5114a\",\"source\":\"WAIT-b4fd333b-6575-47d8-93b8-4c478b19b3fb\",\"target\":\"EMAIL_REPLY-17f282b8-bd4d-47d2-9c6e-00add6b5114a\",\"type\":\"baseEdge\",\"markerEnd\":{\"type\":\"arrow\",\"width\":20,\"height\":20},\"data\":{\"isHovered\":false}},{\"id\":\"eEMAIL_REPLY-17f282b8-bd4d-47d2-9c6e-00add6b5114a-tn-2\",\"source\":\"EMAIL_REPLY-17f282b8-bd4d-47d2-9c6e-00add6b5114a\",\"target\":\"tn-2\",\"type\":\"baseEdge\",\"markerEnd\":{\"type\":\"arrow\",\"width\":20,\"height\":20},\"data\":{\"isHovered\":false}}]",
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.flow.FlowSenderMerge(ctx, flow.Id, &neo4jentity.FlowSenderEntity{
		UserId: &testUser.UserId,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *registrationService) ConfigureTestMailbox(ctx context.Context) (*interfaces.TestUserSetup, error) {
	span, ctx := s.initializeTracing(ctx, "ConfigureTestMailbox", nil)
	defer span.Finish()

	if err := common.ValidateTenant(ctx); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	tenant := common.GetTenantFromContext(ctx)
	testUser, err := s.setupTestUser(ctx, span)
	if err != nil {
		return nil, err
	}

	if err := s.setupTestMailbox(ctx, span, tenant, testUser); err != nil {
		return nil, err
	}

	span.LogKV("result.mailboxAddress", testUser.MailboxAddress)
	return testUser, nil
}

func (s *registrationService) CreatePostmarkServer(ctx context.Context) error {
	span, ctx := s.initializeTracing(ctx, "CreatePostmarkServer", nil)
	defer span.Finish()

	if err := common.ValidateTenant(ctx); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err := s.postmark.CreateServerIfNotExists(ctx); err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}

// Helper functions

func (s *registrationService) initializeTracing(ctx context.Context, operation string, logFields map[string]interface{}) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("RegistrationService.%s", operation))
	tracing.SetDefaultServiceSpanTags(ctx, span)
	for key, value := range logFields {
		span.LogKV(key, value)
	}
	return span, ctx
}

func (s *registrationService) setupTestUser(ctx context.Context, span opentracing.Span) (*interfaces.TestUserSetup, error) {
	existingTestUser, err := s.neo4j.UserReadRepository.FindTestUser(ctx)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "cannot find test user"))
		return nil, err
	}

	var testUserId string
	if existingTestUser == nil {
		testUserId, err = s.user.Save(ctx, nil, nil, data_fields.UserFields{
			FirstName: utils.StringPtr("Test"),
			LastName:  utils.StringPtr("Sender"),
			Test:      utils.BoolPtr(true),
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "cannot create test user"))
			return nil, err
		}
	} else {
		testUserId = mapper.MapDbNodeToUserEntity(existingTestUser).Id
	}

	span.LogKV("result.testUserId", testUserId)
	return &interfaces.TestUserSetup{UserId: testUserId}, nil
}

func (s *registrationService) setupTestMailbox(ctx context.Context, span opentracing.Span, tenant string, testUser *interfaces.TestUserSetup) error {
	mailboxAddress := strings.ToLower(fmt.Sprintf("%s@%s", tenant, mailbov_srv.TEST_MAILBOX_DOMAIN))
	testUser.MailboxAddress = mailboxAddress

	testEmailId, err := s.email.Merge(ctx, nil, tenant, interfaces.EmailFields{
		Email: mailboxAddress,
	}, &common_srv.LinkWith{
		Type: model.USER,
		Id:   testUser.UserId,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to setup test mailbox"))
		return err
	}
	span.LogKV("result.testEmailId", testEmailId)

	return s.createMailboxIfNotExists(ctx, span, tenant, mailboxAddress)
}

func (s *registrationService) createMailboxIfNotExists(ctx context.Context, span opentracing.Span, tenant, mailboxAddress string) error {
	mailbox, err := s.postgres.TenantSettingsMailboxRepository.GetByMailbox(ctx, mailboxAddress)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get by mailbox"))
		return err
	}

	if mailbox == nil {
		if err := s.mailbox.CreateMailbox(ctx, nil, interfaces.CreateMailboxRequest{
			Domain:          mailbov_srv.TEST_MAILBOX_DOMAIN,
			Username:        strings.ToLower(tenant),
			Password:        utils.GenerateLowerAlpha(1) + utils.GenerateKey(11, false),
			LinkedUserEmail: mailboxAddress,
			WebmailEnabled:  true,
			ForwardingTo:    []string{fmt.Sprintf("bcc@%s.customeros.ai", strings.ToLower(tenant))},
		}); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to add mailbox"))
			return err
		}

		mailboxEntity, err := s.postgres.TenantSettingsMailboxRepository.GetByMailbox(ctx, strings.ToLower(tenant)+"@"+mailbov_srv.TEST_MAILBOX_DOMAIN)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get by mailbox"))
			return err
		}

		err = s.events.Publisher.PublishEvent(ctx, mailboxEntity.ID, model.MAILBOX, dto.MailstackProvisionMailbox{})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
