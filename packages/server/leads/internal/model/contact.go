package model

import (
	"time"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
)

type Contact struct {
	ID            string `gorm:"column:id;primaryKey;type:varchar(255)"`
	CompanyID     string `gorm:"column:company_id;type:varchar(255);index;not null"`
	CompanyName   string `gorm:"column:company_name;type:varchar(255);index;not null"`
	PrimaryDomain string `gorm:"column:primary_domain;type:varchar(255);index"`
	FirstName     string `gorm:"column:first_name;type:varchar(100)"`
	LastName      string `gorm:"column:last_name;type:varchar(100)"`
	Email         string `gorm:"column:email;type:varchar(255);uniqueIndex"`
	LinkedIn      string `gorm:"column:linkedin;type:varchar(255)"`
	JobTitle      string `gorm:"column:job_title;type:varchar(255)"`
	Department    string `gorm:"column:department;type:varchar(100)"`
	Seniority     string `gorm:"column:seniority;type:varchar(50)"`
	City          string `gorm:"column:city;type:varchar(255)"`
	Region        string `gorm:"column:region;type:varchar(255)"`
	Country       string `gorm:"column:country;type:varchar(255)"`

	// Engagement metrics
	Stage                enum.CustomerJourneyStage `gorm:"column:stage;type:varchar(50);index;default:'target'"`
	StageEngagementCount int                       `gorm:"column:stage_engagement_count;type:integer;default:0"`
	EngagementScore      float64                   `gorm:"column:engagement_score;type:float;default:0"`
	TotalSessionCount    int                       `gorm:"column:total_session_count;type:integer;default:0"`

	// Marketing action data
	BestOutreachChannel     enum.Channel        `gorm:"column:best_outreach_channel;type:varchar(50)"`
	BestRetargetingPlatform enum.SocialPlatform `gorm:"column:best_retargeting_platform;type:varchar(50)"`

	// Content and interest data
	TopInterestArea       string `gorm:"column:top_interest_area;type:varchar(100)"`
	SecondaryInterestArea string `gorm:"column:secondary_interest_area;type:varchar(100)"`
	PreferredContentType  string `gorm:"column:preferred_content_type;type:varchar(50)"`

	// Communication preferences
	PreferredLanguage string          `gorm:"column:preferred_language;type:varchar(10)"`
	PreferredDevice   enum.DeviceType `gorm:"column:preferred_device;type:varchar(20)"`

	// Timing fields
	FirstEngagementAt *time.Time `gorm:"column:first_engagement_at;type:timestamptz"`
	LastEngagementAt  *time.Time `gorm:"column:last_engagement_at;type:timestamptz"`

	// System fields
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for the Contact model
func (Contact) TableName() string {
	return "contacts"
}
