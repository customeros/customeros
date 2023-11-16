import { atom, selector, useRecoilState } from 'recoil';

interface RelationshipFilterState {
  value: boolean[];
  isActive: boolean;
}

export const defaultState: RelationshipFilterState = {
  value: [true, false],
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
