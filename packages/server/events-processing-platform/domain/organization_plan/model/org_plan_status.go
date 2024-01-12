package model

type OrgPlanStatusString string

const (
	OrgPlanStatusNotStarted OrgPlanStatusString = "NOT_STARTED"
	OrgPlanStatusOnTrack    OrgPlanStatusString = "ON_TRACK"
	OrgPlanStatusLate       OrgPlanStatusString = "LATE"
	OrgPlanStatusDone       OrgPlanStatusString = "DONE"
)

type OrgPlanStatus int32

const (
	NotStarted OrgPlanStatus = iota
	OnTrack
	Late
	Done
)

func (os OrgPlanStatus) String() string {
	switch os {
	case NotStarted:
		return string(OrgPlanStatusNotStarted)
	case OnTrack:
		return string(OrgPlanStatusOnTrack)
	case Late:
		return string(OrgPlanStatusLate)
	case Done:
		return string(OrgPlanStatusDone)
	default:
		return string(OrgPlanStatusNotStarted)
	}
}

/////// milestone status

type OrgPlanMilestoneStatusString string

const (
	OrgPlanMilestoneStatusNotStarted OrgPlanMilestoneStatusString = "NOT_STARTED"
	OrgPlanMilestoneStatusStarted    OrgPlanMilestoneStatusString = "STARTED"
	OrgPlanMilestoneStatusDone       OrgPlanMilestoneStatusString = "DONE"
)

type OrgPlanMilestoneStatus int32

const (
	MilestoneNotStarted OrgPlanMilestoneStatus = iota
	MilestoneStarted
	MilestoneDone
)

func (os OrgPlanMilestoneStatus) String() string {
	switch os {
	case MilestoneNotStarted:
		return string(OrgPlanMilestoneStatusNotStarted)
	case MilestoneStarted:
		return string(OrgPlanMilestoneStatusStarted)
	case MilestoneDone:
		return string(OrgPlanMilestoneStatusDone)
	default:
		return string(OrgPlanMilestoneStatusNotStarted)
	}
}

///// milestone task status

type OrgPlanMilestoneTaskStatusString string

const (
	OrgPlanMilestoneTaskStatusNotDone OrgPlanMilestoneTaskStatusString = "NOT_DONE"
	OrgPlanMilestoneTaskStatusSkipped OrgPlanMilestoneTaskStatusString = "SKIPPED"
	OrgPlanMilestoneTaskStatusDone    OrgPlanMilestoneTaskStatusString = "DONE"
)

type OrgPlanMilestoneTaskStatus int32

const (
	TaskNotDone OrgPlanMilestoneTaskStatus = iota
	TaskSkipped
	TaskDone
)

func (os OrgPlanMilestoneTaskStatus) String() string {
	switch os {
	case TaskNotDone:
		return string(OrgPlanMilestoneTaskStatusNotDone)
	case TaskSkipped:
		return string(OrgPlanMilestoneTaskStatusSkipped)
	case TaskDone:
		return string(OrgPlanMilestoneStatusDone)
	default:
		return string(OrgPlanMilestoneTaskStatusNotDone)
	}
}
