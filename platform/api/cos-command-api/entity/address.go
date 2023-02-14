package entity

import (
	"fmt"
	"time"
)

type PlaceEntity struct {
	Id            string
	Country       string
	State         string
	City          string
	Address       string
	Address2      string
	Zip           string
	Phone         string
	Fax           string
	CreatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
}

func (address PlaceEntity) ToString() string {
	return fmt.Sprintf("id: %s", address.Id)
}

type PlaceEntities []PlaceEntity

func (address PlaceEntity) Labels() []string {
	return []string{"Address"}
}
