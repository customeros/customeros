package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type EmailFields struct {
	Email     string `json:"email"`
	Source    string `json:"source"`
	AppSource string `json:"appSource"`
	Primary   bool   `json:"primary"`
}

type emailService struct {
	services *Services
}

type EmailService interface {
	Merge(ctx context.Context, tenant string, emailFields EmailFields, linkWith *LinkWith) (*string, error)
	ReplaceEmail(ctx context.Context, tenant, previousEmail string, emailFields EmailFields, linkWith LinkWith) (*string, error)
	LinkEmail(ctx context.Context, tenant, emailId, email, appSource string, primary bool, linkWith LinkWith) error
	UnlinkEmail(ctx context.Context, tenant, email, appSource string, linkWith LinkWith) error
	DeleteOrphanEmail(ctx context.Context, tenant, emailId, appSource string) error

	GetAllEmailsForEntityIds(ctx context.Context, tenant string, entityType commonmodel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error)
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
			SourceFields: neo4jmodel.Source{
				AppSource: emailFields.AppSource,
			},
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
					Source:    emailFields.Source,
					AppSource: emailFields.AppSource,
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to upsert email"))
		}

	} else {
		span.LogFields(log.Bool("email.created", false))
	}

	if linkWith != nil && linkWith.Id != "" && linkWith.Type != "" {
		err = s.LinkEmail(ctx, tenant, emailId, emailFields.Email, emailFields.AppSource, emailFields.Primary, *linkWith)
		if err != nil {
			tracing.TraceErr(span, err)
			return &emailId, err
		}
	}

	return &emailId, nil
}

func (s *emailService) ReplaceEmail(ctx context.Context, tenant, previousEmail string, emailFields EmailFields, linkWith LinkWith) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.ReplaceEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", emailFields)
	span.LogKV("previousEmail", previousEmail)

	if previousEmail == emailFields.Email {
		span.LogFields(log.Bool("email.same", true))
		return nil, nil
	}

	if tenant == "" {
		tenant = common.GetTenantFromContext(ctx)
	}

	if previousEmail != "" {
		err := s.UnlinkEmail(ctx, tenant, previousEmail, emailFields.AppSource, linkWith)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unlink email"))
		}
	}

	return s.Merge(ctx, tenant, emailFields, &linkWith)
}

func (s *emailService) LinkEmail(ctx context.Context, tenant, emailId, email, appSource string, primary bool, linkWith LinkWith) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.LinkEmail")
	defer span.Finish()

	if linkWith.Id == "" {
		tracing.TraceErr(span, errors.New("linkWith id is required"))
		return errors.New("linkWith id is required")
	}
	if linkWith.Type == "" {
		tracing.TraceErr(span, errors.New("linkWith type is required"))
		return errors.New("linkWith type is required")
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
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return s.services.GrpcClients.ContactClient.LinkEmailToContact(ctx, &contactpb.LinkEmailToContactGrpcRequest{
				Tenant:         tenant,
				EmailId:        emailId,
				ContactId:      linkWith.Id,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      appSource,
				Primary:        primary,
				Email:          email,
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to link email to contact"))
		}
	case commonmodel.USER.String():
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
			return s.services.GrpcClients.UserClient.LinkEmailToUser(ctx, &userpb.LinkEmailToUserGrpcRequest{
				Tenant:         tenant,
				EmailId:        emailId,
				UserId:         linkWith.Id,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      appSource,
				Primary:        primary,
				Email:          email,
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to link email to user"))
		}
	case commonmodel.ORGANIZATION.String():
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.services.GrpcClients.OrganizationClient.LinkEmailToOrganization(ctx, &organizationpb.LinkEmailToOrganizationGrpcRequest{
				Tenant:         tenant,
				EmailId:        emailId,
				OrganizationId: linkWith.Id,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      appSource,
				Primary:        primary,
				Email:          email,
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to link email to user"))
		}
	default:
		tracing.TraceErr(span, errors.New("unsupported linkWith type %s"+linkWith.Type.String()))
	}

	return err
}

func (s *emailService) UnlinkEmail(ctx context.Context, tenant, email, appSource string, linkWith LinkWith) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.UnlinkEmail")
	defer span.Finish()

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
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return s.services.GrpcClients.ContactClient.UnLinkEmailFromContact(ctx, &contactpb.UnLinkEmailFromContactGrpcRequest{
				Tenant:         tenant,
				ContactId:      linkWith.Id,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      appSource,
				Email:          email,
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to link email to contact"))
		}
	case commonmodel.USER.String():
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
			return s.services.GrpcClients.UserClient.UnLinkEmailFromUser(ctx, &userpb.UnLinkEmailFromUserGrpcRequest{
				Tenant:         tenant,
				UserId:         linkWith.Id,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      appSource,
				Email:          email,
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to link email to user"))
		}
	case commonmodel.ORGANIZATION.String():
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.services.GrpcClients.OrganizationClient.UnLinkEmailFromOrganization(ctx, &organizationpb.UnLinkEmailFromOrganizationGrpcRequest{
				Tenant:         tenant,
				OrganizationId: linkWith.Id,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      appSource,
				Email:          email,
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to link email to user"))
		}
	default:
		tracing.TraceErr(span, errors.New("unsupported linkWith type %s"+linkWith.Type.String()))
	}

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
