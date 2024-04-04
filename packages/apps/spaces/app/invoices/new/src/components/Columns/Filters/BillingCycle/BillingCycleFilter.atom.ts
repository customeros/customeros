import { atom, selector, useRecoilState } from 'recoil';

import { ContractBillingCycle } from '@graphql/types';
import { makeServerToAtomMapper } from '@shared/components/Filters/makeServerToAtomMapper';

export interface BillingCycleFilterState {
  isActive: boolean;
  value: ContractBillingCycle[];
}

export const defaultState: BillingCycleFilterState = {
  value: [],
  isActive: false,
};

export const BillingCycleFilterAtom = atom<BillingCycleFilterState>({
  key: 'billing-cycle-filter',
  default: defaultState,
});

export const BillingCycleFilterSelector = selector({
  key: 'billing-cycle-filter-selector',
  get: ({ get }) => get(BillingCycleFilterAtom),
});

export const useBillingCycleFilter = () => {
  return useRecoilState(BillingCycleFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom BillingCycleFilterState
 */
export const mapBillingCycleToAtom =
  makeServerToAtomMapper<BillingCycleFilterState>(
    {
      filter: {
        property: 'CONTRACT_BILLING_CYCLE',
      },
    },
    ({ filter }) => ({
      isActive: true,
      value: filter?.value,
    }),
    defaultState,
  );
