import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { Skeleton } from '@ui/presentation/Skeleton';
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
  AVATAR: columnHelper.accessor((row) => row, {
    id: 'AVATAR',
    minSize: 42,
    maxSize: 70,
    fixWidth: true,
    enableColumnFilter: false,
    cell: (props) => {
      return (
        <AvatarCell
          id={props.getValue()?.id}
          name={props.getValue()?.name}
          src={props.getValue()?.logoUrl}
        />
      );
    },
    header: () => <></>,
    skeleton: () => (
      <Skeleton
        width='42px'
        height='42px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
  }),
  NAME: columnHelper.accessor((row) => row, {
    id: 'NAME',
    minSize: 200,
    filterFn: filterOrganizationFn,
    enableColumnFilter: false,
    enableSorting: false,
    cell: (props) => {
      const contractName = props.getValue()?.contract?.name ?? 'Unnamed';
      const orgId = props.getValue()?.organization?.id;
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
      <Flex flexDir='column' h='42px' align='flex-start' gap='1'>
        <Skeleton
          width='100px'
          height='18px'
          startColor='gray.300'
          endColor='gray.300'
        />
        <Skeleton
          width='100px'
          height='18px'
          startColor='gray.300'
          endColor='gray.300'
        />
      </Flex>
    ),
  }),
  RENEWAL_LIKELIHOOD: columnHelper.accessor('organization.accountDetails', {
    id: 'RENEWAL_LIKELIHOOD',
    minSize: 200,
    filterFn: filterRenewalLikelihoodFn,
    cell: (props) => {
      const value = props.getValue()?.renewalSummary?.renewalLikelihood;

      return <RenewalLikelihoodCell value={value} />;
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
      <Flex flexDir='column' gap='1'>
        <Skeleton
          width='25%'
          height='18px'
          startColor='gray.300'
          endColor='gray.300'
        />
      </Flex>
    ),
  }),
  RENEWAL_DATE: columnHelper.accessor('organization.accountDetails', {
    id: 'RENEWAL_DATE',
    minSize: 200,
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
    skeleton: () => (
      <Skeleton
        width='50%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
  }),
  FORECAST_ARR: columnHelper.accessor('organization.accountDetails', {
    id: 'FORECAST_ARR',
    minSize: 200,
    filterFn: filterForecastFn,
    enableColumnFilter: false,
    cell: (props) => {
      const value = props.getValue()?.renewalSummary;
      const amount = value?.arrForecast;
      const potentialAmount = value?.maxArrForecast;

      return (
        <RenewalForecastCell
          amount={amount}
          potentialAmount={potentialAmount}
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
      <Flex flexDir='column' gap='1'>
        <Skeleton
          width='50%'
          height='18px'
          startColor='gray.300'
          endColor='gray.300'
        />
        <Skeleton
          width='25%'
          height='18px'
          startColor='gray.300'
          endColor='gray.300'
        />
      </Flex>
    ),
  }),
  OWNER: columnHelper.accessor('contract.owner', {
    id: 'OWNER',
    minSize: 200,
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
    skeleton: () => (
      <Skeleton
        width='75%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
  }),
  LAST_TOUCHPOINT: columnHelper.accessor((row) => row, {
    id: 'LAST_TOUCHPOINT',
    minSize: 250,
    filterFn: filterLastTouchpointFn,
    cell: (props) => (
      <LastTouchpointCell
        lastTouchPointAt={props.row.original.organization.lastTouchPointAt}
        lastTouchPointTimelineEvent={
          props.row.original.organization.lastTouchPointTimelineEvent
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
      <Flex flexDir='column' gap='1'>
        <Skeleton
          width='75%'
          height='18px'
          startColor='gray.300'
          endColor='gray.300'
        />
        <Skeleton
          width='100%'
          height='18px'
          startColor='gray.300'
          endColor='gray.300'
        />
      </Flex>
    ),
  }),
};

export const getColumnsConfig = (tableViewDef?: TableViewDef) => {
  if (!tableViewDef) return [];

  return (tableViewDef.columns ?? []).reduce((acc, curr) => {
    const columnTypeName = curr?.columnType?.name;

    if (!columnTypeName) return acc;

    const column = columns[columnTypeName];

    if (!column) return acc;

    return [...acc, column];
  }, [] as Column[]);
};
