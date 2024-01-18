package model

type OrganizationPlanStatusString string

const (
	OrganizationPlanStatusNotStarted OrganizationPlanStatusString = "NOT_STARTED"
	OrganizationPlanStatusOnTrack    OrganizationPlanStatusString = "ON_TRACK"
	OrganizationPlanStatusLate       OrganizationPlanStatusString = "LATE"
	OrganizationPlanStatusDone       OrganizationPlanStatusString = "DONE"
)

type OrganizationPlanStatus int32

const (
	NotStarted OrganizationPlanStatus = iota
	OnTrack
	Late
	Done
)

func (os OrganizationPlanStatus) String() string {
	switch os {
	case NotStarted:
		return string(OrganizationPlanStatusNotStarted)
	case OnTrack:
		return string(OrganizationPlanStatusOnTrack)
	case Late:
		return string(OrganizationPlanStatusLate)
	case Done:
		return string(OrganizationPlanStatusDone)
	default:
		return string(OrganizationPlanStatusNotStarted)
	}
}

/////// milestone status

type OrganizationPlanMilestoneStatusString string

const (
	OrganizationPlanMilestoneStatusNotStarted OrganizationPlanMilestoneStatusString = "NOT_STARTED"
	OrganizationPlanMilestoneStatusStarted    OrganizationPlanMilestoneStatusString = "STARTED"
	OrganizationPlanMilestoneStatusDone       OrganizationPlanMilestoneStatusString = "DONE"
)

type OrganizationPlanMilestoneStatus int32

const (
	MilestoneNotStarted OrganizationPlanMilestoneStatus = iota
	MilestoneStarted
	MilestoneDone
)

func (os OrganizationPlanMilestoneStatus) String() string {
	switch os {
	case MilestoneNotStarted:
		return string(OrganizationPlanMilestoneStatusNotStarted)
	case MilestoneStarted:
		return string(OrganizationPlanMilestoneStatusStarted)
	case MilestoneDone:
		return string(OrganizationPlanMilestoneStatusDone)
	default:
		return string(OrganizationPlanMilestoneStatusNotStarted)
	}
}

///// milestone task status

type OrganizationPlanMilestoneTaskStatusString string

const (
	OrganizationPlanMilestoneTaskStatusNotDone OrganizationPlanMilestoneTaskStatusString = "NOT_DONE"
	OrganizationPlanMilestoneTaskStatusSkipped OrganizationPlanMilestoneTaskStatusString = "SKIPPED"
	OrganizationPlanMilestoneTaskStatusDone    OrganizationPlanMilestoneTaskStatusString = "DONE"
)

type OrganizationPlanMilestoneTaskStatus int32

const (
	TaskNotDone OrganizationPlanMilestoneTaskStatus = iota
	TaskSkipped
	TaskDone
)

func (os OrganizationPlanMilestoneTaskStatus) String() string {
	switch os {
	case TaskNotDone:
		return string(OrganizationPlanMilestoneTaskStatusNotDone)
	case TaskSkipped:
		return string(OrganizationPlanMilestoneTaskStatusSkipped)
	case TaskDone:
		return string(OrganizationPlanMilestoneStatusDone)
	default:
		return string(OrganizationPlanMilestoneTaskStatusNotDone)
	}
}
