package entity

import (
	"fmt"
	"time"
)

type JobRoleEntity struct {
	Id                  string
	JobTitle            string
	Primary             bool
	ResponsibilityLevel int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Source              DataSource
	SourceOfTruth       DataSource
	AppSource           string
}

func (jobRole JobRoleEntity) ToString() string {
	return fmt.Sprintf("id: %s\njob title: %s", jobRole.Id, jobRole.JobTitle)
}

type JobRoleEntities []JobRoleEntity

func (jobRole JobRoleEntity) Labels(tenant string) []string {
	return []string{"JobRole", "JobRole_" + tenant}
}
