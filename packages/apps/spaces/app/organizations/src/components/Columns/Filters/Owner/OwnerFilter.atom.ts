import { atom, selector, useRecoilState } from 'recoil';

interface OwnerFilterState {
  value: string[];
  isActive: boolean;
}

export const OwnerFilterAtom = atom<OwnerFilterState>({
  key: 'owner-filter',
  default: {
    value: [],
    isActive: false,
  },
});

export const OwnerFilterSelector = selector({
  key: 'owner-filter-selector',
  get: ({ get }) => get(OwnerFilterAtom),
});

export const useOwnerFilter = () => {
  return useRecoilState(OwnerFilterAtom);
};
