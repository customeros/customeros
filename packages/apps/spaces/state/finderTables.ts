import { atom } from 'recoil';
import { SortingDirection } from '@spaces/graphql';

export interface FinderContactTableSortingState {
  column: 'CONTACT' | 'EMAIL' | 'ORGANIZATION' | 'LOCATION' | undefined;
  direction: SortingDirection | undefined;
}

export const finderContactTableSortingState =
  atom<FinderContactTableSortingState>({
    key: 'finderContactTable', // unique ID (with respect to other atoms/selectors)
    default: {
      column: undefined,
      direction: undefined,
    }, // default value (aka initial value)
  });
