import { atom, selector, useRecoilState } from 'recoil';

interface LastTouchpointState {
  value: string[];
  isActive: boolean;
}

export const LastTouchpointAtom = atom<LastTouchpointState>({
  key: 'last-touchpoint-filter',
  default: {
    value: [],
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
