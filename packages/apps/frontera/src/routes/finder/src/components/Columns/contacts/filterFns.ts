import { match } from 'ts-pattern';
import { FilterItem } from '@store/types.ts';
import { isAfter, isBefore } from 'date-fns';
import { ContactStore } from '@store/Contacts/Contact.store.ts';

import {
  Tag,
  User,
  Filter,
  Social,
  ColumnViewType,
  EmailDeliverable,
  ComparisonOperator,
  EmailValidationDetails,
} from '@graphql/types';

import { EmailVerificationStatus } from './filterTypes';

const getFilterFn = (filter: FilterItem | undefined | null) => {
  const noop = (_row: ContactStore) => true;

  if (!filter) return noop;

  return match(filter)
    .with({ property: 'STAGE' }, (filter) => (row: ContactStore) => {
      const filterValues = filter?.value;

      if (!filterValues || !row.value?.organizations.content.length) {
        return false;
      }
      const hasOrgWithMatchingStage = row.value?.organizations.content.every(
        (o) => {
          const stage = row.root?.organizations?.value.get(o.metadata.id)?.value
            ?.stage;

          return filterValues.includes(stage);
        },
      );

      return hasOrgWithMatchingStage;
    })

    .with({ property: 'RELATIONSHIP' }, (filter) => (row: ContactStore) => {
      const filterValues = filter?.value;

      if (!filterValues || !row.value?.organizations.content.length) {
        return false;
      }
      const hasOrgWithMatchingRelationship =
        row.value?.organizations.content.every((o) => {
          const stage = row.root?.organizations?.value.get(o.metadata.id)?.value
            ?.relationship;

          return filterValues.includes(stage);
        });

      return hasOrgWithMatchingRelationship;
    })

    .with(
      { property: ColumnViewType.ContactsName },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const values = row.name;

        if (!values)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeText(filter, values);
      },
    )

    .with(
      { property: ColumnViewType.ContactsOrganization },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const orgs = row.value?.organizations?.content?.map((o) =>
          o.name.toLowerCase().trim(),
        );

        return filterTypeText(filter, orgs?.join(' '));
      },
    )

    .with(
      { property: ColumnViewType.ContactsPrimaryEmail },
      (filter) => (row: ContactStore) => {
        const filterValues = filter?.value;

        if (!filter.active) return true;

        if (!filterValues) return true;
        const emails = row.value.primaryEmail?.email;

        return filterTypeText(filter, emails ?? undefined);
      },
    )

    .with(
      { property: ColumnViewType.ContactsPhoneNumbers },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const value = row.value?.phoneNumbers?.map((p) => p.rawPhoneNumber);

        if (!String(value).length) {
          return ComparisonOperator.IsEmpty === filter.operation;
        }

        return filterTypeText(filter, String(value) ?? undefined);
      },
    )

    .with(
      { property: ColumnViewType.ContactsLinkedin },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;

        const linkedInUrl = row.value.socials?.find(
          (v: { id: string; url: string }) => v.url.includes('linkedin'),
        )?.url;

        return filterTypeText(filter, linkedInUrl);
      },
    )

    .with(
      { property: ColumnViewType.ContactsCity },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const cities = row.value.locations?.map((l) => l?.locality);

        if (!cities)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(
          filter,
          cities?.some((j) => j) ? (cities as string[]) : [],
        );
      },
    )

    .with(
      { property: ColumnViewType.ContactsPersona },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const tags = row.value.tags?.map((l: Tag) => l.id);

        if (!tags)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, tags);
      },
    )

    .with(
      { property: ColumnViewType.ContactsConnections },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const users = row.value.connectedUsers?.map(
          (l: User) => row.root.users.value.get(l.id)?.name,
        );

        if (!users.length)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, users as string[]);
      },
    )

    .with(
      { property: ColumnViewType.ContactsLinkedinFollowerCount },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;

        const followers = row.value?.socials?.find((e: Social) =>
          e?.url?.includes('linkedin'),
        )?.followersCount;

        return filterTypeNumber(filter, followers);
      },
    )

    .with(
      { property: ColumnViewType.ContactsJobTitle },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const jobTitles =
          row.value?.latestOrganizationWithJobRole?.jobRole.jobTitle;

        if (!jobTitles) return false;

        return filterTypeText(filter, jobTitles);
      },
    )

    .with(
      { property: ColumnViewType.ContactsCountry },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;

        const countries = row.value.locations?.map((l) => l.countryCodeA2);

        if (!countries)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, countries as string[]);
      },
    )
    .with({ property: ColumnViewType.ContactsRegion }, (filter) => {
      if (!filter.active) return () => true;

      return (row: ContactStore) => {
        const locations = row.value.locations;
        const region = locations?.[0]?.region;

        if (!region)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, region ? [region] : []);
      };
    })

    .with({ property: ColumnViewType.ContactsFlows }, (filter) => {
      if (!filter.active) return () => true;

      return (row: ContactStore) => {
        const values = row.flows?.map((e) => e?.value?.metadata.id);

        if ((values ?? []).length === 0)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, Array.isArray(values) ? values : []);
      };
    })

    .with({ property: ColumnViewType.ContactsTimeInCurrentRole }, (filter) => {
      if (!filter.active) return () => true;

      return (row: ContactStore) => {
        const timeInCurrentRole =
          row.value?.latestOrganizationWithJobRole?.jobRole.startedAt;

        return filterTypeDate(filter, timeInCurrentRole);
      };
    })

    .with(
      { property: 'EMAIL_VERIFICATION_PRIMARY_EMAIL' },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;

        const filterValues = filter.value;
        const emailValidationData =
          row.value.primaryEmail?.emailValidationDetails;

        if (emailValidationData === undefined) return false;

        return match(filter.operation)
          .with(ComparisonOperator.Contains, () =>
            filterValues?.some(
              (categoryFilter: { value: string; category: string }) =>
                (categoryFilter.category === 'DELIVERABLE' &&
                  isDeliverableV2(categoryFilter.value, emailValidationData)) ||
                (categoryFilter.category === 'UNDELIVERABLE' &&
                  isNotDeliverableV2(
                    categoryFilter?.value,
                    emailValidationData,
                  )) ||
                (categoryFilter.category === 'UNKNOWN' &&
                  isDeliverableUnknownV2(
                    categoryFilter.value,
                    emailValidationData,
                  )),
            ),
          )

          .with(ComparisonOperator.NotContains, () =>
            filterValues.some(
              (categoryFilter: { value: string; category: string }) =>
                !(
                  categoryFilter.category === 'DELIVERABLE' &&
                  isDeliverableV2(categoryFilter.value, emailValidationData)
                ) &&
                !(
                  categoryFilter.category === 'UNDELIVERABLE' &&
                  isNotDeliverableV2(categoryFilter.value, emailValidationData)
                ) &&
                !(
                  categoryFilter.category === 'UNKNOWN' &&
                  isDeliverableUnknownV2(
                    categoryFilter.value,
                    emailValidationData,
                  )
                ),
            ),
          )
          .with(
            ComparisonOperator.IsEmpty,
            () =>
              !emailValidationData ||
              Object.keys(emailValidationData).length === 0,
          )
          .with(
            ComparisonOperator.IsNotEmpty,
            () =>
              !!emailValidationData &&
              Object.keys(emailValidationData).length > 1,
          )
          .otherwise(() => true);
      },
    )

    .with(
      { property: ColumnViewType.ContactsFlowStatus },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;

        if (!row.hasFlows)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return row.flows?.some((e) =>
          filterTypeList(filter, [e?.value?.status]),
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

export const getContactDefaultFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];
  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};

export const getContactFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];
  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};

function isNotDeliverableV2(
  statuses: string,
  data: EmailValidationDetails,
): boolean {
  if (!statuses?.length && data?.deliverable && data?.verified) return true;

  const statusChecks: Record<string, () => boolean> = {
    [EmailVerificationStatus.InvalidMailbox]: () =>
      !data.canConnectSmtp || data.deliverable !== EmailDeliverable.Deliverable,

    [EmailVerificationStatus.MailboxFull]: () => !!data?.isMailboxFull,
    [EmailVerificationStatus.IncorrectFormat]: () => !data.isValidSyntax,
  };

  return statusChecks[statuses]?.() ?? false;
}

function isDeliverableUnknownV2(
  statuses: string,
  data: EmailValidationDetails,
): boolean {
  if (
    !statuses?.length &&
    (!data.verified || data.isCatchAll || data.verifyingCheckAll)
  ) {
    return true;
  }

  const statusChecks: Record<string, () => boolean> = {
    [EmailVerificationStatus.CatchAll]: () =>
      data?.deliverable === EmailDeliverable.Unknown &&
      !!data?.isCatchAll &&
      !!data?.verified,
    [EmailVerificationStatus.NotVerified]: () => !data.verified,
    [EmailVerificationStatus.VerificationInProgress]: () =>
      data.verifyingCheckAll,
  };

  return statusChecks[status]?.() ?? false;
}

function isDeliverableV2(
  statuses: string,
  data: EmailValidationDetails,
): boolean {
  if (data?.deliverable !== EmailDeliverable.Deliverable || !data.verified)
    return false;

  const statusChecks: Record<string, () => boolean> = {
    [EmailVerificationStatus.NoRisk]: () => !data.isRisky,
    [EmailVerificationStatus.FirewallProtected]: () => !!data.isFirewalled,
    [EmailVerificationStatus.FreeAccount]: () => !!data.isFreeAccount,
    [EmailVerificationStatus.GroupMailbox]: () => data.verifyingCheckAll,
  };

  return statusChecks?.[statuses]?.() ?? false;
}
