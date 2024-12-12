package entity

import (
	"gorm.io/datatypes"
	"time"
)

type FlowNode struct {
	ID        string         `gorm:"primaryKey;type:varchar(50)" json:"id"`
	FlowID    string         `gorm:"column:flow_id;type:varchar(50);not null;index" json:"flowId"`
	Type      string         `gorm:"column:type;type:varchar(50);not null" json:"type"`
	Event     string         `gorm:"column:event;type:varchar(255);not null" json:"event"`
	PositionX float64        `gorm:"column:position_x" json:"positionX"`
	PositionY float64        `gorm:"column:position_y" json:"positionY"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Data      datatypes.JSON `gorm:"column:data;type:jsonb" json:"data"` // used to store wait time config etc
}

func (FlowNode) TableName() string {
	return "flow_nodes"
}

func (FlowNode) UniqueIndex() [][]string {
	return [][]string{
		{"flow_id", "id"},
	}
}
