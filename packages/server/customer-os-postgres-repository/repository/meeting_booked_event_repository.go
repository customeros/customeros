package postgres_repository

import (
	"context"
	"errors"

	telemetry "github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type MeetingBookedEventRepository interface {
	Create(ctx context.Context, meetingBookedEvent *postgres_entity.MeetingBookedEvent) error
	GetActiveFirstUpcomingByMeetingBookingEventIDAndClientEmail(ctx context.Context, tenant, meetingBookingEventID, clientEmail string) (*postgres_entity.MeetingBookedEvent, error)
	Cancel(ctx context.Context, id string) error
	DeleteByConditions(ctx context.Context, conditions *postgres_entity.MeetingBookedEvent) error
}

type meetingBookedEventRepository struct {
	gormDb *gorm.DB
}

func NewMeetingBookedEventRepository(gormDb *gorm.DB) MeetingBookedEventRepository {
	return &meetingBookedEventRepository{gormDb: gormDb}
}

func (r *meetingBookedEventRepository) Create(ctx context.Context, meetingBookedEvent *postgres_entity.MeetingBookedEvent) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookedEventRepository.Create")
	defer spans.Finish()

	return r.gormDb.Create(meetingBookedEvent).Error
}

func (r *meetingBookedEventRepository) GetActiveFirstUpcomingByMeetingBookingEventIDAndClientEmail(ctx context.Context, tenant, meetingBookingEventID, clientEmail string) (*postgres_entity.MeetingBookedEvent, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookedEventRepository.GetActiveFirstUpcomingByMeetingBookingEventIDAndClientEmail")
	defer spans.Finish()
	spans.LogKV("meetingBookingEventID", meetingBookingEventID)
	spans.LogKV("clientEmail", clientEmail)

	var meetingBookedEvent postgres_entity.MeetingBookedEvent
	err := r.gormDb.Where("tenant = ? AND meeting_booking_event_id = ? AND client_email = ? AND canceled = false AND start_time > ?", tenant, meetingBookingEventID, clientEmail, utils.Now()).
		Order("start_time ASC").
		First(&meetingBookedEvent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	return &meetingBookedEvent, nil
}

func (r *meetingBookedEventRepository) Cancel(ctx context.Context, id string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookedEventRepository.Cancel")
	defer spans.Finish()
	spans.TagEntity(id)

	return r.gormDb.Model(&postgres_entity.MeetingBookedEvent{}).Where("id = ?", id).Update("canceled", true).Error
}

func (r *meetingBookedEventRepository) DeleteByConditions(ctx context.Context, conditions *postgres_entity.MeetingBookedEvent) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "MeetingBookedEventRepository.DeleteByConditions")
	defer spans.Finish()
	spans.LogObjectAsJson("conditions", conditions)

	query := r.gormDb.Model(&postgres_entity.MeetingBookedEvent{})

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

	return query.Delete(&postgres_entity.MeetingBookedEvent{}).Error
}
