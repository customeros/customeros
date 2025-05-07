package postgres_entity

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"time"
)

type APICallLog struct {
	ID           string         `gorm:"column:id;primaryKey;type:varchar(25)"`
	Vendor       enum.APIVendor `gorm:"column:vendor;type:varchar(255);index;not null"`
	Method       string         `gorm:"column:method;type:varchar(255);not null"`
	URL          string         `gorm:"column:url;type:varchar(255);not null"`
	RequestID    string         `gorm:"column:request_id;type:varchar(55);not null"`
	RequestBody  []byte         `gorm:"column:request_body;type:bytea;not null"`
	Timestamp    time.Time      `gorm:"column:timestamp;type:timestamptz;not null"`
	Duration     int            `gorm:"column:duration;type:int;not null"`
	StatusCode   *int           `gorm:"column:status_code;type:int;not null"`
	ResponseBody *[]byte        `gorm:"column:response_body;type:bytea"`
	ErrorMessage *string        `gorm:"column:error_message;type:text;not null"`
}

func (APICallLog) TableName() string {
	return "api_call_logs"
}
