import { atom, selector, useRecoilState } from 'recoil';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface ForecastFilterState {
  isActive: boolean;
  value: [number, number];
}

export const defaultState: ForecastFilterState = {
  isActive: false,
  value: [0, 10000],
};

export const ForecastFilterAtom = atom<ForecastFilterState>({
  key: 'forecast-filter',
  default: defaultState,
});

export const ForecastFilterSelector = selector({
  key: 'forecast-filter-selector',
  get: ({ get }) => get(ForecastFilterAtom),
});

export const useForecastFilter = () => {
  return useRecoilState(ForecastFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom ForecastFilterState
 */
export const mapForecastToAtom = makeServerToAtomMapper<ForecastFilterState>(
  {
    filter: { property: 'FORECAST_ARR' },
  },
  ({ filter }) => ({
    isActive: true,
    value: filter?.value,
  }),
  defaultState,
);
