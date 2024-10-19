package entity

type OrganizationWithJobRole struct {
	DataLoaderKey
	Organization OrganizationEntity
	JobRole      JobRoleEntity
}

type OrganizationWithJobRoleEntities []OrganizationWithJobRole
