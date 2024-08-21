import { match } from 'ts-pattern';
import { OrganizationStore } from '@store/Organizations/Organization.store';

import {
  Social,
  ColumnViewType,
  OnboardingStatus,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

export const getOrganizationSortFn = (columnId: string) =>
  match(columnId)
    .with(
      'ORGANIZATIONS_NAME',
      () => (row: OrganizationStore) =>
        row.value?.name?.trim().toLocaleLowerCase() || null,
    )
    .with(
      'ORGANIZATIONS_RELATIONSHIP',
      () => (row: OrganizationStore) =>
        row.value?.relationship === OrganizationRelationship.Customer,
    )
    .with(
      'ORGANIZATIONS_ONBOARDING_STATUS',
      () => (row: OrganizationStore) =>
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
      () => (row: OrganizationStore) =>
        match(row.value?.accountDetails?.renewalSummary?.renewalLikelihood)
          .with(OpportunityRenewalLikelihood.HighRenewal, () => 3)
          .with(OpportunityRenewalLikelihood.MediumRenewal, () => 2)
          .with(OpportunityRenewalLikelihood.LowRenewal, () => 1)
          .otherwise(() => null),
    )
    .with('ORGANIZATIONS_RENEWAL_DATE', () => (row: OrganizationStore) => {
      const value = row.value?.accountDetails?.renewalSummary?.nextRenewalDate;

      return value ? new Date(value) : null;
    })
    .with(
      'ORGANIZATIONS_FORECAST_ARR',
      () => (row: OrganizationStore) =>
        row.value?.accountDetails?.renewalSummary?.arrForecast,
    )
    .with('ORGANIZATIONS_OWNER', () => (row: OrganizationStore) => {
      const name = row.value?.owner?.name ?? '';
      const firstName = row.value?.owner?.firstName ?? '';
      const lastName = row.value?.owner?.lastName ?? '';

      const fullName = (name ?? `${firstName} ${lastName}`).trim();

      return fullName.length ? fullName.toLocaleLowerCase() : null;
    })
    .with(
      'ORGANIZATIONS_LEAD_SOURCE',
      () => (row: OrganizationStore) => row.value?.leadSource,
    )
    .with(
      'ORGANIZATIONS_CREATED_DATE',
      () => (row: OrganizationStore) =>
        row.value?.metadata?.created
          ? new Date(row.value?.metadata?.created)
          : null,
    )
    .with(
      'ORGANIZATIONS_YEAR_FOUNDED',
      () => (row: OrganizationStore) => row.value?.yearFounded,
    )
    .with(
      'ORGANIZATIONS_EMPLOYEE_COUNT',
      () => (row: OrganizationStore) => row.value?.employees,
    )
    .with(
      'ORGANIZATIONS_SOCIALS',
      () => (row: OrganizationStore) => row.value?.socialMedia?.[0]?.url,
    )
    .with('ORGANIZATIONS_LAST_TOUCHPOINT', () => (row: OrganizationStore) => {
      const value = row.value?.lastTouchpoint?.lastTouchPointAt;

      if (!value) return null;

      return new Date(value);
    })
    .with(
      'ORGANIZATIONS_LAST_TOUCHPOINT_DATE',
      () => (row: OrganizationStore) => {
        const value = row.value?.lastTouchpoint?.lastTouchPointAt;

        return value ? new Date(value) : null;
      },
    )
    .with('ORGANIZATIONS_CHURN_DATE', () => (row: OrganizationStore) => {
      const value = row.value?.accountDetails?.churned;

      return value ? new Date(value) : null;
    })
    .with(
      'ORGANIZATIONS_LTV',
      () => (row: OrganizationStore) => row.value?.accountDetails?.ltv,
    )
    .with(
      'ORGANIZATIONS_INDUSTRY',
      () => (row: OrganizationStore) => row.value?.industry,
    )
    .with(
      ColumnViewType.OrganizationsContactCount,
      () => (row: OrganizationStore) => row.value?.contacts?.content?.length,
    )
    .with(
      ColumnViewType.OrganizationsLinkedinFollowerCount,
      () => (row: OrganizationStore) =>
        row.value.socialMedia.find((e: Social) => e?.url?.includes('linkedin'))
          ?.followersCount,
    )
    .with(ColumnViewType.OrganizationsCity, () => (row: OrganizationStore) => {
      return row.country?.toLowerCase() || null;
    })
    .with(
      ColumnViewType.OrganizationsIsPublic,
      () => (row: OrganizationStore) => row.value.public,
    )
    .with(
      ColumnViewType.OrganizationsStage,
      () => (row: OrganizationStore) => row.value.stage?.toLowerCase(),
    )
    .with(ColumnViewType.OrganizationsTags, () => (row: OrganizationStore) => {
      return row.value?.tags?.[0]?.name?.trim().toLowerCase() || null;
    })
    .otherwise(() => (_row: OrganizationStore) => false);
