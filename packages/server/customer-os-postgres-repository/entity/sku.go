package postgres_entity

import "time"

type SkuEntity struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);not null" json:"tenant" binding:"required"`
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name" binding:"required"`
	Price     float64   `gorm:"column:price;type:decimal(10,2);not null" json:"price" binding:"required"`
}

func (SkuEntity) TableName() string {
	return "sku"
}
