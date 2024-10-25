import { ContactStore } from '@store/Contacts/Contact.store';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';
import { CountryCell } from '@finder/components/Columns/Cells/country';
import { TextCell } from '@finder/components/Columns/shared/Cells/TextCell';
import { JobTitleCell } from '@finder/components/Columns/contacts/Cells/jobTitle';
import { ContactFlowCell } from '@finder/components/Columns/contacts/Cells/contactFlow';
import { ConnectedToFilter } from '@finder/components/Columns/contacts/Filters/ConnectedToFilter';

import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { User, Social, TableViewDef, ColumnViewType } from '@graphql/types';

import { EmailCell } from './Cells/email';
import { PhoneCell } from './Cells/phone';
import { AvatarCell } from './Cells/avatar';
import { EmailFilter } from './Filters/Email';
import { ContactsTagsCell } from './Cells/tags';
import { FlowsFilter } from './Filters/FlowsFilter';
import { FlowStatusCell } from './Cells/flowStatus';
import { ContactLinkedInCell } from './Cells/socials';
import { ContactNameCell } from './Cells/contactName';
import { ContactAvatarHeader } from './Headers/Avatar';
import { PersonaFilter } from './Filters/PersonaFilter';
import { ConnectedUsers } from './Cells/connectedUsers';
import { SocialsFilter } from '../shared/Filters/Socials';
import { OrganizationCell } from '../shared/Cells/organization';
import { getColumnConfig } from '../shared/util/getColumnConfig';
import { ContactFlowStatusFilter } from './Filters/FlowStatusFilter';
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
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      return (
        <ContactNameCell
          contactId={props.row.id}
          canNavigate={props.getValue()?.hasActiveOrganization}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Name'
        filterWidth='14rem'
        id={ColumnViewType.ContactsName}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            placeholder={'e.g. Isabella Evans'}
            property={ColumnViewType.ContactsName}
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  [ColumnViewType.ContactsOrganization]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsOrganization,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      const lasOrganization =
        props.row.original.value.latestOrganizationWithJobRole?.organization
          .metadata.id;

      if (!lasOrganization) return '-';

      return <OrganizationCell id={lasOrganization} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        filterWidth='14rem'
        title='Organization'
        id={ColumnViewType.ContactsOrganization}
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

  [ColumnViewType.ContactsPrimaryEmail]: columnHelper.accessor(
    'value.primaryEmail',
    {
      id: ColumnViewType.ContactsPrimaryEmail,
      minSize: 230,
      maxSize: 650,
      enableResizing: true,
      enableColumnFilter: true,
      enableSorting: false,
      cell: (props) => {
        const validationDetails =
          props.row.original.value.primaryEmail?.emailValidationDetails;

        return (
          <EmailCell
            contactId={props.row.id}
            validationDetails={validationDetails}
          />
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          filterWidth='15rem'
          title='Primary Email'
          id={ColumnViewType.ContactsPrimaryEmail}
          renderFilter={(initialFocusRef) => (
            <EmailFilter
              initialFocusRef={initialFocusRef}
              property={ColumnViewType.ContactsPersonalEmails}
            />
          )}
          {...getTHeadProps<ContactStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),

  [ColumnViewType.ContactsPhoneNumbers]: columnHelper.accessor(
    'value.phoneNumbers',
    {
      id: ColumnViewType.ContactsPhoneNumbers,
      minSize: 125,
      maxSize: 650,
      enableResizing: true,
      enableColumnFilter: true,
      enableSorting: false,

      header: (props) => (
        <THead<HTMLInputElement>
          title='Mobile Number'
          id={ColumnViewType.ContactsPhoneNumbers}
          renderFilter={(initialFocusRef) => (
            <SearchTextFilter
              initialFocusRef={initialFocusRef}
              placeholder={'e.g. (907) 834-2765'}
              property={ColumnViewType.ContactsPhoneNumbers}
            />
          )}
          {...getTHeadProps<ContactStore>(props)}
        />
      ),
      cell: (props) => {
        const phoneNumber = props.getValue()?.[0];
        const enrichedContact = props.row.original.value.enrichDetails;

        const enrichingStatus =
          !enrichedContact?.enrichedAt &&
          enrichedContact?.requestedAt &&
          !enrichedContact?.failedAt;

        if (!phoneNumber)
          return (
            <p className='text-gray-400'>
              {enrichingStatus ? 'Enriching...' : 'Not set'}
            </p>
          );

        return <PhoneCell phone={phoneNumber?.rawPhoneNumber} />;
      },
      skeleton: () => <Skeleton className='w-[100%] h-[14px]' />,
    },
  ),
  [ColumnViewType.ContactsCity]: columnHelper.accessor('value.locations', {
    id: ColumnViewType.ContactsCity,
    minSize: 125,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      const city = props.getValue()?.[0]?.locality;
      const enrichedContact = props.row.original.value.enrichDetails;
      const enrichingStatus =
        !enrichedContact?.enrichedAt &&
        enrichedContact?.requestedAt &&
        !enrichedContact?.failedAt;

      return <TextCell text={city} enrichingStatus={enrichingStatus} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='City'
        id={ColumnViewType.ContactsCity}
        renderFilter={(initialFocusRef) => (
          <LocationFilter
            type='contacts'
            locationType='locality'
            placeholder={'e.g. New York'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsCity}
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
    minSize: 125,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: false,
    cell: (props) => <ContactLinkedInCell contactId={props.row.original.id} />,
    header: (props) => (
      <THead<HTMLInputElement>
        title='LinkedIn'
        filterWidth='14rem'
        id={ColumnViewType.ContactsLinkedin}
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
    minSize: 120,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: false,
    cell: (props) => {
      return <ContactsTagsCell id={props.row.original.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Persona'
        filterWidth='14rem'
        id={ColumnViewType.ContactsPersona}
        renderFilter={(initialFocusRef) => (
          <PersonaFilter
            placeholder={'e.g. Solo RevOps'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsPersona}
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
  [ColumnViewType.ContactsJobTitle]: columnHelper.accessor('value', {
    id: ColumnViewType.ContactsJobTitle,
    minSize: 120,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: false,
    cell: (props) => {
      const value = props.getValue()?.id;

      return <JobTitleCell contactId={value} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Job Title'
        id={ColumnViewType.ContactsJobTitle}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            placeholder={'e.g. CTO'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsJobTitle}
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
  //       {...getTHeadProps<ContactStore>(props)}
  //     />
  //   ),
  //   skeleton: () => (
  //     <div className='flex flex-col gap-1'>
  //       <Skeleton className='w-[25%] h-[14px]' />
  //     </div>
  //   ),
  // }),
  [ColumnViewType.ContactsTimeInCurrentRole]: columnHelper.accessor(
    'value.jobRoles',
    {
      id: ColumnViewType.ContactsTimeInCurrentRole,
      minSize: 190,
      maxSize: 650,
      enableResizing: true,
      enableColumnFilter: true,
      enableSorting: false,
      cell: (props) => {
        const jobRole =
          props.row.original.value.latestOrganizationWithJobRole?.jobRole;
        const enrichedContact = props.row.original.value.enrichDetails;

        const enrichingStatus =
          !enrichedContact?.enrichedAt &&
          enrichedContact?.requestedAt &&
          !enrichedContact?.failedAt;

        if (!jobRole?.startedAt)
          return (
            <p className='text-gray-400'>
              {enrichingStatus ? 'Enriching...' : 'Not set'}
            </p>
          );

        return <p>{DateTimeUtils.timeAgo(jobRole.startedAt)}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          filterWidth='21rem'
          title='Time In Current Role'
          id={ColumnViewType.ContactsTimeInCurrentRole}
          renderFilter={(initialFocusRef) => (
            <NumericValueFilter
              suffix='month'
              label='time in current role'
              initialFocusRef={initialFocusRef}
              property={ColumnViewType.ContactsTimeInCurrentRole}
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
    minSize: 200,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: false,
    cell: (props) => {
      const value = props.getValue()?.id;

      return <CountryCell id={value} type='contact' />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Country'
        id={ColumnViewType.ContactsCountry}
        renderFilter={(initialFocusRef) => (
          <LocationFilter
            type='contacts'
            locationType='countryCodeA2'
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsCountry}
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
  //  TODO uncomment when data will be available
  // [ColumnViewType.ContactsSkills]: columnHelper.accessor('value.locations', {
  //   id: ColumnViewType.ContactsSkills,
  //   size: 100,
  //   enableSorting: false,
  //   enableColumnFilter: false,
  //   cell: () => {
  //     // TODO implement when data will be available
  //     return <div className='text-gray-400'>Unknown</div>;
  //   },
  //   header: (props) => (
  //     <THead<HTMLInputElement>
  //       id={ColumnViewType.ContactsSkills}
  //       title='Skills'
  //       {...getTHeadProps<ContactStore>(props)}
  //     />
  //   ),
  //   skeleton: () => (
  //     <div className='flex flex-col gap-1'>
  //       <Skeleton className='w-[25%] h-[14px]' />
  //     </div>
  //   ),
  // }),
  //  TODO uncomment when data will be available
  // [ColumnViewType.ContactsSchools]: columnHelper.accessor('value.locations', {
  //   id: ColumnViewType.ContactsSchools,
  //   size: 100,
  //   enableSorting: false,
  //   enableColumnFilter: false,
  //   cell: () => {
  //     // TODO implement when data will be available
  //     return <div className='text-gray-400'>Unknown</div>;
  //   },
  //   header: (props) => (
  //     <THead<HTMLInputElement>
  //       id={ColumnViewType.ContactsSchools}
  //       title='Schools'
  //       {...getTHeadProps<ContactStore>(props)}
  //     />
  //   ),
  //   skeleton: () => (
  //     <div className='flex flex-col gap-1'>
  //       <Skeleton className='w-[25%] h-[14px]' />
  //     </div>
  //   ),
  // }),
  //  TODO uncomment when data will be available
  // [ColumnViewType.ContactsLanguages]: columnHelper.accessor('value.locations', {
  //   id: ColumnViewType.ContactsLanguages,
  //   size: 100,
  //   enableSorting: false,
  //   enableColumnFilter: false,
  //   cell: () => {
  //     // TODO implement when data will be available
  //     return <div className='text-gray-400'>Unknown</div>;
  //   },
  //   header: (props) => (
  //     <THead<HTMLInputElement>
  //       id={ColumnViewType.ContactsLanguages}
  //       title='Languages'
  //       {...getTHeadProps<ContactStore>(props)}
  //     />
  //   ),
  //   skeleton: () => (
  //     <div className='flex flex-col gap-1'>
  //       <Skeleton className='w-[25%] h-[14px]' />
  //     </div>
  //   ),
  // }),
  [ColumnViewType.ContactsLinkedinFollowerCount]: columnHelper.accessor(
    'value',
    {
      id: ColumnViewType.ContactsLinkedinFollowerCount,
      size: 165,
      minSize: 165,
      maxSize: 165,
      enableResizing: true,
      enableColumnFilter: true,
      enableSorting: false,

      cell: (props) => {
        const value = props
          .getValue()
          ?.socials.find((e: Social) =>
            e?.url?.includes('linkedin'),
          )?.followersCount;

        const enrichedContact = props.row.original.value.enrichDetails;
        const enrichingStatus =
          !enrichedContact?.enrichedAt &&
          enrichedContact?.requestedAt &&
          !enrichedContact?.failedAt;

        if (typeof value !== 'number')
          return (
            <div className='text-gray-400'>
              {enrichingStatus ? 'Enriching...' : 'Not set'}
            </div>
          );

        return <div>{value.toLocaleString()}</div>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          filterWidth='17.5rem'
          title='LinkedIn Followers'
          id={ColumnViewType.ContactsLinkedinFollowerCount}
          renderFilter={() => (
            <NumericValueFilter
              label='followers'
              property={ColumnViewType.ContactsLinkedinFollowerCount}
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
    minSize: 150,
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
        renderFilter={() => (
          <NumericValueFilter
            label='followers'
            property={ColumnViewType.ContactsLastInteraction}
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsConnections]: columnHelper.accessor(
    'value.connectedUsers',
    {
      id: ColumnViewType.ContactsConnections,
      minSize: 150,
      maxSize: 600,
      enableColumnFilter: true,
      enableSorting: true,

      cell: (props) => {
        const users = props.row.original.connectedUsers;

        return <ConnectedUsers users={users as User[]} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='LinkedIn Connections'
          id={ColumnViewType.ContactsConnections}
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
    minSize: 150,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      const region = props.getValue()?.[0]?.region;
      const enrichedContact = props.row.original.value.enrichDetails;
      const enrichingStatus =
        !enrichedContact?.enrichedAt &&
        enrichedContact?.requestedAt &&
        !enrichedContact?.failedAt;

      return <TextCell text={region} enrichingStatus={enrichingStatus} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Region'
        id={ColumnViewType.ContactsRegion}
        renderFilter={(initialFocusRef) => (
          <LocationFilter
            type='contacts'
            locationType='region'
            placeholder='e.g. California'
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContactsRegion}
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsFlows]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsFlows,
    minSize: 170,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      return <ContactFlowCell contactId={props.row.original.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Current Flow'
        filterWidth='17.5rem'
        id={ColumnViewType.ContactsFlows}
        renderFilter={(initialFocusRef) => (
          <FlowsFilter
            placeholder='e.g. Aerospace CTO'
            initialFocusRef={initialFocusRef}
          />
        )}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.ContactsFlowStatus]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContactsFlowStatus,
    minSize: 170,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      const value = props.getValue()?.flowContact?.value?.status;

      return <FlowStatusCell value={value} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Status in Flow'
        id={ColumnViewType.ContactsFlowStatus}
        renderFilter={() => <ContactFlowStatusFilter />}
        {...getTHeadProps<ContactStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
};

export const getContactColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
) => getColumnConfig<ColumnDatum>(columns, tableViewDef);
