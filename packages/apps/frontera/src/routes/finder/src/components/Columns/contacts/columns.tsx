import { Contact } from '@store/Contacts/Contact.dto';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';
import { CountryCell } from '@finder/components/Columns/Cells/country';
import { DateCell } from '@finder/components/Columns/shared/Cells/DateCell';
import { TextCell } from '@finder/components/Columns/shared/Cells/TextCell';
import { JobTitleCell } from '@finder/components/Columns/contacts/Cells/jobTitle';
import { ContactFlowCell } from '@finder/components/Columns/contacts/Cells/contactFlow';
import { NextFlowAction } from '@finder/components/Columns/contacts/Cells/nextFlowAction';

import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { User, TableViewDef, ColumnViewType } from '@graphql/types';

import { EmailCell } from './Cells/email';
import { AvatarCell } from './Cells/avatar';
import { ContactsTagsCell } from './Cells/tags';
import { FlowStatusCell } from './Cells/flowStatus';
import { ContactLinkedInCell } from './Cells/socials';
import { ContactNameCell } from './Cells/contactName';
import { ContactAvatarHeader } from './Headers/Avatar';
import { ConnectedUsers } from './Cells/connectedUsers';
import { OrganizationNameCell } from './Cells/organization';
import { getColumnConfig } from '../shared/util/getColumnConfig';

type ColumnDatum = Contact;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  [ColumnViewType.ContactsAvatar]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsAvatar,
    size: 26,
    minSize: 26,
    maxSize: 26,
    enableColumnFilter: false,
    enableResizing: false,
    cell: (props) => {
      const icon = props.getValue()?.value?.icon;
      const logo = props.getValue()?.value?.profilePhotoUrl;

      return (
        <AvatarCell
          icon={icon}
          logo={logo}
          id={props.row.original.id}
          name={props.getValue().name}
          canNavigate={props.getValue()?.hasActiveOrganization}
        />
      );
    },
    header: ContactAvatarHeader,
    skeleton: () => <Skeleton className='size-[24px]' />,
  }),
  [ColumnViewType.ContactsName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsName,
    minSize: 102,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return <ContactNameCell contactId={props.row.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Name'
        id={ColumnViewType.ContactsName}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  [ColumnViewType.ContactsOrganization]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsOrganization,
    minSize: 120,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      const lasOrganizationId = props.row.original.value.primaryOrganizationId;
      const contactId = props.row.original.value.id;

      const org = props.row.original.value.primaryOrganizationName;

      return (
        <OrganizationNameCell
          org={org || ''}
          contactId={contactId}
          orgId={lasOrganizationId || ''}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Company'
        id={ColumnViewType.ContactsOrganization}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),

  [ColumnViewType.ContactsPrimaryEmail]: columnHelper.accessor(
    'value.emails.primary',
    {
      id: ColumnViewType.ContactsPrimaryEmail,
      minSize: 128,
      maxSize: 650,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: false,
      cell: (props) => {
        return <EmailCell contactId={props.row.id} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          filterWidth='15rem'
          title='Primary Email'
          id={ColumnViewType.ContactsPrimaryEmail}
          {...getTHeadProps<Contact>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),

  // [ColumnViewType.ContactsPhoneNumbers]: columnHelper.accessor('value.phones', {
  //   id: ColumnViewType.ContactsPhoneNumbers,
  //   minSize: 138,
  //   maxSize: 650,
  //   enableResizing: true,
  //   enableColumnFilter: false,
  //   enableSorting: false,

  //   header: (props) => (
  //     <THead<HTMLInputElement>
  //       title='Mobile Number'
  //       id={ColumnViewType.ContactsPhoneNumbers}
  //       {...getTHeadProps<Contact>(props)}
  //     />
  //   ),
  //   cell: (props) => {
  //     const phoneNumber = props.getValue()?.[0];

  //     const isEnriching = props.row.original.isEnriching;

  //     if (!phoneNumber)
  //       return (
  //         <p className='text-gray-400'>
  //           {isEnriching ? 'Enriching...' : 'Not set'}
  //         </p>
  //       );

  //     return <PhoneCell phone={phoneNumber?.rawPhoneNumber} />;
  //   },
  //   skeleton: () => <Skeleton className='w-[100%] h-[14px]' />,
  // }),
  [ColumnViewType.ContactsCity]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsCity,
    minSize: 65,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      const city = props.getValue()?.[0]?.locality;
      const isEnriching = props.row.original.isEnriching;

      return <TextCell text={city} isEnriching={isEnriching} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='City'
        id={ColumnViewType.ContactsCity}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[33%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsLinkedin]: columnHelper.accessor(
    'value.linkedInUrl',
    {
      id: ColumnViewType.ContactsLinkedin,
      minSize: 96,
      maxSize: 650,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: false,
      cell: (props) => (
        <ContactLinkedInCell contactId={props.row.original.id} />
      ),
      header: (props) => (
        <THead<HTMLInputElement>
          title='LinkedIn'
          filterWidth='14rem'
          id={ColumnViewType.ContactsLinkedin}
          {...getTHeadProps<Contact>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.ContactsPersona]: columnHelper.accessor('value.tags', {
    id: ColumnViewType.ContactsPersona,
    minSize: 92,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: false,
    cell: (props) => {
      return <ContactsTagsCell id={props.row.original.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Tags'
        filterWidth='14rem'
        id={ColumnViewType.ContactsPersona}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.ContactsJobTitle]: columnHelper.accessor('value', {
    id: ColumnViewType.ContactsJobTitle,
    minSize: 94,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: false,
    cell: (props) => {
      const value = props.getValue()?.id;

      return <JobTitleCell contactId={value} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Job Title'
        id={ColumnViewType.ContactsJobTitle}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  //  TODO uncomment when data will be available
  // [ColumnViewType.ContactsExperience]: columnHelper.accessor('value', {
  //   id: ColumnViewType.ContactsExperience,
  //   size: 100,
  //   enableSorting: false,
  //   enableColumnFilter: false,
  //   cell: () => {
  //     return <div className='text-gray-400'>Unknown</div>;
  //   },
  //   header: (props) => (
  //     <THead<HTMLInputElement>
  //       id={ColumnViewType.ContactsExperience}
  //       title='Experience'
  //       renderFilter={(initialFocusRef) => (
  //         <SearchTextFilter
  //           initialFocusRef={initialFocusRef}
  //           property={ColumnViewType.ContactsExperience}
  //           placeholder={'e.g. CTO'}
  //         />
  //       )}
  //       {...getTHeadProps<Contact>(props)}
  //     />
  //   ),
  //   skeleton: () => (
  //     <div className='flex flex-col gap-1'>
  //       <Skeleton className='w-[25%] h-[14px]' />
  //     </div>
  //   ),
  // }),
  [ColumnViewType.ContactsTimeInCurrentRole]: columnHelper.accessor(
    'value.primaryOrganizationJobRoleStartDate',
    {
      id: ColumnViewType.ContactsTimeInCurrentRole,
      minSize: 171,
      maxSize: 650,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: false,
      cell: (props) => {
        const startedAt =
          props.row.original.value.primaryOrganizationJobRoleStartDate;
        const isEnriching = props.row.original.isEnriching;

        if (!startedAt)
          return (
            <p className='text-gray-400'>
              {isEnriching ? 'Enriching...' : 'Not set'}
            </p>
          );

        return (
          <p className='first-letter:capitalize'>
            {DateTimeUtils.timeAgo(startedAt)}
          </p>
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          filterWidth='21rem'
          title='Time In Current Role'
          id={ColumnViewType.ContactsTimeInCurrentRole}
          {...getTHeadProps<Contact>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[25%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.ContactsCountry]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsCountry,
    minSize: 91,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: false,
    cell: (props) => {
      const value = props.row.original.value.id;

      return <CountryCell id={value} type='contact' />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Country'
        id={ColumnViewType.ContactsCountry}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),

  [ColumnViewType.ContactsLinkedinFollowerCount]: columnHelper.accessor(
    'value.linkedInFollowerCount',
    {
      id: ColumnViewType.ContactsLinkedinFollowerCount,

      minSize: 159,
      maxSize: 165,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: false,

      cell: (props) => {
        const value = props.row.original.value.linkedInFollowerCount;

        const isEnriching = props.row.original.isEnriching;

        if (typeof value !== 'number')
          return (
            <div className='text-gray-400'>
              {isEnriching ? 'Enriching...' : 'Not set'}
            </div>
          );

        return <div>{value.toLocaleString()}</div>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          filterWidth='17.5rem'
          title='LinkedIn Followers'
          id={ColumnViewType.ContactsLinkedinFollowerCount}
          {...getTHeadProps<Contact>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.ContactsLastInteraction]: columnHelper.accessor('value', {
    id: ColumnViewType.ContactsLastInteraction,
    minSize: 139,
    maxSize: 600,
    enableResizing: true,
    cell: (_props) => {
      return <div className='text-gray-400'>None</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        filterWidth='17.5rem'
        title='Last Interaction'
        id={ColumnViewType.ContactsLastInteraction}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsConnections]: columnHelper.accessor(
    'value.connectedUsers',
    {
      id: ColumnViewType.ContactsConnections,
      minSize: 177,
      maxSize: 600,
      enableColumnFilter: false,
      enableSorting: true,
      cell: (props) => {
        const users = props.row.original.connectedUsers;

        return <ConnectedUsers users={users as User[]} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='LinkedIn Connections'
          id={ColumnViewType.ContactsConnections}
          {...getTHeadProps<Contact>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.ContactsRegion]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsRegion,
    minSize: 83,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      const region = props.getValue()?.[0]?.region;
      const isEnriching = props.row.original?.isEnriching;

      return <TextCell text={region} isEnriching={isEnriching} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Region'
        id={ColumnViewType.ContactsRegion}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsUpdatedAt]: columnHelper.accessor('value.updatedAt', {
    id: ColumnViewType.ContactsUpdatedAt,
    minSize: 122,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      const lastUpdatedAt = props.getValue();

      return <DateCell value={lastUpdatedAt} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Last Updated'
        id={ColumnViewType.ContactsUpdatedAt}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsFlows]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsFlows,
    minSize: 128,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return <ContactFlowCell contactId={props.row.original.value.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Current Flows'
        filterWidth='17.5rem'
        id={ColumnViewType.ContactsFlows}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsFlowStatus]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsFlowStatus,
    minSize: 127,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      const value = props.row.original.value.id;

      return <FlowStatusCell contactID={value} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Status in Flow'
        id={ColumnViewType.ContactsFlowStatus}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsFlowNextAction]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsFlowNextAction,
    minSize: 147,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      const value = props.row.original.value.id;

      return <NextFlowAction contactID={value} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Next Flow Action'
        id={ColumnViewType.ContactsFlowNextAction}
        {...getTHeadProps<Contact>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
};

export const getContactColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
) => getColumnConfig<ColumnDatum>(columns, tableViewDef);
