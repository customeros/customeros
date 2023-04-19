package entity

import (
	"fmt"
	"time"
)

type StateEntity struct {
	Id        string
	CountryId string
	Name      string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (country StateEntity) ToString() string {
	return fmt.Sprintf("id: %s\ncountryId: %s\nname: %s\ncode: %s", country.Id, country.CountryId, country.Name, country.Code)
}

func (country StateEntity) StateEntity() []string {
	return []string{"State"}
}

func (StateEntity) Labels() []string {
	return []string{
		"State",
	}
}
