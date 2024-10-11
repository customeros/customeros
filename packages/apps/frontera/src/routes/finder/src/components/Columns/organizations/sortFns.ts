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
      ColumnViewType.OrganizationsName,
      () => (row: OrganizationStore) =>
        row.value?.name?.trim().toLocaleLowerCase() || null,
    )
    .with(
      ColumnViewType.OrganizationsRelationship,
      () => (row: OrganizationStore) =>
        row.value?.relationship === OrganizationRelationship.Customer,
    )
    .with(
      ColumnViewType.OrganizationsOnboardingStatus,
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
      ColumnViewType.OrganizationsRenewalLikelihood,
      () => (row: OrganizationStore) =>
        match(row.value?.accountDetails?.renewalSummary?.renewalLikelihood)
          .with(OpportunityRenewalLikelihood.HighRenewal, () => 3)
          .with(OpportunityRenewalLikelihood.MediumRenewal, () => 2)
          .with(OpportunityRenewalLikelihood.LowRenewal, () => 1)
          .otherwise(() => null),
    )

    .with(
      ColumnViewType.OrganizationsRenewalDate,
      () => (row: OrganizationStore) => {
        const value =
          row.value?.accountDetails?.renewalSummary?.nextRenewalDate;

        return value ? new Date(value) : null;
      },
    )
    .with(
      ColumnViewType.OrganizationsForecastArr,
      () => (row: OrganizationStore) =>
        row.value?.accountDetails?.renewalSummary?.arrForecast,
    )
    .with(ColumnViewType.OrganizationsOwner, () => (row: OrganizationStore) => {
      const name = row.value?.owner?.name ?? '';
      const firstName = row.value?.owner?.firstName ?? '';
      const lastName = row.value?.owner?.lastName ?? '';

      const fullName = (name ?? `${firstName} ${lastName}`).trim();

      return fullName.length ? fullName.toLocaleLowerCase() : null;
    })
    .with(
      ColumnViewType.OrganizationsLeadSource,
      () => (row: OrganizationStore) => row.value?.leadSource,
    )
    .with(
      ColumnViewType.OrganizationsCreatedDate,
      () => (row: OrganizationStore) =>
        row.value?.metadata?.created
          ? new Date(row.value?.metadata?.created)
          : null,
    )
    .with(
      ColumnViewType.OrganizationsYearFounded,
      () => (row: OrganizationStore) => row.value?.yearFounded,
    )
    .with(
      ColumnViewType.OrganizationsEmployeeCount,
      () => (row: OrganizationStore) => row.value?.employees,
    )
    .with(
      ColumnViewType.OrganizationsSocials,
      () => (row: OrganizationStore) => row.value?.socialMedia?.[0]?.url,
    )
    .with(
      ColumnViewType.OrganizationsLastTouchpoint,
      () => (row: OrganizationStore) => {
        const value = row.value?.lastTouchpoint?.lastTouchPointAt;

        if (!value) return null;

        return new Date(value);
      },
    )
    .with(
      ColumnViewType.OrganizationsLastTouchpointDate,
      () => (row: OrganizationStore) => {
        const value = row.value?.lastTouchpoint?.lastTouchPointAt;

        return value ? new Date(value) : null;
      },
    )
    .with(
      ColumnViewType.OrganizationsChurnDate,
      () => (row: OrganizationStore) => {
        const value = row.value?.accountDetails?.churned;

        return value ? new Date(value) : null;
      },
    )
    .with(
      ColumnViewType.OrganizationsLtv,
      () => (row: OrganizationStore) => row.value?.accountDetails?.ltv,
    )
    .with(
      ColumnViewType.OrganizationsIndustry,
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
    .with(
      ColumnViewType.OrganizationsHeadquarters,
      () => (row: OrganizationStore) => {
        return row.country?.toLowerCase() || null;
      },
    )
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
    .with(
      ColumnViewType.OrganizationsParentOrganization,
      () => (row: OrganizationStore) => {
        return (
          row.value?.parentCompanies?.[0]?.organization?.name
            .trim()
            .toLowerCase() || null
        );
      },
    )
    .otherwise(() => (_row: OrganizationStore) => false);
