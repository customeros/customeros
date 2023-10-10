import { Flex } from '@ui/layout/Flex';
import { Plus } from '@ui/media/icons/Plus';
import { Tooltip } from '@ui/overlay/Tooltip';
import { Organization } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { Skeleton } from '@ui/presentation/Skeleton';
import { THead, createColumnHelper } from '@ui/presentation/Table';

import { OwnerCell } from './Cells/owner/OwnerCell';
import { WebsiteCell } from './Cells/website/WebsiteCell';
import { RelationshipStage } from './Cells/stage/RelationshipStage';
import { TimeToRenewalCell } from './Cells/renewal/TimeToRenewalCell';
import { OrganizationCell } from './Cells/organization/OrganizationCell';
import { RenewalForecastCell } from './Cells/renewal/RenewalForecastCell';
import { LastTouchpointCell } from './Cells/touchpoint/LastTouchpointCell';
import { RenewalLikelihoodCell } from './Cells/renewal/RenewalLikelihoodCell';
import { OrganizationRelationship } from './Cells/relationship/OrganizationRelationship';

const columnHelper =
  createColumnHelper<Omit<Organization, 'lastTouchPointTimelineEvent'>>();

interface GetColumnsOptions {
  tabs?: { [key: string]: string } | null;
  createIsLoading?: boolean;
  onCreateOrganization?: () => void;
}

export const getColumns = (options: GetColumnsOptions) => [
  columnHelper.accessor((row) => row, {
    id: 'ORGANIZATION',
    cell: (props) => {
      return (
        <OrganizationCell
          key={props.getValue().id}
          organization={props.getValue()}
          lastPositionParams={options?.tabs?.[props.getValue()?.id]}
        />
      );
    },
    minSize: 200,
    header: (props) => (
      <THead<Organization>
        title='Company'
        icon={
          <Flex w='10' h='10' align='center' justify='center' mr='3'>
            <Tooltip label='Create a company'>
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
        }
        {...props}
      />
    ),
    skeleton: () => (
      <Flex align='center' h='full'>
        <Skeleton
          borderRadius='lg'
          w='40px'
          h='40px'
          startColor='gray.300'
          endColor='gray.300'
        />
        <Flex ml='3' flexDir='column' h='42px' align='center' gap='1'>
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
      </Flex>
    ),
  }),
  columnHelper.accessor('website', {
    id: 'WEBSITE',
    minSize: 200,
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
    header: (props) => <THead<Organization> title='Relationship' {...props} />,
    minSize: 200,
    cell: (props) => {
      const isCustomer = props.getValue();
      const organizationId = props.row.original.id;
      const organization = props.row.original;
      const organizationName = props.row.original.name;

      return (
        <>
          <OrganizationRelationship
            organizationId={organizationId}
            isCustomer={isCustomer ?? false}
            organization={organization}
            organizationName={organizationName}
          />
        </>
      );
    },
    skeleton: () => (
      <Flex gap='1' flexDir='column'>
        <Skeleton
          width='100%'
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
  columnHelper.accessor('accountDetails', {
    id: 'RENEWAL_LIKELIHOOD',
    minSize: 200,
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
    header: (props) => (
      <THead<Organization> title='Renewal Forecast' {...props} />
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
