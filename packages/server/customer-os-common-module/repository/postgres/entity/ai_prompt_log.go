package entity

import "time"

type AiPromptLog struct {
	ID                      string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt               time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	AppSource               string    `gorm:"column:app_source;type:varchar(50);NOT NULL;" json:"appSource" binding:"required"`
	Provider                string    `gorm:"column:provider;type:varchar(50);NOT NULL;" json:"provider" binding:"required"`
	Model                   string    `gorm:"column:model;type:varchar(100);NOT NULL;" json:"model" binding:"required"`
	PromptType              string    `gorm:"column:prompt_type;type:varchar(255);NOT NULL;" json:"promptType" binding:"required"`
	PromptTemplate          *string   `gorm:"column:prompt_template;type:text;" json:"promptTemplate" binding:"required"`
	ParentId                string    `gorm:"column:parent_id;type:text;NOT NULL;" json:"parentId" binding:"required"`
	Prompt                  string    `gorm:"column:prompt;type:text;NOT NULL;" json:"prompt" binding:"required"`
	RawResponse             string    `gorm:"column:raw_response;type:text;NOT NULL;" json:"rawResponse" binding:"required"`
	PostProcessError        bool      `gorm:"column:post_process_error;" json:"postProcessError" binding:"required"`
	PostProcessErrorMessage *string   `gorm:"column:post_process_error_message;type:text;" json:"postProcessErrorMessage" binding:"required"`
}

func (AiPromptLog) TableName() string {
	return "ai_prompt_log"
}
