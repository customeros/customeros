import { atom, selector, useRecoilState } from 'recoil';

import { OpportunityRenewalLikelihood } from '@graphql/types';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface RenewalLiklihoodFilterState {
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
  key: 'renewals-renewal-likelihood-filter',
  default: defaultState,
});

export const RenewalLikelihoodFilterSelector = selector({
  key: 'renewals-renewal-likelihood-filter-selector',
  get: ({ get }) => get(RenewalLikelihoodFilterAtom),
});

export const useRenewalLikelihoodFilter = () => {
  return useRecoilState(RenewalLikelihoodFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom RenewalLiklihoodFilterState
 */
export const mapRenewalLikelihoodToAtom =
  makeServerToAtomMapper<RenewalLiklihoodFilterState>(
    {
      filter: {
        property: 'RENEWAL_LIKELIHOOD',
      },
    },
    ({ filter }) => ({
      isActive: true,
      value: filter?.value,
    }),
    defaultState,
  );
