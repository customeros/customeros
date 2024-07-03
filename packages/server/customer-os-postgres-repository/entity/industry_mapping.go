package entity

import "time"

type IndustryMapping struct {
	ID             uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	InputIndustry  string    `gorm:"column:input_industry;type:varchar(255);NOT NULL" json:"inputIndustry"`
	OutputIndustry string    `gorm:"column:output_industry;type:varchar(255);NOT NULL" json:"outputIndustry"`
}

func (IndustryMapping) TableName() string {
	return "industry_mapping"
}
