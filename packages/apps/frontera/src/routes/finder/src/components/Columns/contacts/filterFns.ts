import { match } from 'ts-pattern';
import { FilterItem } from '@store/types.ts';
import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { EmailVerificationStatus } from '@finder/components/Columns/contacts/Filters/Email/utils.ts';

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
        const filterValue = filter?.value;

        if (!filterValue && filter.active && !filter.includeEmpty) return true;
        if (!row.value?.name?.length && filter.includeEmpty) return true;
        if (!filterValue || !row.value?.name?.length) return false;

        return row.value.name
          ?.toLowerCase()
          .includes(filterValue?.toLowerCase());
      },
    )
    .with(
      { property: ColumnViewType.ContactsOrganization },
      (filter) => (row: ContactStore) => {
        const filterValues = filter?.value;

        if (!filter.active) return true;
        const orgs = row.value?.organizations?.content?.map((o) =>
          o.name.toLowerCase().trim(),
        );

        if (filter.includeEmpty && orgs?.every((org) => !org.length)) {
          return true;
        }

        if (filter.includeEmpty && filterValues.length === 0) {
          return false;
        }

        return orgs?.some((e) => e.includes(filterValues));
      },
    )
    .with(
      { property: ColumnViewType.ContactsEmails },
      (filter) => (row: ContactStore) => {
        const filterValues = filter?.value;

        if (!filter.active) return true;

        if (!filterValues) return true;
        const emails = row.value?.emails?.map((e) => e.email);

        return emails?.some((e) => e?.includes(filterValues));
      },
    )

    .with(
      { property: ColumnViewType.ContactsPhoneNumbers },
      (filter) => (row: ContactStore) => {
        const filterValue = filter?.value;

        if (!filter.active) return true;

        if (!filterValue && filter.active && !filter.includeEmpty) return true;
        if (!row.value?.phoneNumbers?.length && filter.includeEmpty)
          return true;
        if (!filterValue || !row.value?.phoneNumbers?.length) return false;

        return row.value?.phoneNumbers?.[0]?.e164?.includes(filterValue);
      },
    )

    .with(
      { property: ColumnViewType.ContactsLinkedin },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const linkedInUrl = row.value.socials?.find(
          (v: { id: string; url: string }) => v.url.includes('linkedin'),
        )?.url;

        if (!filterValue && filter.active && !filter.includeEmpty) return true;
        if (!linkedInUrl?.length && filter.includeEmpty) return true;

        if (!filterValue || !linkedInUrl?.[0] || filter.includeEmpty) {
          return false;
        }

        return linkedInUrl.includes(filterValue);
      },
    )
    .with(
      { property: ColumnViewType.ContactsCity },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const cities = row.value.locations?.map((l) => l?.locality);

        if (filterValue.length === 0 && filter.active && !filter.includeEmpty)
          return true;

        if (!cities.length && filter.includeEmpty) return true;

        if (filterValue.length === 0 && filter.active && filter.includeEmpty) {
          return cities.forEach((j) => !j);
        }

        return cities.some((c) =>
          filterValue.map((f: string) => f).includes(c),
        );
      },
    )
    .with(
      { property: ColumnViewType.ContactsPersona },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const tags = row.value.tags?.map((l: Tag) => l.name);

        if (!filter.value?.length) return true;

        if (!tags?.length) return false;

        return filter.value.some((f: string) => tags.includes(f));
      },
    )
    .with(
      { property: ColumnViewType.ContactsConnections },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const users = row.value.connectedUsers?.map(
          (l: User) => row.root.users.value.get(l.id)?.name,
        );

        if (!filter.value?.length) return true;

        if (!users?.length) return false;

        return filter.value.some((f: string) => users.includes(f));
      },
    )
    .with(
      { property: ColumnViewType.ContactsLinkedinFollowerCount },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const operator = filter.operation;
        const followers = row.value?.socials?.find((e: Social) =>
          e?.url?.includes('linkedin'),
        )?.followersCount;

        if (operator === ComparisonOperator.Lt) {
          return Number(followers) < Number(filterValue);
        }

        if (operator === ComparisonOperator.Gt) {
          return Number(followers) > Number(filterValue);
        }

        if (operator === ComparisonOperator.Between) {
          const filterValue = filter?.value?.map(Number) as number[];

          return (
            followers >= Number(filterValue[0]) &&
            followers <= Number(filterValue[1])
          );
        }
      },
    )
    .with(
      { property: ColumnViewType.ContactsJobTitle },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const jobTitles = row.value?.jobRoles?.map((j) => j.jobTitle) || [];

        if (!filterValue && filter.active && !filter.includeEmpty) return true;
        if (!jobTitles.length && filter.includeEmpty) return true;
        if (!filterValue && filter.active && filter.includeEmpty)
          return jobTitles.some((j) => !j);

        return jobTitles.some((j) => j?.toLowerCase().includes(filterValue));
      },
    )

    .with(
      { property: ColumnViewType.ContactsCountry },
      (filter) => (row: ContactStore) => {
        const filterValue = filter?.value;

        if (!filter.active) return true;

        const countries = row.value.locations?.map((l) => l.countryCodeA2);

        if (filterValue.length === 0 && filter.active && !filter.includeEmpty)
          return true;

        if (!countries.length && filter.includeEmpty) return true;

        if (filterValue.length === 0 && filter.active && filter.includeEmpty) {
          return countries.some((j) => !j);
        }

        return countries.some((c) =>
          filterValue.map((f: string) => f).includes(c),
        );
      },
    )
    .with({ property: ColumnViewType.ContactsRegion }, (filter) => {
      // Early exit if filter is not active
      if (!filter.active) return () => true;

      const filterValue = filter.value;
      const includeEmpty = filter.includeEmpty;

      return (row: ContactStore) => {
        const locations = row.value.locations;
        const region = locations?.[0]?.region;

        // Check for empty cases
        if (!region) {
          return includeEmpty;
        }

        // If filterValue is empty, return based on includeEmpty
        if (!filterValue.length) {
          return !includeEmpty;
        }

        // Check if country is in filterValue
        return filterValue.includes(region);
      };
    })

    .with({ property: ColumnViewType.ContactsFlows }, (filter) => {
      if (!filter.active) return () => true;

      const filterValue = filter.value;
      const includeEmpty = filter.includeEmpty;

      return (row: ContactStore) => {
        const flow = row.flow;

        if (!flow) {
          return includeEmpty;
        }

        if (!filterValue.length) {
          return !includeEmpty;
        }

        return filterValue.includes(flow.value.name);
      };
    })
    .with(
      { property: 'EMAIL_VERIFICATION' },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;

        const filterValues = filter.value;
        const email = row.value?.emails?.[0];
        const emailValidationData = email?.emailValidationDetails;

        if (!emailValidationData) return false;

        return filterValues.some(
          (categoryFilter: {
            category: string;
            values: EmailVerificationStatus[];
          }) =>
            (categoryFilter.category === EmailDeliverable.Deliverable &&
              isDeliverable(categoryFilter.values, emailValidationData)) ||
            (categoryFilter.category === EmailDeliverable.Undeliverable &&
              isNotDeliverable(categoryFilter.values, emailValidationData)) ||
            (categoryFilter.category === EmailDeliverable.Unknown &&
              isDeliverableUnknown(categoryFilter.values, emailValidationData)),
        );
      },
    )

    .otherwise(() => noop);
};

export const getContactFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];
  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};

function isNotDeliverable(
  statuses: EmailVerificationStatus[],
  data: EmailValidationDetails,
): boolean {
  if (data?.deliverable !== EmailDeliverable.Undeliverable || !data.verified)
    return false;

  if (!statuses.length && data.deliverable && data.verified) return true;

  const statusChecks: Record<string, () => boolean> = {
    [EmailVerificationStatus.InvalidMailbox]: () => !data.canConnectSmtp,
    [EmailVerificationStatus.MailboxFull]: () => !!data?.isMailboxFull,
    [EmailVerificationStatus.IncorrectFormat]: () => !data.isValidSyntax,
  };

  return statuses.some((status) => statusChecks[status]?.() ?? false);
}

function isDeliverableUnknown(
  statuses: EmailVerificationStatus[],
  data: EmailValidationDetails,
): boolean {
  if (
    !statuses.length &&
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

  return statuses.some((status) => statusChecks[status]?.() ?? false);
}

function isDeliverable(
  statuses: EmailVerificationStatus[],
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

  return statuses.some((status) => statusChecks[status]?.() ?? false);
}