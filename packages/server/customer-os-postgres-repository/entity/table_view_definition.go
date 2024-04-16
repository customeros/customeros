package entity

import "time"

type TableViewDefinition struct {
	// tenant, event, webhook, api key
	ID         uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`
	TenantName string    `gorm:"column:tenant_name;type:varchar(255);NOT NULL" json:"tenantName"`
	UserId     string    `gorm:"column:user_id;type:varchar(255)" json:"userId"`
	TableType  string    `gorm:"column:table_type;type:varchar(255);NOT NULL" json:"tableType"`
	Name       string    `gorm:"column:table_name;type:varchar(255);NOT NULL" json:"tableName"`
	Order      int       `gorm:"column:order;type:int;NOT NULL" json:"order"`
	Icon       string    `gorm:"column:icon;type:varchar(255)" json:"icon"`
	Filters    string    `gorm:"column:filters;type:json" json:"filters"`
	Sorting    string    `gorm:"column:sorting;type:json" json:"sorting"`
	Columns    string    `gorm:"column:columns;type:json" json:"columns"`
}

func (TableViewDefinition) TableName() string {
	return "table_view_definition"
}
