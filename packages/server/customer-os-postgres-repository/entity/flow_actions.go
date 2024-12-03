package entity

type FlowAction struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	EventName   string `gorm:"column:event_name;type:varchar(255);not null;index" json:"eventName" binding:"required"`
	ActionType  string `gorm:"column:action_type;type:varchar(255);not null;index" json:"actionType" binding:"required"`
	Description string `gorm:"column:description;type:varchar(255);not null" json:"description" binding:"required"`
	Enabled     bool   `gorm:"column:enabled;type:boolean;default:true" json:"enabled"`
}

func (FlowAction) TableName() string {
	return "flow_actions"
}

func (FlowAction) UniqueIndex() [][]string {
	return [][]string{
		{"event_name", "action_type"},
	}
}
