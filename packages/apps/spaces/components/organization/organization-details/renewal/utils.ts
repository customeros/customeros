import capitalize from 'lodash/capitalize';

import { RenewalLikelihoodProbability } from '@graphql/types';

export function isLikelihoodIncreased(
  curr: RenewalLikelihoodProbability | null = null,
  prev: RenewalLikelihoodProbability | null = null,
) {
  if (!curr) return false;
  if (!prev) return true;

  if (curr === RenewalLikelihoodProbability.High) {
    return true;
  }

  if (
    curr === RenewalLikelihoodProbability.Medium &&
    [
      RenewalLikelihoodProbability.Medium,
      RenewalLikelihoodProbability.Low,
      RenewalLikelihoodProbability.Zero,
    ].includes(prev)
  ) {
    return true;
  }

  if (
    curr === RenewalLikelihoodProbability.Low &&
    [
      RenewalLikelihoodProbability.Low,
      RenewalLikelihoodProbability.Zero,
    ].includes(prev)
  ) {
    return true;
  }

  return false;
}

export function getLikelihoodColor(
  likelihood: RenewalLikelihoodProbability | null = null,
) {
  switch (likelihood) {
    case RenewalLikelihoodProbability.High:
      return 'success.500';
    case RenewalLikelihoodProbability.Medium:
      return 'warning.500';
    case RenewalLikelihoodProbability.Low:
      return 'error.500';
    case RenewalLikelihoodProbability.Zero:
      return 'gray.500';
    default:
      return 'gray.500';
  }
}

export const renewalLikelihoodOptions = [
  {
    label: capitalize(RenewalLikelihoodProbability.High),
    value: RenewalLikelihoodProbability.High,
  },
  {
    label: capitalize(RenewalLikelihoodProbability.Medium),
    value: RenewalLikelihoodProbability.Medium,
  },
  {
    label: capitalize(RenewalLikelihoodProbability.Low),
    value: RenewalLikelihoodProbability.Low,
  },
  {
    label: capitalize(RenewalLikelihoodProbability.Zero),
    value: RenewalLikelihoodProbability.Zero,
  },
];
