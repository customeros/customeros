package entity

// Deprecated
type OpportunityRenewalLikelihood string

const (
	OpportunityRenewalLikelihoodHigh   OpportunityRenewalLikelihood = "HIGH"
	OpportunityRenewalLikelihoodMedium OpportunityRenewalLikelihood = "MEDIUM"
	OpportunityRenewalLikelihoodLow    OpportunityRenewalLikelihood = "LOW"
	OpportunityRenewalLikelihoodZero   OpportunityRenewalLikelihood = "ZERO"
)

// Deprecated
var AllOpportunityRenewalLikelihoods = []OpportunityRenewalLikelihood{
	OpportunityRenewalLikelihoodHigh,
	OpportunityRenewalLikelihoodMedium,
	OpportunityRenewalLikelihoodLow,
	OpportunityRenewalLikelihoodZero,
}

// Deprecated
func GetOpportunityRenewalLikelihood(s string) OpportunityRenewalLikelihood {
	if IsValidOpportunityRenewalLikelihood(s) {
		return OpportunityRenewalLikelihood(s)
	}
	return OpportunityRenewalLikelihoodZero
}

// Deprecated
func IsValidOpportunityRenewalLikelihood(s string) bool {
	for _, ms := range AllOpportunityRenewalLikelihoods {
		if ms == OpportunityRenewalLikelihood(s) {
			return true
		}
	}
	return false
}
