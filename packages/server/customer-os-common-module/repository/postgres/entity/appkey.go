package entity

type AppKey struct {
	ID     uint64 `gorm:"primary_key;autoIncrement:true" json:"id"`
	AppId  string `gorm:"column:app_id;type:varchar(255);NOT NULL" json:"appId" binding:"required"`
	Key    string `gorm:"column:key;type:varchar(255);NOT NULL;index:idx_key,unique" json:"key" binding:"required"`
	Active bool   `gorm:"column:active;type:bool;NOT NULL" json:"active" binding:"required"`
}

func (AppKey) TableName() string {
	return "app_keys"
}
