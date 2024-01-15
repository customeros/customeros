package event

const (
	OrganizationPlanCreateV1           = "V1_ORGANIZATION_PLAN_CREATE"
	OrganizationPlanUpdateV1           = "V1_ORGANIZATION_PLAN_UPDATE"
	OrganizationPlanMilestoneCreateV1  = "V1_ORGANIZATION_PLAN_MILESTONE_CREATE"
	OrganizationPlanMilestoneUpdateV1  = "V1_ORGANIZATION_PLAN_MILESTONE_UPDATE"
	OrganizationPlanMilestoneReorderV1 = "V1_ORGANIZATION_PLAN_MILESTONE_REORDER"
)

const (
	FieldMaskName          = "name"
	FieldMaskRetired       = "retired"
	FieldMaskOrder         = "order"
	FieldMaskOptional      = "optional"
	FieldMaskDurationHours = "durationHours"
	FieldMaskItems         = "items"
)
