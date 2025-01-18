package postgres_entity

import (
	"time"
)

type FlowNode struct {
	ID        string     `gorm:"primaryKey;type:varchar(50)" json:"id"`
	FlowID    string     `gorm:"column:flow_id;type:varchar(50);not null;index" json:"flowId"`
	Type      string     `gorm:"column:type;type:varchar(50);not null" json:"type"`
	AgentID   *string    `gorm:"column:agent_id;type:varchar(50);not null;index" json:"agentId"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// used to store wait time config, email subject & template, etc
	UIData *string `gorm:"column:ui_data;type:varchar(255)" json:"uiData"`
}

func (FlowNode) TableName() string {
	return "flow_nodes"
}

func (FlowNode) UniqueIndex() [][]string {
	return [][]string{
		{"flow_id", "id"},
	}
}
