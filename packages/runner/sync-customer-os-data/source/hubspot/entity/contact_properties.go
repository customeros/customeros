package entity

type ContactProperties struct {
	OwnerId        string `gorm:"column:hubspot_owner_id"`
	LifecycleStage string `gorm:"column:lifecyclestage"`
}

type ContactPropertiesList []ContactProperties

func (ContactProperties) TableName() string {
	return "contacts_properties"
}
