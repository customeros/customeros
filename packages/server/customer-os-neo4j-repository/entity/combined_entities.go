package neo4j_entity

type OrganizationWithJobRole struct {
	DataLoaderKey
	Organization OrganizationEntity
	JobRole      JobRoleEntity
}

type OrganizationWithJobRoleEntities []OrganizationWithJobRole
