package entity

type UserToTenant struct {
	ID       uint64 `gorm:"primary_key;autoIncrement:true" json:"id"`
	Username string `gorm:"column:username;type:varchar(255);NOT NULL" json:"username" binding:"required"`
	Tenant   string `gorm:"column:tenant;type:varchar(255);NOT NULL;" json:"key" binding:"required"`
}

func (UserToTenant) TableName() string {
	return "user_to_tenant"
}
