import { atom, selector, useRecoilState } from 'recoil';

interface OrganizationFilterState {
  value: string;
  isActive: boolean;
  showEmpty: boolean;
}

export const defaultState: OrganizationFilterState = {
  value: '',
  isActive: false,
  showEmpty: false,
};

export const OrganizationFilterAtom = atom<OrganizationFilterState>({
  key: 'organization-filter',
  default: defaultState,
});

export const OrganizationFilterSelector = selector({
  key: 'organization-filter-selector',
  get: ({ get }) => get(OrganizationFilterAtom),
});

export const useOrganizationFilter = () => {
  return useRecoilState(OrganizationFilterAtom);
};
