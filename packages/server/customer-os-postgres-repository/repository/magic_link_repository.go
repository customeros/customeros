package postgres_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type magicLinkRepository struct {
	gormDb *gorm.DB
}

type MagicLinkRepository interface {
	GetByEmail(ctx context.Context, email string) (*postgres_entity.MagicLink, error)
	GetByCode(ctx context.Context, code string) (*postgres_entity.MagicLink, error)
	Create(ctx context.Context, magicLink *postgres_entity.MagicLink) error
	Delete(ctx context.Context, id string) error
}

func NewMagicLinkRepository(db *gorm.DB) MagicLinkRepository {
	return &magicLinkRepository{gormDb: db}
}

func (r *magicLinkRepository) GetByEmail(ctx context.Context, email string) (*postgres_entity.MagicLink, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MagicLinkRepository.GetByEmail")
	defer spans.Finish()
	spans.LogKV("email", email)

	var result *postgres_entity.MagicLink
	err := r.gormDb.
		Where("email = ?", email).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)

	return result, nil
}

func (r *magicLinkRepository) GetByCode(ctx context.Context, code string) (*postgres_entity.MagicLink, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MagicLinkRepository.GetByCode")
	defer spans.Finish()
	spans.LogKV("code", code)

	var result *postgres_entity.MagicLink
	err := r.gormDb.
		Where("code = ?", code).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)

	return result, nil
}

func (r *magicLinkRepository) Create(ctx context.Context, input *postgres_entity.MagicLink) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MagicLinkRepository.Create")
	defer spans.Finish()
	spans.LogObjectAsJson("magicLink", input)

	// Check if the mailbox already exists
	var magicLink postgres_entity.MagicLink
	err := r.gormDb.
		Where("code = ?", input.Code).
		First(&magicLink).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		spans.TraceError(err)
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If not found, create a new mailbox
		magicLink = postgres_entity.MagicLink{
			Email: input.Email,
			Code:  input.Code,
			Url:   input.Url,
		}

		err = r.gormDb.Create(&magicLink).Error
		if err != nil {
			spans.TraceError(err)
			return err
		}
	} else {
		e := errors.New("magic link already exists by code")
		spans.TraceError(e)
		return e
	}

	return nil
}

func (r *magicLinkRepository) Delete(ctx context.Context, id string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MagicLinkRepository.Delete")
	defer spans.Finish()
	spans.LogKV("id", id)

	// Check if the mailbox already exists
	var magicLink postgres_entity.MagicLink
	err := r.gormDb.
		Where("id = ?", id).
		First(&magicLink).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		spans.TraceError(err)
		return err
	}

	if magicLink.ID == "" {
		e := errors.New("magic link not found")
		spans.TraceError(e)
		return e
	}

	err = r.gormDb.Delete(magicLink).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
