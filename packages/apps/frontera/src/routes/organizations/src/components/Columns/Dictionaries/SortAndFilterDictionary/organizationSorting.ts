import { match } from 'ts-pattern';
import { Store } from '@store/store.ts';
import countries from '@assets/countries/countries.json';

import {
  Social,
  Contact,
  Organization,
  ColumnViewType,
  OnboardingStatus,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

export const getOrganizationColumnSortFn = (columnId: string) =>
  match(columnId)
    .with(
      'ORGANIZATIONS_NAME',
      () => (row: Store<Organization>) =>
        row.value?.name?.trim().toLocaleLowerCase() || null,
    )
    .with(
      'ORGANIZATIONS_RELATIONSHIP',
      () => (row: Store<Organization>) => row.value?.isCustomer,
    )
    .with(
      'ORGANIZATIONS_ONBOARDING_STATUS',
      () => (row: Store<Organization>) =>
        match(row.value?.accountDetails?.onboarding?.status)
          .with(OnboardingStatus.NotApplicable, () => null)
          .with(OnboardingStatus.NotStarted, () => 1)
          .with(OnboardingStatus.OnTrack, () => 2)
          .with(OnboardingStatus.Late, () => 3)
          .with(OnboardingStatus.Stuck, () => 4)
          .with(OnboardingStatus.Successful, () => 5)
          .with(OnboardingStatus.Done, () => 6)
          .otherwise(() => null),
    )
    .with(
      'ORGANIZATIONS_RENEWAL_LIKELIHOOD',
      () => (row: Store<Organization>) =>
        match(row.value?.accountDetails?.renewalSummary?.renewalLikelihood)
          .with(OpportunityRenewalLikelihood.HighRenewal, () => 3)
          .with(OpportunityRenewalLikelihood.MediumRenewal, () => 2)
          .with(OpportunityRenewalLikelihood.LowRenewal, () => 1)
          .otherwise(() => null),
    )
    .with('ORGANIZATIONS_RENEWAL_DATE', () => (row: Store<Organization>) => {
      const value = row.value?.accountDetails?.renewalSummary?.nextRenewalDate;

      return value ? new Date(value) : null;
    })
    .with(
      'ORGANIZATIONS_FORECAST_ARR',
      () => (row: Store<Organization>) =>
        row.value?.accountDetails?.renewalSummary?.arrForecast,
    )
    .with('ORGANIZATIONS_OWNER', () => (row: Store<Organization>) => {
      const name = row.value?.owner?.name ?? '';
      const firstName = row.value?.owner?.firstName ?? '';
      const lastName = row.value?.owner?.lastName ?? '';

      const fullName = (name ?? `${firstName} ${lastName}`).trim();

      return fullName.length ? fullName.toLocaleLowerCase() : null;
    })
    .with(
      'ORGANIZATIONS_LEAD_SOURCE',
      () => (row: Store<Organization>) => row.value?.leadSource,
    )
    .with(
      'ORGANIZATIONS_CREATED_DATE',
      () => (row: Store<Organization>) =>
        row.value?.metadata?.created
          ? new Date(row.value?.metadata?.created)
          : null,
    )
    .with(
      'ORGANIZATIONS_YEAR_FOUNDED',
      () => (row: Store<Organization>) => row.value?.yearFounded,
    )
    .with(
      'ORGANIZATIONS_EMPLOYEE_COUNT',
      () => (row: Store<Organization>) => row.value?.employees,
    )
    .with(
      'ORGANIZATIONS_SOCIALS',
      () => (row: Store<Organization>) => row.value?.socialMedia?.[0]?.url,
    )
    .with('ORGANIZATIONS_LAST_TOUCHPOINT', () => (row: Store<Organization>) => {
      const value = row.value?.lastTouchpoint?.lastTouchPointAt;

      if (!value) return null;

      return new Date(value);
    })
    .with(
      'ORGANIZATIONS_LAST_TOUCHPOINT_DATE',
      () => (row: Store<Organization>) => {
        const value = row.value?.lastTouchpoint?.lastTouchPointAt;

        return value ? new Date(value) : null;
      },
    )
    .with('ORGANIZATIONS_CHURN_DATE', () => (row: Store<Organization>) => {
      const value = row.value?.accountDetails?.churned;

      return value ? new Date(value) : null;
    })
    .with(
      'ORGANIZATIONS_LTV',
      () => (row: Store<Organization>) => row.value?.accountDetails?.ltv,
    )
    .with(
      'ORGANIZATIONS_INDUSTRY',
      () => (row: Store<Organization>) => row.value?.industry,
    )
    .with(
      ColumnViewType.OrganizationsContactCount,
      () => (row: Store<Organization>) =>
        row.value?.contacts?.content?.filter(
          (e) => e?.tags?.length && e.tags?.length > 0,
        ).length,
    )
    .with(
      ColumnViewType.OrganizationsLinkedinFollowerCount,
      () => (row: Store<Organization>) =>
        row.value.socialMedia.find((e: Social) => e?.url?.includes('linkedin'))
          ?.followersCount,
    )
    .with(
      ColumnViewType.OrganizationsCity,
      () => (row: Store<Organization>) => {
        const countryName = countries.find(
          (d) =>
            d.alpha2 === row.value.locations?.[0]?.countryCodeA2?.toLowerCase(),
        );

        return countryName?.name?.toLowerCase() || null;
      },
    )
    .with(
      ColumnViewType.OrganizationsIsPublic,
      () => (row: Store<Organization>) => row.value.public,
    )
    .with(ColumnViewType.OrganizationsTags, () => (row: Store<Contact>) => {
      return row.value?.tags?.[0]?.name?.trim().toLowerCase() || null;
    })
    .otherwise(() => (_row: Store<Organization>) => false);
