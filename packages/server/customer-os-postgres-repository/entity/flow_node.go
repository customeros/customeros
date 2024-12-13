package entity

import (
	"gorm.io/datatypes"
	"time"
)

type FlowNode struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	FlowID    string    `gorm:"column:flow_id;type:varchar(50);not null;index" json:"flowId"`
	Type      string    `gorm:"column:type;type:varchar(50);not null" json:"type"`
	Event     *string   `gorm:"column:event;type:varchar(255)" json:"event"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// used to store wait time config, email subject & template, etc
	EventData *datatypes.JSON `gorm:"column:data;type:jsonb" json:"data"`
	UIData    *string         `gorm:"column:ui_data;type:varchar(255)" json:"uiData"`
}

func (FlowNode) TableName() string {
	return "flow_nodes"
}

func (FlowNode) UniqueIndex() [][]string {
	return [][]string{
		{"flow_id", "id"},
	}
}
