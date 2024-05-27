import { atom, selector, useRecoilState } from 'recoil';

import { OrganizationRelationship } from '@graphql/types';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface RelationshipFilterState {
  isActive: boolean;
  value: OrganizationRelationship[];
}

export const defaultState: RelationshipFilterState = {
  value: [OrganizationRelationship.Customer],
  isActive: false,
};

export const RelationshipFilterAtom = atom<RelationshipFilterState>({
  key: 'relationship-filter',
  default: defaultState,
});

export const RelationshipFilterSelector = selector({
  key: 'relationship-filter-selector',
  get: ({ get }) => get(RelationshipFilterAtom),
});

export const useRelationshipFilter = () => {
  return useRecoilState(RelationshipFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom RelationshipFilterState
 */
export const mapRelationshipToAtom =
  makeServerToAtomMapper<RelationshipFilterState>(
    {
      filter: {
        property: 'RELATIONSHIP',
      },
    },
    ({ filter }) => ({
      isActive: true,
      value: filter?.value,
    }),
    defaultState,
  );
