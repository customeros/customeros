import { atom } from 'recoil';

export const finderEverythingTable = atom({
  key: 'finderEverythingTable', // unique ID (with respect to other atoms/selectors)
  default: '', // default value (aka initial value)
});
