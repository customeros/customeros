import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { Skeleton } from '@ui/presentation/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import {
  // SortBy,
  Organization,
  TableViewDef,
} from '@graphql/types';

import { AvatarHeader } from './Headers/Avatar';
import {
  OwnerCell,
  AvatarCell,
  WebsiteCell,
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
    header: AvatarHeader,
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
    cell: (props) => {
      return (
        <OrganizationCell
          id={props.getValue().id}
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
  WEBSITE: columnHelper.accessor('website', {
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
    skeleton: () => (
      <Skeleton
        width='50%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
  }),
  RELATIONSHIP: columnHelper.accessor('isCustomer', {
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

      return <OrganizationRelationship organization={organization} />;
    },
    skeleton: () => (
      <Skeleton
        width='100%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
  }),
  ONBOARDING_STATUS: columnHelper.accessor('accountDetails', {
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
      <Flex flexDir='column' gap='1'>
        <Skeleton
          width='33%'
          height='18px'
          startColor='gray.300'
          endColor='gray.300'
        />
      </Flex>
    ),
  }),
  RENEWAL_LIKELIHOOD: columnHelper.accessor('accountDetails', {
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
        {...getTHeadProps<Organization>(props)}
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
  RENEWAL_DATE: columnHelper.accessor('accountDetails', {
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
    skeleton: () => (
      <Skeleton
        width='50%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
  }),
  FORECAST_ARR: columnHelper.accessor('accountDetails', {
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
  OWNER: columnHelper.accessor('owner', {
    id: 'OWNER',
    minSize: 200,
    filterFn: filterOwnerFn,
    cell: (props) => (
      <OwnerCell id={props.row.original.id} owner={props.getValue()} />
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
        lastTouchPointAt={props.row.original.lastTouchPointAt}
        lastTouchPointTimelineEvent={
          (props.row.original as Organization).lastTouchPointTimelineEvent
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

export const getColumnConfig = (tableViewDef: TableViewDef) => {
  if (!tableViewDef) return [];

  return (tableViewDef.columns ?? []).reduce((acc, curr) => {
    //@ts-expect-error will be fixed
    const columnTypeName = curr?.columnType?.name;

    if (!columnTypeName) return acc;

    const column = columns[columnTypeName];

    if (!column) return acc;

    return [...acc, column];
  }, [] as Column[]);
};

// const getSortConfig = (tableViewDef: TableViewDef) => {
//   // @ts-expect-error remove
//   const sort = JSON.parse(tableViewDef?.sorting);

//   return sort as SortBy;
// };
