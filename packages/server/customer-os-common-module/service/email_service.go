package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type EmailFields struct {
	Email     string                 `json:"email"`
	Source    neo4jentity.DataSource `json:"source"`
	AppSource string                 `json:"appSource"`
	Primary   bool                   `json:"primary"`
}

type emailService struct {
	services *Services
}

type EmailService interface {
	Merge(ctx context.Context, tenant string, emailFields EmailFields, linkWith *LinkWith) (*string, error)
	ReplaceEmail(ctx context.Context, previousEmail string, emailFields EmailFields, linkWith LinkWith) (*string, error)
	UnlinkEmail(ctx context.Context, email, appSource string, linkWith LinkWith) error
	DeleteOrphanEmail(ctx context.Context, tenant, emailId, appSource string) error
	GetAllEmailsForEntityIds(ctx context.Context, tenant string, entityType commonmodel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error)
	SetPrimary(ctx context.Context, email string, forEntity LinkWith) error
	GetPrimaryEmailForEntityId(ctx context.Context, entityType commonmodel.EntityType, entityId string) (*neo4jentity.EmailEntity, error)
	GetPrimaryEmailsForEntityIds(ctx context.Context, entityType commonmodel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error)

	linkEmail(ctx context.Context, emailId, email, appSource string, primary bool, linkWith LinkWith) error
}

func NewEmailService(services *Services) EmailService {
	return &emailService{
		services: services,
	}
}

func (s *emailService) Merge(ctx context.Context, tenant string, emailFields EmailFields, linkWith *LinkWith) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.Merge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", emailFields)

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

	// check if email already exists
	emailId, err = s.services.Neo4jRepositories.EmailReadRepository.GetEmailIdIfExists(ctx, tenant, emailFields.Email)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// email not exist, create one
	if emailId == "" {
		emailId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, commonmodel.NodeLabelEmail)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		err = s.services.Neo4jRepositories.EmailWriteRepository.CreateEmail(ctx, tenant, emailId, neo4jrepository.EmailCreateFields{
			RawEmail:  emailFields.Email,
			CreatedAt: createdAt,
			Source:    emailFields.Source,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		span.LogFields(log.Bool("email.created", true))
		span.LogFields(log.String("email.id", emailId))

		// send event to register email in eventstore
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
			return s.services.GrpcClients.EmailClient.UpsertEmailV2(ctx, &emailpb.UpsertEmailRequest{
				Tenant:         tenant,
				EmailId:        emailId,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				RawEmail:       emailFields.Email,
				CreatedAt:      utils.ConvertTimeToTimestampPtr(&createdAt),
				SourceFields: &commonpb.SourceFields{
					Source:    emailFields.Source.String(),
					AppSource: emailFields.AppSource,
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to upsert email"))
		}

		// send email event to rabbit mq
		err = s.services.RabbitMQService.Publish(ctx, emailId, commonmodel.NodeLabelEmail, dto.NewRegisterEmailEvent(emailFields.Email, emailFields.Source.String()))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddEmailEvent"))
		}
	} else {
		span.LogFields(log.Bool("email.created", false))
	}

	if linkWith != nil && linkWith.Id != "" && linkWith.Type != "" {
		err = s.linkEmail(ctx, emailId, emailFields.Email, emailFields.AppSource, emailFields.Primary, *linkWith)
		if err != nil {
			tracing.TraceErr(span, err)
			return &emailId, err
		}
	}

	return &emailId, nil
}

func (s *emailService) ReplaceEmail(ctx context.Context, previousEmail string, emailFields EmailFields, linkWith LinkWith) (*string, error) {
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
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel())
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

	if previousEmail != "" {
		err := s.UnlinkEmail(ctx, previousEmail, emailFields.AppSource, linkWith)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unlink email"))
		}
	}

	return s.Merge(ctx, tenant, emailFields, &linkWith)
}

func (s *emailService) linkEmail(ctx context.Context, emailId, email, appSource string, primary bool, linkWith LinkWith) error {
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

	// check linked entity exists
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
		return err
	}
	if !exists {
		err = errors.Errorf("linked entity %s with id %s not found", linkWith.Type.String(), linkWith.Id)
		tracing.TraceErr(span, err)
		return err
	}

	// check if email is already linked to entity, if so, skip linking
	alreadyLinked, err := s.services.Neo4jRepositories.EmailReadRepository.IsLinkedToEntityByEmailAddress(ctx, tenant, emailId, linkWith.Id, linkWith.Type)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check if email is already linked to entity"))
	}
	if alreadyLinked {
		span.LogFields(log.Bool("email.alreadyLinked", true))
		return nil
	}

	switch linkWith.Type.String() {
	case commonmodel.CONTACT.String():
		// if contact has no emails yet, set this one as primary
		dbResults, err := s.services.Neo4jRepositories.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, commonmodel.CONTACT, []string{linkWith.Id})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get all emails for contact"))
		} else if len(dbResults) == 0 {
			span.LogFields(log.Bool("firstEmailForContact", true))
			primary = true
		}
		err = s.services.Neo4jRepositories.EmailWriteRepository.LinkWithContact(ctx, tenant, linkWith.Id, emailId, primary)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.LinkWithContact"))
			return err
		}
	case commonmodel.USER.String():
		err = s.services.Neo4jRepositories.EmailWriteRepository.LinkWithUser(ctx, tenant, linkWith.Id, emailId, primary)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.LinkWithUser"))
			return err
		}
	case commonmodel.ORGANIZATION.String():
		err = s.services.Neo4jRepositories.EmailWriteRepository.LinkWithOrganization(ctx, tenant, linkWith.Id, emailId, primary)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.LinkWithOrganization"))
			return err
		}
	default:
		tracing.TraceErr(span, errors.New("unsupported linkWith type "+linkWith.Type.String()))
		return errors.New("unsupported linkWith type " + linkWith.Type.String())
	}

	// publish event to rabbit mq
	err = s.services.RabbitMQService.Publish(ctx, linkWith.Id, linkWith.Type, dto.NewAddEmailEvent(email, primary))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddEmailEvent"))
	}

	// publish event to eventstore for completion
	utils.EventCompleted(ctx, common.GetTenantFromContext(ctx), linkWith.Type.String(), linkWith.Id, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return err
}

func (s *emailService) UnlinkEmail(ctx context.Context, email, appSource string, linkWith LinkWith) error {
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
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
		return err
	}
	if !exists {
		err = errors.Errorf("linked entity %s with id %s not found", linkWith.Type.String(), linkWith.Id)
		tracing.TraceErr(span, err)
		return err
	}

	switch linkWith.Type.String() {
	case commonmodel.CONTACT.String():
		err = s.services.Neo4jRepositories.EmailWriteRepository.UnlinkFromContact(ctx, tenant, linkWith.Id, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.UnlinkFromContact"))
			return err
		}
	case commonmodel.USER.String():
		err = s.services.Neo4jRepositories.EmailWriteRepository.UnlinkFromUser(ctx, tenant, linkWith.Id, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.UnlinkFromUser"))
			return err
		}

	case commonmodel.ORGANIZATION.String():
		err = s.services.Neo4jRepositories.EmailWriteRepository.UnlinkFromOrganization(ctx, tenant, linkWith.Id, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "EmailWriteRepository.UnlinkFromOrganization"))
			return err
		}
	default:
		tracing.TraceErr(span, errors.New("unsupported linkWith type "+linkWith.Type.String()))
		return errors.New("unsupported linkWith type " + linkWith.Type.String())
	}
	// publish event to rabbit mq
	err = s.services.RabbitMQService.Publish(ctx, linkWith.Id, linkWith.Type, dto.NewRemoveEmailEvent(email))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message RemoveEmailEvent"))
	}

	// publish event to eventstore for completion
	utils.EventCompleted(ctx, common.GetTenantFromContext(ctx), linkWith.Type.String(), linkWith.Id, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return err
}

func (s *emailService) DeleteOrphanEmail(ctx context.Context, tenant, emailId, appSource string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.DeleteOrphanEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, emailId)
	span.LogKV("appSource", appSource)

	if tenant == "" {
		tenant = common.GetTenantFromContext(ctx)
	}

	// check if email exists by id
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, emailId, commonmodel.NodeLabelEmail)
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
	isOrphan, err := s.services.Neo4jRepositories.EmailReadRepository.IsOrphanEmail(ctx, tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check if email is orphan"))
		return err
	}
	if !isOrphan {
		err = errors.Errorf("email with id %s is not orphan", emailId)
		tracing.TraceErr(span, err)
		return err
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
		return s.services.GrpcClients.EmailClient.DeleteEmail(ctx, &emailpb.DeleteEmailRequest{
			Tenant:         tenant,
			EmailId:        emailId,
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			AppSource:      appSource,
		})
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to delete email"))
		return err
	}

	return nil
}

func (s *emailService) GetAllEmailsForEntityIds(ctx context.Context, tenant string, entityType commonmodel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetAllEmailsForEntityIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	emailNodes, err := s.services.Neo4jRepositories.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, entityType, entityIds)
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

func (s *emailService) SetPrimary(ctx context.Context, email string, forEntity LinkWith) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.SetPrimary")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("email", email)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if forEntity.Id == "" {
		tracing.TraceErr(span, errors.New("forEntity id is required"))
		return errors.New("forEntity id is required")
	}
	if forEntity.Type == "" {
		tracing.TraceErr(span, errors.New("forEntity type is required"))
		return errors.New("forEntity type is required")
	}

	// check linked entity exists
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), forEntity.Id, forEntity.Type.Neo4jLabel())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
		return err
	}
	if !exists {
		err = errors.Errorf("linked entity %s with id %s not found", forEntity.Type.String(), forEntity.Id)
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.EmailWriteRepository.SetPrimaryForEntity(ctx, common.GetTenantFromContext(ctx), forEntity.Id, email, forEntity.Type)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	utils.EventCompleted(ctx, common.GetTenantFromContext(ctx), forEntity.Type.String(), forEntity.Id, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (s *emailService) GetPrimaryEmailForEntityId(ctx context.Context, entityType commonmodel.EntityType, entityId string) (*neo4jentity.EmailEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetPrimaryEmailForEntityId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("entityId", entityId))

	emailNodes, err := s.services.Neo4jRepositories.EmailReadRepository.GetPrimaryEmailNodesForLinkedEntityIds(ctx, common.GetTenantFromContext(ctx), entityType, []string{entityId})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if len(emailNodes) == 0 {
		return nil, nil
	}

	return mapper.MapDbNodeToEmailEntity(emailNodes[0].Node), nil
}

func (s *emailService) GetPrimaryEmailsForEntityIds(ctx context.Context, entityType commonmodel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetPrimaryEmailsForEntityIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("entityIds", entityIds))

	emailNodes, err := s.services.Neo4jRepositories.EmailReadRepository.GetPrimaryEmailNodesForLinkedEntityIds(ctx, common.GetTenantFromContext(ctx), entityType, entityIds)
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
