package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type BillableEventDetails struct {
	ExternalID    string `json:"externalId"`
	ReferenceData string `json:"referenceData"`
	Subtype       string `json:"subtype"`
	Source        string `json:"source"`
}

type ApiBillableEventRepository interface {
	RegisterEvent(ctx context.Context, tenant string, event postgres_entity.BillableEvent, details BillableEventDetails) (*postgres_entity.ApiBillableEvent, error)
}

type apiBillableEventRepository struct {
	gormDb *gorm.DB
}

func NewApiBillableEventRepository(db *gorm.DB) ApiBillableEventRepository {
	return &apiBillableEventRepository{gormDb: db}
}

// Register creates a new ApiBillableEvent and stores it in the database
func (r *apiBillableEventRepository) RegisterEvent(ctx context.Context, tenant string, event postgres_entity.BillableEvent, details BillableEventDetails) (*postgres_entity.ApiBillableEvent, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BillableEventRepository.RegisterEvent")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("event", event)
	span.LogKV("externalID", details.ExternalID)
	span.LogKV("referenceData", details.ReferenceData)
	span.LogKV("subtype", details.Subtype)
	span.LogKV("source", details.Source)

	// Construct the ApiBillableEvent postgres_entity
	billableEvent := postgres_entity.ApiBillableEvent{
		Tenant:        tenant,
		Event:         event,
		CreatedAt:     utils.Now(),
		ExternalID:    details.ExternalID,
		ReferenceData: details.ReferenceData,
		Subtype:       details.Subtype,
		Source:        details.Source,
	}

	err := r.gormDb.Create(&billableEvent).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to store billable event")
	}

	return &billableEvent, nil
}
