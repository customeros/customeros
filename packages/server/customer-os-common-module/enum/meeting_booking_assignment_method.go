package enum

type MeetingBookingAssignmentMethod string

const (
	MeetingBookingAssignmentMethodRoundRobinMaxFairness     MeetingBookingAssignmentMethod = "ROUND_ROBIN_MAX_FAIRNESS"
	MeetingBookingAssignmentMethodRoundRobinMaxAvailability MeetingBookingAssignmentMethod = "ROUND_ROBIN_MAX_AVAILABILITY"
	MeetingBookingAssignmentMethodCustom                    MeetingBookingAssignmentMethod = "CUSTOM"
)

var AllMeetingBookingAssignmentMethod = []MeetingBookingAssignmentMethod{
	MeetingBookingAssignmentMethodRoundRobinMaxFairness,
	MeetingBookingAssignmentMethodRoundRobinMaxAvailability,
	MeetingBookingAssignmentMethodCustom,
}

var DefaultMeetingBookingAssignmentMethod = MeetingBookingAssignmentMethodRoundRobinMaxAvailability

func GetMeetingBookingAssignmentMethod(s string) MeetingBookingAssignmentMethod {
	if IsValidMeetingBookingAssignmentMethod(s) {
		return MeetingBookingAssignmentMethod(s)
	}
	return DefaultMeetingBookingAssignmentMethod
}

func IsValidMeetingBookingAssignmentMethod(s string) bool {
	for _, ds := range AllMeetingBookingAssignmentMethod {
		if ds == MeetingBookingAssignmentMethod(s) {
			return true
		}
	}
	return false
}

func (m MeetingBookingAssignmentMethod) String() string {
	return string(m)
}
