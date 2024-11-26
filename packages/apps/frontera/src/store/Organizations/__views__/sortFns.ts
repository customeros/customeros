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
    .with(
      ColumnViewType.OrganizationsName,
      () => (row: Organization) =>
        row.value.name?.trim().toLocaleLowerCase() || null,
    )
    .with(
      ColumnViewType.OrganizationsRelationship,
      () => (row: Organization) =>
        row?.value.relationship === OrganizationRelationship.Customer,
    )
    .with(
      ColumnViewType.OrganizationsOnboardingStatus,
      () => (row: Organization) =>
        match(row?.value.accountDetails?.onboarding?.status)
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
        match(row?.value.accountDetails?.renewalSummary?.renewalLikelihood)
          .with(OpportunityRenewalLikelihood.HighRenewal, () => 3)
          .with(OpportunityRenewalLikelihood.MediumRenewal, () => 2)
          .with(OpportunityRenewalLikelihood.LowRenewal, () => 1)
          .otherwise(() => null),
    )
    .with(
      ColumnViewType.OrganizationsUpdatedDate,
      () => (row: Organization) => {
        return row?.value?.metadata?.lastUpdated
          ? new Date(row?.value?.metadata?.lastUpdated)
          : null;
      },
    )
    .with(
      ColumnViewType.OrganizationsRenewalDate,
      () => (row: Organization) => {
        const value =
          row?.value.accountDetails?.renewalSummary?.nextRenewalDate;

        return value ? new Date(value) : null;
      },
    )
    .with(
      ColumnViewType.OrganizationsForecastArr,
      () => (row: Organization) =>
        row?.value.accountDetails?.renewalSummary?.arrForecast,
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
        row?.value.metadata?.created
          ? new Date(row?.value.metadata?.created)
          : null,
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
        const value = row?.value.lastTouchpoint?.lastTouchPointAt;

        if (!value) return null;

        return new Date(value);
      },
    )
    .with(
      ColumnViewType.OrganizationsLastTouchpointDate,
      () => (row: Organization) => {
        const value = row?.value.lastTouchpoint?.lastTouchPointAt;

        return value ? new Date(value) : null;
      },
    )
    .with(ColumnViewType.OrganizationsChurnDate, () => (row: Organization) => {
      const value = row?.value.accountDetails?.churned;

      return value ? new Date(value) : null;
    })
    .with(
      ColumnViewType.OrganizationsLtv,
      () => (row: Organization) => row?.value.accountDetails?.ltv,
    )
    .with(
      ColumnViewType.OrganizationsIndustry,
      () => (row: Organization) => row?.value.industry,
    )
    .with(
      ColumnViewType.OrganizationsContactCount,
      () => (row: Organization) => row?.value.contacts?.content?.length,
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
        return (
          row?.value.parentCompanies?.[0]?.organization?.name
            .trim()
            .toLowerCase() || null
        );
      },
    )
    .otherwise(() => (_row: Organization) => false);
