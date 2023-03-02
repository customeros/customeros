import { atom } from 'recoil';

export const finderSearchTerm = atom({
  key: 'finderSearchTerm', // unique ID (with respect to other atoms/selectors)
  default: '', // default value (aka initial value)
});
