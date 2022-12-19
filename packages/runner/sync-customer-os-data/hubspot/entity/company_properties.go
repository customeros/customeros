package entity

type CompanyProperties struct {
	AirbyteAbId            string `gorm:"column:_airbyte_ab_id"`
	AirbyteCompaniesHashid string `gorm:"column:_airbyte_companies_hashid"`
	Name                   string `gorm:"column:name"`
	Description            string `gorm:"column:description"`
	Domain                 string `gorm:"column:domain"`
	Website                string `gorm:"column:website"`
	Industry               string `gorm:"column:industry"`
	IsPublic               bool   `gorm:"column:is_public"`
}

type CompanyPropertiesList []CompanyProperties

func (CompanyProperties) TableName() string {
	return "companies_properties"
}
