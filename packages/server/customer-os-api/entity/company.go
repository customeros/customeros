package entity

import (
	"fmt"
	"time"
)

type CompanyEntity struct {
	Id          string
	Name        string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Description string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Domain      string `neo4jDb:"property:domain;lookupName:DOMAIN;supportCaseSensitive:true"`
	Website     string `neo4jDb:"property:website;lookupName:WEBSITE;supportCaseSensitive:true"`
	Industry    string `neo4jDb:"property:industry;lookupName:INDUSTRY;supportCaseSensitive:true"`
	IsPublic    bool
	CreatedAt   time.Time
	Readonly    bool
}

func (company CompanyEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", company.Id, company.Name)
}

type CompanyEntities []CompanyEntity

func (company CompanyEntity) Labels() []string {
	return []string{"Company"}
}
