import { atom, selector, useRecoilState } from 'recoil';

import { RenewalLikelihoodProbability } from '@graphql/types';

interface RenewalLiklihoodFilterState {
  isActive: boolean;
  value: RenewalLikelihoodProbability[];
}

export const defaultState: RenewalLiklihoodFilterState = {
  value: [
    RenewalLikelihoodProbability.High,
    RenewalLikelihoodProbability.Medium,
    RenewalLikelihoodProbability.Low,
    RenewalLikelihoodProbability.Zero,
  ],
  isActive: false,
};

export const RenewalLikelihoodFilterAtom = atom<RenewalLiklihoodFilterState>({
  key: 'renewal-likelihood-filter',
  default: defaultState,
});

export const RenewalLikelihoodFilterSelector = selector({
  key: 'renewal-likelihood-filter-selector',
  get: ({ get }) => get(RenewalLikelihoodFilterAtom),
});

export const useRenewalLikelihoodFilter = () => {
  return useRecoilState(RenewalLikelihoodFilterAtom);
};
