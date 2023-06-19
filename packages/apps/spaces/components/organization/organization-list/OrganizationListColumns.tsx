import { createColumnHelper } from '@spaces/ui/presentation/Table/Table';
import { THead } from '@spaces/ui/presentation/Table/THead';
import {
  AddressTableCell,
  OrganizationTableCell,
} from '@spaces/finder/finder-table';
import { ExternalLinkCell } from '@spaces/atoms/table/table-cells/TableCell';
import { Organization } from '@spaces/graphql';
import { OwnerTableCell } from '@spaces/finder/finder-table/OwnerTableCell';
import { Skeleton } from '@spaces/atoms/skeleton';

import { OrganizationRelationship } from '../organization-details/relationship/OrganizationRelationship';
import { RelationshipStage } from '../organization-details/stage/RelationshipStage';

import styles from './organization-list.module.scss';

const columnHelper = createColumnHelper<Organization>();

export const columns = [
  columnHelper.accessor((row) => row, {
    id: 'ORGANIZATION',
    cell: (props) =>
      props.table.options?.meta?.isLoading ? (
        <Skeleton width='100%' height='21px' />
      ) : (
        <OrganizationTableCell organization={props.getValue()} />
      ),
    header: (props) => (
      <THead<Organization>
        title='Company'
        subTitle='Branch'
        columnHasIcon
        {...props}
      />
    ),
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
      if (props.table.options?.meta?.isLoading) {
        return (
          <div className={styles.skeletonWrapper}>
            <Skeleton width='100%' height='21px' />
            <Skeleton width='25%' height='21px' />
          </div>
        );
      }

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
  }),
  columnHelper.accessor('website', {
    id: 'DOMAIN',
    cell: (props) => {
      if (props.table.options?.meta?.isLoading) {
        return <Skeleton width='100%' height='21px' />;
      }

      const url = props.getValue();
      if (!url) return <div />;
      return <ExternalLinkCell url={url} />;
    },
    header: (props) => {
      return (
        <THead<Organization> title='Domain' subTitle='Website' {...props} />
      );
    },
  }),
  columnHelper.accessor('locations', {
    id: 'LOCATION',
    cell: (props) => {
      if (props.table.options?.meta?.isLoading) {
        return <Skeleton width='100%' height='21px' />;
      }
      return <AddressTableCell locations={props.getValue()} />;
    },
    header: (props) => (
      <THead<Organization> title='Location' subTitle='Address' {...props} />
    ),
  }),
  columnHelper.accessor('owner', {
    id: 'OWNER',
    cell: (props) => {
      if (props.table.options?.meta?.isLoading) {
        return <Skeleton width='100%' height='21px' />;
      }

      return (
        <OwnerTableCell
          owner={props.getValue()}
          organizationId={props.row.original.id}
        />
      );
    },
    header: (props) => <THead<Organization> title='Owner' {...props} />,
  }),
];
