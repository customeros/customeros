package repository

import "gorm.io/gorm"

type Repositories struct{}

func InitRepositories(leadsDB, warehouseDB *gorm.DB) *Repositories {
	return &Repositories{}
}
