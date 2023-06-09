import { atom } from 'recoil';
import { FinderContactTableSortingState } from './finderTables';

export const globalCacheData = atom({
  key: 'globalCacheData',
  default: {
    userId: '',
    userEmail: '',
    isOwner: false,
    gCliCache: [],
  },
});
