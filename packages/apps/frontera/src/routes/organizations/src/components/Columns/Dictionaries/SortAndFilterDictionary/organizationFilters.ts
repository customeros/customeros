import { match } from 'ts-pattern';
import { Store } from '@store/store.ts';
import { isAfter } from 'date-fns/isAfter';
import { FilterItem } from '@store/types.ts';

import { Organization, ColumnViewType } from '@graphql/types';

export const getOrganizationFilterFn = (
  filter: FilterItem | undefined | null,
) => {
  const noop = (_row: Store<Organization>) => true;
  if (!filter) return noop;

  return match(filter)
    .with({ property: 'STAGE' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;

      if (!filterValues) return false;

      return filterValues.includes(row.value?.stage);
    })
    .with(
      { property: 'IS_CUSTOMER' },
      (filter) => (row: Store<Organization>) => {
        const filterValues = filter?.value;

        if (!filterValues) return false;

        return filterValues.includes(row.value?.isCustomer);
      },
    )
    .with({ property: 'OWNER_ID' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;

      if (!filterValues) return false;

      return filterValues.includes(row.value?.owner?.id);
    })

    .with(
      { property: 'RELATIONSHIP' },
      (filter) => (row: Store<Organization>) => {
        const filterValues = filter?.value;

        if (!filterValues) return false;

        return filterValues.includes(row.value?.relationship);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsCreatedDate },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return isAfter(
          new Date(row.value.metadata.created),
          new Date(filterValue),
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsName },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (filter.includeEmpty && row.value.name === 'Unnamed') {
          return true;
        }

        return row.value.name.toLowerCase().includes(filterValue.toLowerCase());
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsWebsite },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (filter.includeEmpty && !row.value.website) {
          return true;
        }

        return (
          row.value.website &&
          row.value.website.toLowerCase().includes(filterValue.toLowerCase())
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRelationship },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(row.value.relationship);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsForecastArr },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const forecastValue =
          row.value?.accountDetails?.renewalSummary?.arrForecast;

        if (!forecastValue) return false;

        return (
          forecastValue >= filterValue[0] && forecastValue <= filterValue[1]
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRenewalDate },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const nextRenewalDate =
          row.value?.accountDetails?.renewalSummary?.nextRenewalDate;

        if (!nextRenewalDate) return false;

        return isAfter(new Date(nextRenewalDate), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsOnboardingStatus },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(
          row.value.accountDetails?.onboarding?.status,
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRenewalLikelihood },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(
          row.value.accountDetails?.renewalSummary?.renewalLikelihood,
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsOwner },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (filterValue === '__EMPTY__' && !row.value.owner) {
          return true;
        }

        return filterValue.includes(row.value.owner?.id);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLastTouchpoint },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const lastTouchpoint = row?.value?.lastTouchpoint?.lastTouchPointType;
        const lastTouchpointAt = row?.value?.lastTouchpoint?.lastTouchPointAt;

        const isIncluded = filterValue?.types.length
          ? filterValue?.types?.includes(lastTouchpoint)
          : false;

        const isAfterDate = isAfter(
          new Date(lastTouchpointAt),
          new Date(filterValue?.after),
        );

        return isIncluded && isAfterDate;
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsChurnDate },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const churned = row?.value?.accountDetails?.churned;

        if (!churned) return false;

        return isAfter(new Date(churned), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsSocials },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        // specific logic for linkedin
        const linkedInUrl = row.value.socialMedia?.find((v) =>
          v.url.includes('linkedin'),
        )?.url;

        if (!linkedInUrl && filter.includeEmpty) return true;

        return linkedInUrl && linkedInUrl.includes(filterValue);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLastTouchpointDate },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const lastTouchpointAt = row?.value?.lastTouchpoint?.lastTouchPointAt;

        return isAfter(new Date(lastTouchpointAt), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsEmployeeCount },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value.split('-').map(Number) as number[];
        const employees = row.value.employees;

        if (filterValue.length !== 2) return employees >= filterValue[0];

        return employees >= filterValue[0] && employees <= filterValue[1];
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLeadSource },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(row.value.leadSource);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsIndustry },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(row.value.industry);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLtv },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const ltv = row.value.accountDetails?.ltv;

        if (!ltv) return false;

        if (filterValue.length !== 2) return ltv >= filterValue[0];

        return ltv >= filterValue[0] && ltv <= filterValue[1];
      },
    )
    .otherwise(() => noop);
};
