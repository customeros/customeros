import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { DateTimeUtils } from '@spaces/utils/date.ts';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import {
  // SortBy,
  Organization,
  TableViewDef,
} from '@graphql/types';

import { AvatarHeader } from './Headers/Avatar';
import { LastTouchpointDateCell } from './Cells/touchpointDate';
import {
  OwnerCell,
  AvatarCell,
  WebsiteCell,
  LinkedInCell,
  OnboardingCell,
  OrganizationCell,
  TimeToRenewalCell,
  LastTouchpointCell,
  RenewalForecastCell,
  RenewalLikelihoodCell,
  OrganizationRelationship,
} from './Cells';
import {
  OwnerFilter,
  WebsiteFilter,
  filterOwnerFn,
  ForecastFilter,
  filterWebsiteFn,
  OnboardingFilter,
  filterForecastFn,
  OrganizationFilter,
  RelationshipFilter,
  filterOnboardingFn,
  TimeToRenewalFilter,
  LastTouchpointFilter,
  filterOrganizationFn,
  filterRelationshipFn,
  filterTimeToRenewalFn,
  filterLastTouchpointFn,
  RenewalLikelihoodFilter,
  filterRenewalLikelihoodFn,
} from './Filters';

type ColumnDatum = Omit<Organization, 'lastTouchPointTimelineEvent'>;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  ORGANIZATIONS_AVATAR: columnHelper.accessor((row) => row, {
    id: 'AVATAR',
    minSize: 42,
    maxSize: 70,
    fixWidth: true,
    enableColumnFilter: false,
    cell: (props) => {
      return (
        <AvatarCell
          id={props.getValue()?.metadata.id}
          name={props.getValue()?.name}
          src={props.getValue()?.icon || props.getValue()?.logo}
        />
      );
    },
    header: AvatarHeader,
    skeleton: () => <Skeleton className='size-[42px]' />,
  }),
  ORGANIZATIONS_NAME: columnHelper.accessor((row) => row, {
    id: 'NAME',
    minSize: 200,
    filterFn: filterOrganizationFn,
    cell: (props) => {
      return (
        <OrganizationCell
          id={props.getValue().metadata.id}
          name={props.getValue().name}
          isSubsidiary={!!props.getValue()?.subsidiaryOf?.length}
          parentOrganizationName={
            props.getValue()?.subsidiaryOf?.[0]?.organization.name
          }
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='organization'
        title='Organization'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <OrganizationFilter
            initialFocusRef={initialFocusRef}
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col h-[42px] items-start gap-1'>
        <Skeleton className='w-[100px] h-[18px]' />
        <Skeleton className='w-[100px] h-[18px]' />
      </div>
    ),
  }),
  ORGANIZATIONS_WEBSITE: columnHelper.accessor('website', {
    id: 'WEBSITE',
    minSize: 200,
    enableSorting: false,
    filterFn: filterWebsiteFn,
    cell: (props) => <WebsiteCell website={props.getValue()} />,
    header: (props) => (
      <THead<HTMLInputElement>
        id='website'
        title='Website'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <WebsiteFilter
            initialFocusRef={initialFocusRef}
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[50%] h-[18px]' />,
  }),
  ORGANIZATIONS_RELATIONSHIP: columnHelper.accessor('isCustomer', {
    id: 'RELATIONSHIP',
    minSize: 200,
    filterFn: filterRelationshipFn,
    header: (props) => (
      <THead
        id='relationship'
        title='Relationship'
        renderFilter={() => (
          <RelationshipFilter
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    cell: (props) => {
      const organization = props.row.original;

      // @ts-expect-error will be fixed
      return <OrganizationRelationship organization={organization} />;
    },
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
  ORGANIZATIONS_ONBOARDING_STATUS: columnHelper.accessor('accountDetails', {
    id: 'ONBOARDING_STATUS',
    minSize: 200,
    filterFn: filterOnboardingFn,
    cell: (props) => {
      const status = props.getValue()?.onboarding?.status;
      const updatedAt = props.getValue()?.onboarding?.updatedAt;

      return <OnboardingCell status={status} updatedAt={updatedAt} />;
    },
    header: (props) => (
      <THead
        id='onboardingStatus'
        title='Onboarding'
        renderFilter={() => <OnboardingFilter column={props.column} />}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[33%] h-[18px]' />
      </div>
    ),
  }),
  ORGANIZATIONS_RENEWAL_LIKELIHOOD: columnHelper.accessor('accountDetails', {
    id: 'RENEWAL_LIKELIHOOD',
    minSize: 200,
    filterFn: filterRenewalLikelihoodFn,
    cell: (props) => {
      const value = props.getValue()?.renewalSummary?.renewalLikelihood;

      return (
        <RenewalLikelihoodCell
          value={value}
          id={props.row.original.metadata.id}
        />
      );
    },
    header: (props) => (
      <THead
        id='renewalLikelihood'
        title='Health'
        renderFilter={() => <RenewalLikelihoodFilter column={props.column} />}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[18px]' />
      </div>
    ),
  }),
  ORGANIZATIONS_RENEWAL_DATE: columnHelper.accessor('accountDetails', {
    id: 'RENEWAL_DATE',
    minSize: 200,
    filterFn: filterTimeToRenewalFn,
    cell: (props) => {
      const nextRenewalDate = props.getValue()?.renewalSummary?.nextRenewalDate;

      return <TimeToRenewalCell nextRenewalDate={nextRenewalDate} />;
    },

    header: (props) => (
      <THead
        id='timeToRenewal'
        title='Next Renewal'
        renderFilter={() => (
          <TimeToRenewalFilter
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[50%] h-[18px]' />,
  }),
  ORGANIZATIONS_FORECAST_ARR: columnHelper.accessor('accountDetails', {
    id: 'FORECAST_ARR',
    minSize: 200,
    filterFn: filterForecastFn,
    cell: (props) => {
      const value = props.getValue()?.renewalSummary;
      const amount = value?.arrForecast;
      const potentialAmount = value?.maxArrForecast;

      return (
        <RenewalForecastCell
          amount={amount}
          potentialAmount={potentialAmount}
          id={props.row.original.metadata.id}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='forecast'
        title='ARR Forecast'
        filterWidth='17rem'
        renderFilter={(initialFocusRef) => (
          <ForecastFilter
            initialFocusRef={initialFocusRef}
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[50%] h-[18px]' />
        <Skeleton className='w-[25%] h-[18px]' />
      </div>
    ),
  }),
  ORGANIZATIONS_OWNER: columnHelper.accessor('owner', {
    id: 'OWNER',
    minSize: 200,
    filterFn: filterOwnerFn,
    cell: (props) => (
      <OwnerCell id={props.row.original.metadata.id} owner={props.getValue()} />
    ),
    header: (props) => (
      <THead<HTMLInputElement>
        id='owner'
        title='Owner'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <OwnerFilter
            initialFocusRef={initialFocusRef}
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[18px]' />,
  }),
  ORGANIZATIONS_LEAD_SOURCE: columnHelper.accessor('owner', {
    id: 'LEAD_SOURCE',
    minSize: 200,
    cell: (props) => {
      if (!props.row.original.leadSource) {
        return <p className='text-gray-400'>Unknown</p>;
      }

      return (
        <p className='text-gray-700 cursor-default truncate'>
          {props.row.original.leadSource}
        </p>
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='lead'
        title='Source'
        filterWidth='14rem'
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[18px]' />,
  }),
  ORGANIZATIONS_CREATED_DATE: columnHelper.accessor('metadata', {
    id: 'ORGANIZATIONS_CREATED_DATE',
    minSize: 200,
    cell: (props) => {
      if (!props.row.original.metadata.created) {
        return <p className='text-gray-400'>Unknown</p>;
      }

      return (
        <p className='text-gray-700 cursor-default truncate'>
          {DateTimeUtils.format(
            props.row.original.metadata.created,
            DateTimeUtils.defaultFormatShortString,
          )}
        </p>
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='lead'
        title='Created Date'
        filterWidth='14rem'
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[18px]' />,
  }),
  ORGANIZATIONS_YEAR_FOUNDED: columnHelper.accessor('yearFounded', {
    id: 'ORGANIZATIONS_YEAR_FOUNDED',
    minSize: 200,
    cell: (props) => {
      if (!props.row.original.yearFounded) {
        return <p className='text-gray-400'>Unknown</p>;
      }

      return (
        <p className='text-gray-700 cursor-default truncate'>
          {props.row.original.yearFounded}
        </p>
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='lead'
        title='Year Founded'
        filterWidth='14rem'
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[18px]' />,
  }),
  ORGANIZATIONS_EMPLOYEE_COUNT: columnHelper.accessor('employees', {
    id: 'ORGANIZATIONS_EMPLOYEE_COUNT',
    minSize: 200,
    cell: (props) => {
      if (!props.row.original.employees) {
        return <p className='text-gray-400'>Unknown</p>;
      }

      return (
        <p className='text-gray-700 cursor-default truncate'>
          {props.row.original.employees}
        </p>
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='lead'
        title='Employee Count'
        filterWidth='14rem'
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[18px]' />,
  }),
  ORGANIZATIONS_SOCIALS: columnHelper.accessor('socialMedia', {
    id: 'ORGANIZATIONS_SOCIALS',
    minSize: 200,
    cell: (props) => <LinkedInCell socials={props.row.original.socialMedia} />,
    header: (props) => (
      <THead<HTMLInputElement>
        id='socials'
        title='LinkedIn'
        filterWidth='14rem'
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[18px]' />,
  }),
  ORGANIZATIONS_LAST_TOUCHPOINT: columnHelper.accessor((row) => row, {
    id: 'ORGANIZATIONS_LAST_TOUCHPOINT',
    minSize: 250,
    filterFn: filterLastTouchpointFn,
    cell: (props) => (
      <LastTouchpointCell
        lastTouchPointAt={props.row.original?.lastTouchpoint?.lastTouchPointAt}
        lastTouchPointTimelineEvent={
          props.row.original?.lastTouchpoint?.lastTouchPointTimelineEvent
        }
      />
    ),
    header: (props) => (
      <THead<HTMLInputElement>
        id='lastTouchpoint'
        title='Last Touchpoint'
        renderFilter={() => (
          <LastTouchpointFilter
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[75%] h-[18px]' />
        <Skeleton className='w-[100%] h-[18px]' />
      </div>
    ),
  }),
  ORGANIZATIONS_LAST_TOUCHPOINT_DATE: columnHelper.accessor((row) => row, {
    id: 'ORGANIZATIONS_LAST_TOUCHPOINT_DATE',
    minSize: 200,
    cell: (props) => (
      <LastTouchpointDateCell
        lastTouchPointAt={props.row.original?.lastTouchpoint?.lastTouchPointAt}
      />
    ),
    header: (props) => (
      <THead<HTMLInputElement>
        id='lastTouchpoint'
        title='Last Touchpoint'
        renderFilter={() => (
          <LastTouchpointFilter
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[75%] h-[18px]' />
        <Skeleton className='w-[100%] h-[18px]' />
      </div>
    ),
  }),
};

export const getColumnsConfig = (tableViewDef?: Array<TableViewDef>[0]) => {
  if (!tableViewDef) return [];

  return (tableViewDef.columns ?? []).reduce((acc, curr) => {
    const columnTypeName = curr?.columnType;
    if (!columnTypeName) return acc;

    if (columns[columnTypeName] === undefined) return acc;
    const column = { ...columns[columnTypeName], enableHiding: !curr.visible };

    if (!column) return acc;

    return [...acc, column];
  }, [] as Column[]);
};
