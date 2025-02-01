package enum

import "fmt"

type IntentSignal string

const (
	IntentChurnRisk       IntentSignal = "churn_risk"
	IntentSupportRequired IntentSignal = "support_required"
	IntentIcpCheck        IntentSignal = "icp_check"
)

func (t IntentSignal) String() string {
	return string(t)
}

func GetIntentSignal(s string) (IntentSignal, error) {
	switch IntentSignal(s) {
	case
		IntentIcpCheck,
		IntentChurnRisk,
		IntentSupportRequired:

		return IntentSignal(s), nil

	default:
		return "", fmt.Errorf("invalid IntentSignal: %s", s)
	}
}
