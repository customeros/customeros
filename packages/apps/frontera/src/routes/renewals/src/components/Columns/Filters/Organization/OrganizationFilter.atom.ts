import { atom, selector, useRecoilState } from 'recoil';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface OrganizationFilterState {
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
  key: 'renewals-organization-filter',
  default: defaultState,
});

export const OrganizationFilterSelector = selector({
  key: 'renewals-organization-filter-selector',
  get: ({ get }) => get(OrganizationFilterAtom),
});

export const useOrganizationFilter = () => {
  return useRecoilState(OrganizationFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom OrganizationFilterState
 */
export const mapOrganizationToAtom =
  makeServerToAtomMapper<OrganizationFilterState>(
    {
      filter: {
        property: 'NAME',
      },
    },
    ({ filter }) => ({
      isActive: true,
      value: filter?.value,
      showEmpty: filter?.includeEmpty ?? false,
    }),
    defaultState,
  );
