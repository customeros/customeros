package entity

type AutomationRegistry struct {
	ID          uint64 `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Name        string `gorm:"column:name;type:varchar(255);not null;index" json:"name" binding:"required"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
	Type        string `gorm:"column:type;type:varchar(50);not null;index" json:"type" binding:"required"`
	IsActive    string `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
}

func (AutomationRegistry) TableName() string {
	return "automation_registry"
}
