package entity

import "time"

type ColumnView struct {
	ColumnId      int    `json:"columnId"`
	ColumnType    string `json:"columnType"`
	Width         int    `json:"width"`
	Visible       bool   `json:"visible"`
	Name          string `json:"name"`
	Filter        string `json:"filter"`
	DefaultFilter string `json:"defaultFilter"`
}

type Columns struct {
	Columns []ColumnView `json:"columns"`
}

type TableViewDefinition struct {
	ID             uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Tenant         string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	UserId         string    `gorm:"column:user_id;type:varchar(255)" json:"userId"`
	TableId        string    `gorm:"column:table_id;type:varchar(255);NOT NULL;DEFAULT:''" json:"tableId"`
	TableType      string    `gorm:"column:table_type;type:varchar(255);NOT NULL" json:"tableType"`
	Name           string    `gorm:"column:table_name;type:varchar(255);NOT NULL" json:"tableName"`
	Order          int       `gorm:"column:position;type:int;NOT NULL" json:"order"`
	Icon           string    `gorm:"column:icon;type:varchar(255)" json:"icon"`
	Filters        string    `gorm:"column:filters;type:text" json:"filters"`
	DefaultFilters string    `gorm:"column:default_filters;type:text" json:"defaultFilters"`
	Sorting        string    `gorm:"column:sorting;type:text" json:"sorting"`
	ColumnsJson    string    `gorm:"column:columns;type:text" json:"columns"`
	IsPreset       bool      `gorm:"column:is_preset;type:boolean;NOT NULL;DEFAULT:false" json:"isPreset"`
	IsShared       bool      `gorm:"column:is_shared;type:boolean;NOT NULL;DEFAULT:false" json:"isShared"`
}

func (TableViewDefinition) TableName() string {
	return "table_view_definition"
}
