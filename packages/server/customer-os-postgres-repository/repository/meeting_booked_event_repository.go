package postgres_repository

import (
	"context"
	"errors"

	telemetry "github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type MeetingBookedEventRepository interface {
	Create(ctx context.Context, meetingBookedEvent *postgresEntity.MeetingBookedEvent) error
	Cancel(ctx context.Context, id string) error
	DeleteByConditions(ctx context.Context, conditions *postgresEntity.MeetingBookedEvent) error
	GetByEventID(ctx context.Context, tenant, meetingBookingEventID, eventID string) (*postgresEntity.MeetingBookedEvent, error)
}

type meetingBookedEventRepository struct {
	gormDb *gorm.DB
}

func NewMeetingBookedEventRepository(gormDb *gorm.DB) MeetingBookedEventRepository {
	return &meetingBookedEventRepository{gormDb: gormDb}
}

func (r *meetingBookedEventRepository) Create(ctx context.Context, meetingBookedEvent *postgresEntity.MeetingBookedEvent) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookedEventRepository.Create")
	defer spans.Finish()

	return r.gormDb.Create(meetingBookedEvent).Error
}

func (r *meetingBookedEventRepository) Cancel(ctx context.Context, id string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookedEventRepository.Cancel")
	defer spans.Finish()
	spans.TagEntity(id)

	return r.gormDb.Model(&postgresEntity.MeetingBookedEvent{}).Where("id = ?", id).Update("canceled", true).Error
}

func (r *meetingBookedEventRepository) DeleteByConditions(ctx context.Context, conditions *postgresEntity.MeetingBookedEvent) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookedEventRepository.DeleteByConditions")
	defer spans.Finish()
	spans.LogObjectAsJson("conditions", conditions)

	query := r.gormDb.Model(&postgresEntity.MeetingBookedEvent{})

	// Build where conditions based on provided fields
	if conditions.Tenant != "" {
		query = query.Where("tenant = ?", conditions.Tenant)
	}
	if conditions.MeetingBookingEventID != "" {
		query = query.Where("meeting_booking_event_id = ?", conditions.MeetingBookingEventID)
	}
	if conditions.ClientEmail != "" {
		query = query.Where("client_email = ?", conditions.ClientEmail)
	}
	if !conditions.StartTime.IsZero() {
		query = query.Where("start_time = ?", conditions.StartTime)
	}

	return query.Delete(&postgresEntity.MeetingBookedEvent{}).Error
}

func (r *meetingBookedEventRepository) GetByEventID(ctx context.Context, tenant, meetingBookingEventID, eventID string) (*postgresEntity.MeetingBookedEvent, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookedEventRepository.GetByEventID")
	defer spans.Finish()
	spans.LogKV("meetingBookingEventID", meetingBookingEventID)
	spans.LogKV("eventID", eventID)

	var meetingBookedEvent postgresEntity.MeetingBookedEvent
	err := r.gormDb.Where("tenant = ? AND meeting_booking_event_id = ? AND nylas_event_id = ?", tenant, meetingBookingEventID, eventID).First(&meetingBookedEvent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return &meetingBookedEvent, nil
}
