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

export const callParticipant = atom({
  key: 'callParticipant', // unique ID (with respect to other atoms/selectors)
  default: {
    identity: '',
  }, // default value (aka initial value)
});

interface UserSettingsInterface {
  googleOAuth: boolean,
  microsoftOAuth: boolean,
  userRoles: string[],
}

export const userSettings = atom<UserSettingsInterface>({
  key: 'userSettings', // unique ID (with respect to other atoms/selectors)
  default: {
    googleOAuth: false,
    microsoftOAuth: false,
    userRoles: ['basicUser'],
  }, // default value (aka initial value)
});
