package service

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commongrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

type EmailService interface {
	CreateEmailAddressViaEvents(ctx context.Context, email, appSource string) (string, error)
	GetAllFor(ctx context.Context, entityType commonModel.EntityType, entityId string) (*neo4jentity.EmailEntities, error)
	GetAllForEntityTypeByIds(ctx context.Context, entityType commonModel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error)
	UpdateEmailFor(ctx context.Context, entityType commonModel.EntityType, entityId string, input model.EmailRelationUpdateInput) error
	DetachFromEntity(ctx context.Context, entityType commonModel.EntityType, entityId, email string) (bool, error)
	DetachFromEntityById(ctx context.Context, entityType commonModel.EntityType, entityId, emailId string) (bool, error)
	DeleteById(ctx context.Context, emailId string) (bool, error)
	GetById(ctx context.Context, emailId string) (*neo4jentity.EmailEntity, error)
	GetByEmailAddress(ctx context.Context, email string) (*neo4jentity.EmailEntity, error)
	Update(ctx context.Context, input model.EmailUpdateAddressInput) error
}

type emailService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
	grpcClients  *grpc_client.Clients
}

func NewEmailService(log logger.Logger, repositories *repository.Repositories, services *Services, grpcClients *grpc_client.Clients) EmailService {
	return &emailService{
		log:          log,
		repositories: repositories,
		services:     services,
		grpcClients:  grpcClients,
	}
}

func (s *emailService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *emailService) GetAllFor(ctx context.Context, entityType commonModel.EntityType, entityId string) (*neo4jentity.EmailEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetAllFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId))

	records, err := s.repositories.EmailRepository.GetAllFor(ctx, common.GetContext(ctx).Tenant, entityType, entityId)
	if err != nil {
		return nil, err
	}

	emailEntities := make(neo4jentity.EmailEntities, 0, len(records))
	for _, dbRecord := range records {
		emailEntity := neo4jmapper.MapDbNodeToEmailEntity(utils.ToPtr(dbRecord.Values[0].(dbtype.Node)))
		s.addDbRelationshipToEmailEntity(dbRecord.Values[1].(dbtype.Relationship), emailEntity)
		emailEntities = append(emailEntities, *emailEntity)
	}

	return &emailEntities, nil
}

func (s *emailService) GetAllForEntityTypeByIds(ctx context.Context, entityType commonModel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetAllForEntityTypeByIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("entityIds", entityIds))

	emails, err := s.repositories.Neo4jRepositories.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, common.GetContext(ctx).Tenant, entityType, entityIds)
	if err != nil {
		return nil, err
	}

	emailEntities := make(neo4jentity.EmailEntities, 0, len(emails))
	for _, v := range emails {
		emailEntity := neo4jmapper.MapDbNodeToEmailEntity(v.Node)
		s.addDbRelationshipToEmailEntity(*v.Relationship, emailEntity)
		emailEntity.DataloaderKey = v.LinkedNodeId
		emailEntities = append(emailEntities, *emailEntity)
	}
	return &emailEntities, nil
}

func (s *emailService) UpdateEmailFor(ctx context.Context, entityType commonModel.EntityType, entityId string, input model.EmailRelationUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.UpdateEmailFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", string(entityType)), log.String("entityId", entityId))

	emailEntity, err := s.GetById(ctx, input.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if entityType == commonModel.CONTACT {
		contactExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), entityId, commonModel.NodeLabelContact)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if !contactExists {
			err = errors.New("Contact not found")
			tracing.TraceErr(span, err)
			return err
		}

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return s.grpcClients.ContactClient.LinkEmailToContact(ctx, &contactpb.LinkEmailToContactGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				ContactId:      entityId,
				EmailId:        emailEntity.Id,
				Primary:        utils.IfNotNilBool(input.Primary),
				Label:          utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      constants.AppSourceCustomerOsApi,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			graphql.AddErrorf(ctx, "Failed to add email %s to contact %s", input.ID, entityId)
			return err
		}
	} else if entityType == commonModel.ORGANIZATION {
		organizationExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), entityId, commonModel.NodeLabelOrganization)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if !organizationExists {
			err = errors.New("Organization not found")
			tracing.TraceErr(span, err)
			return err
		}

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.grpcClients.OrganizationClient.LinkEmailToOrganization(ctx, &organizationpb.LinkEmailToOrganizationGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				OrganizationId: entityId,
				EmailId:        emailEntity.Id,
				Primary:        utils.IfNotNilBool(input.Primary),
				Label:          utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      constants.AppSourceCustomerOsApi,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			graphql.AddErrorf(ctx, "Failed to add email %s to organization %s", input.ID, entityId)
			return err
		}
	} else if entityType == commonModel.USER {
		userExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), entityId, commonModel.NodeLabelUser)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if !userExists {
			err = errors.New("User not found")
			tracing.TraceErr(span, err)
			return err
		}

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
			return s.grpcClients.UserClient.LinkEmailToUser(ctx, &userpb.LinkEmailToUserGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				UserId:         entityId,
				EmailId:        emailEntity.Id,
				Primary:        utils.IfNotNilBool(input.Primary),
				Label:          utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      constants.AppSourceCustomerOsApi,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			graphql.AddErrorf(ctx, "Failed to add email %s to user %s", input.ID, entityId)
			return err
		}
	}

	return nil
}

func (s *emailService) DetachFromEntity(ctx context.Context, entityType commonModel.EntityType, entityId, email string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.DetachFromEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email), log.String("entityId", entityId), log.String("entityType", string(entityType)))

	err := s.repositories.EmailRepository.RemoveRelationship(ctx, entityType, common.GetTenantFromContext(ctx), entityId, email)

	if entityType == commonModel.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpoint(ctx, entityId)
	} else if entityType == commonModel.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, entityId)
	}

	return err == nil, err
}

func (s *emailService) DetachFromEntityById(ctx context.Context, entityType commonModel.EntityType, entityId, emailId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.DetachFromEntityById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailId", emailId), log.String("entityId", entityId), log.String("entityType", string(entityType)))

	err := s.repositories.EmailRepository.RemoveRelationshipById(ctx, entityType, common.GetTenantFromContext(ctx), entityId, emailId)

	if entityType == commonModel.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpoint(ctx, entityId)
	} else if entityType == commonModel.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, entityId)
	}

	return err == nil, err
}

func (s *emailService) DeleteById(ctx context.Context, emailId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.DeleteById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailId", emailId))

	err := s.repositories.EmailRepository.DeleteById(ctx, common.GetTenantFromContext(ctx), emailId)
	return err == nil, err
}

func (s *emailService) GetById(ctx context.Context, emailId string) (*neo4jentity.EmailEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailId", emailId))

	emailNode, err := s.repositories.Neo4jRepositories.EmailReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), emailId)
	if err != nil {
		return nil, err
	}
	var emailEntity = neo4jmapper.MapDbNodeToEmailEntity(emailNode)
	return emailEntity, nil
}

func (s *emailService) GetByEmailAddress(ctx context.Context, email string) (*neo4jentity.EmailEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetByEmailAddress")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	emailNode, err := s.repositories.Neo4jRepositories.EmailReadRepository.GetFirstByEmail(ctx, common.GetTenantFromContext(ctx), email)
	if err != nil {
		return nil, err
	}
	var emailEntity = neo4jmapper.MapDbNodeToEmailEntity(emailNode)
	return emailEntity, nil
}

func (s *emailService) CreateEmailAddressViaEvents(ctx context.Context, email, appSource string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.CreateEmailAddressViaEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email), log.String("appSource", appSource))

	email = strings.TrimSpace(email)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
		return s.grpcClients.EmailClient.UpsertEmail(ctx, &emailpb.UpsertEmailGrpcRequest{
			Tenant:   common.GetTenantFromContext(ctx),
			RawEmail: email,
			SourceFields: &commongrpc.SourceFields{
				Source:    string(neo4jentity.DataSourceOpenline),
				AppSource: utils.StringFirstNonEmpty(appSource, constants.AppSourceCustomerOsApi),
			},
			LoggedInUserId: common.GetUserIdFromContext(ctx),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, commonModel.NodeLabelEmail, span)
	return response.Id, nil
}

func (s *emailService) addDbRelationshipToEmailEntity(relationship dbtype.Relationship, emailEntity *neo4jentity.EmailEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	emailEntity.Primary = utils.GetBoolPropOrFalse(props, "primary")
	emailEntity.Label = utils.GetStringPropOrEmpty(props, "label")
}

func (s *emailService) Update(ctx context.Context, input model.EmailUpdateAddressInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	emailEntity, err := s.GetById(ctx, input.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if emailEntity.Email == input.Email || emailEntity.RawEmail == input.Email {
		err = errors.New("Email address is the same as the current one")
		tracing.TraceErr(span, err)
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
		return s.grpcClients.EmailClient.UpsertEmail(ctx, &emailpb.UpsertEmailGrpcRequest{
			Id:       input.ID,
			Tenant:   common.GetTenantFromContext(ctx),
			RawEmail: input.Email,
			SourceFields: &commongrpc.SourceFields{
				Source:    string(neo4jentity.DataSourceOpenline),
				AppSource: constants.AppSourceCustomerOsApi,
			},
			LoggedInUserId: common.GetUserIdFromContext(ctx),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing %s", err.Error())
		return err
	}

	return err
}
