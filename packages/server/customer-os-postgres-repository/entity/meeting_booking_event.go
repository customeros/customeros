package postgres_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type MeetingBookingEvent struct {
	ID        string    `gorm:"column:id;type:varchar(21);primaryKey" json:"id"`
	Tenant    string    `gorm:"column:tenant;size:255;not null;primaryKey" json:"tenant"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`

	Title               string   `gorm:"column:title;size:255;not null" json:"title"`
	DurationMins        int64    `gorm:"column:duration_mins;not null;default:30" json:"durationMins"`
	Description         string   `gorm:"column:description;type:text;not null;default:''" json:"description"`
	AllowedParticipants []string `gorm:"column:allowed_participants;type:text[];default:'{}'" json:"allowedParticipants"`

	BookingFormName  string `gorm:"column:booking_form_name;size:255;not null;default:''" json:"bookingFormName"`
	BookingFormEmail string `gorm:"column:booking_form_email;size:255;not null;default:''" json:"bookingFormEmail"`
	BookingFormPhone string `gorm:"column:booking_form_phone;size:255;not null;default:''" json:"bookingFormPhone"`

	BookOptionEnabled                   bool  `gorm:"column:book_option_enabled;not null;default:false" json:"bookOptionEnabled"`
	BookOptionBufferBetweenMeetingsMins int64 `gorm:"column:book_option_buffer_between_meetings_mins;not null;default:0" json:"bookOptionBufferBetweenMeetingsMins"`
	BookOptionDaysInAdvance             int64 `gorm:"column:book_option_days_in_advance;not null;default:0" json:"bookOptionDaysInAdvance"`
	BookOptionMinNoticeMins             int64 `gorm:"column:book_option_min_notice_mins;not null;default:0" json:"bookOptionMinNoticeMins"`

	EmailNotificationEnabled        bool                                `gorm:"column:email_notification_enabled;not null;default:false" json:"emailNotificationEnabled"`
	Location                        string                              `gorm:"column:location;size:255;not null;default:''" json:"location"`
	AssignmentMethod                enum.MeetingBookingAssignmentMethod `gorm:"column:assignment_method;size:255;not null;default:''" json:"assignmentMethod"`
	ShowLogo                        bool                                `gorm:"column:show_logo;not null;default:false" json:"showLogo"`
	BookingConfirmationRedirectLink string                              `gorm:"column:booking_confirmation_redirect_link;size:255;not null;default:''" json:"bookingConfirmationRedirectLink"`
}

func (MeetingBookingEvent) TableName() string {
	return "meeting_booking_events"
}

// BeforeCreate hook to ensure ID has the correct prefix and validate
func (u *MeetingBookingEvent) BeforeCreate(tx *gorm.DB) error {
	u.ID = utils.GenerateNanoIdWithPrefix("mbe", 16)
	return nil
}

// BeforeUpdate hook to update the UpdatedAt timestamp and validate
func (u *MeetingBookingEvent) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
