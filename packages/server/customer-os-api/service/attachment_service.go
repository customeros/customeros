package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type AttachmentService interface {
	GetAttachmentById(ctx context.Context, id string) (*entity.AttachmentEntity, error)

	Create(ctx context.Context, newAnalysis *entity.AttachmentEntity, source, sourceOfTruth neo4jentity.DataSource) (*entity.AttachmentEntity, error)
	GetAttachmentsForNode(ctx context.Context, linkedWith repository.LinkedWith, linkedNature *repository.LinkedNature, ids []string) (*entity.AttachmentEntities, error)

	LinkNodeWithAttachment(ctx context.Context, linkedWith repository.LinkedWith, linkedNature *repository.LinkedNature, attachmentId, includedById string) (*dbtype.Node, error)
	UnlinkNodeWithAttachment(ctx context.Context, linkedWith repository.LinkedWith, linkedNature *repository.LinkedNature, attachmentId, includedById string) (*dbtype.Node, error)

	MapDbNodeToAttachmentEntity(node dbtype.Node) *entity.AttachmentEntity
}

type attachmentService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewAttachmentService(log logger.Logger, repositories *repository.Repositories) AttachmentService {
	return &attachmentService{
		log:          log,
		repositories: repositories,
	}
}

func (s *attachmentService) LinkNodeWithAttachment(ctx context.Context, linkedWith repository.LinkedWith, linkedNature *repository.LinkedNature, attachmentId, includedById string) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	node, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		return s.repositories.AttachmentRepository.LinkWithXXIncludesAttachmentInTx(ctx, tx, tenant, linkedWith, linkedNature, attachmentId, includedById)
	})
	if err != nil {
		return nil, err
	}
	return node.(*dbtype.Node), err
}

func (s *attachmentService) UnlinkNodeWithAttachment(ctx context.Context, linkedWith repository.LinkedWith, linkedNature *repository.LinkedNature, attachmentId, includedById string) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	node, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		return s.repositories.AttachmentRepository.UnlinkWithXXIncludesAttachmentInTx(ctx, tx, tenant, linkedWith, linkedNature, attachmentId, includedById)
	})
	if err != nil {
		return nil, err
	}
	return node.(*dbtype.Node), err
}
func (s *attachmentService) GetAttachmentsForNode(ctx context.Context, linkedWith repository.LinkedWith, linkedNature *repository.LinkedNature, ids []string) (*entity.AttachmentEntities, error) {
	records, err := s.repositories.AttachmentRepository.GetAttachmentsForXX(ctx, common.GetTenantFromContext(ctx), linkedWith, linkedNature, ids)
	if err != nil {
		return nil, err
	}

	analysisDescribes := s.convertDbNodesToAttachments(records)

	return &analysisDescribes, nil
}

func (s *attachmentService) Create(ctx context.Context, newAnalysis *entity.AttachmentEntity, source, sourceOfTruth neo4jentity.DataSource) (*entity.AttachmentEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	interactionEventDbNode, err := session.ExecuteWrite(ctx, s.createAttachmentInDBTxWork(ctx, newAnalysis, source, sourceOfTruth))
	if err != nil {
		return nil, err
	}
	return s.MapDbNodeToAttachmentEntity(*interactionEventDbNode.(*dbtype.Node)), nil
}

func (s *attachmentService) createAttachmentInDBTxWork(ctx context.Context, newAttachment *entity.AttachmentEntity, source, sourceOfTruth neo4jentity.DataSource) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		analysisDbNode, err := s.repositories.AttachmentRepository.Create(ctx, tx, tenant, *newAttachment, source, sourceOfTruth)
		if err != nil {
			return nil, err
		}
		return analysisDbNode, nil
	}
}

func (s *attachmentService) GetAttachmentById(ctx context.Context, id string) (*entity.AttachmentEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (a:Attachment_%s {id:$id}) RETURN a`,
			common.GetTenantFromContext(ctx)),
			map[string]interface{}{
				"id": id,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.MapDbNodeToAttachmentEntity(queryResult.(dbtype.Node)), nil
}

func (s *attachmentService) MapDbNodeToAttachmentEntity(node dbtype.Node) *entity.AttachmentEntity {
	props := utils.GetPropsFromNode(node)
	createdAt := utils.GetTimePropOrEpochStart(props, "createdAt")
	analysisEntity := entity.AttachmentEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     &createdAt,
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		MimeType:      utils.GetStringPropOrEmpty(props, "mimeType"),
		Extension:     utils.GetStringPropOrEmpty(props, "extension"),
		Size:          utils.GetInt64PropOrZero(props, "size"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &analysisEntity
}

func (s *attachmentService) convertDbNodesToAttachments(records []*utils.DbNodeAndId) entity.AttachmentEntities {
	attachments := entity.AttachmentEntities{}
	for _, v := range records {
		attachment := s.MapDbNodeToAttachmentEntity(*v.Node)
		attachment.DataloaderKey = v.LinkedNodeId
		attachments = append(attachments, *attachment)

	}
	return attachments
}

func (s *attachmentService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}
