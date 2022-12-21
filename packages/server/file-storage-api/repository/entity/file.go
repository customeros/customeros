package entity

type File struct {
	ID        uint64 `gorm:"primary_key;autoIncrement:true" json:"id"`
	TenantId  string `gorm:"column:tenant_id;type:varchar(255);NOT NULL" json:"tenantId" binding:"required"`
	Name      string `gorm:"column:name;type:varchar(255);NOT NULL;" json:"name" binding:"required"`
	Extension string `gorm:"column:extension;type:varchar(255);NOT NULL;" json:"extension" binding:"required"`
	MIME      string `gorm:"column:mime;type:varchar(255);NOT NULL;" json:"mime" binding:"required"`
}

func (File) TableName() string {
	return "files"
}
