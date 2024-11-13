import { Tag01 } from '@ui/media/icons/Tag01';
import { Users03 } from '@ui/media/icons/Users03';
import { Building07 } from '@ui/media/icons/Building07';
import {
  ColumnViewType,
  ComparisonOperator,
  FlowParticipantStatus,
} from '@shared/types/__generated__/graphql.types';

export type FilterType = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  options?: any[];
  icon: JSX.Element;
  filterName: string;
  filterOperators: ComparisonOperator[];
  filterType: 'text' | 'date' | 'number' | 'list';
  filterAccesor: ColumnViewType | 'EMAIL_VERIFICATION_PRIMARY_EMAIL';
  groupOptions?: { label: string; options: { id: string; label: string }[] }[];
};

import { uniqBy } from 'lodash';
import { type RootStore } from '@store/root';

import { Phone } from '@ui/media/icons/Phone';
import { User03 } from '@ui/media/icons/User03';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Globe05 } from '@ui/media/icons/Globe05';
import { Globe06 } from '@ui/media/icons/Globe06';
import { Globe04 } from '@ui/media/icons/Globe04';
import { Shuffle01 } from '@ui/media/icons/Shuffle01';
import { ClockCheck } from '@ui/media/icons/ClockCheck';
import { Certificate01 } from '@ui/media/icons/Certificate01';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';

export const getFilterTypes = (store?: RootStore) => {
  const filterTypes: Partial<
    Record<ColumnViewType | 'EMAIL_VERIFICATION_PRIMARY_EMAIL', FilterType>
  > = {
    [ColumnViewType.ContactsName]: {
      filterType: 'text',
      filterName: 'Contact name',
      filterAccesor: ColumnViewType.ContactsName,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <User03 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.ContactsOrganization]: {
      filterType: 'text',
      filterName: 'Organization',
      filterAccesor: ColumnViewType.ContactsOrganization,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Building07 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },

    [ColumnViewType.ContactsPrimaryEmail]: {
      filterType: 'text',
      filterName: 'Primary email',
      filterAccesor: ColumnViewType.ContactsPrimaryEmail,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Mail01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    ['EMAIL_VERIFICATION_PRIMARY_EMAIL']: {
      filterType: 'list',
      filterName: 'Primary email status',
      filterAccesor: 'EMAIL_VERIFICATION_PRIMARY_EMAIL',
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Mail01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [],
      groupOptions: [
        {
          label: 'Deliverable',
          options: [
            {
              id: EmailVerificationStatus.FirewallProtected,
              label: 'Firewall protected',
            },
            {
              id: EmailVerificationStatus.FreeAccount,
              label: 'Free account',
            },
            {
              id: EmailVerificationStatus.NoRisk,
              label: 'No risk',
            },
          ],
        },
        {
          label: 'Not deliverable',
          options: [
            {
              id: EmailVerificationStatus.IncorrectFormat,
              label: 'Incorrect email format',
            },
            {
              id: EmailVerificationStatus.InvalidMailbox,
              label: 'Mailbox doesn’t exist',
            },
            {
              id: EmailVerificationStatus.MailboxFull,
              label: 'Mailbox full',
            },
          ],
        },
        {
          label: "Don't know",
          options: [
            {
              id: EmailVerificationStatus.CatchAll,
              label: 'Catch all',
            },
            {
              id: EmailVerificationStatus.NotVerified,
              label: 'Not verified yet',
            },
            {
              id: EmailVerificationStatus.VerificationInProgress,
              label: 'Verification in progress',
            },
          ],
        },
      ],
    },
    [ColumnViewType.ContactsPhoneNumbers]: {
      filterType: 'text',
      filterName: 'Mobile number',
      filterAccesor: ColumnViewType.ContactsPhoneNumbers,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Phone className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.ContactsCity]: {
      filterType: 'list',
      filterName: 'City',
      filterAccesor: ColumnViewType.ContactsCity,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Globe05 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: uniqBy(
        store?.contacts
          ?.toArray()
          .flatMap((contact) => contact?.value.locations)
          .filter(
            (l) =>
              l?.locality !== null &&
              l?.locality !== undefined &&
              l?.locality !== '',
          )
          .map((location) => ({
            id: location?.locality,
            label: location?.locality,
          }))
          .sort((a, b) => (a.label ?? '').localeCompare(b.label ?? '')),
        'id',
      ),
    },
    [ColumnViewType.ContactsLinkedin]: {
      filterType: 'text',
      filterName: 'LinkedIn URL',
      filterAccesor: ColumnViewType.ContactsLinkedin,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <LinkedinOutline className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.ContactsPersona]: {
      filterType: 'list',
      filterName: 'Persona',
      filterAccesor: ColumnViewType.ContactsPersona,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Tag01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: store?.tags.toArray().map((tag) => ({
        id: tag?.value?.id,
        label: tag?.value?.name,
      })),
    },
    [ColumnViewType.ContactsJobTitle]: {
      filterType: 'text',
      filterName: 'Job title',
      filterAccesor: ColumnViewType.ContactsJobTitle,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Certificate01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.ContactsTimeInCurrentRole]: {
      filterType: 'date',
      filterName: 'Time in current role',
      filterAccesor: ColumnViewType.ContactsTimeInCurrentRole,
      filterOperators: [ComparisonOperator.Gt, ComparisonOperator.Lt],
      icon: (
        <ClockCheck className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.ContactsLinkedinFollowerCount]: {
      filterType: 'number',
      filterName: 'LinkedIn followers',
      filterAccesor: ColumnViewType.ContactsLinkedinFollowerCount,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: (
        <Users03 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    // [ColumnViewType.ContactsLastInteraction]: {
    //   filterType: 'number',
    //   filterName: 'Last Interaction',
    //   filterAccesor: ColumnViewType.ContactsLastInteraction,
    //   filterOperators: [
    //     ComparisonOperator.Gt,
    //     ComparisonOperator.Lt,
    //     ComparisonOperator.IsEmpty,
    //     ComparisonOperator.IsNotEmpty,
    //   ],
    //   icon: <Calendar  />,
    // },
    [ColumnViewType.ContactsConnections]: {
      filterType: 'list',
      filterName: 'LinkedIn connections',
      filterAccesor: ColumnViewType.ContactsConnections,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <LinkedinOutline className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: store?.users.toArray().map((user) => ({
        id: user?.id,
        label: user?.name,
        avatar: user?.value?.profilePhotoUrl,
      })),
    },
    [ColumnViewType.ContactsCountry]: {
      filterType: 'list',
      filterName: 'Country',
      filterAccesor: ColumnViewType.ContactsCountry,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Globe04 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: uniqBy(
        store?.contacts
          ?.toArray()
          .flatMap((contact) => contact?.value.locations?.[0])
          .map((location) => ({
            id: location?.countryCodeA2,
            label: location?.country,
          }))
          .sort((a, b) => (a.label ?? '').localeCompare(b.label ?? '')),
        'id',
      ),
    },
    [ColumnViewType.ContactsRegion]: {
      filterType: 'list',
      filterName: 'Region',
      filterAccesor: ColumnViewType.ContactsRegion,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Globe06 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: uniqBy(
        store?.contacts
          ?.toArray()
          .flatMap((contact) => contact?.value.locations?.[0])
          .filter(
            (l) =>
              l?.region !== null && l?.region !== undefined && l?.region !== '',
          )
          .map((location) => ({
            id: location.region,
            label: location.region,
          }))
          .sort((a, b) => (a.label ?? '').localeCompare(b.label ?? '')),
        'id',
      ),
    },
    [ColumnViewType.ContactsFlows]: {
      filterType: 'list',
      filterName: 'Current flows',
      filterAccesor: ColumnViewType.ContactsFlows,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Shuffle01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: uniqBy(
        store?.flows.toArray().map((flow) => ({
          id: flow?.id,
          label: flow?.value.name,
          isArchived: flow.value.status,
        })),
        'id',
      ),
    },
    [ColumnViewType.ContactsFlowStatus]: {
      filterType: 'list',
      filterName: 'Status in flow',
      filterAccesor: ColumnViewType.ContactsFlowStatus,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Shuffle01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [
        {
          id: FlowParticipantStatus.OnHold,
          label: 'Blocked',
        },
        {
          id: FlowParticipantStatus.Ready,
          label: 'Ready',
        },
        {
          id: FlowParticipantStatus.Scheduled,
          label: 'Scheduled',
        },
        {
          id: FlowParticipantStatus.InProgress,
          label: 'In progress',
        },
        {
          id: FlowParticipantStatus.Completed,
          label: 'Completed',
        },
        {
          id: FlowParticipantStatus.GoalAchieved,
          label: 'Goal achieved',
        },
        {
          id: 'PENDING',
          label: 'Pending',
        },
      ],
    },
  };

  return filterTypes;
};

import { EmailDeliverable } from '@graphql/types';

export const CategoryHeaderLabel = {
  [EmailDeliverable.Deliverable]: 'Deliverable',
  [EmailDeliverable.Undeliverable]: 'Not Deliverable',
  [EmailDeliverable.Unknown]: 'Don’t know',
};

export enum EmailVerificationStatus {
  NoRisk = 'no_risk',
  FirewallProtected = 'firewall_protected',
  FreeAccount = 'free_account',
  GroupMailbox = 'group_mailbox',
  InvalidMailbox = 'invalid_mailbox',
  MailboxFull = 'mailbox_full',
  IncorrectFormat = 'incorrect_format',
  CatchAll = 'catch_all',
  NotVerified = 'not_verified',
  VerificationInProgress = 'verification_in_progress',
}

export const getOptionsForCategory = (category: EmailDeliverable) => {
  switch (category) {
    case EmailDeliverable.Deliverable:
      return [
        { label: 'No risk', value: EmailVerificationStatus.NoRisk },
        {
          label: 'Firewall protected',
          value: EmailVerificationStatus.FirewallProtected,
        },
        {
          label: 'Free account',
          value: EmailVerificationStatus.FreeAccount,
        },
        {
          disabled: true,
          label: 'Group mailbox',
          value: EmailVerificationStatus.GroupMailbox,
        },
      ];
    case EmailDeliverable.Undeliverable:
      return [
        {
          label: 'Mailbox doesn’t exist',
          value: EmailVerificationStatus.InvalidMailbox,
        },
        {
          label: 'Mailbox full',
          value: EmailVerificationStatus.MailboxFull,
        },
        {
          label: 'Incorrect email format',
          value: EmailVerificationStatus.IncorrectFormat,
        },
      ];
    case EmailDeliverable.Unknown:
      return [
        { label: 'Catch all', value: EmailVerificationStatus.CatchAll },
        {
          label: 'Not verified yet',
          value: EmailVerificationStatus.NotVerified,
        },
        {
          label: 'Verification in progress',
          value: EmailVerificationStatus.VerificationInProgress,
        },
      ];
    default:
      return [];
  }
};
