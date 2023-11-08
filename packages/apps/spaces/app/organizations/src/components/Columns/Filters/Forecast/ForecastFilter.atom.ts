import { atom, selector, useRecoilState } from 'recoil';

interface ForecastFilterState {
  isActive: boolean;
  value: [number, number];
}

export const ForecastFilterAtom = atom<ForecastFilterState>({
  key: 'forecast-filter',
  default: {
    value: [0, 10000],
    isActive: false,
  },
});

export const ForecastFilterSelector = selector({
  key: 'forecast-filter-selector',
  get: ({ get }) => get(ForecastFilterAtom),
});

export const useForecastFilter = () => {
  return useRecoilState(ForecastFilterAtom);
};
