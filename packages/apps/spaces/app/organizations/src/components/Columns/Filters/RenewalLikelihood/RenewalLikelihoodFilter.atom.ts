import { atom, selector, useRecoilState } from 'recoil';

import { OpportunityRenewalLikelihood } from '@graphql/types';

interface RenewalLiklihoodFilterState {
  isActive: boolean;
  value: OpportunityRenewalLikelihood[];
}

export const defaultState: RenewalLiklihoodFilterState = {
  value: [
    OpportunityRenewalLikelihood.HighRenewal,
    OpportunityRenewalLikelihood.MediumRenewal,
    OpportunityRenewalLikelihood.LowRenewal,
    OpportunityRenewalLikelihood.ZeroRenewal,
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
