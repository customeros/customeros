import { Maybe, OpportunityRenewalLikelihood } from '@graphql/types';

export function getLikelihoodColor(
  likelihood: OpportunityRenewalLikelihood | null = null,
) {
  switch (likelihood) {
    case OpportunityRenewalLikelihood.HighRenewal:
      return 'text-greenLight-500';
    case OpportunityRenewalLikelihood.MediumRenewal:
      return 'text-yellow-500';
    case OpportunityRenewalLikelihood.LowRenewal:
      return 'text-orangeDark-800';
    case OpportunityRenewalLikelihood.ZeroRenewal:
      return 'text-gray-500';
    default:
      return 'text-gray-500';
  }
}

export function getRenewalLikelihoodLabel(
  renewalLikelihood?: Maybe<OpportunityRenewalLikelihood> | undefined,
) {
  switch (renewalLikelihood) {
    case OpportunityRenewalLikelihood.HighRenewal:
      return 'High';
    case OpportunityRenewalLikelihood.MediumRenewal:
      return 'Medium';
    case OpportunityRenewalLikelihood.LowRenewal:
      return 'Low';
    case OpportunityRenewalLikelihood.ZeroRenewal:
      return 'Zero';
    default:
      return '';
  }
}
