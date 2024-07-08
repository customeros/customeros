import { DateTimeUtils } from '@utils/date.ts';
import { Social, Contact, Organization, ColumnViewType } from '@graphql/types';

export const csvDataMapper = {
  [ColumnViewType.OrganizationsAvatar]: (d: Organization) => d?.logo,
  [ColumnViewType.OrganizationsName]: (d: Organization) => d.name,
  [ColumnViewType.OrganizationsWebsite]: (d: Organization) => d.website,
  [ColumnViewType.OrganizationsRelationship]: (d: Organization) =>
    d.relationship,

  [ColumnViewType.OrganizationsOnboardingStatus]: (d: Organization) =>
    d?.accountDetails?.onboarding?.status,
  [ColumnViewType.OrganizationsRenewalLikelihood]: (d: Organization) =>
    d?.accountDetails?.renewalSummary?.renewalLikelihood,
  [ColumnViewType.OrganizationsRenewalDate]: (d: Organization) =>
    DateTimeUtils.format(
      d?.accountDetails?.renewalSummary?.nextRenewalDate,
      DateTimeUtils.iso8601,
    ),
  [ColumnViewType.OrganizationsForecastArr]: (d: Organization) =>
    d?.accountDetails?.renewalSummary?.arrForecast,

  [ColumnViewType.OrganizationsOwner]: (d: Organization) => {
    return (
      d.owner?.name ??
      `${d.owner?.firstName ?? ''} ${d.owner?.lastName ?? ''}`?.trim()
    );
  },
  [ColumnViewType.OrganizationsLeadSource]: (d: Organization) => d.leadSource,
  [ColumnViewType.OrganizationsCreatedDate]: (d: Organization) =>
    DateTimeUtils.format(d.metadata.created, DateTimeUtils.iso8601),
  [ColumnViewType.OrganizationsYearFounded]: (d: Organization) => d.yearFounded,
  [ColumnViewType.OrganizationsEmployeeCount]: (data: Organization) =>
    data.employees,
  [ColumnViewType.OrganizationsSocials]: () => 'linkedin/in/nana',
  [ColumnViewType.OrganizationsLastTouchpoint]: (data: Organization) =>
    `${data?.lastTouchpoint?.lastTouchPointType} - ${DateTimeUtils.format(
      data?.lastTouchpoint?.lastTouchPointAt,
      DateTimeUtils.iso8601,
    )}`,
  [ColumnViewType.OrganizationsLastTouchpointDate]: (data: Organization) =>
    data?.lastTouchpoint?.lastTouchPointAt
      ? DateTimeUtils.format(
          data.lastTouchpoint.lastTouchPointAt,
          DateTimeUtils.iso8601,
        )
      : 'Unknown',
  [ColumnViewType.OrganizationsChurnDate]: (data: Organization) =>
    data?.accountDetails?.churned
      ? DateTimeUtils.format(data.accountDetails.churned, DateTimeUtils.iso8601)
      : 'Unknown',
  [ColumnViewType.OrganizationsLtv]: (data: Organization) =>
    data?.accountDetails?.ltv,
  [ColumnViewType.OrganizationsIndustry]: (data: Organization) =>
    data.industry ?? 'Unknown',
  [ColumnViewType.OrganizationsContactCount]: (data: Organization) =>
    data?.contacts?.content?.filter((e: Contact) => e.tags)?.length,
  [ColumnViewType.OrganizationsLinkedinFollowerCount]: (data: Organization) =>
    data?.socialMedia.find((e: Social) => e?.url?.includes('linkedin'))
      ?.followersCount ?? 'Unknown',
  [ColumnViewType.OrganizationsTags]: (data: Organization) =>
    data?.tags?.map((e) => e.name).join('; '),
  [ColumnViewType.OrganizationsIsPublic]: (data: Organization) =>
    data.isPublic ? 'Public' : 'Private',
  [ColumnViewType.OrganizationsStage]: (data: Organization) => data.stage,
  [ColumnViewType.OrganizationsCity]: (data: Organization) =>
    data?.locations?.[0]?.countryCodeA2,
  [ColumnViewType.OrganizationsHeadquarters]: (data: Organization) =>
    data?.locations?.[0]?.countryCodeA2,
};
