import React from 'react';

import { Store } from '@store/store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { Contact, ColumnViewType } from '@graphql/types';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead.tsx';
import { EmailFilter } from '@organizations/components/Columns/Filters/Email';
import { TagsCell } from '@organizations/components/Columns/Cells/tags/TagsCell.tsx';
import { ContactNameCell } from '@organizations/components/Columns/Cells/contactName';
import { PhoneCell } from '@organizations/components/Columns/Cells/phone/PhoneCell.tsx';
import { EmailCell } from '@organizations/components/Columns/Cells/email/EmailCell.tsx';
import { SearchTextFilter } from '@organizations/components/Columns/Filters/SearchTextFilter';
import { ContactLinkedInCell } from '@organizations/components/Columns/Cells/socials/ContactLinkedInCell.tsx';
import { ContactAvatarHeader } from '@organizations/components/Columns/Headers/Avatar/ContactAvatarHeader.tsx';

import { SocialsFilter } from '../Filters';
import { AvatarCell, OrganizationCell } from '../Cells';

type ColumnDatum = Store<Contact>;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

export const contactColumns: Record<string, Column> = {
  [ColumnViewType.ContactsAvatar]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsAvatar,
    size: 26,
    enableColumnFilter: false,
    cell: (props) => {
      const icon = props.getValue()?.value?.icon;
      const logo = props.getValue()?.value?.profilePhotoUrl;

      return (
        <AvatarCell
          icon={icon}
          variant='outlineCircle'
          logo={logo}
          description={''}
          id={props.getValue()?.value?.id}
          name={props.getValue()?.value?.name}
        />
      );
    },
    header: ContactAvatarHeader,
    skeleton: () => <Skeleton className='size-[24px]' />,
  }),
  [ColumnViewType.ContactsName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsName,
    size: 150,
    cell: (props) => {
      return <ContactNameCell contactId={props.row.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsName}
        title='Name'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsName}
            placeholder={'e.g. Isabella Evans'}
          />
        )}
        {...getTHeadProps<Store<Contact>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  [ColumnViewType.ContactsOrganization]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsOrganization,
    size: 150,
    cell: (props) => {
      const organization = props.getValue()?.value?.organizations?.content?.[0];

      if (!organization) return '-';

      return (
        <OrganizationCell
          id={organization.id}
          name={organization.name ?? 'Unnamed'}
          isSubsidiary={false}
          className='font-normal'
          parentOrganizationName=''
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsOrganization}
        title='Organization'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsOrganization}
          />
        )}
        {...getTHeadProps<Store<Contact>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  [ColumnViewType.ContactsEmails]: columnHelper.accessor('value.emails', {
    id: ColumnViewType.ContactsEmails,
    size: 200,
    enableSorting: false,
    cell: (props) => {
      const email = props.getValue()?.[0]?.email;
      const validationDetails = props.getValue()?.[0]?.emailValidationDetails;

      return (
        <EmailCell
          email={email}
          validationDetails={validationDetails}
          contactId={props.row.id}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsEmails}
        title='Email'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <EmailFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsEmails}
          />
        )}
        {...getTHeadProps<Store<Contact>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsPhoneNumbers]: columnHelper.accessor(
    'value.phoneNumbers',
    {
      id: ColumnViewType.ContactsPhoneNumbers,
      size: 125,
      enableSorting: false,

      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.ContactsPhoneNumbers}
          title='Phone'
          renderFilter={(initialFocusRef) => (
            <SearchTextFilter
              initialFocusRef={initialFocusRef}
              property={ColumnViewType.ContactsPhoneNumbers}
              placeholder={'e.g. (907) 834-2765'}
            />
          )}
          {...getTHeadProps<Store<Contact>>(props)}
        />
      ),
      cell: (props) => {
        const phoneNumber = props.getValue()?.[0]?.e164;

        return <PhoneCell phone={phoneNumber} />;
      },
      skeleton: () => <Skeleton className='w-[100%] h-[14px]' />,
    },
  ),
  [ColumnViewType.ContactsCity]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsCity,
    size: 125,
    cell: (props) => {
      const status = props.getValue();

      if (!status?.locality) {
        return <p className='text-gray-400'>Unknown</p>;
      }

      return <div>{status?.locality}</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsCity}
        title='City'
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsCity}
            placeholder={'e.g. New York'}
          />
        )}
        {...getTHeadProps<Store<Contact>>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[33%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsLinkedin]: columnHelper.accessor('value.socials', {
    id: ColumnViewType.ContactsLinkedin,
    size: 125,
    enableSorting: false,
    cell: (props) => <ContactLinkedInCell contactId={props.row.original.id} />,
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsLinkedin}
        title='LinkedIn'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <SocialsFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsLinkedin}
            placeholder={'e.g. linkedin.com/in/isabella-evans'}
          />
        )}
        {...getTHeadProps<Store<Contact>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsPersona]: columnHelper.accessor('value.tags', {
    id: ColumnViewType.ContactsPersona,
    size: 400,
    cell: (props) => {
      return <TagsCell id={props.row.original.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsPersona}
        title='Persona'
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsPersona}
            placeholder={'e.g. Solo RevOps'}
          />
        )}
        {...getTHeadProps<Store<Contact>>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
};
