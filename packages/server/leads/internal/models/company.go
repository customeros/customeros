package models

import (
	"time"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
)

type Company struct {
	ID                    string                    `gorm:"column:id;primaryKey;type:varchar(255)"`
	CompanyName           string                    `gorm:"column:company_name;type:varchar(255);index;not null"`
	PrimaryDomain         string                    `gorm:"column:primary_domain;type:varchar(255);uniqueIndex"`
	ICPFit                bool                      `gorm:"column:icp_fit;type:boolean;default:false"`
	FirstEngagementSource enum.AdPlatform           `gorm:"column:first_engagement_source;type:varchar(50)"`
	LastEngagementSource  enum.AdPlatform           `gorm:"column:last_engagement_source;type:varchar(50)"`
	TopEngagementChannel  enum.Channel              `gorm:"column:top_engagement_channel;type:varchar(50)"`
	IdentifiedContacts    int                       `gorm:"column:identified_contacts;type:varchar(50)"`
	AnonymousVisitors     int                       `gorm:"column:anonymous_visitors:varchar(50)"`
	Stage                 enum.CustomerJourneyStage `gorm:"column:stage;type:varchar(50);index;default:'target'"`
	StageEngagementCount  int                       `gorm:"column:stage_engagement_count;type:integer;default:0"`
	FirstEngagementAt     *time.Time                `gorm:"column:first_engagement_at;type:timestamptz"`
	LastEngagementAt      *time.Time                `gorm:"column:last_engagement_at;type:timestamptz"`

	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Company) TableName() string {
	return "companies"
}
