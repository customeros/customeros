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
	Country                string `gorm:"column:country"`
	State                  string `gorm:"column:state"`
	City                   string `gorm:"column:city"`
	Address                string `gorm:"column:address"`
	Address2               string `gorm:"column:address2"`
	Zip                    string `gorm:"column:zip"`
	Phone                  string `gorm:"column:phone"`
	Employees              int64  `gorm:"column:numberofemployees"`
}

type CompanyPropertiesList []CompanyProperties

func (CompanyProperties) TableName() string {
	return "companies_properties"
}
