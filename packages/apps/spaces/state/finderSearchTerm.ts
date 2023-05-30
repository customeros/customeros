import { atom } from 'recoil';

export const finderContactsGridData = atom<any>({
  key: 'finderContactsSearchTerms', // unique ID (with respect to other atoms/selectors)
  default: {
    searchTerms: [],
    sortBy: {
      column: undefined,
      direction: undefined,
    },
  },
});

export const finderOrganizationsSearchTerms = atom<any[]>({
  key: 'finderOrganizationsSearchTerms', // unique ID (with respect to other atoms/selectors)
  default: [],
});
