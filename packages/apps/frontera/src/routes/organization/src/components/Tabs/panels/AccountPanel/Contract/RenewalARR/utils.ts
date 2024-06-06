import { OpportunityRenewalLikelihood } from '@graphql/types';

export const likelihoodButtons = [
  {
    label: 'Low',
    colorScheme: 'orangeDark',
    likelihood: OpportunityRenewalLikelihood.LowRenewal,
  },
  {
    label: 'Medium',
    colorScheme: 'yellow',
    likelihood: OpportunityRenewalLikelihood.MediumRenewal,
  },
  {
    label: 'High',
    colorScheme: 'greenLight',
    likelihood: OpportunityRenewalLikelihood.HighRenewal,
  },
];
