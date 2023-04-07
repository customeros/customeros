package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateEmailCommand struct {
	eventstore.BaseCommand
	Tenant    string
	Email     string
	Source    models.Source
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type UpsertEmailCommand struct {
	eventstore.BaseCommand
	Tenant    string
	RawEmail  string
	Source    models.Source
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func NewCreateEmailCommand(aggregateID, tenant, rawEmail, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *CreateEmailCommand {
	return &CreateEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
		Email:       rawEmail,
		Source: models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewUpsertEmailCommand(aggregateID, tenant, rawEmail, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *UpsertEmailCommand {
	return &UpsertEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
		RawEmail:    rawEmail,
		Source: models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
