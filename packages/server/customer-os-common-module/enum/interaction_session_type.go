package enum

type InteractionSessionType string

const (
	InteractionSessionTypeThread InteractionSessionType = "THREAD"
)

func DecodeInteractionSessionType(s string) InteractionSessionType {
	switch InteractionSessionType(s) {
	case InteractionSessionTypeThread:
		return InteractionSessionType(s)
	}
	return ""
}

func (i InteractionSessionType) String() string {
	return string(i)
}
