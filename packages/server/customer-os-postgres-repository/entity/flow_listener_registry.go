package entity

type FlowListenerRegistry struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalSystem string `gorm:"column:external_system;type:varchar(255);not null;index" json:"externalSystem" binding:"required"`
	ListenerEvent  string `gorm:"column:listener_event;type:varchar(255);not null;index" json:"listenerEvent" binding:"required"`
	FriendlyName   string `gorm:"column:friendly_name;type:varchar(255);not null;index" json:"friendlyName" binding:"required"`
	Description    string `gorm:"column:description;type:varchar(255);not null" json:"description" binding:"required"`
	Status         string `gorm:"column:status;type:varchar(50)" json:"status"`
}

func (FlowListenerRegistry) TableName() string {
	return "flow_listener_registry"
}
