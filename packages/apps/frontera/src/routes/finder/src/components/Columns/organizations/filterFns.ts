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
  ComparisonOperator,
  OrganizationRelationship,
} from '@graphql/types';

const getFilterV2Fn = (filter: FilterItem | undefined | null) => {
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
          Array.isArray(values) ? values : [values],
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
        const filterValue = filter?.value;

        return isAfter(
          new Date(row.value.metadata.created),
          new Date(filterValue),
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsName },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value?.toLowerCase();
        const filterOperator = filter?.operation;

        const values = row.value.name?.toLowerCase();

        if (
          filter.includeEmpty &&
          filterOperator === ComparisonOperator.IsEmpty
        ) {
          return !values;
        }

        if (filterOperator === ComparisonOperator.IsNotEmpty) {
          return values;
        }

        if (filterOperator === ComparisonOperator.NotContains) {
          return !values.includes(filterValue);
        }

        {
          return values.includes(filterValue);
        }
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsParentOrganization },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value?.toLowerCase();
        const filterOperator = filter?.operation;

        const values =
          row.value.parentCompanies?.[0]?.organization.name?.toLowerCase();

        if (
          filter.includeEmpty &&
          filterOperator === ComparisonOperator.IsEmpty
        ) {
          return !values;
        }

        if (filterOperator === ComparisonOperator.IsNotEmpty) {
          return values;
        }

        if (filterOperator === ComparisonOperator.NotContains) {
          return !values.includes(filterValue);
        }

        {
          return values.includes(filterValue);
        }
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsWebsite },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;

        const websiteUrl = row.value.website || '';
        const filterText = filter.value || '';

        if (filter.includeEmpty && websiteUrl === '') return true;

        if (filterText === '') {
          return !filter.includeEmpty;
        }

        return websiteUrl?.toLowerCase().includes(filterText?.toLowerCase());
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRelationship },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(row.value.relationship);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsStage },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(row.value.stage);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsForecastArr },
      (filter) => (row: OrganizationStore) => {
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
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const nextRenewalDate =
          row.value?.accountDetails?.renewalSummary?.nextRenewalDate?.split(
            'T',
          )[0];

        if (!filterValue) return true;
        if (filterValue?.[1] === null)
          return filterValue?.[0] <= nextRenewalDate;
        if (filterValue?.[0] === null)
          return filterValue?.[1] >= nextRenewalDate;

        return (
          filterValue[0] <= nextRenewalDate && filterValue[1] >= nextRenewalDate
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsOnboardingStatus },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(
          row.value.accountDetails?.onboarding?.status,
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRenewalLikelihood },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(
          row.value.accountDetails?.renewalSummary?.renewalLikelihood,
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsOwner },
      (filter) => (row: OrganizationStore) => {
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
      (filter) => (row: OrganizationStore) => {
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
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const churned = row?.value?.accountDetails?.churned;

        if (!churned) return false;

        return isAfter(new Date(churned), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsSocials },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const linkedInUrl = row.value.socialMedia?.find((v) =>
          v.url.includes('linkedin'),
        )?.url;

        if (!filterValue && filter.active && !filter.includeEmpty) return true;
        if (!linkedInUrl && filter.includeEmpty) return true;
        if (!filterValue || !linkedInUrl) return false;

        return linkedInUrl && linkedInUrl.includes(filterValue);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLastTouchpointDate },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const lastTouchpointAt = row?.value?.lastTouchpoint?.lastTouchPointAt;

        return isAfter(new Date(lastTouchpointAt), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsEmployeeCount },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const operator = filter.operation;
        const employees = row.value.employees;

        if (operator === ComparisonOperator.Lt) {
          return employees < Number(filterValue);
        }

        if (operator === ComparisonOperator.Gt) {
          return employees > Number(filterValue);
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
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const operator = filter.operation;
        const followers = row.value.socialMedia.find((e: Social) =>
          e?.url?.includes('linkedin'),
        )?.followersCount;

        if (operator === ComparisonOperator.Lt) {
          return followers < Number(filterValue);
        }

        if (operator === ComparisonOperator.Gt) {
          return followers > Number(filterValue);
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
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(row.value.leadSource);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsIndustry },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (!filterValue) return false;

        return filterValue.includes(row.value.industry);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLtv },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const ltv = row.value.accountDetails?.ltv;

        if (!ltv) return false;

        if (filterValue.length !== 2) return ltv >= filterValue[0];

        return ltv >= filterValue[0] && ltv <= filterValue[1];
      },
    )
    .with({ property: ColumnViewType.OrganizationsHeadquarters }, (filter) => {
      // Early exit if filter is not active
      if (!filter.active) return () => true;

      const filterValue = filter.value;
      const includeEmpty = filter.includeEmpty;

      return (row: OrganizationStore) => {
        const locations = row.value.locations;
        const country = locations?.[0]?.countryCodeA2;

        if (!country) {
          return includeEmpty;
        }

        if (!filterValue.length) {
          return !includeEmpty;
        }

        return filterValue.includes(country);
      };
    })
    .with(
      { property: ColumnViewType.OrganizationsIsPublic },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const isPublic = row.value.public;

        if (filterValue.includes('public') && isPublic) return true;

        return filterValue.includes('private') && !isPublic;
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsYearFounded },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const operator = filter.operation;

        const yearFounded = row.value.yearFounded;
        const currentYear = new Date().getFullYear();
        const age = currentYear - yearFounded;

        if (!yearFounded) return false;

        if (operator === ComparisonOperator.Lt) {
          return age < Number(filterValue);
        }

        if (operator === ComparisonOperator.Gt) {
          return age > Number(filterValue);
        }

        if (operator === ComparisonOperator.Between) {
          const filterValue = filter?.value?.map(Number) as number[];

          return age >= filterValue[0] && age <= filterValue[1];
        }

        return true;
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsTags },
      (filter) => (row: OrganizationStore) => {
        if (!filter.active) return true;
        const tags = row.value.tags?.map((l: Tag) => l.name);

        if (!filter.value?.length) return true;

        if (!tags?.length) return false;

        return filter.value.some((f: string) => tags.includes(f));
      },
    )

    .otherwise(() => noop);
};

export const getOrganizationDefaultFilterFns = (
  filters: Filter | null,
  isFeatureEnabled: boolean,
) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  if (isFeatureEnabled) {
    return data.map(({ filter }) => getFilterV2Fn(filter));
  }

  return data.map(({ filter }) => getFilterFn(filter));
};

export const getOrganizationFilterFns = (
  filters: Filter | null,
  isFeatureEnabled: boolean,
) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  if (isFeatureEnabled) {
    return data.map(({ filter }) => getFilterV2Fn(filter));
  }

  return data.map(({ filter }) => getFilterFn(filter));
};
