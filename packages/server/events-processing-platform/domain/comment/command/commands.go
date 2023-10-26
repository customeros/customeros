package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/model"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertCommentCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      model.CommentDataFields
	Source          commonmodel.Source
	ExternalSystem  commonmodel.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertCommentCommand(commentId, tenant, userId string, source commonmodel.Source, externalSystem commonmodel.ExternalSystem, dataFields model.CommentDataFields, createdAt, updatedAt *time.Time) *UpsertCommentCommand {
	return &UpsertCommentCommand{
		BaseCommand:    eventstore.NewBaseCommand(commentId, tenant, userId),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
