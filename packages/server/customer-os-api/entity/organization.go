package entity

import (
	"fmt"
	"time"
)

type OrganizationEntity struct {
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

func (organization OrganizationEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", organization.Id, organization.Name)
}

type OrganizationEntities []OrganizationEntity

func (organization OrganizationEntity) Labels() []string {
	return []string{"Organization"}
}
