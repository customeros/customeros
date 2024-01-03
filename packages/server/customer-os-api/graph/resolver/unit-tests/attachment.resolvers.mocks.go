package unit_tests

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockedAttachmentService struct {
	mock.Mock
}

func (s *MockedAttachmentService) GetAttachmentById(ctx context.Context, id string) (*entity.AttachmentEntity, error) {
	return nil, nil
}

func (s *MockedAttachmentService) Create(ctx context.Context, newAnalysis *entity.AttachmentEntity, source, sourceOfTruth neo4jentity.DataSource) (*entity.AttachmentEntity, error) {
	timeToHardcode := time.Date(2023, 9, 28, 12, 0, 0, 0, time.UTC)
	timePointer := &timeToHardcode
	return &entity.AttachmentEntity{
		MimeType:  "text/plain",
		Name:      "readme.txt",
		Size:      123,
		CreatedAt: timePointer,
	}, nil
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
