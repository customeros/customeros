package entity

import "time"

type ColumnView struct {
	ColumnType string `json:"columnType"`
	Width      int    `json:"width"`
	Visible    bool   `json:"visible"`
}

type Columns struct {
	Columns []ColumnView `json:"columns"`
}

type TableViewDefinition struct {
	ID          uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Tenant      string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	UserId      string    `gorm:"column:user_id;type:varchar(255)" json:"userId"`
	TableType   string    `gorm:"column:table_type;type:varchar(255);NOT NULL" json:"tableType"`
	Name        string    `gorm:"column:table_name;type:varchar(255);NOT NULL" json:"tableName"`
	Order       int       `gorm:"column:position;type:int;NOT NULL" json:"order"`
	Icon        string    `gorm:"column:icon;type:varchar(255)" json:"icon"`
	Filters     string    `gorm:"column:filters;type:text" json:"filters"`
	Sorting     string    `gorm:"column:sorting;type:text" json:"sorting"`
	ColumnsJson string    `gorm:"column:columns;type:text" json:"columns"`
}

func (TableViewDefinition) TableName() string {
	return "table_view_definition"
}
