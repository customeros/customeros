import { atom } from 'recoil';

export const searchTermState = atom({
  key: 'searchTerm',
  default: '',
});
