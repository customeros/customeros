package enum

type RenewalLikelihood string

const (
	// still used in opportunities, TODO migrate to V2
	RenewalLikelihoodHigh   RenewalLikelihood = "HIGH"
	RenewalLikelihoodMedium RenewalLikelihood = "MEDIUM"
	RenewalLikelihoodLow    RenewalLikelihood = "LOW"
	RenewalLikelihoodZero   RenewalLikelihood = "ZERO"

	RenewalLikelihoodHighV2   RenewalLikelihood = "HIGH_RENEWAL"
	RenewalLikelihoodMediumV2 RenewalLikelihood = "MEDIUM_RENEWAL"
	RenewalLikelihoodLowV2    RenewalLikelihood = "LOW_RENEWAL"
	RenewalLikelihoodZeroV2   RenewalLikelihood = "ZERO_RENEWAL"
)

func (e RenewalLikelihood) ToV2() RenewalLikelihood {
	switch e {
	case RenewalLikelihoodHigh,
		RenewalLikelihoodHighV2:
		return RenewalLikelihoodHighV2
	case RenewalLikelihoodMedium,
		RenewalLikelihoodMediumV2:
		return RenewalLikelihoodMediumV2
	case RenewalLikelihoodLow,
		RenewalLikelihoodLowV2:
		return RenewalLikelihoodLowV2
	case RenewalLikelihoodZero,
		RenewalLikelihoodZeroV2:
		return RenewalLikelihoodZeroV2
	}
	return ""
}

func (e RenewalLikelihood) IsValid() bool {
	switch e {
	case RenewalLikelihoodHigh, RenewalLikelihoodMedium, RenewalLikelihoodLow, RenewalLikelihoodZero,
		RenewalLikelihoodHighV2, RenewalLikelihoodMediumV2, RenewalLikelihoodLowV2, RenewalLikelihoodZeroV2:
		return true
	}
	return false
}

func (e RenewalLikelihood) String() string {
	return string(e)
}

func DecodeRenewalLikelihood(input string) RenewalLikelihood {
	switch input {
	case "HIGH":
		return RenewalLikelihoodHigh
	case "MEDIUM":
		return RenewalLikelihoodMedium
	case "LOW":
		return RenewalLikelihoodLow
	case "ZERO":
		return RenewalLikelihoodZero
	case "HIGH_RENEWAL":
		return RenewalLikelihoodHighV2
	case "MEDIUM_RENEWAL":
		return RenewalLikelihoodMediumV2
	case "LOW_RENEWAL":
		return RenewalLikelihoodLowV2
	case "ZERO_RENEWAL":
		return RenewalLikelihoodZeroV2
	}
	return ""
}
