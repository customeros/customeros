import { atom, selector, useRecoilState } from 'recoil';

interface RelationshipFilterState {
  value: string[];
  isActive: boolean;
}

export const RelationshipFilterAtom = atom<RelationshipFilterState>({
  key: 'relationship-filter',
  default: {
    value: ['customer', 'prospect'],
    isActive: false,
  },
});

export const getRelationshipFilterValue = (state: RelationshipFilterState) => {
  const value = state.value.map((item) => (item === 'customer' ? true : false));

  return {
    value,
    isActive: state.isActive,
  };
};

export const RelationshipFilterSelector = selector({
  key: 'relationship-filter-selector',
  get: ({ get }) => {
    const state = get(RelationshipFilterAtom);

    return getRelationshipFilterValue(state);
  },
});

export const useRelationshipFilter = () => {
  return useRecoilState(RelationshipFilterAtom);
};
