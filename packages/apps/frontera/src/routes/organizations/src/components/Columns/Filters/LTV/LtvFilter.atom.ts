import { atom, selector, useRecoilState } from 'recoil';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface ForecastFilterState {
  isActive: boolean;
  value: [number, number];
}

export const defaultState: ForecastFilterState = {
  isActive: false,
  value: [0, 1000],
};

export const LtvFilterAtom = atom<ForecastFilterState>({
  key: 'ltv',
  default: defaultState,
});

export const LtvFilterSelector = selector({
  key: 'ltv-filter-selector',
  get: ({ get }) => get(LtvFilterAtom),
});

export const useForecastFilter = () => {
  return useRecoilState(LtvFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom ForecastFilterState
 */
export const mapForecastToAtom = makeServerToAtomMapper<ForecastFilterState>(
  {
    filter: { property: 'LTV' },
  },
  ({ filter }) => ({
    isActive: true,
    value: filter?.value,
  }),
  defaultState,
);
