import { THead, createColumnHelper } from '@ui/presentation/Table';
import { OrganizationTableCell } from '@spaces/finder/finder-table';
import { Organization } from '@graphql/types';
import { OwnerTableCell } from '@spaces/finder/finder-table/OwnerTableCell';
import { LastTouchpointTableCell } from '@spaces/finder/finder-table/LastTouchpointTableCell';
import { Skeleton } from '@spaces/atoms/skeleton/Skeleton';

import { OrganizationRelationship } from '../organization-details/relationship/OrganizationRelationship';
import { RelationshipStage } from '../organization-details/stage/RelationshipStage';
import { RenewalLikelihoodCell } from '../organization-details/renewal/RenewalLikelihoodCell';
import { RenewalForecastCell } from '../organization-details/renewal/RenewalForecastCell';
import { TimeToRenewalCell } from '../organization-details/renewal/TimeToRenewalCell';

import styles from './organization-list.module.scss';

const columnHelper =
  createColumnHelper<Omit<Organization, 'lastTouchPointTimelineEvent'>>();

export const columns = [
  columnHelper.accessor((row) => row, {
    id: 'ORGANIZATION',
    cell: (props) => {
      return (
        <OrganizationTableCell
          key={props.getValue().id}
          organization={props.getValue()}
        />
      );
    },
    minSize: 200,
    header: (props) => <THead<Organization> title='Company' {...props} />,
    skeleton: () => <Skeleton width='100%' height='21px' />,
  }),
  columnHelper.accessor('relationshipStages', {
    id: 'RELATIONSHIP',
    header: (props) => (
      <THead<Organization> title='Relationship | Stage' {...props} />
    ),
    minSize: 200,
    cell: (props) => {
      const relationshipStages = props.getValue();
      const relationship = relationshipStages?.[0]?.relationship;
      const stage = relationshipStages?.[0]?.stage;
      const organizationId = props.row.original.id;

      return (
        <>
          <OrganizationRelationship
            defaultValue={relationship}
            organizationId={organizationId}
          />
          <RelationshipStage
            defaultValue={stage}
            relationship={relationship}
            organizationId={organizationId}
          />
        </>
      );
    },
    skeleton: () => (
      <div className={styles.skeletonWrapper}>
        <Skeleton width='100%' height='21px' />
        <Skeleton width='25%' height='21px' />
      </div>
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
    skeleton: () => <Skeleton width='100%' height='21px' />,
  }),
  columnHelper.accessor('accountDetails', {
    id: 'RENEWAL_CYCLE_NEXT',
    minSize: 200,
    cell: (props) => {
      const values = props.getValue()?.billingDetails;
      const renewalDate = values?.renewalCycleNext;
      const renewalFrequency = values?.frequency;

      return (
        <TimeToRenewalCell
          renewalDate={renewalDate}
          renewalFrequency={renewalFrequency}
        />
      );
    },
    header: (props) => (
      <THead<Organization> title='Time to renewal' {...props} />
    ),
    skeleton: () => <Skeleton width='100%' height='21px' />,
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
    skeleton: () => <Skeleton width='100%' height='21px' />,
  }),
  columnHelper.accessor('owner', {
    id: 'OWNER',
    minSize: 200,
    cell: (props) => (
      <OwnerTableCell
        owner={props.getValue()}
        organizationId={props.row.original.id}
      />
    ),
    header: (props) => <THead<Organization> title='Owner' {...props} />,
    skeleton: () => <Skeleton width='100%' height='21px' />,
  }),
  columnHelper.accessor('market', {
    id: 'LAST_TOUCHPOINT',
    minSize: 200,
    cell: (props) => (
      <LastTouchpointTableCell
        lastTouchPointAt={props.row.original.lastTouchPointAt}
        lastTouchPointTimelineEvent={
          (props.row.original as Organization).lastTouchPointTimelineEvent
        }
      />
    ),
    header: (props) => (
      <THead<Organization> title='Last touchpoint' {...props} />
    ),
    skeleton: () => <Skeleton width='100%' height='21px' />,
  }),
];
