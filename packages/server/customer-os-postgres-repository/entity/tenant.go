package entity

type Tenant struct {
	Name string `gorm:"primary_key;type:varchar(255);not null"`
}

func (Tenant) TableName() string {
	return "tenant"
}
