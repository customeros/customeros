package postgres_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type MeetingBookedEvent struct {
	ID                    string    `gorm:"column:id;type:varchar(25);primaryKey" json:"id"`
	Tenant                string    `gorm:"column:tenant;size:255;not null;primaryKey" json:"tenant"`
	CreatedAt             time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt             time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	MeetingBookingEventID string    `gorm:"column:meeting_booking_event_id;type:varchar(21);not null" json:"meetingBookingEventId"`
	HostCalendarID        string    `gorm:"column:host_calendar_id;size:255" json:"hostCalendarId"`
	HostEmail             string    `gorm:"column:host_email;size:255;not null" json:"hostEmail"`
	HostName              string    `gorm:"column:host_name;type:text" json:"hostName"`
	ClientEmail           string    `gorm:"column:client_email;size:255;not null" json:"clientEmail"`
	ClientName            string    `gorm:"column:client_name;type:text;not null" json:"clientName"`
	ClientPhone           string    `gorm:"column:client_phone;type:text;not null" json:"clientPhone"`
	StartTime             time.Time `gorm:"column:start_time;type:timestamp;not null" json:"startTime"`
	EndTime               time.Time `gorm:"column:end_time;type:timestamp;not null" json:"endTime"`
	DurationMins          int64     `gorm:"column:duration_mins;not null" json:"durationMins"`
	NylasResponse         string    `gorm:"column:nylas_response;type:text;not null" json:"nylasResponse"`
	NylasEventID          string    `gorm:"column:nylas_event_id;type:varchar(255)" json:"nylasEventId"`
	Canceled              bool      `gorm:"column:canceled;type:boolean;not null;default:false" json:"canceled"`
}

func (MeetingBookedEvent) TableName() string {
	return "meeting_booked_events"
}

// BeforeCreate hook to ensure ID has the correct prefix and validate
func (u *MeetingBookedEvent) BeforeCreate(tx *gorm.DB) error {
	u.ID = utils.GenerateNanoIdWithPrefix("mev", 21)
	return nil
}

// BeforeUpdate hook to update the UpdatedAt timestamp and validate
func (u *MeetingBookedEvent) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
