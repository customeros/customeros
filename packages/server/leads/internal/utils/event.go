package utils

func GenerateEventID() string {
	return GenerateNanoIDWithPrefix("event", 21)
}
