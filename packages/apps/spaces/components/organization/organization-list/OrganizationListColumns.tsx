import { createColumnHelper } from '@spaces/ui/presentation/Table/Table';
import { THead } from '@spaces/ui/presentation/Table/THead';
import {
  AddressTableCell,
  OrganizationTableCell,
} from '@spaces/finder/finder-table';
import { ExternalLinkCell } from '@spaces/atoms/table/table-cells/TableCell';
import { Organization } from '@spaces/graphql';
import { OwnerTableCell } from '@spaces/finder/finder-table/OwnerTableCell';
import { Skeleton } from '@spaces/atoms/skeleton/Skeleton';

import { OrganizationRelationship } from '../organization-details/relationship/OrganizationRelationship';
import { RelationshipStage } from '../organization-details/stage/RelationshipStage';

import styles from './organization-list.module.scss';
import {LastTouchpointTableCell} from "@spaces/finder/finder-table/LastTouchpointTableCell";

const columnHelper = createColumnHelper<Organization>();

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
    header: (props) => (
      <THead<Organization>
        title='Company'
        subTitle='Branch'
        columnHasIcon
        {...props}
      />
    ),
    skeleton: () => <Skeleton width='100%' height='21px' />,
  }),
  columnHelper.accessor('relationshipStages', {
    id: 'RELATIONSHIP',
    header: (props) => (
      <THead<Organization>
        title='Relationship'
        subTitle='Stage'
        columnHasIcon
        {...props}
      />
    ),
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
  columnHelper.accessor('website', {
    id: 'DOMAIN',
    cell: (props) => {
      const url = props.getValue();
      if (!url) return <div />;
      return <ExternalLinkCell url={url} />;
    },
    header: (props) => (
      <THead<Organization> title='Domain' subTitle='Website' {...props} />
    ),
    skeleton: () => <Skeleton width='100%' height='21px' />,
  }),
  columnHelper.accessor('locations', {
    id: 'LOCATION',
    cell: (props) => {
      return <AddressTableCell locations={props.getValue()} />;
    },
    header: (props) => (
      <THead<Organization> title='Location' subTitle='Address' {...props} />
    ),
    skeleton: () => <Skeleton width='100%' height='21px' />,
  }),
  columnHelper.accessor('owner', {
    id: 'OWNER',
    cell: (props) => (
      <OwnerTableCell
        owner={props.getValue()}
        organizationId={props.row.original.id}
      />
    ),
    header: (props) => <THead<Organization> title='Owner' {...props} />,
    skeleton: () => <Skeleton width='100%' height='21px'  />,
  }),
  columnHelper.accessor('lastTouchPointAt', {
    id: 'LAST_TOUCHPOINT',
    cell: (props) => (
        <LastTouchpointTableCell
            lastTouchPointAt={props.row.original.lastTouchPointAt}
            lastTouchPointTimelineEvent={props.row.original.lastTouchPointTimelineEvent}
        />
    ),
    header: (props) => <THead<Organization> title='Last touchpoint' subTitle={'How long ago'} {...props} />,
    skeleton: () => <Skeleton width='100%' height='21px' />,
  }),

];
