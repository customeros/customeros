import { atom } from 'recoil';
import { FinderContactTableSortingState } from './finderTables';

export const globalCacheData = atom({
  key: 'globalCacheData',
  default: {
    user: {
      id: undefined,
      firstName: undefined,
      lastName: undefined,
      emails: Array<{ email: string; rawEmail: string; primary: boolean }>,
    } as any | undefined,
    isOwner: undefined,
    gCliCache: undefined,
  },
});
