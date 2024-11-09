package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
)

type AttachmentService interface {
	GetById(ctx context.Context, id string) (*neo4jentity.AttachmentEntity, error)
	GetFor(ctx context.Context, entityType model.EntityType, relation *model.EntityRelation, ids []string) (*neo4jentity.AttachmentEntities, error)

	Create(ctx context.Context, record *neo4jentity.AttachmentEntity) (*neo4jentity.AttachmentEntity, error)
}

type attachmentService struct {
	services *Services
}

func NewAttachmentService(services *Services) AttachmentService {
	return &attachmentService{
		services: services,
	}
}

func (s *attachmentService) GetById(ctx context.Context, id string) (*neo4jentity.AttachmentEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	node, err := s.services.Neo4jRepositories.AttachmentReadRepository.GetById(ctx, tenant, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToAttachmentEntity(node), nil
}

func (s *attachmentService) GetFor(c context.Context, entityType model.EntityType, relation *model.EntityRelation, ids []string) (*neo4jentity.AttachmentEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "AttachmentService.GetFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	records, err := s.services.Neo4jRepositories.AttachmentReadRepository.GetFor(ctx, common.GetTenantFromContext(ctx), entityType, relation, ids)
	if err != nil {
		return nil, err
	}

	attachments := neo4jentity.AttachmentEntities{}
	for _, v := range records {
		attachment := neo4jmapper.MapDbNodeToAttachmentEntity(v.Node)
		attachment.DataloaderKey = v.LinkedNodeId
		attachments = append(attachments, *attachment)

	}
	return &attachments, nil
}

func (s *attachmentService) Create(c context.Context, record *neo4jentity.AttachmentEntity) (*neo4jentity.AttachmentEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "AttachmentService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	interactionEventDbNode, err := session.ExecuteWrite(ctx, s.createAttachmentInDBTxWork(ctx, record))
	if err != nil {
		return nil, err
	}

	return neo4jmapper.MapDbNodeToAttachmentEntity(interactionEventDbNode.(*neo4j.Node)), nil
}

func (s *attachmentService) createAttachmentInDBTxWork(c context.Context, newAttachment *neo4jentity.AttachmentEntity) func(tx neo4j.ManagedTransaction) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "AttachmentService.createAttachmentInDBTxWork")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		dbNode, err := s.services.Neo4jRepositories.AttachmentWriteRepository.Create(ctx, tx, tenant, newAttachment.Id, newAttachment.CdnUrl, newAttachment.BasePath, newAttachment.FileName, newAttachment.MimeType, newAttachment.Size, newAttachment.CreatedAt, newAttachment.Source, newAttachment.SourceOfTruth, newAttachment.AppSource)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		return dbNode, nil
	}
}
