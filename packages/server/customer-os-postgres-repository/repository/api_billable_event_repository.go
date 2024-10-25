package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ApiBillableEventRepository interface {
	RegisterEvent(ctx context.Context, tenant string, event entity.BillableEvent, subtype, externalID, referenceData string) (*entity.ApiBillableEvent, error)
}

type apiBillableEventRepository struct {
	gormDb *gorm.DB
}

func NewApiBillableEventRepository(db *gorm.DB) ApiBillableEventRepository {
	return &apiBillableEventRepository{gormDb: db}
}

// Register creates a new ApiBillableEvent and stores it in the database
func (r *apiBillableEventRepository) RegisterEvent(ctx context.Context, tenant string, event entity.BillableEvent, subtype, externalID, referenceData string) (*entity.ApiBillableEvent, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BillableEventRepository.RegisterEvent")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("event", event)
	span.LogKV("externalID", externalID)
	span.LogKV("referenceData", referenceData)
	span.LogKV("subtype", subtype)

	// Construct the ApiBillableEvent entity
	billableEvent := entity.ApiBillableEvent{
		Tenant:        tenant,
		Event:         event,
		CreatedAt:     utils.Now(),
		ExternalID:    externalID,
		ReferenceData: referenceData,
		Subtype:       subtype,
	}

	err := r.gormDb.Create(&billableEvent).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to store billable event")
	}

	return &billableEvent, nil
}
