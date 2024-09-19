package entity

type BrowserAutomationsRunResult struct {
	Id         uint64 `gorm:"primary_key;autoIncrement:true"`
	RunId      int    `gorm:"column:run_id;type:integer"`
	ResultData string `gorm:"column:result_data;type:text;"`
}

func (BrowserAutomationsRunResult) TableName() string {
	return "browser_automation_run_results"
}
