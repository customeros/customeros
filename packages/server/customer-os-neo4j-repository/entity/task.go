package neo4j_entity

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"time"
)

type TaskProperty string

const (
	TaskPropertyId          TaskProperty = "id"
	TaskPropertySubject     TaskProperty = "subject"
	TaskPropertyDescription TaskProperty = "description"
	TaskPropertyStatus      TaskProperty = "status"
	TaskPropertyDueAt       TaskProperty = "dueAt"
	TaskPropertyCreatedAt   TaskProperty = "createdAt"
	TaskPropertyUpdatedAt   TaskProperty = "updatedAt"
	TaskPropertySource      TaskProperty = "source"
	TaskPropertyAppSource   TaskProperty = "appSource"
)

type TaskEntity struct {
	DataLoaderKey
	Id          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DueAt       *time.Time
	Subject     string
	Description string
	Status      enum.TaskStatus
	Source      DataSource
	AppSource   string
}

type TaskEntities []TaskEntity

func (e *TaskEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *TaskEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
