import { Flex } from '@ui/layout/Flex';
import { Organization } from '@graphql/types';
import { Skeleton } from '@ui/presentation/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { TimeToRenewalCell } from '@organizations/components/Columns/Cells/renewal/TimeToRenewalCell';
import { RenewalLikelihoodCell } from '@organizations/components/Columns/Cells/renewal/RenewalLikelihoodCell';

import { AvatarHeader } from './Headers/Avatar';
import { OwnerCell } from './Cells/owner/OwnerCell';
import { AvatarCell } from './Cells/avatar/AvatarCell';
import { WebsiteCell } from './Cells/website/WebsiteCell';
import { OwnerFilter, filterOwnerFn } from './Filters/Owner';
import { WebsiteFilter, filterWebsiteFn } from './Filters/Website';
import { OrganizationCell } from './Cells/organization/OrganizationCell';
import { RenewalForecastCell } from './Cells/renewal/RenewalForecastCell';
import { LastTouchpointCell } from './Cells/touchpoint/LastTouchpointCell';
import {
  OrganizationFilter,
  filterOrganizationFn,
} from './Filters/Organization';
import {
  RelationshipFilter,
  filterRelationshipFn,
} from './Filters/Relationship';
import { OrganizationRelationship } from './Cells/relationship/OrganizationRelationship';
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
        <AvatarCell id={props.getValue()?.id} name={props.getValue()?.name} />
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
  columnHelper.accessor((row) => row, {
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
    skeleton: () => (
      <Skeleton
        width='50%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
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
  columnHelper.accessor('accountDetails', {
    id: 'RENEWAL_LIKELIHOOD',
    minSize: 200,
    enableColumnFilter: false,
    enableSorting: false,
    // filterFn: filterRenewalLikelihoodFn,
    cell: (props) => {
      const value = props.getValue()?.renewalSummary?.renewalLikelihood;

      return <RenewalLikelihoodCell value={value} />;
    },
    header: (props) => (
      <THead
        id='renewalLikelihood'
        title='Health'
        // renderFilter={() => <RenewalLikelihoodFilter column={props.column} />}
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
  columnHelper.accessor('accountDetails', {
    id: 'RENEWAL_CYCLE_NEXT',
    minSize: 200,
    enableColumnFilter: false,
    enableSorting: false,
    // filterFn: filterTimeToRenewalFn,
    cell: (props) => {
      const nextRenewalDate = props.getValue()?.renewalSummary?.nextRenewalDate;

      return <TimeToRenewalCell nextRenewalDate={nextRenewalDate} />;
    },

    header: (props) => (
      <THead
        id='timeToRenewal'
        title='Next Renewal'
        // renderFilter={() => (
        //   <TimeToRenewalFilter
        //     onFilterValueChange={props.column.setFilterValue}
        //   />
        // )}
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
  columnHelper.accessor('accountDetails', {
    id: 'FORECAST_AMOUNT',
    minSize: 200,
    enableColumnFilter: false,
    enableSorting: false,
    // filterFn: filterForecastFn,
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
        // filterWidth='17rem'
        // renderFilter={(initialFocusRef) => (
        //   <ForecastFilter
        //     initialFocusRef={initialFocusRef}
        //     onFilterValueChange={props.column.setFilterValue}
        //   />
        // )}
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
  columnHelper.accessor('owner', {
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
  columnHelper.accessor((row) => row, {
    id: 'LAST_TOUCHPOINT',
    minSize: 250,
    filterFn: filterLastTouchpointFn,
    // enableColumnFilter: false,
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
];
