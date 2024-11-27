import { match } from 'ts-pattern';
import countries from '@assets/countries/countries.json';
import { ContactStore } from '@store/Contacts/Contact.store';

import {
  User,
  Social,
  JobRole,
  ColumnViewType,
  FlowParticipantStatus,
} from '@graphql/types';

export const getContactSortFn = (columnId: string, flowId?: string) =>
  match(columnId)
    .with(ColumnViewType.ContactsName, () => (row: ContactStore) => {
      return row.name?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.ContactsOrganization,
      () => (row: ContactStore) =>
        row.value?.latestOrganizationWithJobRole?.organization?.name.toLowerCase() ||
        null,
    )
    .with(
      ColumnViewType.ContactsCity,
      () => (row: ContactStore) =>
        row.value?.locations?.[0]?.locality?.trim().toLowerCase() || null,
    )
    .with(ColumnViewType.ContactsPersona, () => (row: ContactStore) => {
      return row.value?.tags?.[0]?.name?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.ContactsLinkedinFollowerCount,
      () => (row: ContactStore) => {
        return row.value.socials.find((e: Social) =>
          e?.url?.includes('linkedin'),
        )?.followersCount;
      },
    )
    .with(ColumnViewType.ContactsJobTitle, () => (row: ContactStore) => {
      return row.value.jobRoles
        .find((e: JobRole) => e.endedAt === null)
        ?.jobTitle?.toLowerCase();
    })
    .with(ColumnViewType.ContactsCountry, () => (row: ContactStore) => {
      const countryName = countries.find(
        (d) =>
          d.alpha2 === row.value.locations?.[0]?.countryCodeA2?.toLowerCase(),
      );

      return countryName?.name?.toLowerCase() || null;
    })
    .with(ColumnViewType.ContactsConnections, () => (row: ContactStore) => {
      return row.value.connectedUsers
        ?.map((l: User) => row.root.users.value.get(l.id)?.name)
        .filter(Boolean)
        .sort((a, b) => (a && b ? a?.localeCompare(b) : -1));
    })

    .with(ColumnViewType.ContactsRegion, () => (row: ContactStore) => {
      return row.value.locations?.[0]?.region?.toLowerCase() || null;
    })
    .with(ColumnViewType.ContactsFlows, () => (row: ContactStore) => {
      return row.flows?.length || null;
    })
    .with(ColumnViewType.ContactsCreatedAt, () => (row: ContactStore) => {
      return row.value.createdAt;
    })
    .with(ColumnViewType.ContactsUpdatedAt, () => (row: ContactStore) => {
      return row.value.updatedAt;
    })
    .with(ColumnViewType.ContactsFlowStatus, () => (row: ContactStore) => {
      if (!flowId) return false;

      const status = row.root.flows.value
        .get(flowId)
        ?.value.participants.find((e) => e.entityId === row.id)?.status;

      const statusOrder = {
        [FlowParticipantStatus.OnHold]: 1,
        [FlowParticipantStatus.Completed]: 2,
        [FlowParticipantStatus.Error]: 3,
        [FlowParticipantStatus.GoalAchieved]: 4,
        [FlowParticipantStatus.InProgress]: 5,
        [FlowParticipantStatus.Ready]: 6,
        [FlowParticipantStatus.Scheduled]: 7,
      };

      return status && statusOrder?.[status];
    })
    .with(ColumnViewType.ContactsFlowNextAction, () => (row: ContactStore) => {
      if (!flowId) return false;

      return row.root.flows.value
        .get(flowId)
        ?.value.participants.find((e) => e.entityId === row.id)
        ?.executions?.find((e) => e.scheduledAt && !e.executedAt)?.scheduledAt;
    })

    .otherwise(() => (_row: ContactStore) => false);
