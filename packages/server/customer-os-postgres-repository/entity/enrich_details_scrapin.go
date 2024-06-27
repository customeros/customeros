package entity

import "time"

type ScrapInFlow string

const (
	ScrapInFlowPersonSearch  ScrapInFlow = "PERSON_SEARCH"
	ScrapInFlowPersonProfile ScrapInFlow = "PERSON_PROFILE"
)

type EnrichDetailsScrapIn struct {
	ID            uint64      `gorm:"primary_key;autoIncrement:true" json:"id"`
	Flow          ScrapInFlow `gorm:"column:flow;type:varchar(255);NOT NULL" json:"flow"`
	Param1        string      `gorm:"column:param1;type:varchar(1000);DEFAULT:'';NOT NULL;index" json:"param1"`
	Param2        string      `gorm:"column:param2;type:varchar(1000);" json:"param2"`
	Param3        string      `gorm:"column:param3;type:varchar(1000);" json:"param3"`
	Param4        string      `gorm:"column:param4;type:varchar(1000);" json:"param4"`
	AllParamsJson string      `gorm:"column:all_params_json;type:text;DEFAULT:'';NOT NULL" json:"allParams"`
	Data          string      `gorm:"column:data;type:text;DEFAULT:'';NOT NULL" json:"data"`
	CreatedAt     time.Time   `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time   `gorm:"column:updated_at;type:timestamp;;DEFAULT:current_timestamp" json:"updatedAt"`
}

func (EnrichDetailsScrapIn) TableName() string {
	return "enrich_details_scrapin"
}
