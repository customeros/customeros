package entity

type BrowserAutomationsRun struct {
	Id              int    `gorm:"primary_key;autoIncrement:true"`
	BrowserConfigId int    `gorm:"column:browser_config_id;type:integer"`
	UserId          string `gorm:"column:user_id;type:varchar(36);"`
	Tenant          string `gorm:"column:tenant;type:varchar(36);"`
	Type            string `gorm:"column:type;type:browser_automation_run_type;"`
	Status          string `gorm:"column:status;type:browser_automation_run_status;"`
	Payload         string `gorm:"column:payload;type:text;"`
}

func (BrowserAutomationsRun) TableName() string {
	return "browser_automation_runs"
}
