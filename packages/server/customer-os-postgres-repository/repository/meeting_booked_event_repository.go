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
