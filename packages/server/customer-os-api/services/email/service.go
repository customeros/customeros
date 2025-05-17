package api_email

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type emailService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewEmailService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.EmailService {
	return &emailService{
		log:          log,
		repositories: repositories,
	}
}

func (s *emailService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *emailService) GetAllFor(ctx context.Context, entityType commonModel.EntityType, entityId string) (*neo4jentity.EmailEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailService.GetAllFor")
	defer spans.Finish()
	spans.LogKV("entityType", entityType.String(), "entityId", entityId)

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
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailService.GetAllForEntityTypeByIds")
	defer spans.Finish()
	spans.LogKV("entityType", entityType.String(), "entityIds", entityIds)

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
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailService.GetById")
	defer spans.Finish()
	spans.LogKV("emailId", emailId)

	emailNode, err := s.repositories.Neo4jRepositories.EmailReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), emailId)
	if err != nil {
		return nil, err
	}
	emailEntity := neo4jmapper.MapDbNodeToEmailEntity(emailNode)
	return emailEntity, nil
}

func (s *emailService) GetByEmailAddress(ctx context.Context, email string) (*neo4jentity.EmailEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailService.GetByEmailAddress")
	defer spans.Finish()
	spans.LogKV("email", email)

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
