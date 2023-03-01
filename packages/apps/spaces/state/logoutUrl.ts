import { atom } from 'recoil';

export const logoutUrlState = atom({
  key: 'logoutUrlState', // unique ID (with respect to other atoms/selectors)
  default: '#', // default value (aka initial value)
});
