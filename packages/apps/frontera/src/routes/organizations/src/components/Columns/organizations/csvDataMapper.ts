import { OrganizationStore } from '@store/Organizations/Organization.store';

import { DateTimeUtils } from '@utils/date.ts';
import { Social, ColumnViewType } from '@graphql/types';

export const csvDataMapper = {
  [ColumnViewType.OrganizationsAvatar]: (d: OrganizationStore) =>
    d?.value?.logo,
  [ColumnViewType.OrganizationsName]: (d: OrganizationStore) => d.value?.name,
  [ColumnViewType.OrganizationsWebsite]: (d: OrganizationStore) =>
    d.value?.website,
  [ColumnViewType.OrganizationsRelationship]: (d: OrganizationStore) =>
    d.value?.relationship,

  [ColumnViewType.OrganizationsOnboardingStatus]: (d: OrganizationStore) =>
    d?.value?.accountDetails?.onboarding?.status,
  [ColumnViewType.OrganizationsRenewalLikelihood]: (d: OrganizationStore) =>
    d?.value?.accountDetails?.renewalSummary?.renewalLikelihood,
  [ColumnViewType.OrganizationsRenewalDate]: (d: OrganizationStore) =>
    DateTimeUtils.format(
      d?.value?.accountDetails?.renewalSummary?.nextRenewalDate,
      DateTimeUtils.iso8601,
    ),
  [ColumnViewType.OrganizationsForecastArr]: (d: OrganizationStore) =>
    d?.value?.accountDetails?.renewalSummary?.arrForecast,

  [ColumnViewType.OrganizationsOwner]: (d: OrganizationStore) => d.owner,
  [ColumnViewType.OrganizationsLeadSource]: (d: OrganizationStore) =>
    d.value?.leadSource,
  [ColumnViewType.OrganizationsCreatedDate]: (d: OrganizationStore) =>
    DateTimeUtils.format(d.value?.metadata.created, DateTimeUtils.iso8601),
  [ColumnViewType.OrganizationsYearFounded]: (d: OrganizationStore) =>
    d.value?.yearFounded,
  [ColumnViewType.OrganizationsEmployeeCount]: (d: OrganizationStore) =>
    d.value?.employees,
  [ColumnViewType.OrganizationsSocials]: (d: OrganizationStore) =>
    d.value?.socialMedia.find((e) => e?.url?.includes('linkedin'))?.url,

  [ColumnViewType.OrganizationsLastTouchpoint]: (d: OrganizationStore) =>
    `${d?.value?.lastTouchpoint?.lastTouchPointType} - ${DateTimeUtils.format(
      d?.value?.lastTouchpoint?.lastTouchPointAt,
      DateTimeUtils.iso8601,
    )}`,
  [ColumnViewType.OrganizationsLastTouchpointDate]: (d: OrganizationStore) =>
    d?.value?.lastTouchpoint?.lastTouchPointAt
      ? DateTimeUtils.format(
          d.value?.lastTouchpoint.lastTouchPointAt,
          DateTimeUtils.iso8601,
        )
      : 'Unknown',
  [ColumnViewType.OrganizationsChurnDate]: (d: OrganizationStore) =>
    d?.value?.accountDetails?.churned
      ? DateTimeUtils.format(
          d.value?.accountDetails.churned,
          DateTimeUtils.iso8601,
        )
      : 'Unknown',
  [ColumnViewType.OrganizationsLtv]: (d: OrganizationStore) =>
    d?.value?.accountDetails?.ltv,
  [ColumnViewType.OrganizationsIndustry]: (d: OrganizationStore) =>
    d.value?.industry ?? 'Unknown',
  [ColumnViewType.OrganizationsContactCount]: (d: OrganizationStore) =>
    d?.contacts?.length,
  [ColumnViewType.OrganizationsLinkedinFollowerCount]: (d: OrganizationStore) =>
    d?.value?.socialMedia.find((e: Social) => e?.url?.includes('linkedin'))
      ?.followersCount ?? 'Unknown',
  [ColumnViewType.OrganizationsTags]: (d: OrganizationStore) =>
    d?.value?.tags?.map((e) => e.name).join('; '),
  [ColumnViewType.OrganizationsIsPublic]: (d: OrganizationStore) =>
    d?.value?.isPublic ? 'Public' : 'Private',
  [ColumnViewType.OrganizationsStage]: (d: OrganizationStore) => d.value?.stage,
  [ColumnViewType.OrganizationsCity]: (d: OrganizationStore) =>
    d?.value?.locations?.[0]?.countryCodeA2,
  [ColumnViewType.OrganizationsHeadquarters]: (d: OrganizationStore) =>
    d?.value?.locations?.[0]?.countryCodeA2,
};
