import addDays from 'date-fns/addDays';
import { atom, selector, useRecoilState } from 'recoil';

interface LastTouchpointState {
  value: string[];
  before?: string;
  isActive: boolean;
}

const defaultValue = addDays(new Date(), 7).toISOString().split('T')[0];

export const LastTouchpointAtom = atom<LastTouchpointState>({
  key: 'last-touchpoint-filter',
  default: {
    value: [],
    before: defaultValue,
    isActive: false,
  },
});

export const LastTouchpointSelector = selector({
  key: 'last-touchpoint-filter-selector',
  get: ({ get }) => get(LastTouchpointAtom),
});

export const useLastTouchpointFilter = () => {
  return useRecoilState(LastTouchpointAtom);
};
