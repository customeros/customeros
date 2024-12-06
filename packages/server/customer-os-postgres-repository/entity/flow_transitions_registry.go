package entity

type FlowTransitionsRegistry struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	FromNodeType string `gorm:"column:from_node_type;type:varchar(255);not null;index" json:"fromNodeType" binding:"required"`
	FromNode     string `gorm:"column:from_node;type:varchar(255);not null;index" json:"fromNode" binding:"required"`
	ToNodeType   string `gorm:"column:to_node_type;type:varchar(255);not null;index" json:"toNodeType" binding:"required"`
	ToNode       string `gorm:"column:to_node;type:varchar(255);not null;index" json:"toNode" binding:"required"`
	Enabled      bool   `gorm:"column:enabled;type:boolean;default:true" json:"enabled"`

	// Add uniqueness constraint
	UniqueTransition string `gorm:"uniqueIndex:idx_from_to_nodes;type:VIRTUAL AS (CONCAT(from_node,'-',to_node)) STORED"`
}

func (FlowTransitionsRegistry) TableName() string {
	return "flow_transitions_registry"
}
