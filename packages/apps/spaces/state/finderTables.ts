import { atom } from 'recoil';
import { SortingDirection } from '../graphQL/__generated__/generated';

export interface FinderOrganizationTableSortingState {
  column: 'ORGANIZATION' | 'DOMAIN' | 'LOCATION' | undefined;
  direction: SortingDirection | undefined;
}
export interface FinderContactTableSortingState {
  column: 'CONTACT' | 'EMAIL' | 'ORGANIZATION' | 'LOCATION' | undefined;
  direction: SortingDirection | undefined;
}

export const finderOrganizationTableSortingState =
  atom<FinderOrganizationTableSortingState>({
    key: 'finderOrganizationTable', // unique ID (with respect to other atoms/selectors)
    default: {
      column: undefined,
      direction: undefined,
    }, // default value (aka initial value)
  });

export const finderContactTableSortingState =
  atom<FinderContactTableSortingState>({
    key: 'finderContactTable', // unique ID (with respect to other atoms/selectors)
    default: {
      column: undefined,
      direction: undefined,
    }, // default value (aka initial value)
  });
