package entity

type AutomationEventRegistry struct {
	ID          string `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Source      string `gorm:"column:source;type:varchar(50);not null;index" json:"source" binding:"required"`
	Name        string `gorm:"column:name;type:varchar(255);not null;index" json:"Name" binding:"required"`
	Description string `gorm:"column:description;type:varchar(255);not null" json:"description" binding:"required"`
	Type        string `gorm:"column:type;type:varchar(255);not null" json:"type" binding:"required"`
	IsActive    string `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
}

func (AutomationEventRegistry) TableName() string {
	return "automation_event_registry"
}
