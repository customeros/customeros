package models

import (
	"time"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
)

type SessionEngagement struct {
	// Primary identification
	ID        string `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	LeadID    string `gorm:"column:lead_id;type:varchar(255);index;not null"`
	SessionID string `gorm:"column:session_id;type:varchar(255);uniqueIndex"`

	// Session timing
	SessionStartedAt time.Time `gorm:"column:session_started_at;type:timestamptz;not null;index"`
	SessionEndedAt   time.Time `gorm:"column:session_ended_at;type:timestamptz"`
	SessionDuration  int       `gorm:"column:session_duration;type:integer"` // in seconds

	// Session metrics
	PageviewCount   int     `gorm:"column:pageview_count;type:integer;default:0"`
	EngagementDepth float64 `gorm:"column:engagement_depth;type:float;default:0"` // Composite score
	FormSubmissions int     `gorm:"column:form_submissions;type:integer;default:0"`
	GoalCompletions int     `gorm:"column:goal_completions;type:integer;default:0"`

	// Entry point attribution
	EntryPage        string              `gorm:"column:entry_page;type:varchar(255)"`
	ExitPage         string              `gorm:"column:exit_page;type:varchar(255)"`
	Channel          enum.Channel        `gorm:"column:channel;type:varchar(50)"`
	Source           enum.AdPlatform     `gorm:"column:source;type:varchar(50)"`
	SourcePlatform   enum.SocialPlatform `gorm:"column:source_platform;type:varchar(50)"`
	ViewedOnPlatform enum.SocialPlatform `gorm:"column:viewed_on_platform;type:varchar(50)"`
	GCLID            string              `gorm:"column:gclid;type:varchar(255)"`
	ReferrerDomain   string              `gorm:"column:referrer_domain;type:varchar(255)"`

	// Aggregated UTM for entry page
	UTMSource   string `gorm:"column:utm_source;type:varchar(100)"`
	UTMMedium   string `gorm:"column:utm_medium;type:varchar(100)"`
	UTMCampaign string `gorm:"column:utm_campaign;type:varchar(100)"`
	UTMContent  string `gorm:"column:utm_content;type:varchar(100)"`
	UTMTerm     string `gorm:"column:utm_term;type:varchar(100)"`

	// Journey context
	CustomerJourneyStage enum.CustomerJourneyStage `gorm:"column:journey_stage;type:varchar(50)"`
	IsSignificantSession bool                      `gorm:"column:is_significant_session;type:boolean;default:false"`

	// Context info
	DeviceType enum.DeviceType `gorm:"column:device_type;type:varchar(20)"`
	Language   string          `gorm:"column:language;type:varchar(10)"`

	// Page/content engagement summaries

	// System fields
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
}
