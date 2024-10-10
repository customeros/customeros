package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type EmailService interface {
	GetAllFor(ctx context.Context, entityType commonModel.EntityType, entityId string) (*neo4jentity.EmailEntities, error)
	GetAllForEntityTypeByIds(ctx context.Context, entityType commonModel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error)
	GetById(ctx context.Context, emailId string) (*neo4jentity.EmailEntity, error)
	GetByEmailAddress(ctx context.Context, email string) (*neo4jentity.EmailEntity, error)
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

	if emailNode == nil {
		return nil, nil
	}

	return neo4jmapper.MapDbNodeToEmailEntity(emailNode), nil
}

func (s *emailService) addDbRelationshipToEmailEntity(relationship dbtype.Relationship, emailEntity *neo4jentity.EmailEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	emailEntity.Primary = utils.GetBoolPropOrFalse(props, "primary")
}
