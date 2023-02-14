package entity

import (
	"fmt"
	"time"
)

// Deprecated
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
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

func (place PlaceEntity) ToString() string {
	return fmt.Sprintf("id: %s", place.Id)
}

type PlaceEntities []PlaceEntity

func (place PlaceEntity) Labels(tenant string) []string {
	return []string{"Place", "Place_" + tenant}
}
