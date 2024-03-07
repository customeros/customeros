package enum

type RenewalLikelihood string

const (
	RenewalLikelihoodHigh   RenewalLikelihood = "HIGH"
	RenewalLikelihoodMedium RenewalLikelihood = "MEDIUM"
	RenewalLikelihoodLow    RenewalLikelihood = "LOW"
	RenewalLikelihoodZero   RenewalLikelihood = "ZERO"
)

var AllRenewalLikelihood = []RenewalLikelihood{
	RenewalLikelihoodHigh,
	RenewalLikelihoodMedium,
	RenewalLikelihoodLow,
	RenewalLikelihoodZero,
}

func (e RenewalLikelihood) IsValid() bool {
	switch e {
	case RenewalLikelihoodHigh, RenewalLikelihoodMedium, RenewalLikelihoodLow, RenewalLikelihoodZero:
		return true
	}
	return false
}

func (e RenewalLikelihood) String() string {
	return string(e)
}
