import { match } from 'ts-pattern';
import { Store } from '@store/store';
import { isAfter } from 'date-fns/isAfter';
import { Filter, FilterItem } from '@store/types';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date.ts';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import {
  Organization,
  TableViewDef,
  ColumnViewType,
  OnboardingStatus,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import { AvatarHeader } from './Headers/Avatar';
import { LastTouchpointDateCell } from './Cells/touchpointDate';
import {
  OwnerCell,
  AvatarCell,
  WebsiteCell,
  IndustryCell,
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
  LtvFilter,
  OwnerFilter,
  SourceFilter,
  WebsiteFilter,
  ChurnedFilter,
  SocialsFilter,
  ForecastFilter,
  IndustryFilter,
  EmployeesFilter,
  OnboardingFilter,
  CreatedDateFilter,
  OrganizationFilter,
  RelationshipFilter,
  TimeToRenewalFilter,
  LastTouchpointFilter,
  LastInteractedFilter,
  RenewalLikelihoodFilter,
} from './Filters';

type ColumnDatum = Store<Organization>;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  [ColumnViewType.OrganizationsAvatar]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.OrganizationsAvatar,
    size: 26,
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
  [ColumnViewType.OrganizationsName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.OrganizationsName,
    size: 150,
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
        id={ColumnViewType.OrganizationsName}
        title='Organization'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <OrganizationFilter initialFocusRef={initialFocusRef} />
        )}
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsWebsite]: columnHelper.accessor(
    'value.website',
    {
      id: ColumnViewType.OrganizationsWebsite,
      size: 125,
      enableSorting: false,
      cell: (props) => {
        const organizationId = props.row.original.value.metadata.id;

        return <WebsiteCell organizationId={organizationId} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsWebsite}
          title='Website'
          filterWidth='14rem'
          renderFilter={(initialFocusRef) => (
            <WebsiteFilter initialFocusRef={initialFocusRef} />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsRelationship]: columnHelper.accessor(
    'value.relationship',
    {
      id: ColumnViewType.OrganizationsRelationship,
      size: 125,
      header: (props) => (
        <THead
          id={ColumnViewType.OrganizationsRelationship}
          title='Relationship'
          renderFilter={() => <RelationshipFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      cell: (props) => {
        const id = props.row.original.value.metadata?.id;

        return <OrganizationRelationshipCell id={id} />;
      },
      skeleton: () => <Skeleton className='w-[100%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsOnboardingStatus]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsOnboardingStatus,
      size: 125,
      cell: (props) => {
        const status = props.getValue()?.onboarding?.status;
        const updatedAt = props.getValue()?.onboarding?.updatedAt;

        return <OnboardingCell status={status} updatedAt={updatedAt} />;
      },
      header: (props) => (
        <THead
          id={ColumnViewType.OrganizationsOnboardingStatus}
          title='Onboarding'
          renderFilter={() => <OnboardingFilter />}
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
  [ColumnViewType.OrganizationsRenewalLikelihood]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsRenewalLikelihood,
      size: 100,
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
          id={ColumnViewType.OrganizationsRenewalLikelihood}
          title='Health'
          renderFilter={() => <RenewalLikelihoodFilter />}
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
  [ColumnViewType.OrganizationsRenewalDate]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsRenewalDate,
      size: 125,
      cell: (props) => {
        const nextRenewalDate =
          props.getValue()?.renewalSummary?.nextRenewalDate;

        return <TimeToRenewalCell nextRenewalDate={nextRenewalDate} />;
      },

      header: (props) => (
        <THead
          id={ColumnViewType.OrganizationsRenewalDate}
          title='Next Renewal'
          renderFilter={() => <TimeToRenewalFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsForecastArr]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsForecastArr,
      size: 150,
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
          id={ColumnViewType.OrganizationsForecastArr}
          title='ARR Forecast'
          filterWidth='17rem'
          renderFilter={(initialFocusRef) => (
            <ForecastFilter initialFocusRef={initialFocusRef} />
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
    },
  ),
  [ColumnViewType.OrganizationsOwner]: columnHelper.accessor('value.owner', {
    id: ColumnViewType.OrganizationsOwner,
    size: 150,
    cell: (props) => {
      return (
        <OwnerCell
          id={props.row.original.value.metadata?.id}
          owner={props.getValue()}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.OrganizationsOwner}
        title='Owner'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <OwnerFilter initialFocusRef={initialFocusRef} />
        )}
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsLeadSource]: columnHelper.accessor(
    'value.leadSource',
    {
      id: ColumnViewType.OrganizationsLeadSource,
      size: 100,
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
          id={ColumnViewType.OrganizationsLeadSource}
          title='Source'
          renderFilter={() => <SourceFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsCreatedDate]: columnHelper.accessor(
    'value.metadata.created',
    {
      id: ColumnViewType.OrganizationsCreatedDate,
      size: 125,
      cell: (props) => {
        const value = props.getValue();

        if (!value) {
          return <p className='text-gray-400'>Unknown</p>;
        }

        return (
          <p className='text-gray-700 cursor-default truncate'>
            {DateTimeUtils.format(
              value,
              DateTimeUtils.defaultFormatShortString,
            )}
          </p>
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsCreatedDate}
          title='Created Date'
          filterWidth='14rem'
          renderFilter={() => <CreatedDateFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsYearFounded]: columnHelper.accessor(
    'value.yearFounded',
    {
      id: ColumnViewType.OrganizationsYearFounded,
      size: 100,
      enableColumnFilter: false,
      cell: (props) => {
        const value = props.getValue();

        if (!value) {
          return <p className='text-gray-400'>Unknown</p>;
        }

        return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsForecastArr}
          title='Founded'
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsEmployeeCount]: columnHelper.accessor(
    'value.employees',
    {
      id: ColumnViewType.OrganizationsEmployeeCount,
      size: 125,
      cell: (props) => {
        const value = props.getValue();

        if (!value) {
          return <p className='text-gray-400'>Unknown</p>;
        }

        return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsEmployeeCount}
          title='Employees'
          renderFilter={() => <EmployeesFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsSocials]: columnHelper.accessor(
    'value.socialMedia',
    {
      id: ColumnViewType.OrganizationsSocials,
      size: 125,
      enableSorting: false,
      cell: (props) => <LinkedInCell organizationId={props.row.original.id} />,
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsSocials}
          title='LinkedIn'
          filterWidth='14rem'
          renderFilter={(initialFocusRef) => (
            <SocialsFilter initialFocusRef={initialFocusRef} />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsLastTouchpoint]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.OrganizationsLastTouchpoint,
      size: 200,
      cell: (props) => (
        <LastTouchpointCell
          lastTouchPointAt={
            props.row.original?.value?.lastTouchpoint?.lastTouchPointAt
          }
          lastTouchPointTimelineEvent={
            props.row.original?.value?.lastTouchpoint
              ?.lastTouchPointTimelineEvent
          }
          lastTouchPointType={
            props.row.original?.value?.lastTouchpoint?.lastTouchPointType
          }
        />
      ),
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsLastTouchpoint}
          title='Last Touchpoint'
          renderFilter={() => <LastTouchpointFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[75%] h-[14px]' />
          <Skeleton className='w-[100%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OrganizationsLastTouchpointDate]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.OrganizationsLastTouchpointDate,
      size: 150,
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
          id={ColumnViewType.OrganizationsLastTouchpointDate}
          title='Last Interacted'
          renderFilter={() => <LastInteractedFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[75%] h-[14px]' />
          <Skeleton className='w-[100%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OrganizationsChurnDate]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsChurnDate,
      size: 115,
      cell: (props) => {
        const value = props.row.original.value.accountDetails?.churned;

        return (
          <p
            className={cn(
              'text-gray-700 cursor-default',
              !value && 'text-gray-400',
            )}
          >
            {DateTimeUtils.format(
              value,
              DateTimeUtils.defaultFormatShortString,
            ) || 'Unknown'}
          </p>
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsChurnDate}
          title='Churn Date'
          renderFilter={() => {
            return <ChurnedFilter />;
          }}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsLtv]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsLtv,
      size: 100,
      cell: (props) => {
        const value = props.row.original.value.accountDetails?.ltv;

        return (
          <p
            className={cn(
              'text-gray-700 cursor-default',
              !value && 'text-gray-400',
            )}
          >
            {value || 'Unknown'}
          </p>
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsLtv}
          title='LTV'
          filterWidth='14rem'
          renderFilter={(initialFocusRef) => {
            return <LtvFilter initialFocusRef={initialFocusRef} />;
          }}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsIndustry]: columnHelper.accessor(
    'value.industry',
    {
      id: ColumnViewType.OrganizationsIndustry,
      size: 200,
      cell: (props) => {
        const value = props.getValue();

        return <IndustryCell value={value} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsIndustry}
          title='Industry'
          filterWidth='auto'
          renderFilter={(initialFocusRef) => (
            <IndustryFilter initialFocusRef={initialFocusRef} />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
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
      () => (row: Store<Organization>) =>
        row.value?.name?.trim().toLocaleLowerCase() || null,
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

      return fullName.length ? fullName.toLocaleLowerCase() : null;
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
    .with('ORGANIZATIONS_CHURN_DATE', () => (row: Store<Organization>) => {
      const value = row.value?.accountDetails?.churned;

      return value ? new Date(value) : null;
    })
    .with(
      'ORGANIZATIONS_LTV',
      () => (row: Store<Organization>) => row.value?.accountDetails?.ltv,
    )
    .with(
      'ORGANIZATIONS_INDUSTRY',
      () => (row: Store<Organization>) => row.value?.industry,
    )
    .otherwise(() => (_row: Store<Organization>) => false);

export const getFilterFn = (filter: FilterItem | undefined | null) => {
  const noop = (_row: Store<Organization>) => true;
  if (!filter) return noop;

  return match(filter)
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

    .with(
      { property: 'RELATIONSHIP' },
      (filter) => (row: Store<Organization>) => {
        const filterValues = filter?.value;

        if (!filterValues) return false;

        return filterValues.includes(row.value?.relationship);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsCreatedDate },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return isAfter(
          new Date(row.value.metadata.created),
          new Date(filterValue),
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsName },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (filter.includeEmpty && row.value.name === 'Unnamed') {
          return true;
        }

        return row.value.name.toLowerCase().includes(filterValue.toLowerCase());
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsWebsite },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (filter.includeEmpty && !row.value.website) {
          return true;
        }

        return (
          row.value.website &&
          row.value.website.toLowerCase().includes(filterValue.toLowerCase())
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRelationship },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(row.value.relationship);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsForecastArr },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const forecastValue =
          row.value?.accountDetails?.renewalSummary?.arrForecast;

        if (!forecastValue) return false;

        return (
          forecastValue >= filterValue[0] && forecastValue <= filterValue[1]
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRenewalDate },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const nextRenewalDate =
          row.value?.accountDetails?.renewalSummary?.nextRenewalDate;

        if (!nextRenewalDate) return false;

        return isAfter(new Date(nextRenewalDate), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsOnboardingStatus },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(
          row.value.accountDetails?.onboarding?.status,
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsRenewalLikelihood },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(
          row.value.accountDetails?.renewalSummary?.renewalLikelihood,
        );
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsOwner },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (filterValue === '__EMPTY__' && !row.value.owner) {
          return true;
        }

        return filterValue.includes(row.value.owner?.id);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLastTouchpoint },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const lastTouchpoint = row?.value?.lastTouchpoint?.lastTouchPointType;
        const lastTouchpointAt = row?.value?.lastTouchpoint?.lastTouchPointAt;

        const isIncluded = filterValue?.types.length
          ? filterValue?.types?.includes(lastTouchpoint)
          : false;

        const isAfterDate = isAfter(
          new Date(lastTouchpointAt),
          new Date(filterValue?.after),
        );

        return isIncluded && isAfterDate;
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsChurnDate },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const churned = row?.value?.accountDetails?.churned;

        if (!churned) return false;

        return isAfter(new Date(churned), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsSocials },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        // specific logic for linkedin
        const linkedInUrl = row.value.socialMedia?.find((v) =>
          v.url.includes('linkedin'),
        )?.url;

        if (!linkedInUrl && filter.includeEmpty) return true;

        return linkedInUrl && linkedInUrl.includes(filterValue);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLastTouchpointDate },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const lastTouchpointAt = row?.value?.lastTouchpoint?.lastTouchPointAt;

        return isAfter(new Date(lastTouchpointAt), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsEmployeeCount },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value.split('-').map(Number) as number[];
        const employees = row.value.employees;

        if (filterValue.length !== 2) return employees >= filterValue[0];

        return employees >= filterValue[0] && employees <= filterValue[1];
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLeadSource },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(row.value.leadSource);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsIndustry },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(row.value.industry);
      },
    )
    .with(
      { property: ColumnViewType.OrganizationsLtv },
      (filter) => (row: Store<Organization>) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const ltv = row.value.accountDetails?.ltv;

        if (!ltv) return false;

        if (filterValue.length !== 2) return ltv >= filterValue[0];

        return ltv >= filterValue[0] && ltv <= filterValue[1];
      },
    )
    .otherwise(() => noop);
};

export const getAllFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};
