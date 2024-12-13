package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowActionRegistryRepository interface {
	Initialize(ctx context.Context) error
	Create(ctx context.Context, flowAction entity.FlowActionRegistry) (*entity.FlowActionRegistry, error)
	Find(ctx context.Context, flowAction entity.FlowActionRegistry) (*entity.FlowActionRegistry, error)
	FindAll(ctx context.Context) (*[]entity.FlowActionRegistry, error)
}

type flowActionRegistryRepository struct {
	gormDb *gorm.DB
}

func NewFlowActionRegistryRepository(gormDb *gorm.DB) FlowActionRegistryRepository {
	return &flowActionRegistryRepository{gormDb: gormDb}
}

func (r *flowActionRegistryRepository) FindAll(ctx context.Context) (*[]entity.FlowActionRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var actions []entity.FlowActionRegistry
	err := r.gormDb.WithContext(ctx).
		Where("enabled = true").
		Order("action DESC").
		Find(&actions).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &actions, nil
}

func (r *flowActionRegistryRepository) Find(ctx context.Context, flowAction entity.FlowActionRegistry) (*entity.FlowActionRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var action entity.FlowActionRegistry
	query := r.gormDb.WithContext(ctx).Where("enabled = ?", true)

	// Add additional filters based on non-zero fields in flowAction
	if flowAction != (entity.FlowActionRegistry{}) {
		query = query.Where(&flowAction)
	}

	err := query.First(&action).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &action, nil
}

func (r *flowActionRegistryRepository) Create(ctx context.Context, flowAction entity.FlowActionRegistry) (*entity.FlowActionRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRegistryRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.WithContext(ctx).Create(&flowAction).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &flowAction, nil
}

func (r *flowActionRegistryRepository) Initialize(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRegistryRepository.Initialize")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredActions := []entity.FlowActionRegistry{
		{
			Action:       enum.ActionContactCreate.String(),
			FriendlyName: "Create a Contact",
			Description:  "Creates a new Contact",
			Enabled:      true,
		},
		{
			Action:       enum.ActionEmailSendNew.String(),
			FriendlyName: "Send an email",
			Description:  "Begins a new email thread",
			Enabled:      true,
		},
		{
			Action:       enum.ActionEmailSendReply.String(),
			FriendlyName: "Reply to an email thread",
			Description:  "Sends a reply to an existing email thread",
			Enabled:      true,
		},
		{
			Action:       enum.ActionLinkedinConnect.String(),
			FriendlyName: "Send Linkedin Connection Request",
			Description:  "Sends a connection request on LinkedIn",
			Enabled:      true,
		},
		{
			Action:       enum.ActionLinkedinMessage.String(),
			FriendlyName: "Send Linkedin Message",
			Description:  "Sends a direct message to a LinkedIn connection",
			Enabled:      true,
		},
		{
			Action:       enum.ActionOrganizationCreate.String(),
			FriendlyName: "Create an Organization",
			Description:  "Creates a new Organization",
			Enabled:      true,
		},
		{
			Action:       enum.ActionTimelineEventCreate.String(),
			FriendlyName: "Add Event to Timeline",
			Description:  "Adds a new event to the Organization Timeline",
			Enabled:      true,
		},
		// ... add more here
	}

	for _, action := range requiredActions {
		existingAction, err := r.Find(ctx, action)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if existingAction.Action != "" {
			continue
		}

		_, createErr := r.Create(ctx, action)
		if err != nil {
			tracing.TraceErr(span, createErr)
			return createErr
		}
	}

	return nil
}
