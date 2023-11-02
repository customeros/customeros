import { Flex } from '@ui/layout/Flex';
import { Plus } from '@ui/media/icons/Plus';
import { Tooltip } from '@ui/overlay/Tooltip';
import { Organization } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { Skeleton } from '@ui/presentation/Skeleton';
import { THead, createColumnHelper } from '@ui/presentation/Table';

import { OwnerCell } from './Cells/owner/OwnerCell';
import { AvatarCell } from './Cells/avatar/AvatarCell';
import { WebsiteCell } from './Cells/website/WebsiteCell';
import { TimeToRenewalCell } from './Cells/renewal/TimeToRenewalCell';
import { OrganizationCell } from './Cells/organization/OrganizationCell';
import { RenewalForecastCell } from './Cells/renewal/RenewalForecastCell';
import { LastTouchpointCell } from './Cells/touchpoint/LastTouchpointCell';
import { RenewalLikelihoodCell } from './Cells/renewal/RenewalLikelihoodCell';
import { OrganizationRelationship } from './Cells/relationship/OrganizationRelationship';

const columnHelper =
  createColumnHelper<Omit<Organization, 'lastTouchPointTimelineEvent'>>();

interface GetColumnsOptions {
  createIsLoading?: boolean;
  onCreateOrganization?: () => void;
  tabs?: { [key: string]: string } | null;
}

export const getColumns = (options: GetColumnsOptions) => [
  columnHelper.accessor((row) => row, {
    id: 'AVATAR',
    minSize: 42,
    maxSize: 70,
    fixWidth: true,
    enableColumnFilter: false,
    cell: (props) => {
      return (
        <AvatarCell
          organization={props.getValue()}
          lastPositionParams={options?.tabs?.[props.getValue()?.id]}
        />
      );
    },
    header: (props) => {
      return (
        <Flex w='42px' align='center' justify='center'>
          <Tooltip label='Create an organization'>
            <IconButton
              size='sm'
              variant='ghost'
              aria-label='create organization'
              isLoading={options?.createIsLoading}
              onClick={options?.onCreateOrganization}
              icon={<Plus color='gray.400' boxSize='5' />}
            />
          </Tooltip>
        </Flex>
      );
    },
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
    id: 'ORGANIZATION',
    minSize: 200,
    enableColumnFilter: false,
    cell: (props) => {
      return (
        <OrganizationCell
          key={props.getValue().id}
          organization={props.getValue()}
          lastPositionParams={options?.tabs?.[props.getValue()?.id]}
        />
      );
    },
    header: (props) => <THead<Organization> title='Organization' {...props} />,
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
    enableColumnFilter: false,
    cell: (props) => <WebsiteCell website={props.getValue()} />,
    header: (props) => <THead<Organization> title='Website' {...props} />,
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
    filterFn: 'equalsString',
    header: (props) => (
      <THead<Organization>
        title='Relationship'
        renderFilter={(column) => {
          return <p>Hello World</p>;
        }}
        {...props}
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
    cell: (props) => {
      const organizationId = props.row.original.id;
      const value = props.getValue()?.renewalLikelihood;
      const currentProbability = value?.probability;
      const previousProbability = value?.previousProbability;
      const updatedAt = value?.updatedAt;

      return (
        <RenewalLikelihoodCell
          updatedAt={updatedAt}
          organizationId={organizationId}
          currentProbability={currentProbability}
          previousProbability={previousProbability}
        />
      );
    },
    header: (props) => (
      <THead<Organization> title='Renewal Likelihood' {...props} />
    ),
    skeleton: () => (
      <Flex flexDir='column' gap='1'>
        <Skeleton
          width='25%'
          height='18px'
          startColor='gray.300'
          endColor='gray.300'
        />
        <Skeleton
          width='75%'
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
    cell: (props) => {
      const values = props.getValue()?.billingDetails;
      const renewalDate = values?.renewalCycleNext;
      const renewalFrequency = values?.renewalCycle;

      return (
        <TimeToRenewalCell
          renewalDate={renewalDate}
          renewalFrequency={renewalFrequency}
        />
      );
    },
    header: (props) => (
      <THead<Organization> title='Time to Renewal' {...props} />
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
    cell: (props) => {
      const value = props.getValue()?.renewalForecast;
      const amount = value?.amount;
      const potentialAmount = value?.potentialAmount;

      return (
        <RenewalForecastCell
          amount={amount}
          potentialAmount={potentialAmount}
          isUpdatedByUser={!!value?.updatedById}
        />
      );
    },
    header: (props) => <THead<Organization> title='ARR Forecast' {...props} />,
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
    enableColumnFilter: false,
    cell: (props) => (
      <OwnerCell id={props.row.original.id} owner={props.getValue()} />
    ),
    header: (props) => <THead<Organization> title='Owner' {...props} />,
    skeleton: () => (
      <Skeleton
        width='75%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
  }),
  columnHelper.accessor('market', {
    id: 'LAST_TOUCHPOINT',
    minSize: 250,
    enableColumnFilter: false,
    cell: (props) => (
      <LastTouchpointCell
        lastTouchPointAt={props.row.original.lastTouchPointAt}
        lastTouchPointTimelineEvent={
          (props.row.original as Organization).lastTouchPointTimelineEvent
        }
      />
    ),
    header: (props) => (
      <THead<Organization> title='Last Touchpoint' {...props} />
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
