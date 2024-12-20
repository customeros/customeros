package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type magicLinkRepository struct {
	gormDb *gorm.DB
}

type MagicLinkRepository interface {
	GetByEmail(ctx context.Context, email string) (*entity.MagicLink, error)
	GetByCode(ctx context.Context, code string) (*entity.MagicLink, error)
	Create(ctx context.Context, magicLink *entity.MagicLink) error
	Delete(ctx context.Context, id string) error
}

func NewMagicLinkRepository(db *gorm.DB) MagicLinkRepository {
	return &magicLinkRepository{gormDb: db}
}

func (r *magicLinkRepository) GetByEmail(ctx context.Context, email string) (*entity.MagicLink, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "MagicLinkRepository.GetByEmail")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("email", email))

	var result *entity.MagicLink
	err := r.gormDb.
		Where("email = ?", email).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			span.LogFields(tracingLog.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return result, nil
}

func (r *magicLinkRepository) GetByCode(ctx context.Context, code string) (*entity.MagicLink, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "MagicLinkRepository.GetByCode")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("code", code))

	var result *entity.MagicLink
	err := r.gormDb.
		Where("code = ?", code).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			span.LogFields(tracingLog.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return result, nil
}

func (r *magicLinkRepository) Create(ctx context.Context, input *entity.MagicLink) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "MagicLinkRepository.Create")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.Object("magicLink", input))

	// Check if the mailbox already exists
	var magicLink entity.MagicLink
	err := r.gormDb.
		Where("code = ?", input.Code).
		First(&magicLink).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If not found, create a new mailbox
		magicLink = entity.MagicLink{
			Email: input.Email,
			Code:  input.Code,
			Url:   input.Url,
		}

		err = r.gormDb.Create(&magicLink).Error
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		e := errors.New("magic link already exists by code")
		tracing.TraceErr(span, e)
		return e
	}

	return nil
}

func (r *magicLinkRepository) Delete(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "MagicLinkRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.Object("id", id))

	// Check if the mailbox already exists
	var magicLink entity.MagicLink
	err := r.gormDb.
		Where("id = ?", id).
		First(&magicLink).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return err
	}

	if magicLink.ID == "" {
		e := errors.New("magic link not found")
		tracing.TraceErr(span, e)
		return e
	}

	err = r.gormDb.Delete(magicLink).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
