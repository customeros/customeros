package entity

import neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

type DashboardViewResultEntity struct {
	Organization *neo4jentity.OrganizationEntity
	Contact      *neo4jentity.ContactEntity
}
