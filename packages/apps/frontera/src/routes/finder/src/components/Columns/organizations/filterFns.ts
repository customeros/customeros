import { match } from 'ts-pattern';
import { isBefore } from 'date-fns';
import { FilterItem } from '@store/types';
import { isAfter } from 'date-fns/isAfter';
import { OrganizationStore } from '@store/Organizations/Organization.store';

import {
  Tag,
  Filter,
  Social,
  ColumnViewType,
  OnboardingStatus,
  ComparisonOperator,
  OrganizationRelationship,
} from '@graphql/types';

const getFilterFn = (filter: FilterItem | undefined | null) => {
  const noop = (_row: OrganizationStore) => true;

  if (!filter) return noop;

  return match(filter)
    .with({ property: 'STAGE' }, (filter) => (row: OrganizationStore) => {
      const filterValues = filter?.value;

      if (!filterValues) return false;

      return filterValues.includes(row.value?.stage);
    })
    .with({ property: 'IS_CUSTOMER' }, (filter) => (row: OrganizationStore) => {
      const filterValues = filter?.value;

      if (!filterValues) return false;

      return filterValues.includes(
        row.value?.relationship === OrganizationRelationship.Customer,
      );
    })
    .with({ property: 'OWNER_ID' }, (filter) => (row: OrganizationStore) => {
      const filterValues = filter?.value;

      if (!filterValues) return false;

      return filterValues.includes(row.value?.owner?.id);
    })

    .with(
      { property: 'RELATIONSHIP' },
      (filter) => (row: OrganizationStore) => {
        const filterValues = filter?.value;

        if (!filterValues) return false;

        return filterValues.includes(row.value?.relationship);
      },
    )

    .with(
      { property: ColumnViewType.OrganizationsCreatedDate },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const value = row.value.metadata.created;

        return filterTypeDate(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsName },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const values = row.value.name?.toLowerCase();

        return filterTypeText(filter, values);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsParentOrganization },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const values =
          row.value.parentCompanies?.[0]?.organization.name?.toLowerCase();

        return filterTypeText(filter, values);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsWebsite },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const value = row.value.website || '';

        return filterTypeText(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRelationship },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const values = row.value.relationship;

        if (!values)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(
          filter,
          Array.isArray(values) ? values.map(String) : [String(values)],
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsStage },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const values = row.value.stage;

        if (!values)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(
          filter,
          Array.isArray(values) ? values : [values],
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsForecastArr },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const forecastValue =
          row.value?.accountDetails?.renewalSummary?.arrForecast;

        if (!forecastValue) return false;

        return filterTypeNumber(filter, forecastValue);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRenewalDate },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const nextRenewalDate =
          row.value?.accountDetails?.renewalSummary?.nextRenewalDate?.split(
            'T',
          )[0];

        return filterTypeDate(filter, nextRenewalDate);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsOnboardingStatus },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const values = row.value.accountDetails?.onboarding?.status;

        if (!values) return false;
        if (values === OnboardingStatus.NotApplicable)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(
          filter,
          Array.isArray(values) ? values : [values],
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRenewalLikelihood },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const values =
          row.value.accountDetails?.renewalSummary?.renewalLikelihood;

        if (!values) return filter.operation === ComparisonOperator.IsEmpty;

        return filterTypeList(
          filter,
          Array.isArray(values) ? values : [values],
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsOwner },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const values = row.value.owner?.id;

        if (!values)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(
          filter,
          Array.isArray(values) ? values : [values],
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLastTouchpoint },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const lastTouchpoint = row?.value?.lastTouchpoint?.lastTouchPointType;

        if (!lastTouchpoint)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(
          filter,
          Array.isArray(lastTouchpoint) ? lastTouchpoint : [lastTouchpoint],
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsChurnDate },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const churned = row?.value?.accountDetails?.churned;

        if (!churned) return false;

        return filterTypeDate(filter, churned);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsSocials },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const linkedInUrl = row.value.socialMedia?.find((v) =>
          v.url.includes('linkedin'),
        )?.url;

        return filterTypeText(filter, linkedInUrl);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLastTouchpointDate },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const lastTouchpointAt = row?.value?.lastTouchpoint?.lastTouchPointAt;

        return filterTypeDate(filter, lastTouchpointAt);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsEmployeeCount },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const employees = row.value?.employees;

        return filterTypeNumber(filter, employees);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsContactCount },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const contactsCount = row.value.contacts.content.length;

        return filterTypeNumber(filter, contactsCount);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLinkedinFollowerCount },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const followers = row.value.socialMedia.find((e: Social) =>
          e?.url?.includes('linkedin'),
        )?.followersCount;

        return filterTypeNumber(filter, followers);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLeadSource },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const value = row.value.leadSource;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, Array.isArray(value) ? value : [value]);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsIndustry },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const value = row.value.industry;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, Array.isArray(value) ? value : [value]);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLtv },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const ltv = row.value.accountDetails?.ltv;

        if (!ltv) return false;

        return filterTypeNumber(filter, ltv);
      },
    )
    .with({ property: ColumnViewType.OrganizationsHeadquarters }, (filter) => {
      return (row: OrganizationStore) => {
        if (!filter.active) return true;
        const locations = row.value.locations;
        const country = locations?.[0]?.countryCodeA2;

        if (!country)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(
          filter,
          Array.isArray(country) ? country : [country],
        );
      };
    })
    .with(
      { property: ColumnViewType.OrganizationsIsPublic },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const isPublic = row.value.public === true ? 'Public' : 'Private';

        return filterTypeList(
          filter,
          Array.isArray(isPublic) ? isPublic.map(String) : [String(isPublic)],
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsYearFounded },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const yearFounded = row.value.yearFounded;

        if (!yearFounded) return false;

        return filterTypeNumber(filter, yearFounded);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsTags },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const values = row.value.tags?.map((l: Tag) => l.id);

        if (!values)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(
          filter,
          Array.isArray(values) ? values : [values],
        );
      },
    )

    .otherwise(() => noop);
};

const filterTypeText = (filter: FilterItem, value: string | undefined) => {
  const filterValue = filter?.value?.toLowerCase();
  const filterOperator = filter?.operation;
  const valueLower = value?.toLowerCase();

  return match(filterOperator)
    .with(ComparisonOperator.IsEmpty, () => !value)
    .with(ComparisonOperator.IsNotEmpty, () => value)
    .with(
      ComparisonOperator.NotContains,
      () => !valueLower?.includes(filterValue),
    )
    .with(ComparisonOperator.Contains, () => valueLower?.includes(filterValue))
    .otherwise(() => false);
};

const filterTypeNumber = (filter: FilterItem, value: number | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  if (value === undefined || value === null) return false;

  return match(filterOperator)
    .with(ComparisonOperator.Lt, () => value < Number(filterValue))
    .with(ComparisonOperator.Gt, () => value > Number(filterValue))
    .with(ComparisonOperator.Eq, () => value === Number(filterValue))
    .with(ComparisonOperator.NotEqual, () => value !== Number(filterValue))
    .otherwise(() => true);
};

const filterTypeList = (filter: FilterItem, value: string[] | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  return match(filterOperator)
    .with(ComparisonOperator.IsEmpty, () => !value?.length)
    .with(ComparisonOperator.IsNotEmpty, () => value?.length)
    .with(
      ComparisonOperator.NotContains,
      () =>
        !value?.length ||
        (value?.length && !value.some((v) => filterValue?.includes(v))),
    )
    .with(
      ComparisonOperator.Contains,
      () => value?.length && value.some((v) => filterValue?.includes(v)),
    )
    .otherwise(() => false);
};

const filterTypeDate = (filter: FilterItem, value: string | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  if (!value) return false;

  return match(filterOperator)
    .with(ComparisonOperator.Lt, () =>
      isBefore(new Date(value), new Date(filterValue)),
    )
    .with(ComparisonOperator.Gt, () =>
      isAfter(new Date(value), new Date(filterValue)),
    )

    .otherwise(() => true);
};

export const getOrganizationDefaultFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};

export const getOrganizationFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};
