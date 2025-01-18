package postgres_entity

import "time"

type DMARCMonitoring struct {
	ID            string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant        string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	EmailProvider string    `gorm:"column:email_provider;type:varchar(255)" json:"emailProvider"`
	Domain        string    `gorm:"column:domain;type:varchar(255)" json:"domain"`
	ReportStart   time.Time `gorm:"column:report_start;type:timestamp" json:"reportStart"`
	ReportEnd     time.Time `gorm:"column:report_end;type:timestamp" json:"reportEnd"`
	MessageCount  int       `gorm:"column:message_count;type:integer" json:"messageCount"`
	SPFPass       int       `gorm:"column:spf_pass;type:integer" json:"spfPass"`
	DKIMPass      int       `gorm:"column:dkim_pass;type:integer" json:"dkimPass"`
	DMARCPass     int       `gorm:"column:dmarc_pass;type:integer" json:"dmarcPass"`
	Data          string    `gorm:"type:text"`
}

func (DMARCMonitoring) TableName() string {
	return "mailstack_dmarc_monitoring"
}
