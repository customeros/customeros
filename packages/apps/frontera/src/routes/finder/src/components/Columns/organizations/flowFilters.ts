import { match } from 'ts-pattern';
import { Organization } from '@store/Organizations/Organization.dto';

import {
  Filter,
  FilterItem,
  ColumnViewType,
  ComparisonOperator,
} from '@graphql/types';

export const getFlowFilters = (filter: FilterItem | undefined | null) => {
  const noop = (_row: Organization) => true;

  if (!filter) return noop;

  return match(filter)
    .with(
      { property: ColumnViewType.OrganizationsIndustry },
      (filter) => (row: Organization) => {
        const filterValues = filter?.value;

        if (!filterValues) return false;

        return filterValues.includes(row.value?.industryName);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsIsPublic },
      (filter) => (row: Organization) => {
        const filterValues = filter?.value;

        return row.value?.public === filterValues;
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsEmployeeCount },
      (filter) => (row: Organization) => {
        const filterValues = filter?.value;
        const filterType = filter?.operation;
        const employees = row.value?.employees;

        if (!filterValues) return false;

        if (
          filterValues.length === 1 &&
          !!row.value?.employees &&
          filterType === ComparisonOperator.Lt
        ) {
          return employees < filterValues[0];
        } else if (
          filterValues.length === 1 &&
          !!row.value?.employees &&
          filterType === ComparisonOperator.Gt
        ) {
          return employees > filterValues[0];
        } else {
          return employees > filterValues[0] && employees < filterValues[1];
        }
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsTags },
      (filter) => (row: Organization) => {
        const filterValues = filter?.value;

        if (!filterValues) return false;

        return filterValues.every((value: string) =>
          row.value.tags?.some((obj) => obj.metadata.id === value),
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsYearFounded },
      (filter) => (row: Organization) => {
        const filterValues = filter?.value;
        const filterType = filter?.operation;

        const age = row.value?.yearFounded;

        if (!filterValues || row.value.yearFounded === null) return false;

        if (typeof filterValues === 'object' && filterValues.length === 0)
          return false;

        if (typeof filterValues === 'object' && filterValues.length > 1)
          return age > filterValues[0] && age < filterValues[1];

        if (
          filterValues &&
          !!row.value?.yearFounded &&
          filterType === ComparisonOperator.Lt
        ) {
          return age < filterValues;
        } else if (filterValues && !!row.value?.yearFounded) {
          return age > filterValues;
        }
      },
    )

    .with(
      { property: ColumnViewType.OrganizationsLinkedinFollowerCount },
      (filter) => (row: Organization) => {
        const filterValues = filter?.value;
        const filterType = filter?.operation;
        const followersCount = row.value?.socialMedia.find((s) =>
          s.url.includes('linkedin'),
        )?.followersCount;

        if (!filterValues) return false;

        if (filterValues.length === 1 && filterType === ComparisonOperator.Lt) {
          return followersCount < filterValues[0];
        } else if (
          filterValues.length === 1 &&
          filterType === ComparisonOperator.Gt
        ) {
          return followersCount > filterValues[0];
        } else {
          return (
            followersCount > filterValues[0] && followersCount < filterValues[1]
          );
        }
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsHeadquarters },
      (filter) => (row: Organization) => {
        const filterValues = filter?.value;

        if (!filterValues) return false;

        return filterValues.includes(row.value.locations?.[0]?.countryCodeA2);
      },
    )

    .otherwise(() => noop);
};

export const getFlowFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getFlowFilters(filter));
};
