package contact

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/customeros/mailsherpa/emailparser"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	common_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type contactService struct {
	log          logger.Logger
	neo4j        *neoRepo.Repositories
	events       *events.EventsService
	domain       interfaces.DomainService
	email        interfaces.EmailService
	organization interfaces.OrganizationService
	jobrole      interfaces.JobRoleService
	social       interfaces.SocialService
}

func NewContactService(log logger.Logger, neo4j *neoRepo.Repositories, events *events.EventsService, domain interfaces.DomainService, email interfaces.EmailService, organization interfaces.OrganizationService, jobrole interfaces.JobRoleService, social interfaces.SocialService) interfaces.ContactService {
	return &contactService{
		log:          log,
		neo4j:        neo4j,
		events:       events,
		domain:       domain,
		email:        email,
		organization: organization,
		jobrole:      jobrole,
		social:       social,
	}
}

func (s *contactService) CreateContactWithOrganizationByEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, email string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.CreateContactWithOrganizationByEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("email", email)

	contactId := ""
	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		var innerErr error
		contactId, innerErr = s.CreateContactByEmail(ctx, txWithPostCommit, email)
		if innerErr != nil {
			return nil, innerErr
		}

		mailValidation := mailsherpa.ValidateEmailSyntax(email)
		validDomainForOrganization := s.domain.IsAcceptedDomainForOrganization(ctx, mailValidation.Domain)

		if validDomainForOrganization {
			organizationId, innerErr := s.organization.Save(ctx, txWithPostCommit, nil, data_fields.OrganizationFields{
				Domains: []string{mailValidation.Domain},
			})
			if innerErr != nil {
				return nil, innerErr
			}

			innerErr = s.LinkContactWithOrganization(ctx, txWithPostCommit, contactId, organizationId, "", "", constants.AppSourceCustomerOsApi, true, nil, nil)
			if innerErr != nil {
				return nil, innerErr
			}
		}

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return contactId, nil
}

func (s *contactService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, contactFields data_fields.ContactFields, updateOnlyIfEmpty bool, options ...common_srv.ServiceOptions) (string, error) {
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

	if utils.IfNotNilString(id) == "" {
		createFlow = true
		span.LogKV("flow", "create")

		// generate id
		contactId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelContact)
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
		span.LogKV("flow", "update")

		contactId = *id
		// validate contact exists
		exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, contactId, model.NodeLabelContact)
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

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		innerErr := s.neo4j.ContactWriteRepository.SaveContactInTx(ctx, txWithPostCommit.Tx, tenant, contactId, contactFields, updateOnlyIfEmpty)
		if innerErr != nil {
			s.log.Errorf("Error while saving contact %s: %s", contactId, err.Error())
			return nil, innerErr
		}

		if contactFields.ExternalSystemAvailable() {
			innerErr = s.neo4j.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, contactId, model.NodeLabelContact, *contactFields.ExternalSystem)
			if err != nil {
				s.log.Errorf("Error while link contact %s with external system %s: %s", contactId, contactFields.ExternalSystem.ExternalSystemId, err.Error())
				return nil, innerErr
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if createFlow {
				err = s.events.Publisher.PublishEvent(ctx, contactId, model.CONTACT, dto.CreateContact{contactFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateContact"))
				}
				if common_srv.PublishCompletedEvents(options...) {
					s.events.Publisher.PublishEventCompleted(ctx, tenant, contactId, model.CONTACT, utils.NewEventCompletedDetails().WithCreate())
				}
			} else {
				err = s.events.Publisher.PublishEvent(ctx, contactId, model.CONTACT, dto.UpdateContact{contactFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateContact"))
				}
				if common_srv.PublishCompletedEvents(options...) && common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
					s.events.Publisher.PublishEventCompleted(ctx, tenant, contactId, model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())
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

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		contactFields := data_fields.ContactFields{Hide: utils.BoolPtr(true)}
		err = s.neo4j.ContactWriteRepository.SaveContactInTx(ctx, txWithPostCommit.Tx, tenant, contactId, contactFields, false)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to save contact"))
			s.log.Errorf("error while hiding contact %s: %s", contactId, err.Error())
			return nil, err
		}

		err = s.neo4j.OrganizationWriteRepository.RefreshContactCountByContactId(ctx, txWithPostCommit.Tx, tenant, contactId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to refresh contact count by contact id"))
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.events.Publisher.PublishEvent(ctx, contactId, model.CONTACT, dto.HideContact{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message HideContact"))
			}

			s.events.Publisher.PublishEventCompleted(ctx, tenant, contactId, model.CONTACT, utils.NewEventCompletedDetails().WithDelete())
			return nil
		})
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

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

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		contactFields := data_fields.ContactFields{Hide: utils.BoolPtr(false)}
		err = s.neo4j.ContactWriteRepository.SaveContactInTx(ctx, txWithPostCommit.Tx, tenant, contactId, contactFields, false)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to save contact"))
			s.log.Errorf("error while showing contact %s: %s", contactId, err.Error())
			return nil, err
		}
		err = s.neo4j.OrganizationWriteRepository.RefreshContactCountByContactId(ctx, txWithPostCommit.Tx, tenant, contactId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to refresh contact count by contact id"))
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.events.Publisher.PublishEvent(ctx, contactId, model.CONTACT, dto.ShowContact{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message ShowContact"))
			}

			s.events.Publisher.PublishEventCompleted(ctx, tenant, contactId, model.CONTACT, utils.NewEventCompletedDetails().WithCreate())

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

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

	contactDbNode, err := s.neo4j.ContactReadRepository.GetContact(ctx, tenant, contactId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while getting contact %s: %s", contactId, err.Error())
		return nil, err
	}

	return neo4jmapper.MapDbNodeToContactEntity(contactDbNode), nil
}

func (s *contactService) GetContactsByIds(ctx context.Context, contactIds []string) ([]*neo4jentity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactsByIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactIds", fmt.Sprintf("%v", contactIds)))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	contactDbNodes, err := s.neo4j.ContactReadRepository.GetContacts(ctx, tenant, contactIds)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while getting contacts %s", err.Error())
		return nil, err
	}

	contactEntities := make([]*neo4jentity.ContactEntity, 0)
	for _, contactDbNode := range contactDbNodes {
		contactEntities = append(contactEntities, neo4jmapper.MapDbNodeToContactEntity(contactDbNode))
	}

	return contactEntities, nil
}

func (s *contactService) LinkContactWithOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId, organizationId, jobTitle, description, source string, primary bool, startedAt, endedAt *time.Time) error {
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

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		var innerErr error
		// validate contact exists
		exists, innerErr := s.neo4j.CommonReadRepository.ExistsByIdInTx(ctx, *txWithPostCommit.Tx, tenant, contactId, model.NodeLabelContact)
		if innerErr != nil || !exists {
			innerErr = errors.New("contact not found")
			tracing.TraceErr(span, innerErr)
			return nil, innerErr
		}

		// validate organization exists
		exists, innerErr = s.neo4j.CommonReadRepository.ExistsByIdInTx(ctx, *txWithPostCommit.Tx, tenant, organizationId, model.NodeLabelOrganization)
		if innerErr != nil || !exists {
			innerErr = errors.New("organization not found")
			tracing.TraceErr(span, innerErr)
			return nil, innerErr
		}

		if startedAt != nil && startedAt.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
			startedAt = nil
		}
		if endedAt != nil && endedAt.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
			endedAt = nil
		}

		jobRoleData := data_fields.JobRoleFields{
			Description: utils.StringPtr(description),
			JobTitle:    utils.StringPtr(jobTitle),
			Primary:     utils.BoolPtr(primary),
			StartedAt:   startedAt,
			EndedAt:     endedAt,
			Source:      utils.StringPtr(neo4jmodel.GetSource(source)),
			AppSource:   utils.StringPtr(common.GetAppSourceFromContext(ctx)),
		}

		_, innerErr = s.jobrole.Save(ctx, txWithPostCommit, nil, utils.StringPtr(contactId), utils.StringPtr(organizationId), jobRoleData)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to link contact with organization"))
			return nil, innerErr
		}

		// set primary job role
		primaryOrgId := &organizationId
		if endedAt != nil {
			primaryOrgId = nil
		}
		err = s.SetPrimaryJobRole(ctx, txWithPostCommit, contactId, primaryOrgId)
		if err != nil {
			return nil, err
		}

		err = s.neo4j.OrganizationWriteRepository.RefreshContactCountByOrgId(ctx, txWithPostCommit.Tx, tenant, organizationId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to refresh contact count by organization id"))
		}

		// reset contact enrich attempts
		_ = s.neo4j.ContactWriteRepository.ResetEnrichAttempts(ctx, txWithPostCommit.Tx, tenant, contactId)

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			s.events.Publisher.PublishEventCompleted(ctx, tenant, contactId, model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())
			s.events.Publisher.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

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

			innerErr = s.events.Publisher.PublishEvent(ctx, contactId, model.CONTACT, dtoData)
			if innerErr != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddContactToOrganization for contact"))
			}
			innerErr = s.events.Publisher.PublishEvent(ctx, organizationId, model.ORGANIZATION, dtoData)
			if innerErr != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddContactToOrganization for organization"))
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
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
	contacts, err := s.neo4j.ContactReadRepository.GetContactsByLinkedIn(ctx, common.GetTenantFromContext(ctx), url, alias, externalId)
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

	contacts, err := s.neo4j.ContactReadRepository.GetContactsWithEmail(ctx, common.GetTenantFromContext(ctx), email)
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

func (s *contactService) CreateContactByLinkedIn(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkedInUrl string, options ...common_srv.ServiceOptions) (string, error) {
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
			linkedInAlreadyUsed, existingContactId, err := s.CheckContactExistsWithLinkedIn(ctx, linkedInUrl, "", "")
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

				_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
					if contactByLinkedInEntity.IsHidden() {
						err = s.ShowContact(ctx, txWithPostCommit, existingContactId)
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "unable to show contact"))
							return "", err
						}
					} else {
						// just update contact' updatedAt
						err = s.neo4j.CommonWriteRepository.TouchEntity(ctx, txWithPostCommit.Tx, tenant, model.NodeLabelContact, existingContactId)
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "error on updating contact updatedAt"))
						}
						txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
							if common_srv.PublishCompletedEvents(options...) {
								s.events.Publisher.PublishEventCompleted(ctx, tenant, existingContactId, model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())
							}
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

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		createdContactId, err = s.Save(ctx, txWithPostCommit, nil, data_fields.ContactFields{}, false, options...)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create contact"))
			return "", err
		}
		if (neo4jentity.SocialEntity{Url: linkedInUrl}).IsLinkedin() {
			_, err := s.social.AddSocialToEntity(ctx, txWithPostCommit,
				common_srv.LinkWith{
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

func (s *contactService) CreateContactByEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, email string, options ...common_srv.ServiceOptions) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.CreateContactByEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("email", email)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// check email is valid
	if email == "" {
		err = errors.New("email is required")
		span.LogKV("response.error", err.Error())
		return "", err
	}
	mailvalidate := mailsherpa.ValidateEmailSyntax(email)
	if !mailvalidate.IsValid {
		err = errors.New("email is not valid")
		span.LogKV("response.error", err.Error())
		return "", err
	}
	if mailvalidate.IsRoleAccount {
		err = errors.New("email is role account")
		span.LogKV("response.error", err.Error())
		return "", err
	}
	if mailvalidate.IsSystemGenerated {
		err = errors.New("email is system generated")
		span.LogKV("response.error", err.Error())
		return "", err
	}

	// Reject contact creation if email is already used by another contact, return existing contact id
	emailAlreadyUsed, existingContactId, err := s.CheckContactExistsWithEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to check contact exists by email"))
		return "", err
	}
	if emailAlreadyUsed {
		contactByEmailEntity, err := s.GetContactById(ctx, existingContactId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to get contact by id"))
			return "", err
		}
		_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
			if contactByEmailEntity.IsHidden() {
				err = s.ShowContact(ctx, txWithPostCommit, existingContactId)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to show contact"))
					return "", err
				}
			} else {
				// just update contact' updatedAt
				err = s.neo4j.CommonWriteRepository.TouchEntity(ctx, txWithPostCommit.Tx, tenant, model.NodeLabelContact, existingContactId)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "error on updating contact updatedAt"))
				}
				txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
					if common_srv.PublishCompletedEvents(options...) {
						s.events.Publisher.PublishEventCompleted(ctx, tenant, existingContactId, model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())
					}
					return nil
				})
			}
			return nil, nil
		})
		return existingContactId, nil
	}

	createdContactId := ""

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		var innerErr error

		contactFields := data_fields.ContactFields{}
		parsedEmail, innerErr := emailparser.Parse(email)
		if innerErr != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to parse email"))
		}
		if parsedEmail.FirstName != "" {
			contactFields.FirstName = utils.StringPtr(utils.CleanName(parsedEmail.FirstName))
		}
		if parsedEmail.LastName != "" {
			contactFields.LastName = utils.StringPtr(utils.CleanName(parsedEmail.LastName))
		}

		createdContactId, innerErr = s.Save(ctx, txWithPostCommit, nil, contactFields, false, options...)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create contact"))
			return "", innerErr
		}

		_, innerErr = s.email.Merge(ctx, txWithPostCommit, tenant, interfaces.EmailFields{Email: email}, &common_srv.LinkWith{
			Id:   createdContactId,
			Type: model.CONTACT,
		})
		if innerErr != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to add email to contact"))
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create contact and link with email"))
		return createdContactId, err
	}

	return createdContactId, nil
}

func (s *contactService) SetPrimaryJobRole(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId string, primaryOrganizationId *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.SetPrimaryJobRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, contactId)
	span.LogFields(log.String("primaryOrganizationId", utils.IfNotNilString(primaryOrganizationId)))

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// get job roles with org id for contact
		jobRolesWithOrgId, err := s.neo4j.JobRoleReadRepository.GetAllForContactWithOrganizationId(ctx, txWithPostCommit.Tx, common.GetTenantFromContext(ctx), contactId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		// no data found, return
		if len(jobRolesWithOrgId) == 0 {
			return nil, nil
		}

		primaryJobRoleEntities := make([]*neo4jentity.JobRoleEntity, 0)
		allJobRoleEntities := make([]*neo4jentity.JobRoleEntity, 0)
		var jobRoleEntityForPrimaryOrganization *neo4jentity.JobRoleEntity
		for _, jobRoleWithOrgId := range jobRolesWithOrgId {
			jobRoleEntity := neo4jmapper.MapDbNodeToJobRoleEntity(jobRoleWithOrgId.Node)
			if primaryOrganizationId != nil && jobRoleWithOrgId.LinkedNodeId == *primaryOrganizationId {
				jobRoleEntityForPrimaryOrganization = jobRoleEntity
			}
			if jobRoleEntity.Primary {
				primaryJobRoleEntities = append(primaryJobRoleEntities, jobRoleEntity)
			}
			allJobRoleEntities = append(allJobRoleEntities, jobRoleEntity)
		}

		jobRoleIdToBeSetPrimary := ""

		// check 1 - if job role for primary organization found, set it primary
		if jobRoleEntityForPrimaryOrganization != nil {
			if !jobRoleEntityForPrimaryOrganization.Primary {
				jobRoleIdToBeSetPrimary = jobRoleEntityForPrimaryOrganization.Id
			}
		} else {
			if len(primaryJobRoleEntities) == 1 {
				// single primary job role already exists, return
				return nil, nil
			}
			var selectedJobRoleEntity *neo4jentity.JobRoleEntity
			for _, jobRoleEntity := range allJobRoleEntities {
				if jobRoleIdToBeSetPrimary == "" {
					jobRoleIdToBeSetPrimary = jobRoleEntity.Id
					selectedJobRoleEntity = jobRoleEntity
				} else {
					if utils.IsAfter(jobRoleEntity.EndedAt, selectedJobRoleEntity.EndedAt) {
						jobRoleIdToBeSetPrimary = jobRoleEntity.Id
						selectedJobRoleEntity = jobRoleEntity
					}
				}
			}

		}

		if jobRoleIdToBeSetPrimary != "" {
			// set it primary,
			err = s.neo4j.JobRoleWriteRepository.SetJobRolePrimaryInTx(ctx, txWithPostCommit.Tx, common.GetTenantFromContext(ctx), jobRoleIdToBeSetPrimary)
			if err != nil {
				return nil, err
			}
			// set other non-primary job roles as non-primary
			if len(allJobRoleEntities) > 1 {
				err = s.neo4j.JobRoleWriteRepository.SetOtherJobRolesForContactNonPrimaryInTx(ctx, txWithPostCommit.Tx, common.GetTenantFromContext(ctx), contactId, jobRoleIdToBeSetPrimary)
				if err != nil {
					return nil, err
				}
			}
		}

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *contactService) GetFirstContactByEmail(ctx context.Context, email string) (*neo4jentity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetFirstContactByEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, nil
	}

	contactDbNodes, err := s.neo4j.ContactReadRepository.GetContactsWithEmail(ctx, common.GetContext(ctx).Tenant, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if len(contactDbNodes) == 0 {
		return nil, nil
	}
	return neo4jmapper.MapDbNodeToContactEntity(contactDbNodes[0]), nil
}
