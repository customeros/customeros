import { atom } from 'recoil';

export const finderContactsSearchTerms = atom<any[]>({
  key: 'finderContactsSearchTerms', // unique ID (with respect to other atoms/selectors)
  default: [],
});

export const finderOrganizationsSearchTerms = atom<any[]>({
  key: 'finderOrganizationsSearchTerms', // unique ID (with respect to other atoms/selectors)
  default: [],
});
