package entity

import "time"

type StatsApiCalls struct {
	ID        uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL;index:idx_stats_api_calls_unique,unique" json:"tenant"`
	Api       string    `gorm:"column:api;type:varchar(255);NOT NULL;index:idx_stats_api_calls_unique,unique" json:"api"`
	Day       time.Time `gorm:"column:day;type:date;NOT NULL;index:idx_stats_api_calls_unique,unique" json:"day"`
	Calls     uint64    `gorm:"column:calls;type:bigint;NOT NULL" json:"calls"`
}

func (StatsApiCalls) TableName() string {
	return "stats_api_calls"
}
