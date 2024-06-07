import { subDays } from 'date-fns/subDays';
import { atom, selector, useRecoilState } from 'recoil';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface ChurnedState {
  value: string;
  isActive: boolean;
}

const defaultValue = subDays(new Date(), 30).toISOString().split('T')[0];

export const defaultState: ChurnedState = {
  value: defaultValue,
  isActive: false,
};

export const ChurnedAtom = atom<ChurnedState>({
  key: 'churned-filter',
  default: defaultState,
});

export const ChurnedFilterSelector = selector({
  key: 'churned-filter-selector',
  get: ({ get }) => get(ChurnedAtom),
});

export const useChurnedFilter = () => {
  return useRecoilState(ChurnedAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom ChurnedState
 */
export const mapChurnedToAtom = makeServerToAtomMapper<ChurnedState>(
  {
    filter: {
      property: 'CHURNED_AT',
    },
  },
  ({ filter }) => ({
    isActive: true,
    value: filter?.value,
  }),
  defaultState,
);
