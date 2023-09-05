import { atom } from 'recoil';

export const ownerListData = atom({
  key: 'ownerListData',
  default: {
    ownerList: [] as any[],
  },
});

export const userData = atom({
  key: 'userData', // unique ID (with respect to other atoms/selectors)
  default: {
    identity: '',
    id: '',
  }, // default value (aka initial value)
});
export const tenantName = atom({
  key: 'tenantName', // unique ID (with respect to other atoms/selectors)
  default: '', // default value (aka initial value)
});
