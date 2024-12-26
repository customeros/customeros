package service

import (
	"context"
	"fmt"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

const AppSource = constants.AppSourceSyncEmail

func (s *mailService) GetEmailsForProcessingForUser(ctx context.Context, tenant, userEmailAddress string) {
	span, ctx := s.initializeTracing(ctx, "MailService.SyncEmailsForUser")
	defer span.Finish()
	span.LogKV("userEmailAddress", userEmailAddress)

	rawEmailsIdsForProcess, err := s.services.PostgresRepositories.RawEmailRepository.GetEmailsIdsForUserForSync(tenant, userEmailAddress)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get emails for sync"))
		return
	}

	if len(rawEmailsIdsForProcess) == 0 {
		err = fmt.Errorf("no emails found for sync")
		tracing.TraceErr(span, err)
		return
	}

	distinctExternalSystems := make([]string, 0)
	for _, email := range rawEmailsIdsForProcess {
		if !utils.Contains(distinctExternalSystems, email.ExternalSystem) {
			distinctExternalSystems = append(distinctExternalSystems, email.ExternalSystem)
		}
	}

	externalSystemStr := ""
	if len(distinctExternalSystems) > 0 {
		externalSystemStr = distinctExternalSystems[0]
	}

	// Create email node in neo4j
	err = s.createUserEmailAddressAsNode(ctx, tenant, userEmailAddress, externalSystemStr, span)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create user source as email node"))
		return
	}

	for _, externalSystem := range distinctExternalSystems {
		// TODO alexb add caching for each tenant of external systems
		err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(
			ctx, tenant, externalSystem, externalSystem,
		)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to merge external system"))
			return
		}
	}

	s.processRawEmails(ctx, tenant, rawEmailsIdsForProcess, span)
}

func (s *mailService) ProcessEmail(ctx context.Context, tenant string, rawEmailId uuid.UUID) entity.UpdateRawEmailTable {
	var db entity.UpdateRawEmailTable

	span, ctx := s.initializeTracing(ctx, "MailService.ProcessEmail")
	defer span.Finish()
	span.LogFields(log.String("rawEmailId", rawEmailId.String()))

	rawEmail, err := s.services.PostgresRepositories.RawEmailRepository.GetEmailForProcess(rawEmailId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get email for process"))
		db.EmailProcessingStatus = postgresentity.ERROR
		db.Error = err
		return db
	}

	emailMessageData, err := s.LoadEmail(ctx, rawEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to load email"))
		db.EmailProcessingStatus = postgresentity.ERROR
		db.Error = err
		return db
	}

	if emailMessageData.Identifiers.MessageId == "" {
		db.EmailProcessingStatus = postgresentity.ERROR
		db.Error = fmt.Errorf("email message id is empty")
		return db
	}

	if len(emailMessageData.Participants.AllEmails) == 0 {
		reason := "no email address belongs to a workspace domain"
		db.EmailProcessingStatus = postgresentity.SKIPPED
		db.Reason = &reason
		return db
	}

	check := s.ProcessEmailCheck(ctx, tenant, &emailMessageData)
	if !check.ProcessEmail {
		db.EmailProcessingStatus = postgresentity.SKIPPED
		db.Reason = &check.SkipReason
		db.BouncedEmails = &check.BouncedEmails
		// set all bounced emails to undeliverable
		for _, e := range *db.BouncedEmails {
			err := s.services.Neo4jRepositories.EmailWriteRepository.SetDeliverableByEmailForAllTenants(ctx, e, "false")
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to set deliverable by email for all tenants"))
			}
		}
		return db
	}

	// set all non-bounced emails to deliverable
	for _, e := range emailMessageData.Participants.AllEmails {
		err := s.services.Neo4jRepositories.EmailWriteRepository.SetDeliverableByEmailForAllTenants(ctx, e, "true")
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set deliverable by email for all tenants"))
		}
	}

	sentAt, err := convertToUTC(emailMessageData.Content.SentDate)
	emailMessageData.CreatedAt = sentAt
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to convert email sent date to UTC"))
		db.EmailProcessingStatus = postgresentity.ERROR
		db.Error = err
		return db
	}

	interactionEventId, err := s.services.Neo4jRepositories.InteractionEventReadRepository.GetInteractionEventIdByExternalId(
		ctx, tenant, rawEmail.ExternalSystem, rawEmail.MessageId,
	)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check if interaction event exists"))
		db.EmailProcessingStatus = postgresentity.ERROR
		db.Error = err
		return db
	}

	if interactionEventId != "" {
		tracing.TraceErr(span, errors.Wrap(err, "interaction event already exists"))
		reason := "interaction event already exists"
		db.EmailProcessingStatus = postgresentity.SKIPPED
		db.Reason = &reason
		return db
	}

	chanErr := s.buildChannelData(&emailMessageData, span)
	if chanErr != nil {
		tracing.TraceErr(span, errors.Wrap(chanErr, "failed to build channel data"))
		db.EmailProcessingStatus = postgresentity.ERROR
		db.Error = err
		return db
	}

	return s.processInboundEmail(ctx, tenant, &emailMessageData, rawEmail)
}

func (s *mailService) ProcessEmailByMessageId(ctx context.Context, tenant, usernameSource, messageId string) entity.UpdateRawEmailTable {
	var db entity.UpdateRawEmailTable

	span, ctx := s.initializeTracing(ctx, "MailService.ProcessEmailByMessageId")
	defer span.Finish()
	span.LogFields(
		log.String("userSource", usernameSource),
		log.String("messageId", messageId))

	rawEmail, err := s.services.PostgresRepositories.RawEmailRepository.GetEmailForSyncByMessageId(tenant, usernameSource, messageId)
	if err != nil {
		err = fmt.Errorf("failed to get emails for sync: %v", err)
		tracing.TraceErr(span, err)
		db.EmailProcessingStatus = postgresentity.ERROR
		db.Error = err
		return db
	}

	if rawEmail == nil {
		db.EmailProcessingStatus = postgresentity.ERROR
		db.Error = fmt.Errorf("email with message id %v not found", messageId)
		return db
	}

	return s.ProcessEmail(ctx, tenant, rawEmail.ID)
}

func (s *mailService) ProcessEmailByEmailRawId(ctx context.Context, tenant string, emailId uuid.UUID) entity.UpdateRawEmailTable {
	return s.ProcessEmail(ctx, tenant, emailId)
}

func (s *mailService) processRawEmails(ctx context.Context, tenant string, rawEmails []postgresentity.RawEmail, span opentracing.Span) {
	for _, rawEmail := range rawEmails {
		dbUpdateRecord := s.ProcessEmail(ctx, tenant, rawEmail.ID)

		err := s.services.PostgresRepositories.RawEmailRepository.UpdateRawEmailTable(rawEmail.ID, dbUpdateRecord)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to mark raw email as processed in postgres"))
		}
	}
}

func (s *mailService) processInboundEmail(ctx context.Context, tenant string, email *EmailMessageData, rawEmail *postgresentity.RawEmail) entity.UpdateRawEmailTable {
	span, ctx := s.initializeTracing(ctx, "MailService.processInboundEmail")
	defer span.Finish()
	tracing.TagTenant(span, tenant)

	var db entity.UpdateRawEmailTable
	var txWithPostCommit *utils.TxWithPostCommit

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		return nil, s.processSessionAndEvents(ctx, txWithPostCommit, tenant, email, rawEmail)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		db.EmailProcessingStatus = postgresentity.ERROR
		db.Error = fmt.Errorf("mail with message id %s failed processing", email.Identifiers.MessageId)
		return db
	}

	db.EmailProcessingStatus = postgresentity.PROCESSED
	return db
}

func (s *mailService) processSessionAndEvents(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant string, emailMessageData *EmailMessageData, rawEmail *postgresentity.RawEmail) error {
	span, ctx := s.initializeTracing(ctx, "MailService.processSessionAndEvents")
	defer span.Finish()
	tracing.TagTenant(span, tenant)

	syncDate := utils.Now()

	// get EmailForCustomerOS
	cosEmail := s.buildEmailForCustomerOS(emailMessageData, rawEmail.ExternalSystem)

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// step 1: Get or create session
		sessionId, err := s.services.Neo4jRepositories.InteractionSessionWriteRepository.MergeByIdentifierAndChannel(
			ctx, txWithPostCommit.Tx, tenant, emailMessageData.Identifiers.EmailThreadId, syncDate, cosEmail, commonenum.InteractionSessionTypeThread, commonenum.InteractionSessionChannelEmail, rawEmail.ExternalSystem, AppSource,
		)
		if err != nil {
			err = fmt.Errorf("failed merge interaction session: %w", err)
			return nil, err
		}

		// step 2: Get or create interaction event
		interactionEventId, err := s.services.Neo4jRepositories.InteractionEventWriteRepository.MergeByExternalSystem(
			ctx, txWithPostCommit.Tx, tenant, syncDate, cosEmail, rawEmail.ExternalSystem, AppSource,
		)
		if err != nil {
			err = fmt.Errorf("failed merge interaction event: %w", err)
			return nil, err
		}

		// step 3: Link event to session
		if err = s.services.Neo4jRepositories.InteractionEventWriteRepository.LinkInteractionEventToSession(
			ctx, txWithPostCommit.Tx, tenant, interactionEventId, sessionId,
		); err != nil {
			err = fmt.Errorf("failed to link event to session: %w", err)
			return nil, err
		}

		// step 4: Process participants
		if err = s.linkParticipants(ctx, txWithPostCommit, tenant, interactionEventId, &emailMessageData.Participants, syncDate, rawEmail.ExternalSystem); err != nil {
			err = fmt.Errorf("failed to link participants: %v", err)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *mailService) buildEmailForCustomerOS(email *EmailMessageData, externalSystem string) model.SaveEmailMessage {
	save := model.SaveEmailMessage{
		Html:           email.Content.Html,
		Text:           email.Content.Text,
		Subject:        email.Content.Subject,
		CreatedAt:      email.CreatedAt,
		ExternalSystem: externalSystem,
		ExternalId:     email.Identifiers.MessageId,
		EmailThreadId:  email.Identifiers.EmailThreadId,
		Channel:        "EMAIL",
		ChannelData:    email.ChannelData,
	}
	return save
}

func (s *mailService) linkParticipants(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant string, interactionEventId string, participants *EmailParticipants, syncDate time.Time, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.linkParticipants")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogKV("interactionEventId", interactionEventId)

	emailIds := make(map[string]string)

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// Link From participant
		fromId, err := s.getOrCreateEmailId(ctx, txWithPostCommit, tenant, participants.From.Email, syncDate, externalSystem, emailIds)
		if err != nil {
			err = fmt.Errorf("failed to get or create email id: %w", err)
			return nil, err
		}
		if err = s.services.Neo4jRepositories.InteractionEventWriteRepository.InteractionEventSentByEmail(ctx, txWithPostCommit.Tx, tenant, interactionEventId, fromId); err != nil {
			err = fmt.Errorf("failed to create interaction event sent by email: %w", err)
			return nil, err
		}

		// Link To participants
		if err = s.linkEmailGroup(ctx, txWithPostCommit, tenant, interactionEventId, "TO", participants.GetToEmailAddresses(), syncDate, externalSystem, emailIds); err != nil {
			err = fmt.Errorf("failed to link email group for TO: %w", err)
			return nil, err
		}

		// Link CC participants
		if err = s.linkEmailGroup(ctx, txWithPostCommit, tenant, interactionEventId, "CC", participants.GetCcEmailAddresses(), syncDate, externalSystem, emailIds); err != nil {
			err = fmt.Errorf("failed to link email group for CC: %w", err)
			return nil, err
		}

		// Link BCC participants
		if err = s.linkEmailGroup(ctx, txWithPostCommit, tenant, interactionEventId, "BCC", participants.GetBccEmailAddresses(), syncDate, externalSystem, emailIds); err != nil {
			err = fmt.Errorf("failed to link email group for BCC: %w", err)
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *mailService) linkEmailGroup(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant string, interactionEventId string, groupType string, emails []string,
	syncDate time.Time, externalSystem string, emailIds map[string]string) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.linkEmailGroup")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, interactionEventId)
	span.LogFields(log.String("groupType", groupType), log.String("emails", strings.Join(emails, ",")), log.String("externalSystem", externalSystem))

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		var groupEmailIds []string
		for _, email := range emails {
			if email == "" {
				continue
			}

			emailId, err := s.getOrCreateEmailId(ctx, txWithPostCommit, tenant, email, syncDate, externalSystem, emailIds)
			if err != nil {
				err = fmt.Errorf("failed to get or create email ID: %w", err)
				return nil, err
			}

			if !utils.Contains(groupEmailIds, emailId) {
				groupEmailIds = append(groupEmailIds, emailId)
			}
		}

		if len(groupEmailIds) > 0 {
			return nil, s.services.Neo4jRepositories.InteractionEventWriteRepository.InteractionEventSentToEmails(
				ctx, txWithPostCommit.Tx, tenant, interactionEventId, groupType, groupEmailIds,
			)
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *mailService) getOrCreateEmailId(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant string, email string, syncDate time.Time, externalSystem string, emailIds map[string]string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.getOrCreateEmailId")
	defer span.Finish()
	tracing.TagTenant(span, tenant)

	if id, exists := emailIds[email]; exists {
		return id, nil
	}

	id, err := s.GetEmailIdForEmail(ctx, txWithPostCommit, tenant, email, externalSystem)
	if err != nil {
		err = fmt.Errorf("failed to get email ID for %s: %v", email, err)
		tracing.TraceErr(span, err)
		return "", err
	}
	if id == "" {
		err = fmt.Errorf("no email ID found for %s", email)
		tracing.TraceErr(span, err)
		return "", err
	}

	emailIds[email] = id
	return id, nil
}

func (s *mailService) buildChannelData(email *EmailMessageData, span opentracing.Span) error {
	channelData, err := neo4jentity.BuildEmailChannelData(
		email.Identifiers.ProviderMessageId,
		email.Identifiers.EmailThreadId,
		email.Content.Subject,
		strings.Join(email.Participants.GetReplyToEmailAddresses(), " "),
		strings.Join(email.Identifiers.References, " "),
	)
	if err != nil {
		err = fmt.Errorf("building channel data failed: %v", err)
		tracing.TraceErr(span, err)
		return err
	}
	email.Channel = "EMAIL"
	email.ChannelData = channelData

	return nil
}

func (s *mailService) createUserEmailAddressAsNode(ctx context.Context, tenant, userEmailAddress, externalSystem string, span opentracing.Span) error {
	_, err := s.services.EmailService.Merge(ctx, nil, tenant, EmailFields{
		Email:     userEmailAddress,
		AppSource: AppSource,
		Source:    neo4jentity.DecodeDataSource(externalSystem),
	}, nil)
	if err != nil {
		err = fmt.Errorf("unable to create email: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *mailService) mapDbNodeToEmailEntity(node dbtype.Node) *neo4jentity.EmailEntity {
	props := utils.GetPropsFromNode(node)
	result := neo4jentity.EmailEntity{
		Id:       utils.GetStringPropOrEmpty(props, "id"),
		Email:    utils.GetStringPropOrEmpty(props, "email"),
		RawEmail: utils.GetStringPropOrEmpty(props, "rawEmail"),
	}
	return &result
}

func (s *mailService) GetEmailIdForEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant, email string, source string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.GetEmailIdForEmail")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogKV("email", email)

	if email == "" {
		return "", nil
	}
	emailSyntax := mailsherpa.ValidateEmailSyntax(email)
	if !emailSyntax.IsValid {
		return "", nil
	}

	output, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		emailId, err := s.services.Neo4jRepositories.EmailReadRepository.GetEmailIdIfExists(ctx, txWithPostCommit.Tx, tenant, email)
		if err != nil {
			return "", fmt.Errorf("unable to retrieve email id for tenant: %w", err)
		}
		if emailId != "" {
			return emailId, nil
		}

		//if it's a personal email, we create just the email node in tenant
		domain := utils.ExtractDomainFromEmail(email)
		if domain == "" {
			err = errors.New("unable to extract domain from email: " + email)
			return "", err
		}
		if utils.Contains(s.services.Cache.GetPersonalEmailProviders(), domain) {
			emailIdPtr, err := s.services.EmailService.Merge(ctx, txWithPostCommit, tenant, EmailFields{
				Email:     email,
				Source:    neo4jentity.DecodeDataSource(source),
				AppSource: constants.AppSourceSyncEmail,
			}, nil)
			if err != nil {
				return "", err
			}
			emailId = utils.IfNotNilString(emailIdPtr)
			return emailId, nil
		}

		var organizationNode *neo4j.Node
		var organizationId string

		organizationNode, err = s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, txWithPostCommit.Tx, tenant, domain)
		if err != nil {
			return "", fmt.Errorf("unable to retrieve organization for tenant: %v", err)
		}

		if organizationNode == nil {
			stage := neo4jenum.Lead
			leadSource := ""

			if source == neo4jentity.DataSourceGmail.String() {
				leadSource = "Gmail"
			} else if source == neo4jentity.DataSourceOutlook.String() {
				leadSource = "Outlook"
			} else if source == neo4jentity.DataSourceMailstack.String() {
				leadSource = "Mailstack"
				stage = neo4jenum.Target
			} else {
				leadSource = "Email"
			}

			organizationFields := data_fields.OrganizationFields{
				LeadSource:   utils.StringPtr(leadSource),
				Stage:        utils.ToPtr(stage),
				Relationship: utils.ToPtr(neo4jenum.OrganizationRelationshipProspect),
				Domains:      []string{domain},
				Source:       utils.StringPtr(source),
			}

			organizationId, err = s.services.OrganizationService.Save(ctx, txWithPostCommit, nil, organizationFields)
			if err != nil {
				return "", fmt.Errorf("failed to create organization: %w", err)
			}
		} else {
			organizationId = utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id")
		}

		contactId, err := s.services.ContactService.CreateContactByEmail(ctx, txWithPostCommit, email)
		if err != nil {
			return "", fmt.Errorf("unable to create contact: %w", err)
		}

		err = s.services.ContactService.LinkContactWithOrganization(ctx, txWithPostCommit, contactId, organizationId, "", "", source, false, nil, nil)

		emailId, err = s.services.Neo4jRepositories.EmailReadRepository.GetEmailIdIfExists(ctx, txWithPostCommit.Tx, tenant, email)
		if err != nil {
			return "", fmt.Errorf("unable to retrieve email id for tenant: %w", err)
		}

		return emailId, nil
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get email id for email"))
		return "", err
	}

	return output.(string), nil
}
