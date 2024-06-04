import { match } from 'ts-pattern';
import { Store } from '@store/store';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { DateTimeUtils } from '@utils/date.ts';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import {
  Filter,
  Organization,
  TableViewDef,
  OnboardingStatus,
  OpportunityRenewalLikelihood,
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
  OrganizationRelationshipCell,
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

type ColumnDatum = Store<Organization>;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  ORGANIZATIONS_AVATAR: columnHelper.accessor((row) => row, {
    id: 'ORGANIZATIONS_AVATAR',
    minSize: 24,
    maxSize: 50,
    fixWidth: true,
    enableColumnFilter: false,
    cell: (props) => {
      const icon = props.getValue()?.value?.icon;
      const logo = props.getValue()?.value?.logo;
      const description = props.getValue()?.value?.valueProposition;

      return (
        <AvatarCell
          icon={icon}
          logo={logo}
          description={description}
          id={props.getValue()?.value?.metadata?.id}
          name={props.getValue()?.value?.name}
        />
      );
    },
    header: AvatarHeader,
    skeleton: () => <Skeleton className='size-[24px]' />,
  }),
  ORGANIZATIONS_NAME: columnHelper.accessor((row) => row, {
    id: 'ORGANIZATIONS_NAME',
    minSize: 200,
    filterFn: filterOrganizationFn,
    cell: (props) => {
      return (
        <OrganizationCell
          id={props.getValue().value.metadata?.id}
          name={props.getValue().value.name}
          isSubsidiary={!!props.getValue()?.value?.subsidiaryOf?.length}
          parentOrganizationName={
            props.getValue()?.value?.subsidiaryOf?.[0]?.organization.name
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
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  ORGANIZATIONS_WEBSITE: columnHelper.accessor('value.website', {
    id: 'ORGANIZATIONS_WEBSITE',
    minSize: 150,
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
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
  }),
  ORGANIZATIONS_RELATIONSHIP: columnHelper.accessor('value.relationship', {
    id: 'ORGANIZATIONS_RELATIONSHIP',
    minSize: 150,
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
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    cell: (props) => {
      const id = props.row.original.value.metadata?.id;

      return <OrganizationRelationshipCell id={id} />;
    },
    skeleton: () => <Skeleton className='w-[100%] h-[14px]' />,
  }),
  ORGANIZATIONS_ONBOARDING_STATUS: columnHelper.accessor(
    'value.accountDetails',
    {
      id: 'ORGANIZATIONS_ONBOARDING_STATUS',
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
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[33%] h-[14px]' />
        </div>
      ),
    },
  ),
  ORGANIZATIONS_RENEWAL_LIKELIHOOD: columnHelper.accessor(
    'value.accountDetails',
    {
      id: 'ORGANIZATIONS_RENEWAL_LIKELIHOOD',
      minSize: 150,
      filterFn: filterRenewalLikelihoodFn,
      cell: (props) => {
        const value = props.getValue()?.renewalSummary?.renewalLikelihood;

        return (
          <RenewalLikelihoodCell
            value={value}
            id={props.row.original.value.metadata?.id}
          />
        );
      },
      header: (props) => (
        <THead
          id='renewalLikelihood'
          title='Health'
          renderFilter={() => <RenewalLikelihoodFilter column={props.column} />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[25%] h-[14px]' />
        </div>
      ),
    },
  ),
  ORGANIZATIONS_RENEWAL_DATE: columnHelper.accessor('value.accountDetails', {
    id: 'ORGANIZATIONS_RENEWAL_DATE',
    minSize: 150,
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
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
  }),
  ORGANIZATIONS_FORECAST_ARR: columnHelper.accessor('value.accountDetails', {
    id: 'ORGANIZATIONS_FORECAST_ARR',
    minSize: 150,
    filterFn: filterForecastFn,
    cell: (props) => {
      const value = props.getValue()?.renewalSummary;
      const amount = value?.arrForecast;
      const potentialAmount = value?.maxArrForecast;

      return (
        <RenewalForecastCell
          amount={amount}
          potentialAmount={potentialAmount}
          id={props.row.original.value.metadata?.id}
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
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[50%] h-[14px]' />
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  ORGANIZATIONS_OWNER: columnHelper.accessor('value.owner', {
    id: 'ORGANIZATIONS_OWNER',
    minSize: 200,
    filterFn: filterOwnerFn,
    cell: (props) => (
      <OwnerCell
        id={props.row.original.value.metadata?.id}
        owner={props.getValue()}
      />
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
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  ORGANIZATIONS_LEAD_SOURCE: columnHelper.accessor('value.leadSource', {
    id: 'ORGANIZATIONS_LEAD_SOURCE',
    minSize: 200,
    cell: (props) => {
      if (!props.getValue()) {
        return <p className='text-gray-400'>Unknown</p>;
      }

      return (
        <p className='text-gray-700 cursor-default truncate'>
          {props.getValue()}
        </p>
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='lead'
        title='Source'
        filterWidth='14rem'
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  ORGANIZATIONS_CREATED_DATE: columnHelper.accessor('value.metadata.created', {
    id: 'ORGANIZATIONS_CREATED_DATE',
    minSize: 200,
    cell: (props) => {
      const value = props.getValue();

      if (!value) {
        return <p className='text-gray-400'>Unknown</p>;
      }

      return (
        <p className='text-gray-700 cursor-default truncate'>
          {DateTimeUtils.format(value, DateTimeUtils.defaultFormatShortString)}
        </p>
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='lead'
        title='Created Date'
        filterWidth='14rem'
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  ORGANIZATIONS_YEAR_FOUNDED: columnHelper.accessor('value.yearFounded', {
    id: 'ORGANIZATIONS_YEAR_FOUNDED',
    minSize: 200,
    cell: (props) => {
      const value = props.getValue();

      if (!value) {
        return <p className='text-gray-400'>Unknown</p>;
      }

      return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='lead'
        title='Year Founded'
        filterWidth='14rem'
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  ORGANIZATIONS_EMPLOYEE_COUNT: columnHelper.accessor('value.employees', {
    id: 'ORGANIZATIONS_EMPLOYEE_COUNT',
    minSize: 200,
    cell: (props) => {
      const value = props.getValue();

      if (!value) {
        return <p className='text-gray-400'>Unknown</p>;
      }

      return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='lead'
        title='Employee Count'
        filterWidth='14rem'
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  ORGANIZATIONS_SOCIALS: columnHelper.accessor('value.socialMedia', {
    id: 'ORGANIZATIONS_SOCIALS',
    minSize: 200,
    cell: (props) => {
      const organizationId = props.row.original.value.metadata.id;

      return <LinkedInCell organizationId={organizationId} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id='socials'
        title='LinkedIn'
        filterWidth='14rem'
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  ORGANIZATIONS_LAST_TOUCHPOINT: columnHelper.accessor((row) => row, {
    id: 'ORGANIZATIONS_LAST_TOUCHPOINT',
    minSize: 350,
    filterFn: filterLastTouchpointFn,
    cell: (props) => (
      <LastTouchpointCell
        lastTouchPointAt={
          props.row.original?.value?.lastTouchpoint?.lastTouchPointAt
        }
        lastTouchPointTimelineEvent={
          props.row.original?.value?.lastTouchpoint?.lastTouchPointTimelineEvent
        }
        lastTouchPointType={
          props.row.original?.value?.lastTouchpoint?.lastTouchPointType
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
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[75%] h-[14px]' />
        <Skeleton className='w-[100%] h-[14px]' />
      </div>
    ),
  }),
  ORGANIZATIONS_LAST_TOUCHPOINT_DATE: columnHelper.accessor((row) => row, {
    id: 'ORGANIZATIONS_LAST_TOUCHPOINT_DATE',
    minSize: 200,
    enableSorting: true,
    cell: (props) => (
      <LastTouchpointDateCell
        lastTouchPointAt={
          props.row.original?.value?.lastTouchpoint?.lastTouchPointAt
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
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[75%] h-[14px]' />
        <Skeleton className='w-[100%] h-[14px]' />
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

export const getColumnSortFn = (columnId: string) =>
  match(columnId)
    .with(
      'ORGANIZATIONS_NAME',
      () => (row: Store<Organization>) => row.value?.name?.trim() || null,
    )
    .with(
      'ORGANIZATIONS_RELATIONSHIP',
      () => (row: Store<Organization>) => row.value?.isCustomer,
    )
    .with(
      'ORGANIZATIONS_ONBOARDING_STATUS',
      () => (row: Store<Organization>) =>
        match(row.value?.accountDetails?.onboarding?.status)
          .with(OnboardingStatus.NotApplicable, () => null)
          .with(OnboardingStatus.NotStarted, () => 1)
          .with(OnboardingStatus.OnTrack, () => 2)
          .with(OnboardingStatus.Late, () => 3)
          .with(OnboardingStatus.Stuck, () => 4)
          .with(OnboardingStatus.Successful, () => 5)
          .with(OnboardingStatus.Done, () => 6)
          .otherwise(() => null),
    )
    .with(
      'ORGANIZATIONS_RENEWAL_LIKELIHOOD',
      () => (row: Store<Organization>) =>
        match(row.value?.accountDetails?.renewalSummary?.renewalLikelihood)
          .with(OpportunityRenewalLikelihood.HighRenewal, () => 3)
          .with(OpportunityRenewalLikelihood.MediumRenewal, () => 2)
          .with(OpportunityRenewalLikelihood.LowRenewal, () => 1)
          .otherwise(() => null),
    )
    .with('ORGANIZATIONS_RENEWAL_DATE', () => (row: Store<Organization>) => {
      const value = row.value?.accountDetails?.renewalSummary?.nextRenewalDate;

      return value ? new Date(value) : null;
    })
    .with(
      'ORGANIZATIONS_FORECAST_ARR',
      () => (row: Store<Organization>) =>
        row.value?.accountDetails?.renewalSummary?.arrForecast,
    )
    .with('ORGANIZATIONS_OWNER', () => (row: Store<Organization>) => {
      const name = row.value?.owner?.name ?? '';
      const firstName = row.value?.owner?.firstName ?? '';
      const lastName = row.value?.owner?.lastName ?? '';

      const fullName = (name ?? `${firstName} ${lastName}`).trim();

      return fullName.length ? fullName : null;
    })
    .with(
      'ORGANIZATIONS_LEAD_SOURCE',
      () => (row: Store<Organization>) => row.value?.leadSource,
    )
    .with(
      'ORGANIZATIONS_CREATED_DATE',
      () => (row: Store<Organization>) =>
        row.value?.metadata?.created
          ? new Date(row.value?.metadata?.created)
          : null,
    )
    .with(
      'ORGANIZATIONS_YEAR_FOUNDED',
      () => (row: Store<Organization>) => row.value?.yearFounded,
    )
    .with(
      'ORGANIZATIONS_EMPLOYEE_COUNT',
      () => (row: Store<Organization>) => row.value?.employees,
    )
    .with(
      'ORGANIZATIONS_SOCIALS',
      () => (row: Store<Organization>) => row.value?.socialMedia?.[0]?.url,
    )
    .with('ORGANIZATIONS_LAST_TOUCHPOINT', () => (row: Store<Organization>) => {
      const value = row.value?.lastTouchpoint?.lastTouchPointAt;

      if (!value) return null;

      return new Date(value);
    })
    .with(
      'ORGANIZATIONS_LAST_TOUCHPOINT_DATE',
      () => (row: Store<Organization>) => {
        const value = row.value?.lastTouchpoint?.lastTouchPointAt;

        return value ? new Date(value) : null;
      },
    )
    .otherwise(() => (_row: Store<Organization>) => null);

export const getPredefinedFilterFn = (serverFilter: Filter | null) => {
  if (!serverFilter) return null;

  const data = serverFilter?.AND?.[0];

  return match(data?.filter)
    .with({ property: 'STAGE' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;

      if (!filterValues) return false;

      return filterValues.includes(row.value?.stage);
    })
    .with(
      { property: 'IS_CUSTOMER' },
      (filter) => (row: Store<Organization>) => {
        const filterValues = filter?.value;

        if (!filterValues) return false;

        return filterValues.includes(row.value?.isCustomer);
      },
    )
    .with({ property: 'OWNER_ID' }, (filter) => (row: Store<Organization>) => {
      const filterValues = filter?.value;

      if (!filterValues) return false;

      return filterValues.includes(row.value?.owner?.id);
    })
    .otherwise(() => null);
};
