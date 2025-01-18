package postgres_entity

import (
	"time"

	"gorm.io/datatypes"
)

type FlowEdge struct {
	ID         string     `gorm:"primaryKey;type:varchar(50)" json:"id"`
	FlowID     string     `gorm:"column:flow_id;type:varchar(50);not null;index" json:"flowId"`
	FromNodeID string     `gorm:"column:from_node_id;type:varchar(50);not null" json:"fromNodeId"`
	ToNodeID   string     `gorm:"column:to_node_id;type:varchar(50);not null" json:"toNodeId"`
	Condition  *string    `gorm:"column:condition;type:varchar(255)" json:"condition"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// used to store conditional logic
	Data *datatypes.JSON `gorm:"column:data;type:jsonb" json:"data"`
}

func (FlowEdge) TableName() string {
	return "flow_edges"
}

func (FlowEdge) UniqueIndex() [][]string {
	return [][]string{
		{"flow_id", "from_node_id", "to_node_id"},
	}
}
