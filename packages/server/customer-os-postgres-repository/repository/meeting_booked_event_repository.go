package postgres_repository

import (
	"context"

	telemetry "github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type MeetingBookedEventRepository interface {
	Create(ctx context.Context, meetingBookedEvent *postgres_entity.MeetingBookedEvent) error
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
