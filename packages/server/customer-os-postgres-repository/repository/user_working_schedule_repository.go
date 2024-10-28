package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type UserWorkingScheduleRepository interface {
	GetForUser(ctx context.Context, tenant, userId string) ([]*entity.UserWorkingSchedule, error)
	Store(ctx context.Context, tenant string, input *entity.UserWorkingSchedule) error
}

type userWorkingScheduleRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewUserWorkingScheduleRepository(gormDb *gorm.DB) UserWorkingScheduleRepository {
	return &userWorkingScheduleRepositoryImpl{gormDb: gormDb}
}

func (repo *userWorkingScheduleRepositoryImpl) GetForUser(ctx context.Context, tenant, userId string) ([]*entity.UserWorkingSchedule, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserWorkingScheduleRepository.GetForUser")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	span.LogFields(log.String("tenant", tenant), log.String("userId", userId))

	var e []*entity.UserWorkingSchedule
	err := repo.gormDb.Where("tenant = ? and user_id = ?", tenant, userId).Find(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return e, nil
}

func (repo *userWorkingScheduleRepositoryImpl) Store(ctx context.Context, tenant string, input *entity.UserWorkingSchedule) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserWorkingScheduleRepository.Store")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	span.LogFields(log.Object("input", input))

	input.Tenant = tenant
	input.CreatedAt = utils.Now()

	err := repo.gormDb.Save(&input).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
