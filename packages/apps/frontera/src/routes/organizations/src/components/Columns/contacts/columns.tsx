import React from 'react';

import { ContactStore } from '@store/Contacts/Contact.store';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { CountryCell } from '@organizations/components/Columns/Cells/country';
import { Social, JobRole, TableViewDef, ColumnViewType } from '@graphql/types';
import { TextCell } from '@organizations/components/Columns/shared/Cells/TextCell';
import { ConnectedToFilter } from '@organizations/components/Columns/contacts/Filters/ConnectedToFilter';

import { EmailCell } from './Cells/email';
import { PhoneCell } from './Cells/phone';
import { AvatarCell } from './Cells/avatar';
import { EmailFilter } from './Filters/Email';
import { ContactsTagsCell } from './Cells/tags';
import { ContactLinkedInCell } from './Cells/socials';
import { ContactNameCell } from './Cells/contactName';
import { ContactAvatarHeader } from './Headers/Avatar';
import { OrganizationCell } from './Cells/organization';
import { PersonaFilter } from './Filters/PersonaFilter';
import { ConnectedUsers } from './Cells/connectedUsers';
import { SocialsFilter } from '../shared/Filters/Socials';
import { getColumnConfig } from '../shared/util/getColumnConfig';
import { SearchTextFilter } from '../shared/Filters/SearchTextFilter';
import { NumericValueFilter } from '../shared/Filters/NumericValueFilter';
import { LocationFilter } from '../shared/Filters/LocationFilter/LocationFilter';

type ColumnDatum = ContactStore;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
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
          logo={logo}
          id={props.row.original.organizationId}
          name={props.getValue().name}
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
        {...getTHeadProps<ContactStore>(props)}
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
          name={organization.name || 'Unknown'}
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
        {...getTHeadProps<ContactStore>(props)}
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
        {...getTHeadProps<ContactStore>(props)}
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
          {...getTHeadProps<ContactStore>(props)}
        />
      ),
      cell: (props) => {
        const phoneNumber = props.getValue()?.[0];
        if (!phoneNumber) return <p className='text-gray-400'>Unknown</p>;

        return <PhoneCell phone={phoneNumber?.rawPhoneNumber} />;
      },
      skeleton: () => <Skeleton className='w-[100%] h-[14px]' />,
    },
  ),
  [ColumnViewType.ContactsCity]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsCity,
    size: 125,
    cell: (props) => {
      const city = props.getValue()?.[0]?.locality;

      return <TextCell text={city} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsCity}
        title='City'
        renderFilter={(initialFocusRef) => (
          <LocationFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsCity}
            locationType='locality'
            placeholder={'e.g. New York'}
            type='contacts'
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
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
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsPersona]: columnHelper.accessor('value.tags', {
    id: ColumnViewType.ContactsPersona,
    size: 200,
    cell: (props) => {
      return <ContactsTagsCell id={props.row.original.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsPersona}
        title='Persona'
        renderFilter={(initialFocusRef) => (
          <PersonaFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsPersona}
            placeholder={'e.g. Solo RevOps'}
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsJobTitle]: columnHelper.accessor('value.jobRoles', {
    id: ColumnViewType.ContactsJobTitle,
    size: 250,
    cell: (props) => {
      const value = props.getValue()?.[0]?.jobTitle;

      return <TextCell text={value} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsJobTitle}
        title='Job Title'
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsJobTitle}
            placeholder={'e.g. CTO'}
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsExperience]: columnHelper.accessor('value', {
    id: ColumnViewType.ContactsExperience,
    size: 100,
    enableSorting: false,
    enableColumnFilter: false,
    cell: () => {
      // TODO implement when data will be available
      return <div className='text-gray-400'>Unknown</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsExperience}
        title='Experience'
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsExperience}
            placeholder={'e.g. CTO'}
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsTimeInCurrentRole]: columnHelper.accessor(
    'value.jobRoles',
    {
      id: ColumnViewType.ContactsTimeInCurrentRole,
      size: 170,
      cell: (props) => {
        const jobRole = props.getValue()?.find((role: JobRole) => {
          return role?.endedAt !== null;
        });
        if (!jobRole?.startedAt)
          return <p className='text-gray-400'>Unknown</p>;

        return <p>{DateTimeUtils.timeAgo(jobRole.startedAt)}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.ContactsTimeInCurrentRole}
          title='Time in current role'
          filterWidth='21rem'
          renderFilter={(initialFocusRef) => (
            <NumericValueFilter
              initialFocusRef={initialFocusRef}
              property={ColumnViewType.ContactsTimeInCurrentRole}
              label='time in current role'
              suffix='month'
            />
          )}
          {...getTHeadProps<ContactStore>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[25%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.ContactsCountry]: columnHelper.accessor('value.metadata', {
    id: ColumnViewType.ContactsCountry,
    size: 200,
    cell: (props) => {
      const value = props.getValue()?.id;

      return <CountryCell id={value} type='contact' />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsCountry}
        title='Country'
        renderFilter={(initialFocusRef) => (
          <LocationFilter
            type='contacts'
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsCountry}
            locationType='countryCodeA2'
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsSkills]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsSkills,
    size: 100,
    enableSorting: false,
    enableColumnFilter: false,
    cell: () => {
      // TODO implement when data will be available
      return <div className='text-gray-400'>Unknown</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsSkills}
        title='Skills'
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsSchools]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsSchools,
    size: 100,
    enableSorting: false,
    enableColumnFilter: false,
    cell: () => {
      // TODO implement when data will be available
      return <div className='text-gray-400'>Unknown</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsSchools}
        title='Schools'
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsLanguages]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsLanguages,
    size: 100,
    enableSorting: false,
    enableColumnFilter: false,
    cell: () => {
      // TODO implement when data will be available
      return <div className='text-gray-400'>Unknown</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsLanguages}
        title='Languages'
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsLinkedinFollowerCount]: columnHelper.accessor(
    'value',
    {
      id: ColumnViewType.ContactsLinkedinFollowerCount,
      size: 165,

      cell: (props) => {
        const value = props
          .getValue()
          ?.socials.find((e: Social) =>
            e?.url?.includes('linkedin'),
          )?.followersCount;
        if (typeof value !== 'number')
          return <div className='text-gray-400'>Unknown</div>;

        return <div>{value.toLocaleString()}</div>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.ContactsLinkedinFollowerCount}
          title='LinkedIn Followers'
          filterWidth='17.5rem'
          renderFilter={() => (
            <NumericValueFilter
              property={ColumnViewType.ContactsLinkedinFollowerCount}
              label='followers'
            />
          )}
          {...getTHeadProps<ContactStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.ContactsLastInteraction]: columnHelper.accessor('value', {
    id: ColumnViewType.ContactsLastInteraction,
    size: 150,

    cell: (_props) => {
      return <div className='text-gray-400'>Unknown</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsLastInteraction}
        title='Last Interaction'
        filterWidth='17.5rem'
        renderFilter={() => (
          <NumericValueFilter
            property={ColumnViewType.ContactsLastInteraction}
            label='followers'
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsConnections]: columnHelper.accessor(
    'value.locations',
    {
      id: ColumnViewType.ContactsConnections,
      size: 150,
      enableColumnFilter: true,
      enableSorting: true,
      cell: (props) => {
        const users = props.getValue().map((v: { id: string }) => v?.id);

        return <ConnectedUsers users={users} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.ContactsConnections}
          title='Connected To'
          renderFilter={(initialFocusRef) => (
            <ConnectedToFilter initialFocusRef={initialFocusRef} />
          )}
          {...getTHeadProps<ContactStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.ContactsRegion]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsRegion,
    size: 150,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      const region = props.getValue()?.[0]?.region;

      return <TextCell text={region} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.ContactsRegion}
        title='Region'
        renderFilter={(initialFocusRef) => (
          <LocationFilter
            type='contacts'
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsRegion}
            locationType='region'
            placeholder='e.g. California'
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
};

export const getContactColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
) => getColumnConfig<ColumnDatum>(columns, tableViewDef);
