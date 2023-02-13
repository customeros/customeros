package entity

import (
	"fmt"
	"time"
)

type LocationEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    DataSource
	AppSource string
}

func (location LocationEntity) ToString() string {
	return fmt.Sprintf("id: %s name: %s", location.Id, location.Name)
}

type LocationEntities []LocationEntity

func (location LocationEntity) Labels(tenant string) []string {
	return []string{"Location", "Location_" + tenant}
}
