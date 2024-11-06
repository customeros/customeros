package entity

type BrowserConfig struct {
	Id     int    `gorm:"primary_key;autoIncrement:true"`
	UserId string `gorm:"column:user_id;type:varchar(36);"`
	Tenant string `gorm:"column:tenant;type:varchar(36);"`
	Status string `gorm:"column:session_status;type:browser_config_session_status;"`
}

func (BrowserConfig) TableName() string {
	return "browser_configs"
}
