package resolver

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"

	//"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	//"github.com/mrdulin/gqlgen-cnode/graph/model"
	"github.com/stretchr/testify/mock"
)

type MockedAttachmentService struct {
	mock.Mock
}

//func (s *MockedAttachmentService) AttachmentCreate(input model.AttachmentInput) *model.Attachment {
//	args := s.Called(input)
//	return args.Get(0).(*model.Attachment)
//}
//
//func (s *MockedAttachmentService) Attachment(id string) *model.Attachment {
//	args := s.Called(id)
//	return args.Get(0).(*model.Attachment)
//}

func (s *MockedAttachmentService) GetAttachmentById(ctx context.Context, id string) (*entity.AttachmentEntity, error) {
	return nil, nil
}

func (s *MockedAttachmentService) Create(ctx context.Context, newAnalysis *entity.AttachmentEntity, source, sourceOfTruth entity.DataSource) (*entity.AttachmentEntity, error) {
	return nil, nil
}
func (s *MockedAttachmentService) GetAttachmentsForNode(ctx context.Context, linkedWith repository.LinkedWith, linkedNature *repository.LinkedNature, ids []string) (*entity.AttachmentEntities, error) {
	return nil, nil
}

func (s *MockedAttachmentService) LinkNodeWithAttachment(ctx context.Context, linkedWith repository.LinkedWith, linkedNature *repository.LinkedNature, attachmentId, includedById string) (*dbtype.Node, error) {
	return nil, nil
}
func (s *MockedAttachmentService) UnlinkNodeWithAttachment(ctx context.Context, linkedWith repository.LinkedWith, linkedNature *repository.LinkedNature, attachmentId, includedById string) (*dbtype.Node, error) {
	return nil, nil
}

func (s *MockedAttachmentService) MapDbNodeToAttachmentEntity(node dbtype.Node) *entity.AttachmentEntity {
	return nil
}
