import { match } from 'ts-pattern';
import { FilterItem } from '@store/types.ts';
import { isAfter, isBefore } from 'date-fns';
import { EmailVerificationStatus } from '@finder/components/Columns/contacts/filterTypes';

import {
  Filter,
  ColumnViewType,
  EmailDeliverable,
  ComparisonOperator,
  EmailValidationDetails,
} from '@graphql/types';

import type { Contact } from '../Contact.dto';

const getFilterFn = (
  filter: FilterItem | undefined | null,
  flowId?: string,
) => {
  const noop = (_row: Contact) => true;

  if (!filter) return noop;

  return (
    match(filter)
      .with({ property: 'STAGE' }, (filter) => (row: Contact) => {
        const filterValues = filter?.value;

        if (!filterValues || !row.value?.primaryOrganizationName) {
          return false;
        }
        const hasOrgWithMatchingStage =
          row.store.root.organizations.value.forEach((o) => {
            const stage = row.store.root?.organizations?.value.get(o.id)?.value
              ?.stage;

            return filterValues.includes(stage);
          });

        return hasOrgWithMatchingStage;
      })

      .with({ property: 'RELATIONSHIP' }, (filter) => (row: Contact) => {
        const filterValues = filter?.value;

        if (!filterValues || !row.value?.primaryOrganizationName) {
          return false;
        }
        const hasOrgWithMatchingRelationship =
          row.store.root?.organizations.value.forEach((o) => {
            const stage = row.store.root?.organizations?.value.get(o.id)?.value
              ?.relationship;

            return filterValues.includes(stage);
          });

        return hasOrgWithMatchingRelationship;
      })

      .with(
        { property: ColumnViewType.ContactsName },
        (filter) => (row: Contact) => {
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
        (filter) => (row: Contact) => {
          if (!filter.active) return true;
          const orgs = row.value.primaryOrganizationName;

          return filterTypeText(filter, orgs);
        },
      )

      .with(
        { property: ColumnViewType.ContactsPrimaryEmail },
        (filter) => (row: Contact) => {
          if (!filter.active) return true;

          const emails = row.value.emails
            .filter((e) => e.primary)
            .map((e) => e.email);

          return filterTypeText(filter, emails.join('; '));
        },
      )

      // .with(
      //   { property: ColumnViewType.ContactsPhoneNumbers },
      //   (filter) => (row: Contact) => {
      //     if (!filter.active) return true;
      //     const value = row.value?.phoneNumbers?.map((p) => p.rawPhoneNumber);

      //     if (!String(value).length) {
      //       return ComparisonOperator.IsEmpty === filter.operation;
      //     }

      //     return filterTypeText(filter, String(value) ?? undefined);
      //   },
      // )

      .with(
        { property: ColumnViewType.ContactsLinkedin },
        (filter) => (row: Contact) => {
          if (!filter.active) return true;

          const linkedInUrl = row.value.linkedInUrl;

          return filterTypeText(filter, linkedInUrl);
        },
      )

      .with(
        { property: ColumnViewType.ContactsCity },
        (filter) => (row: Contact) => {
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
        (filter) => (row: Contact) => {
          if (!filter.active) return true;
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          const tags = row.value.tags?.map((l: any) => l.metadata.id);

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
        (filter) => (row: Contact) => {
          if (!filter.active) return true;
          const users = Array.from(row.value.connectedUsers.values()).map(
            (u) => row.store.root.users.value.get(u)?.id ?? '',
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
        (filter) => (row: Contact) => {
          if (!filter.active) return true;

          const followers = row.value?.linkedInFollowerCount;

          return filterTypeNumber(filter, followers);
        },
      )

      .with(
        { property: ColumnViewType.ContactsJobTitle },
        (filter) => (row: Contact) => {
          if (!filter.active) return true;
          const jobTitles = row.value?.primaryOrganizationJobRoleTitle;

          return filterTypeText(filter, jobTitles);
        },
      )

      .with(
        { property: ColumnViewType.ContactsCountry },
        (filter) => (row: Contact) => {
          if (!filter.active) return true;
          const locations = row.value.locations;
          const country = locations?.[0]?.countryCodeA2;

          if (!country)
            return (
              filter.operation === ComparisonOperator.IsEmpty ||
              filter.operation === ComparisonOperator.NotContains
            );

          return filterTypeList(
            filter,
            Array.isArray(country) ? country : [country],
          );
        },
      )
      .with({ property: ColumnViewType.ContactsRegion }, (filter) => {
        if (!filter.active) return () => true;

        return (row: Contact) => {
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

        return (row: Contact) => {
          const values = row.flows?.map((e) => e?.value?.metadata.id);

          if ((values ?? []).length === 0)
            return (
              filter.operation === ComparisonOperator.IsEmpty ||
              filter.operation === ComparisonOperator.NotContains
            );

          return filterTypeList(filter, Array.isArray(values) ? values : []);
        };
      })

      .with(
        { property: ColumnViewType.ContactsTimeInCurrentRole },
        (filter) => {
          if (!filter.active) return () => true;

          return (row: Contact) => {
            const timeInCurrentRole =
              row.value?.primaryOrganizationJobRoleStartDate;

            return filterTypeDate(filter, timeInCurrentRole);
          };
        },
      )

      .with(
        { property: 'EMAIL_VERIFICATION_PRIMARY_EMAIL' },
        (filter) => (row: Contact) => {
          if (!filter.active) return true;

          const filterValues = filter.value;
          const emailValidationData = row.value.emails?.find(
            (e) => e.primary,
          )?.emailValidationDetails;

          if (emailValidationData === undefined) return false;

          return match(filter.operation)
            .with(ComparisonOperator.In, () =>
              filterValues?.some(
                (id: string) =>
                  isDeliverableV2(id, emailValidationData) ||
                  isNotDeliverableV2(id, emailValidationData) ||
                  isDeliverableUnknownV2(id, emailValidationData),
              ),
            )
            .with(ComparisonOperator.NotIn, () =>
              filterValues?.every(
                (id: string) =>
                  !isDeliverableV2(id, emailValidationData) &&
                  !isNotDeliverableV2(id, emailValidationData) &&
                  !isDeliverableUnknownV2(id, emailValidationData),
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
        (filter) => (row: Contact) => {
          if (!filter.active) return true;

          if (!row.hasFlows)
            return (
              filter.operation === ComparisonOperator.IsEmpty ||
              filter.operation === ComparisonOperator.NotContains
            );
          if (!flowId) return false;

          const participantStatus = row.store.root.flows.value
            .get(flowId)
            ?.value.participants.find((e) => e.entityId === row.id)?.status;

          if (!participantStatus) return false;

          return filterTypeList(filter, [participantStatus]);
        },
      )
      .with({ property: 'FLOW_ID' }, (filter) => (row: Contact) => {
        if (!filter.active || !filter.value[0]) return true;

        return row?.flowsIds.includes(filter.value[0]);
      })

      .otherwise(() => noop)
  );
};

const filterTypeText = (
  filter: FilterItem,
  value: string | undefined | null,
) => {
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
    .with(ComparisonOperator.Equals, () => value === Number(filterValue))
    .with(ComparisonOperator.NotEquals, () => value !== Number(filterValue))
    .otherwise(() => true);
};

const filterTypeList = (filter: FilterItem, value: string[] | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  return match(filterOperator)
    .with(ComparisonOperator.IsEmpty, () => !value?.length)
    .with(ComparisonOperator.IsNotEmpty, () => value?.length)
    .with(
      ComparisonOperator.NotIn,
      () =>
        !value?.length ||
        (value?.length && !value.some((v) => filterValue?.includes(v))),
    )
    .with(ComparisonOperator.In, () => {
      return value?.length && value?.some((v) => filterValue?.includes(v));
    })

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

export const getContactFilterFns = (
  filters: Filter | null,
  flowId?: string,
) => {
  if (!filters || !filters.AND) return [];
  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter, flowId));
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

  return statusChecks[statuses]?.() ?? false;
}

function isDeliverableV2(
  statuses: string,
  data: EmailValidationDetails,
): boolean {
  if (data?.deliverable !== EmailDeliverable.Deliverable || !data.verified)
    return false;

  const statusChecks: Record<string, () => boolean> = {
    [EmailVerificationStatus.NoRisk]: () =>
      data?.deliverable === EmailDeliverable.Deliverable, // validation rules https://www.figma.com/design/uWolaNIV9vDhQ5PfQGNC7o/Views?node-id=3127-7449&p=f&t=at3ByiH8XZGopHZV-0
    [EmailVerificationStatus.FirewallProtected]: () => !!data.isFirewalled,
    [EmailVerificationStatus.FreeAccount]: () => !!data.isFreeAccount,
    [EmailVerificationStatus.GroupMailbox]: () => data.verifyingCheckAll,
  };

  return statusChecks?.[statuses]?.() ?? false;
}
