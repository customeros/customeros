import subDays from 'date-fns/subDays';
import { atom, selector, useRecoilState } from 'recoil';

import { touchpoints } from './util';

interface LastTouchpointState {
  after?: string;
  value: string[];
  isActive: boolean;
}

const defaultAfter = subDays(new Date(), 7).toISOString().split('T')[0];

export const defaultState: LastTouchpointState = {
  value: touchpoints.map((t) => t.value),
  after: defaultAfter,
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
