import { atom } from 'recoil';
import { SortingDirection } from '../graphQL/__generated__/generated';

interface FinderOrganizationTableSortingState {
  column: 'NAME' | undefined;
  direction: SortingDirection;
}

export const finderOrganizationTableSortingState =
  atom<FinderOrganizationTableSortingState>({
    key: 'finderOrganizationTable', // unique ID (with respect to other atoms/selectors)
    default: {
      column: 'NAME',
      direction: SortingDirection.Asc,
    }, // default value (aka initial value)
  });

export const finderContactTableSortingState =
  atom<FinderOrganizationTableSortingState>({
    key: 'finderContactTable', // unique ID (with respect to other atoms/selectors)
    default: {
      column: 'NAME',
      direction: SortingDirection.Asc,
    }, // default value (aka initial value)
  });
