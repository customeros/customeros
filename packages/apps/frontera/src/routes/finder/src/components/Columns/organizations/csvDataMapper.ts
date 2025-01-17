import { Organization } from '@store/Organizations/Organization.dto';

import { DateTimeUtils } from '@utils/date.ts';
import { ColumnViewType } from '@graphql/types';

export const csvDataMapper = {
  [ColumnViewType.OrganizationsAvatar]: (d: Organization) => d?.value?.logoUrl,
  [ColumnViewType.OrganizationsName]: (d: Organization) => d.value?.name,
  [ColumnViewType.OrganizationsRelationship]: (d: Organization) =>
    d.value?.relationship,
  [ColumnViewType.OrganizationsUpdatedDate]: (d: Organization) =>
    d?.value?.updatedAt
      ? DateTimeUtils.format(d.value?.updatedAt, DateTimeUtils.iso8601)
      : 'Unknown',
  [ColumnViewType.OrganizationsLeadSource]: (d: Organization) =>
    d?.value?.leadSource || 'Unknown',

  [ColumnViewType.OrganizationsOnboardingStatus]: (d: Organization) =>
    d?.value?.onboardingStatus,
  [ColumnViewType.OrganizationsRenewalLikelihood]: (d: Organization) =>
    d?.value?.renewalSummaryRenewalLikelihood,
  [ColumnViewType.OrganizationsRenewalDate]: (d: Organization) =>
    DateTimeUtils.format(
      d?.value?.renewalSummaryNextRenewalAt,
      DateTimeUtils.iso8601,
    ),
  [ColumnViewType.OrganizationsForecastArr]: (d: Organization) =>
    d?.value?.renewalSummaryArrForecast,

  [ColumnViewType.OrganizationsOwner]: (d: Organization) => {
    return d.owner?.name ?? 'No owner';
  },

  [ColumnViewType.OrganizationsCreatedDate]: (d: Organization) =>
    DateTimeUtils.format(d.value?.createdAt, DateTimeUtils.iso8601),
  [ColumnViewType.OrganizationsYearFounded]: (d: Organization) =>
    d.value?.yearFounded,
  [ColumnViewType.OrganizationsEmployeeCount]: (d: Organization) =>
    d.value?.employees,
  [ColumnViewType.OrganizationsSocials]: (d: Organization) =>
    d.value?.socialMedia.find((e) => e?.url?.includes('linkedin'))?.url,
  [ColumnViewType.OrganizationsPrimaryDomains]: (d: Organization) =>
    d.value?.domainsDetails
      .filter((e) => e.primary)
      ?.map((e) => e.domain)
      .join('; '),

  [ColumnViewType.OrganizationsLastTouchpoint]: (d: Organization) =>
    `${d?.value?.lastTouchPointType} - ${DateTimeUtils.format(
      d?.value?.lastTouchPointAt,
      DateTimeUtils.iso8601,
    )}`,
  [ColumnViewType.OrganizationsLastTouchpointDate]: (d: Organization) =>
    d?.value?.lastTouchPointAt
      ? DateTimeUtils.format(d.value?.lastTouchPointAt, DateTimeUtils.iso8601)
      : 'Unknown',
  [ColumnViewType.OrganizationsChurnDate]: (d: Organization) =>
    d?.value?.churnedAt
      ? DateTimeUtils.format(d.value?.churnedAt, DateTimeUtils.iso8601)
      : 'Unknown',
  [ColumnViewType.OrganizationsLtv]: (d: Organization) => d?.value?.ltv,
  [ColumnViewType.OrganizationsIndustry]: (d: Organization) =>
    d.value?.industryName ?? 'Unknown',
  [ColumnViewType.OrganizationsContactCount]: (d: Organization) =>
    d?.contacts?.length,
  [ColumnViewType.OrganizationsLinkedinFollowerCount]: (d: Organization) =>
    d?.value?.socialMedia.find((e) => e?.url?.includes('linkedin'))
      ?.followersCount ?? 'Unknown',
  [ColumnViewType.OrganizationsTags]: (d: Organization) =>
    d?.value?.tags?.map((e) => e.name).join('; '),
  [ColumnViewType.OrganizationsIsPublic]: (d: Organization) =>
    d?.value?.public ? 'Public' : 'Private',
  [ColumnViewType.OrganizationsStage]: (d: Organization) => d.value?.stage,
  [ColumnViewType.OrganizationsCity]: (d: Organization) =>
    d?.value?.locations?.[0]?.countryCodeA2,
  [ColumnViewType.OrganizationsHeadquarters]: (d: Organization) =>
    d?.value?.locations?.[0]?.countryCodeA2,
};
