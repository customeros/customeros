import { Tag01 } from '@ui/media/icons/Tag01';
import { Users03 } from '@ui/media/icons/Users03';
import { Building07 } from '@ui/media/icons/Building07';
import {
  ColumnViewType,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

export type FilterType = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  options?: any[];
  icon: JSX.Element;
  filterName: string;
  filterOperators: ComparisonOperator[];
  filterType: 'text' | 'date' | 'number' | 'list';
  groupOptions?: { label: string; options: { id: string; label: string }[] }[];
  filterAccesor:
    | ColumnViewType
    | 'EMAIL_VERIFICATION_WORK_EMAIL'
    | 'EMAIL_VERIFICATION_PERSONAL_EMAIL';
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

import { EmailVerificationStatus } from './Filters/Email/utils';

export const getFilterTypes = (store?: RootStore) => {
  const filterTypes: Partial<
    Record<
      | ColumnViewType
      | 'EMAIL_VERIFICATION_WORK_EMAIL'
      | 'EMAIL_VERIFICATION_PERSONAL_EMAIL',
      FilterType
    >
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
      icon: <User03 />,
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
      icon: <Building07 />,
    },
    [ColumnViewType.ContactsEmails]: {
      filterType: 'text',
      filterName: 'Work email',
      filterAccesor: ColumnViewType.ContactsEmails,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Mail01 />,
    },
    [ColumnViewType.ContactsPersonalEmails]: {
      filterType: 'text',
      filterName: 'Personal Email',
      filterAccesor: ColumnViewType.ContactsPersonalEmails,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Mail01 />,
    },
    [ColumnViewType.ContactsPhoneNumbers]: {
      filterType: 'text',
      filterName: 'Mobile Number',
      filterAccesor: ColumnViewType.ContactsPhoneNumbers,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Phone />,
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
      icon: <Globe05 />,
      options: uniqBy(
        store?.contacts
          ?.toArray()
          .flatMap((contact) => contact?.value.locations?.[0])
          .filter(
            (l) =>
              l?.locality !== null &&
              l?.locality !== undefined &&
              l?.locality !== '',
          )
          .map((location) => ({
            id: location.locality,
            label: location.locality,
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
      icon: <LinkedinOutline />,
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
      icon: <Tag01 />,
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
      icon: <Certificate01 />,
    },
    [ColumnViewType.ContactsTimeInCurrentRole]: {
      filterType: 'number',
      filterName: 'Time In Current Role',
      filterAccesor: ColumnViewType.ContactsTimeInCurrentRole,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: <ClockCheck />,
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
      icon: <Users03 />,
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
      filterName: 'LinkedIn Connections',
      filterAccesor: ColumnViewType.ContactsConnections,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <LinkedinOutline />,
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
      icon: <Globe04 />,
      options: uniqBy(
        store?.contacts
          ?.toArray()
          .flatMap((contact) => contact?.value.locations?.[0])
          .filter(
            (l) =>
              l?.country !== null &&
              l?.country !== undefined &&
              l?.country !== '',
          )
          .map((location) => ({
            id: location.countryCodeA2,
            label: location.country,
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
      icon: <Globe06 />,
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
      filterType: 'text',
      filterName: 'Current Flow',
      filterAccesor: ColumnViewType.ContactsFlows,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Shuffle01 />,
    },
    [ColumnViewType.ContactsFlowStatus]: {
      filterType: 'list',
      filterName: 'Flow status',
      filterAccesor: ColumnViewType.ContactsFlowStatus,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Shuffle01 />,
    },
    ['EMAIL_VERIFICATION_WORK_EMAIL']: {
      filterType: 'list',
      filterName: 'Email status work email',
      filterAccesor: 'EMAIL_VERIFICATION_WORK_EMAIL',
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Mail01 />,
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
    ['EMAIL_VERIFICATION_PERSONAL_EMAIL']: {
      filterType: 'list',
      filterName: 'Email status personal email',
      filterAccesor: 'EMAIL_VERIFICATION_PERSONAL_EMAIL',
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Mail01 />,
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
  };

  return filterTypes;
};
