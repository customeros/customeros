package entity

type AgentRegistry struct {
	ID           string  `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name         string  `gorm:"column:name;type:varchar(255);not null;index" json:"name" binding:"required"`
	Capabilities string  `gorm:"column:capabilities;type:text" json:"capabilities"`
	ConfigSchema *string `gorm:"column:config_schema;type:text" json:"config"`
	Goal         string  `gorm:"column:goal;type:text" json:"goal"`
	IsActive     string  `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
	Icom         string  `gorm:"column:icon;type:text" json:"icon"`
}

func (AgentRegistry) TableName() string {
	return "agent_registry"
}
