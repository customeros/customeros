package repository

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type GlobalOrganizationRepository interface {
	GetById(ctx context.Context, id uint64) (*entity.GlobalOrganization, error)
	GetByPrimaryDomain(ctx context.Context, domain string) (*entity.GlobalOrganization, error)
	GetByPrimaryDomains(ctx context.Context, domains []string) ([]*entity.GlobalOrganization, error)
	Create(ctx context.Context, organization *entity.GlobalOrganization) (*entity.GlobalOrganization, error)
	Update(ctx context.Context, organization *entity.GlobalOrganization) (*entity.GlobalOrganization, error)
	Search(ctx context.Context, searchTerm string, limit int) ([]*entity.GlobalOrganization, error)
}

type globalOrganizationRepository struct {
	db *gorm.DB
}

func NewGlobalOrganizationRepository(gormDb *gorm.DB) GlobalOrganizationRepository {
	return &globalOrganizationRepository{db: gormDb}
}

func (r *globalOrganizationRepository) GetById(ctx context.Context, id uint64) (*entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Uint64("id", id))

	organization := &entity.GlobalOrganization{}
	result := r.db.WithContext(ctx).Where("id = ?", id).First(organization)
	if result.Error != nil {
		span.LogFields(tracingLog.Bool("found", false))
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Bool("found", true))
	return organization, nil
}

func (r *globalOrganizationRepository) GetByPrimaryDomains(ctx context.Context, domains []string) ([]*entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetByPrimaryDomains")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Object("domains", domains))

	organizations := make([]*entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).Where("primary_domain IN ?", domains).Find(&organizations)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Int("result.count", len(organizations)))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetByPrimaryDomain(ctx context.Context, domain string) (*entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetByPrimaryDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("domain", domain))

	organization := &entity.GlobalOrganization{}
	result := r.db.WithContext(ctx).Where("primary_domain = ?", domain).First(organization)
	if result.Error != nil {
		span.LogFields(tracingLog.Bool("found", false))
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Bool("found", true))
	return organization, nil
}

func (r *globalOrganizationRepository) Create(ctx context.Context, organization *entity.GlobalOrganization) (*entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "organization", organization)

	result := r.db.WithContext(ctx).Create(&organization)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	return organization, nil
}

func (r *globalOrganizationRepository) Update(ctx context.Context, organization *entity.GlobalOrganization) (*entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "organization", organization)

	organization.UpdatedAt = utils.Now()
	result := r.db.WithContext(ctx).Save(organization)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	return organization, nil
}

func (r *globalOrganizationRepository) Search(ctx context.Context, searchTerm string, limit int) ([]*entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.SearchOrganizations")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("searchTerm", searchTerm))

	organizations := make([]*entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Select("id, name, primary_domain, website, logo_url, icon_url, domains").
		Where("name ILIKE ? OR primary_domain ILIKE ? OR other_domains ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Int("found", len(organizations)))
	return organizations, nil
}
