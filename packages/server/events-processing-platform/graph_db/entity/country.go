package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
)

// Deprecated
type CountryEntity struct {
	Name      string
	CodeA2    string
	CodeA3    string
	PhoneCode string

	DataloaderKey string
}

type CountryEntities []CountryEntity

func (CountryEntity) Labels(string) []string {
	return []string{
		neo4jutil.NodeLabelCountry,
	}
}
