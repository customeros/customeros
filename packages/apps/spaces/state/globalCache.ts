import { atom } from 'recoil';

export const globalCacheData = atom({
  key: 'globalCacheData',
  default: {
    user: {
      id: undefined,
      firstName: undefined,
      lastName: undefined,
      emails: Array<{ email: string; rawEmail: string; primary: boolean }>,
    } as any | undefined,
    isOwner: false,
    gmailOauthTokenNeedsManualRefresh: false,
    gCliCache: [],
  },
});
