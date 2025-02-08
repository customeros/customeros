package enum

import "fmt"

type IntentSignal string

const (
	IntentChurnRisk       IntentSignal = "churn_risk"
	IntentSupportRequired IntentSignal = "support_required"
)

func (t IntentSignal) String() string {
	return string(t)
}

func GetIntentSignal(s string) (IntentSignal, error) {
	switch IntentSignal(s) {
	case
		IntentChurnRisk,
		IntentSupportRequired:

		return IntentSignal(s), nil

	default:
		return "", fmt.Errorf("invalid IntentSignal: %s", s)
	}
}
