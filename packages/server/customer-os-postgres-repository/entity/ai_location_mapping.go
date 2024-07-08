package entity

import "time"

type AiLocationMapping struct {
	ID            uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	Input         string    `gorm:"column:input;type:text;NOT NULL;" json:"input"`
	ResponseJson  string    `gorm:"column:response_json;type:text;NOT NULL;" json:"responseJson"`
	AiPromptLogId string    `gorm:"column:ai_prompt_log_id;type:uuid;" json:"aiPromptLogId"`
}

func (AiLocationMapping) TableName() string {
	return "ai_location_mapping"
}
