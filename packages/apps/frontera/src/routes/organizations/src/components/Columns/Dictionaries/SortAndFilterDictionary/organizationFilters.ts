import { match } from 'ts-pattern';
import { Store } from '@store/store.ts';
import { isAfter } from 'date-fns/isAfter';
import { FilterItem } from '@store/types.ts';

import {
  Social,
  Organization,
  ColumnViewType,
  ComparisonOperator,
} from '@graphql/types';

function checkCommonStrings(
  array1: (string | null | undefined)[],
  array2: Array<string>,
) {
  const set1 = new Set(array1);
  const set2 = new Set(array2);

  return [...set1].filter((item) => set2.has(<string>item));
}
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
        const filterValue = filter?.value;

        const operator = filter.operation;
        const employees = row.value.employees;

        if (operator === ComparisonOperator.Lte) {
          return employees <= filterValue[0];
        }
        if (operator === ComparisonOperator.Gte) {
          return employees >= filterValue[0];
        }

        if (operator === ComparisonOperator.Between) {
          const filterValue = filter?.value?.map(Number) as number[];

          return employees >= filterValue[0] && employees <= filterValue[1];
        }

        return true;
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLinkedinFollowerCount },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const operator = filter.operation;
        const followers = row.value.socialMedia.find((e: Social) =>
          e?.url?.includes('linkedin'),
        )?.followersCount;

        if (operator === ComparisonOperator.Lte) {
          return followers <= filterValue[0];
        }
        if (operator === ComparisonOperator.Gte) {
          return followers >= filterValue[0];
        }

        if (operator === ComparisonOperator.Between) {
          const filterValue = filter?.value?.map(Number) as number[];

          return followers >= filterValue[0] && followers <= filterValue[1];
        }

        return true;
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

        if (!filterValue) return false;

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
    .with(
      { property: ColumnViewType.OrganizationsCity },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const countries = row.value.locations
          .map((l) => l.countryCodeA2)
          .filter((l) => !!l?.length);

        return checkCommonStrings(countries, filterValue).length > 0;
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsIsPublic },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const isPublic = row.value.public;

        if (filterValue.includes('public') && isPublic) return true;

        return filterValue.includes('private') && !isPublic;
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsYearFounded },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const operator = filter.operation;

        const yearFounded = row.value.yearFounded;
        const currentYear = new Date().getFullYear();
        const age = currentYear - yearFounded;
        if (!yearFounded) return false;
        if (operator === ComparisonOperator.Lte) {
          return age <= Number(filterValue);
        }
        if (operator === ComparisonOperator.Gte) {
          return age >= Number(filterValue);
        }

        if (operator === ComparisonOperator.Between) {
          const filterValue = filter?.value?.map(Number) as number[];

          return age >= filterValue[0] && age <= filterValue[1];
        }

        return true;
      },
    )

    .otherwise(() => noop);
};
