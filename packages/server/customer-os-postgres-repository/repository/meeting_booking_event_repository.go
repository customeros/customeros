package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type MeetingBookingEventRepository interface {
	GetById(ctx context.Context, tenant, id string) (*postgresEntity.MeetingBookingEvent, error)
	GetAll(ctx context.Context, tenant string) ([]*postgresEntity.MeetingBookingEvent, error)
	Save(ctx context.Context, meetingBookingEvent *postgresEntity.MeetingBookingEvent) (*postgresEntity.MeetingBookingEvent, error)
	Delete(ctx context.Context, tenant, id string) error
}

type meetingBookingEventRepository struct {
	db *gorm.DB
}

func NewMeetingBookingEventRepository(db *gorm.DB) MeetingBookingEventRepository {
	return &meetingBookingEventRepository{
		db: db,
	}
}

func (r *meetingBookingEventRepository) GetById(ctx context.Context, tenant, id string) (*postgresEntity.MeetingBookingEvent, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookingEventRepository.GetById")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)
	spans.LogKV("id", id)

	var meetingBookingEvent postgresEntity.MeetingBookingEvent
	err := r.db.Where("tenant = ? AND id = ?", tenant, id).First(&meetingBookingEvent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return &meetingBookingEvent, nil
}

func (r *meetingBookingEventRepository) GetAll(ctx context.Context, tenant string) ([]*postgresEntity.MeetingBookingEvent, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookingEventRepository.GetAll")
	defer spans.Finish()

	var meetings []*postgresEntity.MeetingBookingEvent
	err := r.db.Where("tenant = ?", tenant).Order("created_at DESC").Find(&meetings).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(meetings))
	return meetings, nil
}

func (r *meetingBookingEventRepository) Save(ctx context.Context, meetingBookingEvent *postgresEntity.MeetingBookingEvent) (*postgresEntity.MeetingBookingEvent, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookingEventRepository.Save")
	defer spans.Finish()
	spans.LogObjectAsJson("meetingBookingEvent", meetingBookingEvent)

	if meetingBookingEvent.Tenant == "" {
		meetingBookingEvent.Tenant = common.GetTenantFromContext(ctx)
	}

	err := r.db.Save(meetingBookingEvent).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return meetingBookingEvent, nil
}

func (r *meetingBookingEventRepository) Delete(ctx context.Context, tenant, id string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookingEventRepository.Delete")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)
	spans.LogKV("id", id)

	err := r.db.Where("tenant = ? AND id = ?", tenant, id).Delete(&postgresEntity.MeetingBookingEvent{}).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}
