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
    .with({ property: 'ownership' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;

      return row.value?.public === filterValues;
    })
    .with({ property: 'employees' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;
      const filterType = filter?.operation;

      const employees = row.value?.employees;

      if (!filterValues) return false;

      if (
        filterValues.length === 1 &&
        !!row.value?.employees &&
        filterType === 'LTE'
      ) {
        return employees <= filterValues[0];
      } else if (
        filterValues.length === 1 &&
        !!row.value?.employees &&
        filterType === 'GTE'
      ) {
        return employees >= filterValues[0];
      } else {
        return employees >= filterValues[0] && employees <= filterValues[1];
      }
    })
    .with({ property: 'tags' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;
      if (!filterValues) return false;

      return filterValues.every((value: string) =>
        row.value.tags?.some((obj) => obj.id === value),
      );
    })
    .with({ property: 'age' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;
      const filterType = filter?.operation;

      const age = new Date().getFullYear() - row.value?.yearFounded;

      if (!filterValues) return false;

      if (
        filterValues.length === 1 &&
        !!row.value?.yearFounded &&
        filterType === 'LTE'
      ) {
        return age <= filterValues[0];
      } else if (
        filterValues.length === 1 &&
        !!row.value?.yearFounded &&
        filterType === 'GTE'
      ) {
        return age >= filterValues[0];
      } else {
        return age >= filterValues[0] && age <= filterValues[1];
      }
    })

    .with({ property: 'followers' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;
      const filterType = filter?.operation;
      const followersCount = row.value?.socialMedia.find((s) =>
        s.url.includes('linkedin'),
      )?.followersCount;

      if (!filterValues) return false;

      if (filterValues.length === 1 && filterType === 'LTE') {
        return followersCount <= filterValues[0];
      } else if (filterValues.length === 1 && filterType === 'GTE') {
        return followersCount >= filterValues[0];
      } else {
        return (
          followersCount >= filterValues[0] && followersCount <= filterValues[1]
        );
      }
    })
    .with(
      { property: 'headquarters' },
      (filter) => (row: Store<Organization>) => {
        const filterValues = filter?.value;
        if (!filterValues) return false;

        return filterValues.includes(row.value.locations?.[0]?.countryCodeA2);
      },
    )

    .otherwise(() => noop);
};
