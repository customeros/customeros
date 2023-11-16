import { atom, selector, useRecoilState } from 'recoil';

interface OwnerFilterState {
  value: string[];
  isActive: boolean;
  showEmpty: boolean;
}

export const defaultState: OwnerFilterState = {
  value: [],
  isActive: false,
  showEmpty: false,
};

export const OwnerFilterAtom = atom<OwnerFilterState>({
  key: 'owner-filter',
  default: defaultState,
});

export const OwnerFilterSelector = selector({
  key: 'owner-filter-selector',
  get: ({ get }) => {
    const state = get(OwnerFilterAtom);

    return {
      value: state.value.filter((v) => v !== '__EMPTY__'),
      isActive: state.isActive,
      showEmpty: state.value.includes('__EMPTY__'),
    };
  },
});

export const useOwnerFilter = () => {
  return useRecoilState(OwnerFilterAtom);
};
