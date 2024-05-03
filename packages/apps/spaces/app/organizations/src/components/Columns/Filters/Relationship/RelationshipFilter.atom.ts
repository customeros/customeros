import { atom, selector, useRecoilState } from 'recoil';

import { makeServerToAtomMapper } from '../shared/makeServerToAtomMapper';

export interface RelationshipFilterState {
  value: boolean[];
  isActive: boolean;
}

export const defaultState: RelationshipFilterState = {
  value: [true],
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
        property: 'IS_CUSTOMER',
      },
    },
    ({ filter }) => ({
      isActive: true,
      value: filter?.value,
    }),
    defaultState,
  );
