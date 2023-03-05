import { atom } from 'recoil';

export const userData = atom({
  key: 'userData', // unique ID (with respect to other atoms/selectors)
  default: {
    identity: '',
  }, // default value (aka initial value)
});
