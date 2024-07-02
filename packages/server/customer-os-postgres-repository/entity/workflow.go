package entity

import "time"

type WorkflowType string

const (
	WorkflowTypeIdealCustomerProfile WorkflowType = "IDEAL_CUSTOMER_PROFILE"
	WorkflowTypeIdealContactPersona  WorkflowType = "IDEAL_CONTACT_PERSONA"
)

type Workflow struct {
	ID           uint64       `gorm:"primary_key;autoIncrement:true" json:"id"`
	CreatedAt    time.Time    `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt    time.Time    `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Tenant       string       `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	WorkflowType WorkflowType `gorm:"column:workflow_type;type:varchar(255);NOT NULL" json:"workflowType"`
	Name         string       `gorm:"column:name;type:varchar(255)" json:"name"`
	Condition    string       `gorm:"column:condition;type:text" json:"condition"`
	Live         bool         `gorm:"column:live;type:boolean;DEFAULT:false" json:"live"`
}

func (Workflow) TableName() string {
	return "workflow"
}
