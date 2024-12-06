package entity

type FlowActionRegistry struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Action       string `gorm:"column:action;type:varchar(255);not null;index" json:"action" binding:"required"`
	FriendlyName string `gorm:"column:friendly_name;type:varchar(255);not null;index" json:"FriendlyName" binding:"required"`
	Description  string `gorm:"column:description;type:varchar(255);not null" json:"description" binding:"required"`
	Enabled      bool   `gorm:"column:enabled;type:boolean;default:true" json:"enabled"`
}

func (FlowActionRegistry) TableName() string {
	return "flow_action_registry"
}
