import { atom } from 'recoil';
import {UserSettingsInterface} from "../services/settings/settingsService";

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

export const userSettings = atom<UserSettingsInterface>({
  key: 'userSettings', // unique ID (with respect to other atoms/selectors)
  default: {
    id: '',
    tenantName: '',
    username: '',
    googleOAuthAllScopesEnabled: false,
    googleOAuthUserAccessToken: '',
  }, // default value (aka initial value)
});
