import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { TableViewDef, RenewalRecord } from '@graphql/types';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';

import {
  OwnerCell,
  AvatarCell,
  OrganizationCell,
  TimeToRenewalCell,
  LastTouchpointCell,
  RenewalForecastCell,
  RenewalLikelihoodCell,
} from './Cells';
import {
  OwnerFilter,
  filterOwnerFn,
  ForecastFilter,
  filterForecastFn,
  OrganizationFilter,
  TimeToRenewalFilter,
  LastTouchpointFilter,
  filterOrganizationFn,
  filterTimeToRenewalFn,
  filterLastTouchpointFn,
  RenewalLikelihoodFilter,
  filterRenewalLikelihoodFn,
} from './Filters';

type ColumnDatum = RenewalRecord;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  RENEWALS_AVATAR: columnHelper.accessor((row) => row, {
    id: 'AVATAR',
    size: 24,
    fixWidth: true,
    enableColumnFilter: false,
    cell: (props) => {
      const orgId = props.getValue()?.organization?.metadata?.id;
      const orgName = props.getValue()?.organization?.name ?? 'Unnamed';

      return (
        <AvatarCell id={orgId} name={orgName} src={props.getValue()?.logo} />
      );
    },
    header: () => <div className='' />,
    skeleton: () => <Skeleton className='w-[24px] h-[24px] bg-gray-300' />,
  }),
  RENEWALS_NAME: columnHelper.accessor((row) => row, {
    id: 'NAME',
    size: 200,
    filterFn: filterOrganizationFn,
    enableColumnFilter: false,
    enableSorting: false,
    cell: (props) => {
      const contractName =
        props.getValue()?.contract?.contractName ?? 'Unnamed';
      const orgId = props.getValue()?.organization?.metadata?.id;
      const orgName = props.getValue()?.organization?.name ?? 'Unnamed';

      return (
        <OrganizationCell
          id={orgId}
          name={contractName}
          isSubsidiary={true}
          parentOrganizationName={orgName}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='contractName'
        title='Contract name'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <OrganizationFilter
            initialFocusRef={initialFocusRef}
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps<RenewalRecord>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col h-[42px] items-start gap-1'>
        <Skeleton className='w-[100px] h-[18px] bg-gray-300' />
      </div>
    ),
  }),
  RENEWALS_RENEWAL_LIKELIHOOD: columnHelper.accessor('opportunity', {
    id: 'RENEWAL_LIKELIHOOD',
    size: 150,
    filterFn: filterRenewalLikelihoodFn,
    cell: (props) => {
      const opportunityId = props.getValue().id;
      const value = props.getValue().renewalLikelihood;

      return (
        <RenewalLikelihoodCell
          value={value}
          opportunityId={opportunityId}
          id={props.row.original.organization.metadata.id}
        />
      );
    },
    header: (props) => (
      <THead
        id='renewalLikelihood'
        title='Health'
        renderFilter={() => <RenewalLikelihoodFilter column={props.column} />}
        {...getTHeadProps<RenewalRecord>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[50%] h-[18px] bg-gray-300' />
      </div>
    ),
  }),
  RENEWALS_RENEWAL_DATE: columnHelper.accessor('organization.accountDetails', {
    id: 'RENEWAL_DATE',
    size: 150,
    filterFn: filterTimeToRenewalFn,
    enableColumnFilter: false,
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
        {...getTHeadProps<RenewalRecord>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[50%] h-[18px] bg-gray-300' />,
  }),
  RENEWALS_FORECAST_ARR: columnHelper.accessor('opportunity', {
    id: 'FORECAST_ARR',
    size: 150,
    filterFn: filterForecastFn,
    enableColumnFilter: false,
    cell: (props) => {
      const value = props.getValue();
      const amount = value?.amount;
      const potentialAmount = value?.maxAmount;
      const opportunityId = value?.id;
      const renewalAdjustedRate = value?.renewalAdjustedRate;

      return (
        <RenewalForecastCell
          amount={amount}
          opportunityId={opportunityId}
          potentialAmount={potentialAmount}
          adjustedRate={renewalAdjustedRate}
          id={props.row.original.organization?.metadata.id}
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
        {...getTHeadProps<RenewalRecord>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[18px] bg-gray-300' />
      </div>
    ),
  }),
  RENEWALS_OWNER: columnHelper.accessor('organization.owner', {
    id: 'OWNER',
    size: 150,
    filterFn: filterOwnerFn,
    cell: (props) => (
      <OwnerCell id={props.getValue()?.id} owner={props.getValue()} />
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
        {...getTHeadProps<RenewalRecord>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[18px] bg-gray-300' />,
  }),
  RENEWALS_LAST_TOUCHPOINT: columnHelper.accessor((row) => row, {
    id: 'LAST_TOUCHPOINT',
    size: 200,
    filterFn: filterLastTouchpointFn,
    cell: (props) => (
      <LastTouchpointCell
        lastTouchPointType={
          props.row.original.organization?.lastTouchpoint?.lastTouchPointType
        }
        lastTouchPointAt={
          props.row.original?.organization?.lastTouchpoint?.lastTouchPointAt
        }
        lastTouchPointTimelineEvent={
          props.row.original.organization?.lastTouchpoint
            ?.lastTouchPointTimelineEvent
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
        {...getTHeadProps<RenewalRecord>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[75%] h-[18px] bg-gray-300' />
      </div>
    ),
  }),
};

export const getColumnsConfig = (tableViewDef?: TableViewDef) => {
  if (!tableViewDef) return [];

  return (tableViewDef.columns ?? []).reduce((acc, curr) => {
    const columnTypeName = curr?.columnType;

    if (!columnTypeName) return acc;

    const column = { ...columns[columnTypeName], enableHiding: !curr.visible };

    if (!column) return acc;

    return [...acc, column];
  }, [] as Column[]);
};
