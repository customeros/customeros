import { match } from 'ts-pattern';

import {
  ColumnViewType,
  OnboardingStatus,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import type { Organization } from '../Organization.dto';

export const getOrganizationSortFn = (columnId: string) =>
  match(columnId)
    .with(ColumnViewType.OrganizationsName, () => (row: Organization) => {
      if (row.name === null || row.name.trim() === '') return null;

      return row.name.trim().toLowerCase();
    })
    .with(
      ColumnViewType.OrganizationsRelationship,
      () => (row: Organization) =>
        row?.value.relationship === OrganizationRelationship.Customer,
    )
    .with(
      ColumnViewType.OrganizationsOnboardingStatus,
      () => (row: Organization) =>
        match(row?.value.onboardingStatus)
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
      () => (row: Organization) =>
        match(row?.value.renewalSummaryRenewalLikelihood)
          .with(OpportunityRenewalLikelihood.HighRenewal, () => 3)
          .with(OpportunityRenewalLikelihood.MediumRenewal, () => 2)
          .with(OpportunityRenewalLikelihood.LowRenewal, () => 1)
          .otherwise(() => null),
    )
    .with(
      ColumnViewType.OrganizationsUpdatedDate,
      () => (row: Organization) => {
        return row?.value?.updatedAt ? new Date(row?.value?.updatedAt) : null;
      },
    )
    .with(
      ColumnViewType.OrganizationsRenewalDate,
      () => (row: Organization) => {
        const value = row?.value?.renewalSummaryNextRenewalAt;

        return value ? new Date(value) : null;
      },
    )
    .with(
      ColumnViewType.OrganizationsForecastArr,
      () => (row: Organization) => row?.value?.renewalSummaryArrForecast,
    )
    .with(ColumnViewType.OrganizationsOwner, () => (row: Organization) => {
      const name = row?.owner?.name ?? '';
      const firstName = row?.owner?.value.firstName ?? '';
      const lastName = row?.owner?.value.lastName ?? '';

      const fullName = (name ?? `${firstName} ${lastName}`).trim();

      return fullName.length ? fullName.toLocaleLowerCase() : null;
    })
    .with(
      ColumnViewType.OrganizationsLeadSource,
      () => (row: Organization) => row?.value.leadSource,
    )
    .with(
      ColumnViewType.OrganizationsCreatedDate,
      () => (row: Organization) =>
        row?.value.createdAt ? new Date(row?.value.createdAt) : null,
    )
    .with(
      ColumnViewType.OrganizationsYearFounded,
      () => (row: Organization) => row?.value.yearFounded,
    )
    .with(
      ColumnViewType.OrganizationsEmployeeCount,
      () => (row: Organization) => row?.value.employees,
    )
    .with(
      ColumnViewType.OrganizationsSocials,
      () => (row: Organization) => row?.value.socialMedia?.[0]?.url,
    )
    .with(
      ColumnViewType.OrganizationsLastTouchpoint,
      () => (row: Organization) => {
        const value = row?.value?.lastTouchPointAt;

        if (!value) return null;

        return new Date(value);
      },
    )
    .with(
      ColumnViewType.OrganizationsLastTouchpointDate,
      () => (row: Organization) => {
        const value = row?.value?.lastTouchPointAt;

        return value ? new Date(value) : null;
      },
    )
    .with(ColumnViewType.OrganizationsChurnDate, () => (row: Organization) => {
      const value = row?.value.churnedAt;

      return value ? new Date(value) : null;
    })
    .with(
      ColumnViewType.OrganizationsLtv,
      () => (row: Organization) => row?.value?.ltv,
    )
    .with(
      ColumnViewType.OrganizationsIndustry,
      () => (row: Organization) => row?.value.industryName,
    )
    .with(
      ColumnViewType.OrganizationsContactCount,
      () => (row: Organization) => row?.value.contacts?.length,
    )
    .with(
      ColumnViewType.OrganizationsLinkedinFollowerCount,
      () => (row: Organization) =>
        row.value.socialMedia.find((e) => e?.url?.includes('linkedin'))
          ?.followersCount,
    )
    .with(
      ColumnViewType.OrganizationsHeadquarters,
      () => (row: Organization) => {
        return row.country?.toLowerCase() || null;
      },
    )
    .with(
      ColumnViewType.OrganizationsIsPublic,
      () => (row: Organization) => row.value.public,
    )
    .with(
      ColumnViewType.OrganizationsStage,
      () => (row: Organization) => row.value.stage?.toLowerCase(),
    )
    .with(ColumnViewType.OrganizationsTags, () => (row: Organization) => {
      return row?.value.tags?.[0]?.name?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.OrganizationsParentOrganization,
      () => (row: Organization) => {
        return row?.value?.parentName?.trim().toLowerCase() || null;
      },
    )
    .otherwise(() => (_row: Organization) => false);
