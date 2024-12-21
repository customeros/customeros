import { match } from 'ts-pattern';
import countries from '@assets/countries/countries.json';

import { ColumnViewType, FlowParticipantStatus } from '@graphql/types';

import type { Contact } from '../Contact.dto';

export const getContactSortFn = (columnId: string, flowId?: string) =>
  match(columnId)
    .with(ColumnViewType.ContactsName, () => (row: Contact) => {
      return row.name?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.ContactsOrganization,
      () => (row: Contact) =>
        row.value?.primaryOrganizationName?.toLowerCase() || null,
    )
    .with(
      ColumnViewType.ContactsCity,
      () => (row: Contact) =>
        row.value?.locations?.[0]?.locality?.trim().toLowerCase() || null,
    )
    .with(ColumnViewType.ContactsPersona, () => (row: Contact) => {
      return row.value?.tags?.[0]?.name?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.ContactsLinkedinFollowerCount,
      () => (row: Contact) => {
        return row.value.linkedInFollowerCount;
      },
    )
    .with(ColumnViewType.ContactsJobTitle, () => (row: Contact) => {
      return row.value.primaryOrganizationJobRoleTitle;
    })
    .with(ColumnViewType.ContactsCountry, () => (row: Contact) => {
      const countryName = countries.find(
        (d) =>
          d.alpha2 === row.value.locations?.[0]?.countryCodeA2?.toLowerCase(),
      );

      return countryName?.name?.toLowerCase() || null;
    })
    .with(ColumnViewType.ContactsConnections, () => (row: Contact) => {
      return (
        row.value.connectedUsers
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          ?.map((l: any) => row.store.root.users.value.get(l.id)?.name)
          .filter(Boolean)
          .sort((a, b) => (a && b ? a?.localeCompare(b) : -1))
      );
    })

    .with(ColumnViewType.ContactsRegion, () => (row: Contact) => {
      return row.value.locations?.[0]?.region?.toLowerCase() || null;
    })
    .with(ColumnViewType.ContactsFlows, () => (row: Contact) => {
      return row.flows?.length || null;
    })
    .with(ColumnViewType.ContactsCreatedAt, () => (row: Contact) => {
      return row.value.createdAt;
    })
    .with(ColumnViewType.ContactsUpdatedAt, () => (row: Contact) => {
      return row.value.updatedAt;
    })
    .with(ColumnViewType.ContactsFlowStatus, () => (row: Contact) => {
      if (!flowId) return false;

      const status = row.store.root.flows.value
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

      return status && statusOrder[status as keyof typeof statusOrder];
    })
    .with(ColumnViewType.ContactsFlowNextAction, () => (row: Contact) => {
      if (!flowId) return false;

      return row.store.root.flows.value
        .get(flowId)
        ?.value.participants.find((e) => e.entityId === row.id)
        ?.executions?.find((e) => e.scheduledAt && !e.executedAt)?.scheduledAt;
    })

    .otherwise(() => (_row: Contact) => false);
