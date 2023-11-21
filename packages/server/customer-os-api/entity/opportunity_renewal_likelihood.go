package entity

type OpportunityRenewalLikelihood string

const (
	OpportunityRenewalLikelihoodHigh   OpportunityRenewalLikelihood = "HIGH_RENEWAL"
	OpportunityRenewalLikelihoodMedium OpportunityRenewalLikelihood = "MEDIUM_RENEWAL"
	OpportunityRenewalLikelihoodLow    OpportunityRenewalLikelihood = "LOW_RENEWAL"
	OpportunityRenewalLikelihoodZero   OpportunityRenewalLikelihood = "ZERO_RENEWAL"
)

var AllOpportunityRenewalLikelihoods = []OpportunityRenewalLikelihood{
	OpportunityRenewalLikelihoodHigh,
	OpportunityRenewalLikelihoodMedium,
	OpportunityRenewalLikelihoodLow,
	OpportunityRenewalLikelihoodZero,
}

func GetOpportunityRenewalLikelihood(s string) OpportunityRenewalLikelihood {
	if IsValidOpportunityRenewalLikelihood(s) {
		return OpportunityRenewalLikelihood(s)
	}
	return OpportunityRenewalLikelihoodZero
}

func IsValidOpportunityRenewalLikelihood(s string) bool {
	for _, ms := range AllOpportunityRenewalLikelihoods {
		if ms == OpportunityRenewalLikelihood(s) {
			return true
		}
	}
	return false
}
