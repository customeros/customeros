package enum

type InteractionSessionStatus string

const (
	InteractionSessionStatusActive   InteractionSessionStatus = "ACTIVE"
	InteractionSessionStatusInactive InteractionSessionStatus = "INACTIVE"
)

func DecodeInteractionSessionStatus(s string) InteractionSessionStatus {
	switch InteractionSessionStatus(s) {
	case InteractionSessionStatusActive, InteractionSessionStatusInactive:
		return InteractionSessionStatus(s)
	}
	return ""
}

func (i InteractionSessionStatus) String() string {
	return string(i)
}
