package entity

type RawState string

const (
	PENDING RawState = "PENDING"
	SENT    RawState = "SENT"
	SKIPPED RawState = "SKIPPED"
	ERROR   RawState = "ERROR"
)

func (rawState RawState) String() string {
	return string(rawState)
}

func DecodeRawState(rawState string) RawState {
	switch rawState {
	case "PENDING":
		return PENDING
	case "SENT":
		return SENT
	case "SKIPPED":
		return SKIPPED
	case "ERROR":
		return ERROR
	default:
		return PENDING
	}
}
