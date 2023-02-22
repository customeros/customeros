package entity

import (
	"fmt"
	"time"
)

type CountryEntity struct {
	Id        string
	Name      string
	CodeA2    string
	CodeA3    string
	PhoneCode string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (country CountryEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s\ncodeA2: %s\ncodeA3: %s\nphoneCode: %s", country.Id, country.Name, country.CodeA2, country.CodeA3, country.PhoneCode)
}

func (country CountryEntity) CountryEntity() []string {
	return []string{"Country"}
}
