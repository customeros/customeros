import { match } from 'ts-pattern';
import { Store } from '@store/store';

import {
  FilterItem,
  Organization,
} from '@shared/types/__generated__/graphql.types';

export const getFlowFilters = (filter: FilterItem | undefined | null) => {
  const noop = (_row: Store<Organization>) => true;
  if (!filter) return noop;

  return match(filter)
    .with({ property: 'industry' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;

      if (!filterValues) return false;

      return filterValues.includes(row.value?.industry);
    })
    .with({ property: 'employees' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;

      const employees = row.value?.employees;

      if (!filterValues) return false;

      return employees >= filterValues[0] && employees <= filterValues[1];
    })
    .with({ property: 'tags' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;

      if (!filterValues) return false;

      return filterValues.includes(row.value?.tags);
    })
    .with({ property: 'age' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;

      const age = row.value?.yearFounded;
      if (!filterValues) return false;

      if (filterValues.length === 2) {
        return age >= filterValues[0] && age <= filterValues[1];
      } else {
        return filterValues.includes(age);
      }
    })
    .otherwise(() => noop);

  //   .with({property:'followers'}, (filter) => (row: Store<Organization>) => {
  //     const filterValues = filter?.value;

  //     if (!filterValues) return false;

  //     return filterValues.includes(row.value?.);
  //   }
  // )

  // .with(
  //   { property: 'ownership' },
  //   (filter) => (row: Store<Organization>) => {
  //     const filterValues = filter?.value;

  //     if (!filterValues) return false;

  //     return filterValues.includes(row.value?.ownership);
  //   },
  // )
};
