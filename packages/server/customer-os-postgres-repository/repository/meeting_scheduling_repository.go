package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type MeetingSchedulingRepository interface {
	GetById(ctx context.Context, tenant, id string) (*postgres_entity.MeetingScheduling, error)
	GetAll(ctx context.Context, tenant string) ([]*postgres_entity.MeetingScheduling, error)
	Save(ctx context.Context, meetingScheduling *postgres_entity.MeetingScheduling) (*postgres_entity.MeetingScheduling, error)
	Delete(ctx context.Context, tenant, id string) error
}

type meetingSchedulingRepository struct {
	db *gorm.DB
}

func NewMeetingSchedulingRepository(db *gorm.DB) MeetingSchedulingRepository {
	return &meetingSchedulingRepository{
		db: db,
	}
}

func (r *meetingSchedulingRepository) GetById(ctx context.Context, tenant, id string) (*postgres_entity.MeetingScheduling, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingSchedulingRepository.GetById")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)
	spans.LogKV("id", id)

	var meetingScheduling postgres_entity.MeetingScheduling
	err := r.db.Where("tenant = ? AND id = ?", tenant, id).First(&meetingScheduling).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return &meetingScheduling, nil
}

func (r *meetingSchedulingRepository) GetAll(ctx context.Context, tenant string) ([]*postgres_entity.MeetingScheduling, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingSchedulingRepository.GetAll")
	defer spans.Finish()

	var meetings []*postgres_entity.MeetingScheduling
	err := r.db.Where("tenant = ?", tenant).Order("created_at DESC").Find(&meetings).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(meetings))
	return meetings, nil
}

func (r *meetingSchedulingRepository) Save(ctx context.Context, meetingScheduling *postgres_entity.MeetingScheduling) (*postgres_entity.MeetingScheduling, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingSchedulingRepository.Save")
	defer spans.Finish()
	spans.LogObjectAsJson("meetingScheduling", meetingScheduling)

	if meetingScheduling.Tenant == "" {
		meetingScheduling.Tenant = common.GetTenantFromContext(ctx)
	}

	err := r.db.Save(meetingScheduling).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return meetingScheduling, nil
}

func (r *meetingSchedulingRepository) Delete(ctx context.Context, tenant, id string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingSchedulingRepository.Delete")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)
	spans.LogKV("id", id)

	err := r.db.Where("tenant = ? AND id = ?", tenant, id).Delete(&postgres_entity.MeetingScheduling{}).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}
