import { atom } from 'recoil';

export const finderContactsSearchTerms = atom({
  key: 'finderContactsSearchTerms', // unique ID (with respect to other atoms/selectors)
  default: [],
});

export const finderOrganizationsSearchTerms = atom({
  key: 'finderOrganizationsSearchTerms', // unique ID (with respect to other atoms/selectors)
  default: [],
});
