package attachment

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type attachmentService struct {
	neo4j *neo4j_repository.Repositories
}

func NewAttachmentService(neo4j *neo4j_repository.Repositories) interfaces.AttachmentService {
	return &attachmentService{
		neo4j: neo4j,
	}
}

func (s *attachmentService) GetById(ctx context.Context, id string) (*neo4j_entity.AttachmentEntity, error) {
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

	node, err := s.neo4j.AttachmentReadRepository.GetById(ctx, tenant, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToAttachmentEntity(node), nil
}

func (s *attachmentService) GetFor(c context.Context, entityType model.EntityType, relation *model.EntityRelation, ids []string) (*neo4j_entity.AttachmentEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "AttachmentService.GetFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	records, err := s.neo4j.AttachmentReadRepository.GetFor(ctx, common.GetTenantFromContext(ctx), entityType, relation, ids)
	if err != nil {
		return nil, err
	}

	attachments := neo4j_entity.AttachmentEntities{}
	for _, v := range records {
		attachment := mapper.MapDbNodeToAttachmentEntity(v.Node)
		attachment.DataloaderKey = v.LinkedNodeId
		attachments = append(attachments, *attachment)

	}
	return &attachments, nil
}

func (s *attachmentService) Create(c context.Context, record *neo4j_entity.AttachmentEntity) (*neo4j_entity.AttachmentEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "AttachmentService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *s.neo4j.Neo4jDriver)
	defer session.Close(ctx)

	interactionEventDbNode, err := session.ExecuteWrite(ctx, s.createAttachmentInDBTxWork(ctx, record))
	if err != nil {
		return nil, err
	}

	return mapper.MapDbNodeToAttachmentEntity(interactionEventDbNode.(*neo4j.Node)), nil
}

func (s *attachmentService) createAttachmentInDBTxWork(c context.Context, newAttachment *neo4j_entity.AttachmentEntity) func(tx neo4j.ManagedTransaction) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "AttachmentService.createAttachmentInDBTxWork")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		dbNode, err := s.neo4j.AttachmentWriteRepository.Create(ctx, tx, tenant, newAttachment.Id, newAttachment.CdnUrl, newAttachment.BasePath, newAttachment.FileName, newAttachment.MimeType, newAttachment.Size, newAttachment.CreatedAt, newAttachment.Source, newAttachment.SourceOfTruth, newAttachment.AppSource)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		return dbNode, nil
	}
}
