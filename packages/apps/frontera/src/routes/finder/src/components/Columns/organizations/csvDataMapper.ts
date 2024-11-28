import { Organization } from '@store/Organizations/Organization.dto';

import { DateTimeUtils } from '@utils/date.ts';
import { ColumnViewType } from '@graphql/types';

export const csvDataMapper = {
  [ColumnViewType.OrganizationsAvatar]: (d: Organization) => d?.value?.logo,
  [ColumnViewType.OrganizationsName]: (d: Organization) => d.value?.name,
  [ColumnViewType.OrganizationsWebsite]: (d: Organization) => d.value?.website,
  [ColumnViewType.OrganizationsRelationship]: (d: Organization) =>
    d.value?.relationship,

  [ColumnViewType.OrganizationsOnboardingStatus]: (d: Organization) =>
    d?.value?.accountDetails?.onboarding?.status,
  [ColumnViewType.OrganizationsRenewalLikelihood]: (d: Organization) =>
    d?.value?.accountDetails?.renewalSummary?.renewalLikelihood,
  [ColumnViewType.OrganizationsRenewalDate]: (d: Organization) =>
    DateTimeUtils.format(
      d?.value?.accountDetails?.renewalSummary?.nextRenewalDate,
      DateTimeUtils.iso8601,
    ),
  [ColumnViewType.OrganizationsForecastArr]: (d: Organization) =>
    d?.value?.accountDetails?.renewalSummary?.arrForecast,

  [ColumnViewType.OrganizationsOwner]: (d: Organization) => {
    return d.owner?.name ?? 'No owner';
  },
  [ColumnViewType.OrganizationsLeadSource]: (d: Organization) =>
    d.value?.leadSource,
  [ColumnViewType.OrganizationsCreatedDate]: (d: Organization) =>
    DateTimeUtils.format(d.value?.metadata.created, DateTimeUtils.iso8601),
  [ColumnViewType.OrganizationsYearFounded]: (d: Organization) =>
    d.value?.yearFounded,
  [ColumnViewType.OrganizationsEmployeeCount]: (d: Organization) =>
    d.value?.employees,
  [ColumnViewType.OrganizationsSocials]: (d: Organization) =>
    d.value?.socialMedia.find((e) => e?.url?.includes('linkedin'))?.url,

  [ColumnViewType.OrganizationsLastTouchpoint]: (d: Organization) =>
    `${d?.value?.lastTouchpoint?.lastTouchPointType} - ${DateTimeUtils.format(
      d?.value?.lastTouchpoint?.lastTouchPointAt,
      DateTimeUtils.iso8601,
    )}`,
  [ColumnViewType.OrganizationsLastTouchpointDate]: (d: Organization) =>
    d?.value?.lastTouchpoint?.lastTouchPointAt
      ? DateTimeUtils.format(
          d.value?.lastTouchpoint.lastTouchPointAt,
          DateTimeUtils.iso8601,
        )
      : 'Unknown',
  [ColumnViewType.OrganizationsChurnDate]: (d: Organization) =>
    d?.value?.accountDetails?.churned
      ? DateTimeUtils.format(
          d.value?.accountDetails.churned,
          DateTimeUtils.iso8601,
        )
      : 'Unknown',
  [ColumnViewType.OrganizationsLtv]: (d: Organization) =>
    d?.value?.accountDetails?.ltv,
  [ColumnViewType.OrganizationsIndustry]: (d: Organization) =>
    d.value?.industry ?? 'Unknown',
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
