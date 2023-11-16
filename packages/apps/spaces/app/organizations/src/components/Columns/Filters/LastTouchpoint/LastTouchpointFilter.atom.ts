import subDays from 'date-fns/subDays';
import { atom, selector, useRecoilState } from 'recoil';

interface LastTouchpointState {
  after?: string;
  value: string[];
  isActive: boolean;
}

const defaultValue = subDays(new Date(), 7).toISOString().split('T')[0];

export const defaultState: LastTouchpointState = {
  value: [],
  after: defaultValue,
  isActive: false,
};

export const LastTouchpointAtom = atom<LastTouchpointState>({
  key: 'last-touchpoint-filter',
  default: defaultState,
});

export const LastTouchpointSelector = selector({
  key: 'last-touchpoint-filter-selector',
  get: ({ get }) => get(LastTouchpointAtom),
});

export const useLastTouchpointFilter = () => {
  return useRecoilState(LastTouchpointAtom);
};
