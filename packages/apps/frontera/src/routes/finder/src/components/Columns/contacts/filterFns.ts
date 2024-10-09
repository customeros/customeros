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
        const values = row.value.name;

        if (!values) return false;

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
      { property: ColumnViewType.ContactsEmails },
      (filter) => (row: ContactStore) => {
        const filterValues = filter?.value;

        if (!filter.active) return true;

        if (!filterValues) return true;
        const emails = row.value?.emails
          .filter((e) => e.work)
          .map((e) => e.email);

        return filterTypeText(filter, emails?.join(' '));
      },
    )

    .with(
      { property: ColumnViewType.ContactsPersonalEmails },
      (filter) => (row: ContactStore) => {
        const filterValues = filter?.value;

        if (!filter.active) return true;

        if (!filterValues) return true;
        const emails = row.value?.emails
          .filter((e) => !e.work)
          .map((e) => e.email);

        if (!emails) return false;

        return filterTypeText(filter, emails?.join(' '));
      },
    )

    .with(
      { property: ColumnViewType.ContactsPhoneNumbers },
      (filter) => (row: ContactStore) => {
        const value = row.value?.phoneNumbers
          ?.map((p) => p.rawPhoneNumber)
          ?.join(' ');

        if (!filter.active) return true;
        if (value.length < 1) return false;

        return filterTypeText(filter, value ?? undefined);
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

        return filterTypeList(
          filter,
          cities?.filter((city) => city !== undefined) as string[],
        );
      },
    )
    .with(
      { property: ColumnViewType.ContactsPersona },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const tags = row.value.tags?.map((l: Tag) => l.id);

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
        const jobTitles = row.value?.jobRoles?.map((j) => j.jobTitle) || [];

        return filterTypeText(filter, jobTitles.join(' '));
      },
    )

    .with(
      { property: ColumnViewType.ContactsCountry },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;

        const countries = row.value.locations?.map((l) => l.countryCodeA2);

        return filterTypeList(filter, countries as string[]);
      },
    )
    .with({ property: ColumnViewType.ContactsRegion }, (filter) => {
      if (!filter.active) return () => true;

      return (row: ContactStore) => {
        const locations = row.value.locations;
        const region = locations?.[0]?.region;

        return filterTypeList(filter, region ? [region] : []);
      };
    })

    .with({ property: ColumnViewType.ContactsFlows }, (filter) => {
      if (!filter.active) return () => true;

      return (row: ContactStore) => {
        const flow = row.flow;

        return filterTypeText(filter, flow?.value?.name);
      };
    })
    .with(
      { property: 'EMAIL_VERIFICATION_WORK_EMAIL' },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;

        const filterValues = filter.value;
        const email = row.value?.emails?.find((e) => e.work === true);
        const emailValidationData = email?.emailValidationDetails;

        if (emailValidationData === undefined) return false;

        return match(filter.operation)
          .with(ComparisonOperator.Contains, () =>
            filterValues?.some(
              (categoryFilter: { value: string; category: string }) =>
                (categoryFilter.category === 'DELIVERABLE' &&
                  isDeliverable(categoryFilter.value, emailValidationData)) ||
                (categoryFilter.category === 'UNDELIVERABLE' &&
                  isNotDeliverable(
                    categoryFilter?.value,
                    emailValidationData,
                  )) ||
                (categoryFilter.category === 'UNKNOWN' &&
                  isDeliverableUnknown(
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
                  isDeliverable(categoryFilter.value, emailValidationData)
                ) &&
                !(
                  categoryFilter.category === 'UNDELIVERABLE' &&
                  isNotDeliverable(categoryFilter.value, emailValidationData)
                ) &&
                !(
                  categoryFilter.category === 'UNKNOWN' &&
                  isDeliverableUnknown(
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
      { property: 'EMAIL_VERIFICATION_PERSONAL_EMAIL' },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;

        const filterValues = filter.value;
        const email = row.value?.emails?.find((e) => e.work === false);
        const emailValidationData = email?.emailValidationDetails;

        if (emailValidationData === undefined) return false;

        return match(filter.operation)
          .with(ComparisonOperator.Contains, () =>
            filterValues?.some(
              (categoryFilter: { value: string; category: string }) =>
                (categoryFilter.category === 'DELIVERABLE' &&
                  isDeliverable(categoryFilter.value, emailValidationData)) ||
                (categoryFilter.category === 'UNDELIVERABLE' &&
                  isNotDeliverable(
                    categoryFilter?.value,
                    emailValidationData,
                  )) ||
                (categoryFilter.category === 'UNKNOWN' &&
                  isDeliverableUnknown(
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
                  isDeliverable(categoryFilter.value, emailValidationData)
                ) &&
                !(
                  categoryFilter.category === 'UNDELIVERABLE' &&
                  isNotDeliverable(categoryFilter.value, emailValidationData)
                ) &&
                !(
                  categoryFilter.category === 'UNKNOWN' &&
                  isDeliverableUnknown(
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

        const filterValues = filter.value;
        const currentFlowStatus = row.flowContact?.value?.status;

        if (!currentFlowStatus) return false;

        return filterValues.includes(currentFlowStatus);
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

export const getContactFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];
  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};

function isNotDeliverable(
  statuses: string,
  data: EmailValidationDetails,
): boolean {
  if (data?.deliverable !== EmailDeliverable.Undeliverable || !data.verified)
    return false;

  if (!statuses?.length && data?.deliverable && data?.verified) return true;

  const statusChecks: Record<string, () => boolean> = {
    [EmailVerificationStatus.InvalidMailbox]: () => !data.canConnectSmtp,
    [EmailVerificationStatus.MailboxFull]: () => !!data?.isMailboxFull,
    [EmailVerificationStatus.IncorrectFormat]: () => !data.isValidSyntax,
  };

  return statusChecks[statuses]?.() ?? false;
}

function isDeliverableUnknown(
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

function isDeliverable(
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
