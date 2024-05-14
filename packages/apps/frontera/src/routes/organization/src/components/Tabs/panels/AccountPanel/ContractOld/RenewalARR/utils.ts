import { OpportunityRenewalLikelihood } from '@graphql/types';

export const getButtonStyles = (
  likelihood: OpportunityRenewalLikelihood | null | undefined,
  variant: OpportunityRenewalLikelihood,
) => ({
  w: 'full',
  variant: 'outline',

  bg: likelihood === variant ? 'white' : 'gray.50',

  '&:hover': {
    bg: likelihood === variant ? 'white' : 'gray.50',
    color: 'gray.500',
  },
});

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
