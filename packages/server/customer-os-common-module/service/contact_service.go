package service

import (
	"context"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type ContactService interface {
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, contactFields data_fields.ContactFields, updateOnlyIfEmpty bool) (string, error)
	CreateContactByLinkedIn(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkedInUrl string) (string, error)
	HideContact(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId string) error
	ShowContact(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId string) error
	GetContactById(ctx context.Context, contactId string) (*neo4jentity.ContactEntity, error)
	LinkContactWithOrganization(ctx context.Context, contactId, organizationId, jobTitle, description, source string, primary bool, startedAt, endedAt *time.Time) error
	CheckContactExistsWithLinkedIn(ctx context.Context, url, alias, externalId string) (bool, string, error)
	CheckContactExistsWithEmail(ctx context.Context, email string) (bool, string, error)
}

type contactService struct {
	log      logger.Logger
	services *Services
}

func NewContactService(log logger.Logger, services *Services) ContactService {
	return &contactService{
		log:      log,
		services: services,
	}
}

func (s *contactService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, contactFields data_fields.ContactFields, updateOnlyIfEmpty bool) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "contactFields", contactFields)
	span.LogFields(log.Bool("updateOnlyIfEmpty", updateOnlyIfEmpty))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	contactId := ""

	if id == nil || *id == "" {
		createFlow = true
		span.LogKV("flow", "create")
	} else {
		span.LogKV("flow", "update")
	}

	if createFlow {
		// generate id
		contactId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelContact)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}

		// prepare missing fields
		if contactFields.CreatedAt == nil {
			contactFields.CreatedAt = utils.NowPtr()
		} else {
			contactFields.CreatedAt = utils.TimePtr(utils.NowIfZero(*contactFields.CreatedAt))
		}
		if utils.IfNotNilString(contactFields.Source) == "" {
			contactFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if utils.IfNotNilString(contactFields.AppSource) == "" {
			contactFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
		if contactFields.Hide == nil {
			contactFields.Hide = utils.BoolPtr(false)
		}
	} else {
		contactId = *id
		// validate contact exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, contactId, model.NodeLabelContact)
		if err != nil || !exists {
			err = errors.New("contact not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, contactId)

	// Clean and update contact names if not updated manually
	if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
		if contactFields.Name != nil {
			contactFields.Name = utils.StringPtr(utils.CleanName(*contactFields.Name))
		}
		if contactFields.FirstName != nil {
			contactFields.FirstName = utils.StringPtr(utils.CleanName(*contactFields.FirstName))
		}
		if contactFields.LastName != nil {
			contactFields.LastName = utils.StringPtr(utils.CleanName(*contactFields.LastName))
		}
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		innerErr := s.services.Neo4jRepositories.ContactWriteRepository.SaveContactInTx(ctx, txWithPostCommit.Tx, tenant, contactId, contactFields, updateOnlyIfEmpty)
		if innerErr != nil {
			s.log.Errorf("Error while saving contact %s: %s", contactId, err.Error())
			return nil, innerErr
		}

		if contactFields.ExternalSystemAvailable() {
			innerErr = s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, contactId, model.NodeLabelContact, *contactFields.ExternalSystem)
			if err != nil {
				s.log.Errorf("Error while link contact %s with external system %s: %s", contactId, contactFields.ExternalSystem.ExternalSystemId, err.Error())
				return nil, innerErr
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if createFlow {
				err = s.services.RabbitMQService.PublishEvent(ctx, contactId, model.CONTACT, dto.CreateContact{contactFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateContact"))
				}
				utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithCreate())
			} else {
				err = s.services.RabbitMQService.PublishEvent(ctx, contactId, model.CONTACT, dto.UpdateContact{contactFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateContact"))
				}
				if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
					utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
				}
			}

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createFlow {
		span.LogFields(log.Bool("response.contactCreated", true))
	} else {
		span.LogFields(log.Bool("response.contactUpdated", true))
	}
	return contactId, nil
}

func (s *contactService) HideContact(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.HideContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, contactId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		contactFields := data_fields.ContactFields{Hide: utils.BoolPtr(false)}
		err = s.services.Neo4jRepositories.ContactWriteRepository.SaveContactInTx(ctx, txWithPostCommit.Tx, tenant, contactId, contactFields, false)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("error while hiding contact %s: %s", contactId, err.Error())
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.services.RabbitMQService.PublishEvent(ctx, contactId, model.CONTACT, dto.HideContact{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message HideContact"))
			}

			utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithDelete())
			return nil
		})
		return nil, nil
	})

	return nil
}

func (s *contactService) ShowContact(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.ShowContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, contactId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		contactFields := data_fields.ContactFields{Hide: utils.BoolPtr(true)}
		err = s.services.Neo4jRepositories.ContactWriteRepository.SaveContactInTx(ctx, txWithPostCommit.Tx, tenant, contactId, contactFields, false)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("error while showing contact %s: %s", contactId, err.Error())
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.services.RabbitMQService.PublishEvent(ctx, contactId, model.CONTACT, dto.ShowContact{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message ShowContact"))
			}

			utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithCreate())

			return nil
		})

		return nil, nil
	})

	return nil
}

func (s *contactService) GetContactById(ctx context.Context, contactId string) (*neo4jentity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, contactId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	contactDbNode, err := s.services.Neo4jRepositories.ContactReadRepository.GetContact(ctx, tenant, contactId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while getting contact %s: %s", contactId, err.Error())
		return nil, err
	}

	return neo4jmapper.MapDbNodeToContactEntity(contactDbNode), nil
}

func (s *contactService) LinkContactWithOrganization(ctx context.Context, contactId, organizationId, jobTitle, description, source string, primary bool, startedAt, endedAt *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.LinkContactWithOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, contactId)
	span.LogFields(log.String("organizationId", organizationId), log.String("jobTitle", jobTitle), log.String("description", description), log.Bool("primary", primary))
	if startedAt != nil {
		span.LogFields(log.Object("startedAt", startedAt))
	}
	if endedAt != nil {
		span.LogFields(log.Object("endedAt", endedAt))
	}

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate contact exists
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, contactId, model.NodeLabelContact)
	if err != nil || !exists {
		err = errors.New("contact not found")
		tracing.TraceErr(span, err)
		return err
	}

	// validate organization exists
	exists, err = s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, organizationId, model.NodeLabelOrganization)
	if err != nil || !exists {
		err = errors.New("organization not found")
		tracing.TraceErr(span, err)
		return err
	}

	if startedAt != nil && startedAt.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
		startedAt = nil
	}
	if endedAt != nil && endedAt.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
		endedAt = nil
	}

	jobRoleData := neo4jrepository.JobRoleFields{
		Description: description,
		JobTitle:    jobTitle,
		Primary:     primary,
		StartedAt:   startedAt,
		EndedAt:     endedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:    neo4jmodel.GetSource(source),
			AppSource: common.GetAppSourceFromContext(ctx),
		},
	}

	err = s.services.Neo4jRepositories.JobRoleWriteRepository.LinkContactWithOrganization(ctx, tenant, contactId, organizationId, jobRoleData)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to link contact with organization"))
		return err
	}

	// reset contact enrich attempts
	_ = s.services.Neo4jRepositories.ContactWriteRepository.ResetEnrichAttempts(ctx, nil, tenant, contactId)

	utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())

	// send 2 events, for contact another for organization
	dtoData := dto.AddContactToOrganization{
		ContactId:      contactId,
		OrganizationId: organizationId,
		JobTitle:       jobTitle,
		Description:    description,
		Primary:        primary,
		StartedAt:      startedAt,
		EndedAt:        endedAt,
	}

	err = s.services.RabbitMQService.PublishEvent(ctx, contactId, model.CONTACT, dtoData)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddContactToOrganization for contact"))
	}
	err = s.services.RabbitMQService.PublishEvent(ctx, organizationId, model.ORGANIZATION, dtoData)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddContactToOrganization for organization"))
	}

	return nil
}

func (s *contactService) CheckContactExistsWithLinkedIn(ctx context.Context, url, alias, externalId string) (bool, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.CheckContactExistsWithLinkedIn")
	defer span.Finish()

	if alias == "" {
		// use identifier as alias
		alias = neo4jentity.SocialEntity{Url: url}.ExtractLinkedinPersonIdentifierFromUrl()
	}
	contacts, err := s.services.Neo4jRepositories.ContactReadRepository.GetContactsByLinkedIn(ctx, common.GetTenantFromContext(ctx), url, alias, externalId)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, "", err
	}
	contactId := ""
	if len(contacts) > 0 {
		contactId = contacts[0].Props["id"].(string)
	}
	return len(contacts) > 0, contactId, nil
}

func (s *contactService) CheckContactExistsWithEmail(ctx context.Context, email string) (bool, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.CheckContactExistsWithEmail")
	defer span.Finish()

	if email == "" {
		return false, "", nil
	}
	// if not a valid email, return false
	syntaxValidation := mailsherpa.ValidateEmailSyntax(email)
	if !syntaxValidation.IsValid {
		return false, "", nil
	}

	contacts, err := s.services.Neo4jRepositories.ContactReadRepository.GetContactsWithEmail(ctx, common.GetTenantFromContext(ctx), email)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, "", err
	}
	contactId := ""
	if len(contacts) > 0 {
		contactId = contacts[0].Props["id"].(string)
	}
	return len(contacts) > 0, contactId, nil
}

func (s *contactService) CreateContactByLinkedIn(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkedInUrl string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.CreateContactByLinkedIn")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("linkedInUrl", linkedInUrl)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// check if linkedInUrl is valid
	socialEntity := neo4jentity.SocialEntity{Url: linkedInUrl}
	if !socialEntity.IsLinkedin() {
		err := errors.New("not a valid linkedin url")
		tracing.TraceErr(span, err)
		return "", err
	}

	// Reject contact creation if linked-in url is already used by another contact
	if utils.IfNotNilString(linkedInUrl) != "" {
		if (neo4jentity.SocialEntity{Url: linkedInUrl}).IsLinkedin() {
			linkedInAlreadyUsed, existingContactId, err := s.services.ContactService.CheckContactExistsWithLinkedIn(ctx, linkedInUrl, "", "")
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to check contact exists with linkedin"))
				return "", err
			}
			if linkedInAlreadyUsed {
				contactByLinkedInEntity, err := s.GetContactById(ctx, existingContactId)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to get contact by id"))
					return "", err
				}

				_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
					if contactByLinkedInEntity.Hide {
						err = s.ShowContact(ctx, txWithPostCommit, existingContactId)
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "unable to show contact"))
							return "", err
						}
					} else {
						// just update contact' updatedAt
						err = s.services.Neo4jRepositories.CommonWriteRepository.TouchEntity(ctx, txWithPostCommit.Tx, tenant, model.NodeLabelContact, existingContactId)
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "error on updating contact updatedAt"))
						}
						txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
							utils.EventCompleted(ctx, tenant, model.CONTACT.String(), existingContactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
							return nil
						})
					}
					return nil, nil
				})
				return existingContactId, nil
			}
		}
	}

	createdContactId := ""

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		createdContactId, err = s.Save(ctx, txWithPostCommit, nil, data_fields.ContactFields{}, false)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create contact"))
			return "", err
		}
		if (neo4jentity.SocialEntity{Url: linkedInUrl}).IsLinkedin() {
			_, err := s.services.SocialService.AddSocialToEntity(ctx, txWithPostCommit,
				LinkWith{
					Id:   createdContactId,
					Type: model.CONTACT,
				},
				neo4jentity.SocialEntity{
					Url:       linkedInUrl,
					Source:    neo4jentity.DecodeDataSource(neo4jentity.DataSourceOpenline.String()),
					AppSource: common.GetAppSourceFromContext(ctx),
				})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to merge social with contact"))
			}
		}
		return nil, nil
	})

	return createdContactId, nil
}
