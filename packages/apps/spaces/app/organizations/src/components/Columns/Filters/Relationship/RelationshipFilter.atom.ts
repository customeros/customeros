import { atom, selector, useRecoilState } from 'recoil';

interface RelationshipFilterState {
  value: boolean[];
  isActive: boolean;
}

export const RelationshipFilterAtom = atom<RelationshipFilterState>({
  key: 'relationship-filter',
  default: {
    value: [true, false],
    isActive: false,
  },
});

export const RelationshipFilterSelector = selector({
  key: 'relationship-filter-selector',
  get: ({ get }) => get(RelationshipFilterAtom),
});

export const useRelationshipFilter = () => {
  return useRecoilState(RelationshipFilterAtom);
};
