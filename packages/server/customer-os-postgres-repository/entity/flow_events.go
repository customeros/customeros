package entity

type FlowEvent struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalSystem string `gorm:"column:external_system;type:varchar(255);not null;index" json:"externalSystem" binding:"required"`
	Resource       string `gorm:"column:resource;type:varchar(255);not null;index" json:"resource" binding:"required"`
	Action         string `gorm:"column:action;type:varchar(255);not null;index" json:"action" binding:"required"`
	EventName      string `gorm:"column:event_name;type:varchar(255);not null" json:"eventName" binding:"required"`
	Description    string `gorm:"column:description;type:varchar(255);not null" json:"description" binding:"required"`
	Enabled        bool   `gorm:"column:enabled;type:boolean;default:true" json:"enabled"`
}

func (FlowEvent) TableName() string {
	return "flow_events"
}

func (FlowEvent) UniqueIndex() [][]string {
	return [][]string{
		{"external_system", "resource", "action"},
	}
}
