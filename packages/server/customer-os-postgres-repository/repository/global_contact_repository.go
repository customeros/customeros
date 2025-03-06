package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type GlobalContactRepository interface {
	Create(ctx context.Context, contact *postgres_entity.GlobalContact) (*postgres_entity.GlobalContact, error)
	Update(ctx context.Context, contact *postgres_entity.GlobalContact) (*postgres_entity.GlobalContact, error)
	GetByLinkedInIdentifier(ctx context.Context, linkedInIdentifier string) (*postgres_entity.GlobalContact, error)
	GetByLinkedInIdentifierAndDomain(ctx context.Context, linkedInIdentifier string, primaryDomain string) (*postgres_entity.GlobalContact, error)
	GetByWorkEmail(ctx context.Context, workEmail string) (*postgres_entity.GlobalContact, error)
	GetByWorkEmailAndDomain(ctx context.Context, workEmail string, primaryDomain string) (*postgres_entity.GlobalContact, error)
	GetByPersonalEmail(ctx context.Context, personalEmail string) (*postgres_entity.GlobalContact, error)
	GetByPersonalEmailAndDomain(ctx context.Context, personalEmail string, primaryDomain string) (*postgres_entity.GlobalContact, error)
}

type globalContactRepository struct {
	db *gorm.DB
}

func NewGlobalContactRepository(gormDb *gorm.DB) GlobalContactRepository {
	return &globalContactRepository{db: gormDb}
}

func (r *globalContactRepository) addSortingClauses(query *gorm.DB) *gorm.DB {
	return query.
		Order("CASE WHEN job_ended_at IS NULL THEN 0 ELSE 1 END"). // Prioritize records with null job_ended_at
		Order("job_ended_at DESC").                                // Then most recent job_ended_at
		Order("job_started_at DESC")                               // Then most recent job_started_at
}

func (r *globalContactRepository) Create(ctx context.Context, contact *postgres_entity.GlobalContact) (*postgres_entity.GlobalContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalContactRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	result := r.db.WithContext(ctx).Create(contact)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}

	return contact, nil
}

func (r *globalContactRepository) Update(ctx context.Context, contact *postgres_entity.GlobalContact) (*postgres_entity.GlobalContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalContactRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	result := r.db.WithContext(ctx).Save(contact)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}

	return contact, nil
}

func (r *globalContactRepository) GetByLinkedInIdentifier(ctx context.Context, linkedInIdentifier string) (*postgres_entity.GlobalContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalContactRepository.GetByLinkedInIdentifier")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("linkedInIdentifier", linkedInIdentifier))

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("linkedin_identifier = ?", linkedInIdentifier),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetByLinkedInIdentifierAndDomain(ctx context.Context, linkedInIdentifier string, primaryDomain string) (*postgres_entity.GlobalContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalContactRepository.GetByLinkedInIdentifierAndDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(
		tracingLog.String("linkedInIdentifier", linkedInIdentifier),
		tracingLog.String("primaryDomain", primaryDomain),
	)

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("linkedin_identifier = ? AND primary_domain = ?", linkedInIdentifier, primaryDomain),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetByWorkEmail(ctx context.Context, workEmail string) (*postgres_entity.GlobalContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalContactRepository.GetByWorkEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("workEmail", workEmail))

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("work_email = ?", workEmail),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetByWorkEmailAndDomain(ctx context.Context, workEmail string, primaryDomain string) (*postgres_entity.GlobalContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalContactRepository.GetByWorkEmailAndDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(
		tracingLog.String("workEmail", workEmail),
		tracingLog.String("primaryDomain", primaryDomain),
	)

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("work_email = ? AND primary_domain = ?", workEmail, primaryDomain),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetByPersonalEmail(ctx context.Context, personalEmail string) (*postgres_entity.GlobalContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalContactRepository.GetByPersonalEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("personalEmail", personalEmail))

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("personal_email = ?", personalEmail),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetByPersonalEmailAndDomain(ctx context.Context, personalEmail string, primaryDomain string) (*postgres_entity.GlobalContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalContactRepository.GetByPersonalEmailAndDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(
		tracingLog.String("personalEmail", personalEmail),
		tracingLog.String("primaryDomain", primaryDomain),
	)

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("personal_email = ? AND primary_domain = ?", personalEmail, primaryDomain),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}

	return &contact, nil
}
