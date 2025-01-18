package command

import (
	"github.com/customeros/customeros/packages/server/events/event/common"
	"time"

	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization/model"
	"github.com/customeros/customeros/packages/server/events/eventstore"
)

type UpsertCustomFieldCommand struct {
	eventstore.BaseCommand
	Source          common.Source
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	CustomFieldData model.CustomField
}

func NewUpsertCustomFieldCommand(organizationId, tenant, source, sourceOfTruth, appSource, userId string,
	createdAt, updatedAt *time.Time, customField model.CustomField) *UpsertCustomFieldCommand {
	return &UpsertCustomFieldCommand{
		BaseCommand: eventstore.NewBaseCommand(organizationId, tenant, userId),
		Source: common.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		CustomFieldData: customField,
	}
}
