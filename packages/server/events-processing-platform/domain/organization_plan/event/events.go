package event

const (
	OrgPlanCreateV1           = "V1_ORG_PLAN_CREATE"
	OrgPlanUpdateV1           = "V1_ORG_PLAN_UPDATE"
	OrgPlanMilestoneCreateV1  = "V1_ORG_PLAN_MILESTONE_CREATE"
	OrgPlanMilestoneUpdateV1  = "V1_ORG_PLAN_MILESTONE_UPDATE"
	OrgPlanMilestoneReorderV1 = "V1_ORG_PLAN_MILESTONE_REORDER"
)

const (
	FieldMaskName          = "name"
	FieldMaskRetired       = "retired"
	FieldMaskOrder         = "order"
	FieldMaskOptional      = "optional"
	FieldMaskDurationHours = "durationHours"
	FieldMaskItems         = "items"
)
