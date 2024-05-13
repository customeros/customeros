import { Organization } from '@graphql/types';
import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { OrganizationRelationship } from '@organizations/components/Columns/Cells';
import { TimeToRenewalCell } from '@organizations/components/Columns/Cells/renewal/TimeToRenewalCell';
import {
  ForecastFilter,
  filterForecastFn,
} from '@organizations/components/Columns/Filters/Forecast';
import {
  TimeToRenewalFilter,
  filterTimeToRenewalFn,
} from '@organizations/components/Columns/Filters/TimeToRenewal';

import { AvatarHeader } from './Headers/Avatar';
import { OwnerCell } from './Cells/owner/OwnerCell';
import { AvatarCell } from './Cells/avatar/AvatarCell';
import { WebsiteCell } from './Cells/website/WebsiteCell';
import { OwnerFilter, filterOwnerFn } from './Filters/Owner';
import { WebsiteFilter, filterWebsiteFn } from './Filters/Website';
import { OnboardingCell } from './Cells/onboarding/OnboardingCell';
import { OrganizationCell } from './Cells/organization/OrganizationCell';
import { RenewalForecastCell } from './Cells/renewal/RenewalForecastCell';
import { LastTouchpointCell } from './Cells/touchpoint/LastTouchpointCell';
import { OnboardingFilter, filterOnboardingFn } from './Filters/Onboarding';
import {
  OrganizationFilter,
  filterOrganizationFn,
} from './Filters/Organization';
import {
  RelationshipFilter,
  filterRelationshipFn,
} from './Filters/Relationship';
import {
  LastTouchpointFilter,
  filterLastTouchpointFn,
} from './Filters/LastTouchpoint';

const columnHelper =
  createColumnHelper<Omit<Organization, 'lastTouchPointTimelineEvent'>>();

export const columns = [
  columnHelper.accessor((row) => row, {
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
          src={props.getValue()?.logo}
        />
      );
    },
    header: AvatarHeader,
    skeleton: () => <Skeleton className='size-[42px]' />,
  }),
  columnHelper.accessor((row) => row, {
    id: 'NAME',
    minSize: 200,
    filterFn: filterOrganizationFn,
    cell: (props) => {
      return (
        <OrganizationCell
          id={props.getValue().metadata.id}
          name={props.getValue().name}
          isSubsidiary={!!props.getValue()?.parentCompanies?.length}
          parentOrganizationName={
            props.getValue()?.parentCompanies?.[0]?.organization.name
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
  columnHelper.accessor('website', {
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
  columnHelper.accessor('isCustomer', {
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

      // @ts-expect-error - fixme
      return <OrganizationRelationship organization={organization} />;
    },
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
  columnHelper.accessor('accountDetails', {
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
  columnHelper.accessor('accountDetails', {
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
  columnHelper.accessor('accountDetails', {
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
  columnHelper.accessor('accountDetails', {
    id: 'FORECAST_ARR',
    minSize: 200,
    filterFn: filterForecastFn,
    cell: (props) => {
      const value = props.getValue()?.renewalSummary;
      const amount = value?.arrForecast;
      const potentialAmount = value?.maxArrForecast;

      return (
        <RenewalForecastCell
          id={props.row.original.metadata.id}
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
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[50%] h-[18px]' />
        <Skeleton className='w-[25%] h-[18px]' />
      </div>
    ),
  }),
  columnHelper.accessor('owner', {
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
  columnHelper.accessor((row) => row, {
    id: 'LAST_TOUCHPOINT',
    minSize: 250,
    filterFn: filterLastTouchpointFn,
    cell: (props) => (
      <LastTouchpointCell
        lastTouchPointAt={props.row.original.lastTouchpoint?.lastTouchPointAt}
        lastTouchPointTimelineEvent={
          (props.row.original as Organization).lastTouchpoint
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
];
