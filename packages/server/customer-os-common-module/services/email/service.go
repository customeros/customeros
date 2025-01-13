package email

import (
	"context"
	"fmt"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	common_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type emailService struct {
	neo4j   *neoRepo.Repositories
	events  *events.EventsService
	contact interfaces.ContactService
	org     interfaces.OrganizationService
}

func NewEmailService(neo4j *neoRepo.Repositories, events *events.EventsService, contact interfaces.ContactService, org interfaces.OrganizationService) interfaces.EmailService {
	return &emailService{
		neo4j:   neo4j,
		events:  events,
		contact: contact,
		org:     org,
	}
}

func (s *emailService) Merge(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant string, emailFields interfaces.EmailFields, linkWith *common_srv.LinkWith) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.Merge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", emailFields)

	if common.GetTenantFromContext(ctx) == "" {
		tracing.TraceErr(span, errors.New("tenant is missing in context"))
	}

	if tenant == "" {
		tenant = common.GetTenantFromContext(ctx)
	}
	if common.GetTenantFromContext(ctx) == "" {
		ctx = common.SetTenantInContext(ctx, tenant)
	}

	emailId := ""
	var err error
	createdAt := utils.Now()

	if emailFields.Email == "" {
		return nil, nil
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// check if email already exists
		emailId, err = s.neo4j.EmailReadRepository.GetEmailIdIfExists(ctx, txWithPostCommit.Tx, tenant, emailFields.Email)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		// email not exist, create one
		if emailId == "" {
			emailId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, commonmodel.NodeLabelEmail)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			err = s.neo4j.EmailWriteRepository.CreateEmail(ctx, txWithPostCommit.Tx, tenant, emailId, neo4jrepository.EmailCreateFields{
				RawEmail:  emailFields.Email,
				CreatedAt: createdAt,
				Source:    emailFields.Source,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
				// send email event to rabbit mq
				err = s.events.Publisher.PublishEvent(ctx, emailId, commonmodel.NodeLabelEmail, dto.NewRegisterEmailEvent(emailFields.Email, emailFields.Source.String()))
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddEmailEvent"))
				}
				return nil
			})
		}

		if linkWith != nil && linkWith.Id != "" && linkWith.Type != "" {
			err = s.linkEmail(ctx, txWithPostCommit, emailId, emailFields.Email, emailFields.AppSource, emailFields.Primary, *linkWith)
			if err != nil {
				tracing.TraceErr(span, err)
				return &emailId, err
			}
		}

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.String("result.emailId", emailId))

	return &emailId, nil
}

func (s *emailService) ReplaceEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, previousEmail string, emailFields interfaces.EmailFields, linkWith common_srv.LinkWith) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.ReplaceEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", emailFields)
	span.LogKV("previousEmail", previousEmail)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	// check if linkWith is valid
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
		return nil, err
	}
	if !exists {
		err = errors.Errorf("linked entity %s with id %s not found", linkWith.Type.String(), linkWith.Id)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if previousEmail == emailFields.Email {
		span.LogFields(log.Bool("email.same", true))
		return nil, nil
	}

	// check if email is already linked to other entity of the same type
	if linkWith.Type == commonmodel.CONTACT {
		emailUsed, existingContactId, err := s.contact.CheckContactExistsWithEmail(ctx, emailFields.Email)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if emailUsed && existingContactId != linkWith.Id {
			span.LogFields(log.String("result.error", fmt.Sprintf("email %s already used by contact %s", emailFields.Email, existingContactId)))
			err = coserrors.ErrEmailUsed
			return nil, err
		}
	} else if linkWith.Type == commonmodel.ORGANIZATION {
		emailUsed, existingOrganizationId, err := s.org.CheckOrganizationExistsWithEmail(ctx, emailFields.Email)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if emailUsed && existingOrganizationId != linkWith.Id {
			err = errors.Errorf("email %s already used by organization %s", emailFields.Email, existingOrganizationId)
			return nil, err
		}
	}

	var emailId *string
	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		if previousEmail != "" {
			err = s.UnlinkEmail(ctx, txWithPostCommit, previousEmail, emailFields.AppSource, linkWith)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to unlink email"))
			}
		}

		emailId, err = s.Merge(ctx, txWithPostCommit, tenant, emailFields, &linkWith)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to merge email"))
		}
		return nil, err
	})

	return emailId, err
}

func (s *emailService) linkEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, emailId, email, appSource string, primary bool, linkWith common_srv.LinkWith) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.LinkEmail")
	defer span.Finish()

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	tenant := common.GetTenantFromContext(ctx)

	if linkWith.Id == "" {
		tracing.TraceErr(span, errors.New("linkWith id is required"))
		return errors.New("linkWith id is required")
	}
	if linkWith.Type == "" {
		tracing.TraceErr(span, errors.New("linkWith type is required"))
		return errors.New("linkWith type is required")
	}

	// set default values
	if appSource != "" {
		common.SetAppSourceInContext(ctx, appSource)
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// check linked entity exists
		exists, err := s.neo4j.CommonReadRepository.ExistsByIdInTx(ctx, *txWithPostCommit.Tx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
			return nil, err
		}
		if !exists {
			err = errors.Errorf("linked entity %s with id %s not found", linkWith.Type.String(), linkWith.Id)
			tracing.TraceErr(span, err)
			return nil, err
		}

		// check if email is already linked to entity, if so, skip linking
		alreadyLinked, err := s.neo4j.EmailReadRepository.IsLinkedToEntityByEmailAddress(ctx, txWithPostCommit.Tx, tenant, emailId, linkWith.Id, linkWith.Type)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to check if email is already linked to entity"))
		}
		if alreadyLinked {
			span.LogFields(log.Bool("email.alreadyLinked", true))
			return nil, nil
		}

		// check if email is already linked to other entity of same type
		if linkWith.Type == commonmodel.CONTACT {
			emailUsed, existingContactId, err := s.contact.CheckContactExistsWithEmail(ctx, email)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			if emailUsed && existingContactId != linkWith.Id {
				span.LogFields(log.String("result.error", fmt.Sprintf("email %s already used by contact %s", email, existingContactId)))
				return nil, coserrors.ErrEmailUsed
			}
		} else if linkWith.Type == commonmodel.ORGANIZATION {
			emailUsed, existingOrganizationId, err := s.org.CheckOrganizationExistsWithEmail(ctx, email)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			if emailUsed && existingOrganizationId != linkWith.Id {
				return nil, errors.Errorf("email %s already used by organization %s", email, existingOrganizationId)
			}
		}

		switch linkWith.Type.String() {
		case commonmodel.CONTACT.String():
			// if contact has no emails yet, set this one as primary
			dbResults, err := s.neo4j.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, commonmodel.CONTACT, []string{linkWith.Id})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to get all emails for contact"))
			} else if len(dbResults) == 0 {
				span.LogFields(log.Bool("firstEmailForContact", true))
				primary = true
			}
			err = s.neo4j.EmailWriteRepository.LinkWithContact(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, emailId, primary)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			// reset contact enrich attempts
			_ = s.neo4j.ContactWriteRepository.ResetEnrichAttempts(ctx, txWithPostCommit.Tx, tenant, linkWith.Id)
		case commonmodel.USER.String():
			err = s.neo4j.EmailWriteRepository.LinkWithUser(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, emailId, primary)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.LinkWithUser"))
				return nil, err
			}
		case commonmodel.ORGANIZATION.String():
			err = s.neo4j.EmailWriteRepository.LinkWithOrganization(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, emailId, primary)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.LinkWithOrganization"))
				return nil, err
			}
		default:
			tracing.TraceErr(span, errors.New("unsupported linkWith type "+linkWith.Type.String()))
			return nil, errors.New("unsupported linkWith type " + linkWith.Type.String())
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// publish event to rabbit mq
			err = s.events.Publisher.PublishEvent(ctx, linkWith.Id, linkWith.Type, dto.NewAddEmailEvent(email, primary))
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddEmailEvent"))
			}

			// publish completion event for linked entity
			s.events.Publisher.PublishEventCompleted(ctx, tenant, linkWith.Id, linkWith.Type, utils.NewEventCompletedDetails().WithUpdate())
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (s *emailService) UnlinkEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, email, appSource string, linkWith common_srv.LinkWith) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.UnlinkEmail")
	defer span.Finish()

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	tenant := common.GetTenantFromContext(ctx)

	// set default values
	if appSource != "" {
		common.SetAppSourceInContext(ctx, appSource)
	}

	if linkWith.Id == "" {
		tracing.TraceErr(span, errors.New("linkWith id is required"))
		return errors.New("linkWith id is required")
	}
	if linkWith.Type == "" {
		tracing.TraceErr(span, errors.New("linkWith type is required"))
		return errors.New("linkWith type is required")
	}

	// check linked entity exists
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
		return err
	}
	if !exists {
		err = errors.Errorf("linked entity %s with id %s not found", linkWith.Type.String(), linkWith.Id)
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		switch linkWith.Type.String() {
		case commonmodel.CONTACT.String():
			err = s.neo4j.EmailWriteRepository.UnlinkFromContact(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, email)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.UnlinkFromContact"))
				return nil, err
			}
		case commonmodel.USER.String():
			err = s.neo4j.EmailWriteRepository.UnlinkFromUser(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, email)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.UnlinkFromUser"))
				return nil, err
			}

		case commonmodel.ORGANIZATION.String():
			err = s.neo4j.EmailWriteRepository.UnlinkFromOrganization(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, email)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.UnlinkFromOrganization"))
				return nil, err
			}
		default:
			tracing.TraceErr(span, errors.New("unsupported linkWith type "+linkWith.Type.String()))
			return nil, errors.New("unsupported linkWith type " + linkWith.Type.String())
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// publish event to rabbit mq
			err = s.events.Publisher.PublishEvent(ctx, linkWith.Id, linkWith.Type, dto.NewRemoveEmailEvent(email))
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message RemoveEmailEvent"))
			}

			// publish event for completion
			s.events.Publisher.PublishEventCompleted(ctx, tenant, linkWith.Id, linkWith.Type, utils.NewEventCompletedDetails().WithUpdate())

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (s *emailService) DeleteOrphanEmail(ctx context.Context, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.DeleteOrphanEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, emailId)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// check if email exists by id
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, emailId, commonmodel.NodeLabelEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check if email exists by id"))
		return err
	}
	if !exists {
		err = errors.Errorf("email with id %s not found", emailId)
		tracing.TraceErr(span, err)
		return err
	}

	// check if email is orphan
	isOrphan, err := s.neo4j.EmailReadRepository.IsOrphanEmail(ctx, tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check if email is orphan"))
		return err
	}
	if !isOrphan {
		err = errors.Errorf("email with id %s is not orphan", emailId)
		tracing.TraceErr(span, err)
		return err
	}

	// delete email node
	err = s.neo4j.EmailWriteRepository.DeleteOrphanEmail(ctx, tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to delete orphan email"))
		return err
	}

	// check if email exists by id
	exists, err = s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, emailId, commonmodel.NodeLabelEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check if email exists by id"))
		return err
	}
	if exists {
		span.LogFields(log.Bool("result.deleted", false))
	} else {
		span.LogFields(log.Bool("result.deleted", true))
		err = s.events.Publisher.PublishEvent(ctx, emailId, commonmodel.EMAIL, dto.Delete{})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish event Delete"))
		}
	}

	return nil
}

func (s *emailService) GetAllEmailsForEntityIds(ctx context.Context, tenant string, entityType commonmodel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetAllEmailsForEntityIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	emailNodes, err := s.neo4j.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, entityType, entityIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	emailEntities := make(neo4jentity.EmailEntities, 0, len(emailNodes))
	for _, v := range emailNodes {
		emailEntity := mapper.MapDbNodeToEmailEntity(v.Node)
		emailEntity.DataloaderKey = v.LinkedNodeId
		relationshipProps := utils.GetPropsFromRelationship(*v.Relationship)
		emailEntity.Primary = utils.GetBoolPropOrFalse(relationshipProps, "primary")
		emailEntities = append(emailEntities, *emailEntity)
	}
	return &emailEntities, nil
}

func (s *emailService) SetPrimary(ctx context.Context, email string, forEntity common_srv.LinkWith) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.SetPrimary")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("email", email)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	if forEntity.Id == "" {
		tracing.TraceErr(span, errors.New("forEntity id is required"))
		return errors.New("forEntity id is required")
	}
	if forEntity.Type == "" {
		tracing.TraceErr(span, errors.New("forEntity type is required"))
		return errors.New("forEntity type is required")
	}

	// check linked entity exists
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), forEntity.Id, forEntity.Type.Neo4jLabel())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
		return err
	}
	if !exists {
		err = errors.Errorf("linked entity %s with id %s not found", forEntity.Type.String(), forEntity.Id)
		tracing.TraceErr(span, err)
		return err
	}

	err = s.neo4j.EmailWriteRepository.SetPrimaryForEntity(ctx, common.GetTenantFromContext(ctx), forEntity.Id, email, forEntity.Type)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	s.events.Publisher.PublishEventCompleted(ctx, tenant, forEntity.Id, forEntity.Type, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (s *emailService) GetPrimaryEmailForEntityId(ctx context.Context, entityType commonmodel.EntityType, entityId string) (*neo4jentity.EmailEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetPrimaryEmailForEntityId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("entityId", entityId))

	emailNodes, err := s.neo4j.EmailReadRepository.GetPrimaryEmailNodesForLinkedEntityIds(ctx, common.GetTenantFromContext(ctx), entityType, []string{entityId})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if len(emailNodes) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))
	return mapper.MapDbNodeToEmailEntity(emailNodes[0].Node), nil
}

func (s *emailService) GetPrimaryEmailsForEntityIds(ctx context.Context, entityType commonmodel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetPrimaryEmailsForEntityIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("entityIds", entityIds))

	emailNodes, err := s.neo4j.EmailReadRepository.GetPrimaryEmailNodesForLinkedEntityIds(ctx, common.GetTenantFromContext(ctx), entityType, entityIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	emailEntities := make(neo4jentity.EmailEntities, 0, len(emailNodes))
	for _, v := range emailNodes {
		emailEntity := mapper.MapDbNodeToEmailEntity(v.Node)
		emailEntity.DataloaderKey = v.LinkedNodeId
		emailEntities = append(emailEntities, *emailEntity)
	}
	return &emailEntities, nil
}

func (s *emailService) UpdateEmailValidationDetails(ctx context.Context, emailId string, validationFields data_fields.EmailValidationFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.UpdateEmailValidationDetails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, emailId)
	tracing.LogObjectAsJson(span, "validationFields", validationFields)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	err = s.neo4j.EmailWriteRepository.EmailValidated(ctx, tenant, emailId, validationFields)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (s *emailService) RequestEmailValidation(ctx context.Context, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.RequestEmailValidation")
	defer span.Finish()
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, emailId)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.events.Publisher.PublishEvent(ctx, emailId, commonmodel.EMAIL, dto.RequestValidateEmail{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error publishing email validation request"))
	}
	return err
}
